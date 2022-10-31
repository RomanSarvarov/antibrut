package antibrut

import (
	"context"
	"errors"
	"net"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/romsar/antibrut/clock"
)

// ErrMaxAttemptsExceeded это ошибка, которая говорит о том,
// что количество попыток больше максимально-допустимого значения.
var ErrMaxAttemptsExceeded = errors.New("max attempts exceeded")

// ErrIPInBlackList это ошибка, которая говорит о том,
// что данный IP адрес в черном списке.
var ErrIPInBlackList = errors.New("ip address in black list")

// Service предоставляет API для работы с Анти-брутфорсом.
// Он не содержит в себе алгоритм тротлинга,
// а принимает его как зависимость.
type Service struct {
	// repo предоставляет доступ к хранилищу.
	repo repository

	// rl содержит алгоритм для тротлинга.
	rl rateLimiter

	// pruneDuration количество времени, после которого
	// очищать неактуальные записи в хранилищах.
	pruneDuration clock.Duration
}

// rateLimiter содержит алгоритм для тротлинга
// и предоставляет API для работы с ним.
type rateLimiter interface {
	// Check проверяет "хороший" ли запрос или его стоить заблокировать.
	Check(ctx context.Context, c LimitationCode, val string) error

	// Reset удаляет бакеты из хранилища.
	Reset(ctx context.Context, filter ResetFilter) error
}

// repository предоставляет доступ к хранилищу.
type repository interface {
	// FindIPRulesByIP находит особые правила для IP адреса на основе IP адреса.
	FindIPRulesByIP(ctx context.Context, ip IP) ([]*IPRule, error)

	// FindIPRuleBySubnet находит особые правила для IP адреса на основе подсети.
	FindIPRuleBySubnet(ctx context.Context, subnet Subnet) (*IPRule, error)

	// CreateIPRule создает особое правило для IP адреса.
	CreateIPRule(ctx context.Context, ipRule *IPRule) (*IPRule, error)

	// UpdateIPRule обновляет особое правило для IP адреса.
	UpdateIPRule(ctx context.Context, id IPRuleID, upd *IPRuleUpdate) (*IPRule, error)

	// DeleteIPRules удаляет особые правила для IP адресов.
	DeleteIPRules(ctx context.Context, filter IPRuleFilter) (int64, error)
}

// Option возвращает функцию, модифицирующую Service.
type Option func(s *Service)

// WithPruneDuration возвращает функцию, которая
// устанавливает pruneDuration у Service.
func WithPruneDuration(d clock.Duration) Option {
	return func(s *Service) {
		s.pruneDuration = d
	}
}

// NewService создает Service.
func NewService(repo repository, rl rateLimiter, opts ...Option) *Service {
	s := &Service{
		repo: repo,
		rl:   rl,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Check проверяет "хороший" ли запрос или его стоить заблокировать.
// Первым делом ищет особое правило для IP адреса (проверка на черный/белый список).
// Если особое правило не найдено - выполнит проверку через алгоритм.
func (s *Service) Check(ctx context.Context, login Login, pass Password, ip IP) error {
	errGrp, ctx := errgroup.WithContext(ctx)

	// white/black list check
	inWl, inBl, err := s.checkWhiteBlackList(ctx, ip)
	if err != nil {
		return err
	}
	if inWl {
		return nil
	}
	if inBl {
		return ErrIPInBlackList
	}

	// algo checks
	errGrp.Go(func() error {
		if login.IsZero() {
			return nil
		}

		return s.rl.Check(ctx, LoginLimitation, login.String())
	})

	errGrp.Go(func() error {
		if pass.IsZero() {
			return nil
		}

		return s.rl.Check(ctx, PasswordLimitation, pass.String())
	})

	errGrp.Go(func() error {
		if ip.IsZero() {
			return nil
		}

		return s.rl.Check(ctx, IPLimitation, ip.String())
	})

	if err := errGrp.Wait(); err != nil {
		return err
	}

	return nil
}

// Reset удаляет бакеты из хранилища.
func (s *Service) Reset(ctx context.Context, login Login, ip IP) error {
	errGrp, ctx := errgroup.WithContext(ctx)

	errGrp.Go(func() error {
		if login.IsZero() {
			return nil
		}

		return s.rl.Reset(ctx, ResetFilter{
			LimitationCode: LoginLimitation,
			Value:          login.String(),
		})
	})

	errGrp.Go(func() error {
		if ip.IsZero() {
			return nil
		}

		return s.rl.Reset(ctx, ResetFilter{
			LimitationCode: IPLimitation,
			Value:          ip.String(),
		})
	})

	if err := errGrp.Wait(); err != nil {
		return err
	}

	return nil
}

// AddIPToWhiteList добавляет IP адрес в белый список.
func (s *Service) AddIPToWhiteList(ctx context.Context, subnet Subnet) error {
	return s.createOrUpdateIPRule(ctx, WhiteList, subnet)
}

// DeleteIPFromWhiteList удаляет IP адрес из белого списка.
func (s *Service) DeleteIPFromWhiteList(ctx context.Context, subnet Subnet) error {
	_, err := s.repo.DeleteIPRules(ctx, IPRuleFilter{
		Type:   WhiteList,
		Subnet: subnet,
	})
	return err
}

// AddIPToBlackList добавляет IP адрес в чёрный список.
func (s *Service) AddIPToBlackList(ctx context.Context, subnet Subnet) error {
	return s.createOrUpdateIPRule(ctx, BlackList, subnet)
}

// DeleteIPFromBlackList удаляет IP адрес из чёрного списка.
func (s *Service) DeleteIPFromBlackList(ctx context.Context, subnet Subnet) error {
	_, err := s.repo.DeleteIPRules(ctx, IPRuleFilter{
		Type:   BlackList,
		Subnet: subnet,
	})
	return err
}

// Work запускает бесконечный цикл с полезной работой,
// которую необходимо выполнять для корректной работы Service.
func (s *Service) Work(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	if err := s.work(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := s.work(ctx); err != nil {
				return err
			}
		}
	}
}

// work выполняет полезную работу, которую запускает Work.
func (s *Service) work(ctx context.Context) error {
	// Удалить неактуальные бакеты.
	if s.pruneDuration.ToDuration() > 0 {
		return s.rl.Reset(ctx, ResetFilter{
			CreatedAtTo: clock.Now().Add(-s.pruneDuration.ToDuration()),
		})
	}
	return nil
}

// createOrUpdateIPRule находит или обновляет особое правило для IP адресов.
func (s *Service) createOrUpdateIPRule(ctx context.Context, t IPRuleType, subnet Subnet) error {
	rule, err := s.repo.FindIPRuleBySubnet(ctx, subnet)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return err
		}

		_, err := s.repo.CreateIPRule(ctx, &IPRule{
			Type:   t,
			Subnet: subnet,
		})
		if err != nil {
			return err
		}

		return nil
	}

	if rule.Type != t {
		_, err = s.repo.UpdateIPRule(ctx, rule.ID, &IPRuleUpdate{
			Type:   t,
			Subnet: subnet,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// checkWhiteBlackList проверяет IP адрес на нахождение в белом/черном списках.
func (s *Service) checkWhiteBlackList(ctx context.Context, ip IP) (wl bool, bl bool, err error) {
	if ip.IsZero() {
		return false, false, nil
	}

	rules, err := s.repo.FindIPRulesByIP(ctx, ip)
	if err != nil {
		return false, false, nil
	}

	for _, rule := range rules {
		if ok, _ := rule.Subnet.Contains(ip); ok {
			if rule.IsWhiteList() {
				return true, false, nil
			}

			if rule.IsBlackList() {
				return false, true, nil
			}
		}
	}

	return false, false, nil
}

// Login это логин пользователя.
type Login string

// String это текстовое представление Login.
func (l Login) String() string {
	return string(l)
}

// IsZero это флаг, который говорит о том,
// что значение является zero-value.
func (l Login) IsZero() bool {
	return l == ""
}

// Password это пароль пользователя.
type Password string

// String это текстовое представление Password.
func (p Password) String() string {
	return string(p)
}

// IsZero это флаг, который говорит о том,
// что значение является zero-value.
func (p Password) IsZero() bool {
	return p == ""
}

// IP это IP адрес.
type IP string

// String это текстовое представление IP.
func (ip IP) String() string {
	return string(ip)
}

// IsZero это флаг, который говорит о том,
// что значение является zero-value.
func (ip IP) IsZero() bool {
	return ip == ""
}

// Subnet это подсеть.
type Subnet string

// String это текстовое представление Subnet.
func (s Subnet) String() string {
	return string(s)
}

// Contains проверяет, что подсеть адрес содержит IP адрес.
func (s Subnet) Contains(ip IP) (bool, error) {
	_, ipNet, err := net.ParseCIDR(s.String())
	if err != nil {
		return false, err
	}

	return ipNet.Contains(net.ParseIP(ip.String())), nil
}
