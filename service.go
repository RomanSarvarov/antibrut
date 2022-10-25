package antibrut

import (
	"context"
	"net"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/romsar/antibrut/clock"
)

// ErrMaxAttemptsExceeded это ошибка, которая говорит о том,
// что количество попыток больше максимально-допустимого значения.
var ErrMaxAttemptsExceeded = errors.New("max attempts exceeded")

// ErrIPInBlackList это ошибка, которая говорит о том,
// что данный IP адрес в черном списке.
var ErrIPInBlackList = errors.New("ip address in black list")

type Service struct {
	repo repository
	rl   rateLimiter
	cfg  Config
}

type rateLimiter interface {
	Check(ctx context.Context, c LimitationCode, val string) error
	Reset(ctx context.Context, filter ResetFilter) error
}

type repository interface {
	FindIPRulesByIP(ctx context.Context, ip IP) ([]*IPRule, error)
	FindIPRuleBySubnet(ctx context.Context, subnet Subnet) (*IPRule, error)
	CreateIPRule(ctx context.Context, ipRule *IPRule) (*IPRule, error)
	UpdateIPRule(ctx context.Context, id IPRuleID, ipRule *IPRule) (*IPRule, error)
	DeleteIPRuleBySubnet(ctx context.Context, subnet Subnet) (int64, error)
}

type Config struct {
	PruneDuration clock.Duration
}

func NewService(repo repository, rl rateLimiter, cfg Config) *Service {
	return &Service{
		repo: repo,
		rl:   rl,
		cfg:  cfg,
	}
}

func (s *Service) Check(ctx context.Context, login Login, pass Password, ip IP) error {
	errGrp, ctx := errgroup.WithContext(ctx)

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

func (s *Service) AddIPToWhiteList(ctx context.Context, subnet Subnet) error {
	return s.createOrUpdateIPRule(ctx, WhiteList, subnet)
}

func (s *Service) DeleteIPFromWhiteList(ctx context.Context, subnet Subnet) error {
	_, err := s.repo.DeleteIPRuleBySubnet(ctx, subnet)
	return err
}

func (s *Service) AddIPToBlackList(ctx context.Context, subnet Subnet) error {
	return s.createOrUpdateIPRule(ctx, BlackList, subnet)
}

func (s *Service) DeleteIPFromBlackList(ctx context.Context, subnet Subnet) error {
	_, err := s.repo.DeleteIPRuleBySubnet(ctx, subnet)
	return err
}

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

func (s *Service) work(ctx context.Context) error {
	// Удалить неактуальные бакеты.
	if s.cfg.PruneDuration.ToDuration() > 0 {
		return s.rl.Reset(ctx, ResetFilter{
			DateTo: clock.Now().Add(-s.cfg.PruneDuration.ToDuration()),
		})
	}
	return nil
}

func (s *Service) createOrUpdateIPRule(ctx context.Context, t IPRuleType, subnet Subnet) error {
	rule, err := s.repo.FindIPRuleBySubnet(ctx, subnet)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return err
		}

		rule, err = s.repo.CreateIPRule(ctx, &IPRule{
			Type:   t,
			Subnet: subnet,
		})
		if err != nil {
			return err
		}

		return nil
	}

	if rule.Type != t {
		_, err = s.repo.UpdateIPRule(ctx, rule.ID, &IPRule{
			Type:   t,
			Subnet: subnet,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

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

type Login string

func (l Login) String() string {
	return string(l)
}

func (l Login) IsZero() bool {
	return l == ""
}

type Password string

func (p Password) String() string {
	return string(p)
}

func (p Password) IsZero() bool {
	return p == ""
}

type IP string

func (ip IP) String() string {
	return string(ip)
}

func (ip IP) IsZero() bool {
	return ip == ""
}

type Subnet string

func (s Subnet) String() string {
	return string(s)
}

func (s Subnet) Contains(ip IP) (bool, error) {
	_, ipNet, err := net.ParseCIDR(s.String())
	if err != nil {
		return false, err
	}

	return ipNet.Contains(net.ParseIP(ip.String())), nil
}
