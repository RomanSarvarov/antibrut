package antibrut

import (
	"context"
	"net"

	"golang.org/x/sync/errgroup"
)

type Service struct {
	repo repository
	rl   rateLimiter
}

func NewService(repo repository, rl rateLimiter) *Service {
	return &Service{
		repo: repo,
		rl:   rl,
	}
}

type rateLimiter interface {
	Check(ctx context.Context, c LimitationCode, val string) error
}

type repository interface {
	FindIPRulesByIP(ctx context.Context, ip IP) ([]*IPRule, error)
	DeleteBuckets(ctx context.Context, filter BucketFilter) (int64, error)
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

		_, err := s.repo.DeleteBuckets(ctx, BucketFilter{
			LimitationCode: LoginLimitation,
			Value:          login.String(),
		})
		return err
	})

	errGrp.Go(func() error {
		if ip.IsZero() {
			return nil
		}

		_, err := s.repo.DeleteBuckets(ctx, BucketFilter{
			LimitationCode: IPLimitation,
			Value:          ip.String(),
		})
		return err
	})

	if err := errGrp.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *Service) AddIPToWhiteList(ctx context.Context, subnet Subnet) error {
	return nil
}

func (s *Service) DeleteIPFromWhiteList(ctx context.Context, subnet Subnet) error {
	return nil
}

func (s *Service) AddIPToBlackList(ctx context.Context, subnet Subnet) error {
	return nil
}

func (s *Service) DeleteIPFromBlackList(ctx context.Context, subnet Subnet) error {
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
