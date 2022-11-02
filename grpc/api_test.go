package grpc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	proto "github.com/romsar/antibrut/proto/antibrut/v1"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	dockerContainer testcontainers.Container
	grpcConn        *grpc.ClientConn

	abClient proto.AntiBrutServiceClient
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}

func (s *APISuite) SetupSuite() {
	if testing.Short() {
		s.T().Skip("skipping e2e test")
	}
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../",
			Dockerfile: "Dockerfile",
		},
		Env: map[string]string{
			"ANTIBRUT_GRPC_ADDRESS":        ":9090",
			"ANTIBRUT_PRUNE_DURATION":      "1m",
			"ANTIBRUT_RATE_LIMITER_DRIVER": "inmem",
			"ANTIBRUT_SQLITE_DSN":          "file::memory:?cache=shared&_foreign_keys=on",
		},
		ExposedPorts: []string{"9090/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Migrations done"),
			wait.ForExposedPort(),
		),
		Name: "testantibrut",
	}

	var err error
	s.dockerContainer, err = testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	ip, err := s.dockerContainer.Host(s.ctx)
	s.Require().NoError(err)

	port, err := s.dockerContainer.MappedPort(s.ctx, "9090")
	s.Require().NoError(err)

	s.grpcConn, err = grpc.Dial(ip+":"+port.Port(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.abClient = proto.NewAntiBrutServiceClient(s.grpcConn)
}

func (s *APISuite) TearDownTest() {
	err := s.dockerContainer.Terminate(context.Background())
	s.NoError(err)
}

func (s *APISuite) TestCheck() {
	s.Run("limit exceeded", func() {
		check := func(req, good *proto.CheckRequest, max int) func() {
			return func() {
				f := func(req *proto.CheckRequest) (*proto.CheckResponse, error) {
					return s.abClient.Check(s.ctx, req)
				}

				for i := 0; i < max; i++ {
					resp, err := f(req)
					s.Require().NoError(err)
					s.Require().True(resp.GetOk())
				}

				resp, err := f(req)
				s.Require().NoError(err)
				s.Require().False(resp.GetOk())
			}
		}

		s.Run("login", check(&proto.CheckRequest{Login: "foo"}, &proto.CheckRequest{Login: "bar"}, 10))
		s.Run("password", check(&proto.CheckRequest{Password: "foo"}, &proto.CheckRequest{Login: "bar"}, 100))
		s.Run("ip", check(&proto.CheckRequest{Ip: "192.168.5.15"}, &proto.CheckRequest{Ip: "192.168.10.150"}, 1000))
	})

	s.Run("no data passed", func() {
		_, err := s.abClient.Check(s.ctx, &proto.CheckRequest{})
		s.Require().Equal(codes.InvalidArgument, status.Code(err))
	})
}

func (s *APISuite) TestReset() {
	s.Run("success", func() {
		reset := func(checkReq *proto.CheckRequest, resetReq *proto.ResetRequest) func() {
			return func() {
				check := func() (*proto.CheckResponse, error) {
					return s.abClient.Check(s.ctx, checkReq)
				}

				i := 0
				for {
					resp, err := check()
					s.Require().NoError(err)

					if resp.GetOk() == false {
						break
					}

					i++

					if i >= 10000 {
						s.Fail("cannot reach max limit")
					}
				}

				_, err := s.abClient.Reset(s.ctx, resetReq)
				s.Require().NoError(err)

				resp, err := check()
				s.Require().NoError(err)
				s.Require().True(resp.GetOk())
			}
		}

		s.Run(
			"login",
			reset(
				&proto.CheckRequest{Login: "foo"},
				&proto.ResetRequest{Login: "foo"},
			),
		)

		s.Run(
			"ip",
			reset(
				&proto.CheckRequest{Ip: "192.168.5.15"},
				&proto.ResetRequest{Ip: "192.168.5.15"},
			),
		)
	})

	s.Run("no data passed", func() {
		_, err := s.abClient.Reset(s.ctx, &proto.ResetRequest{})
		s.Require().Equal(codes.InvalidArgument, status.Code(err))
	})
}

func (s *APISuite) TestAddIPToWhiteList() {
	s.Run("success", func() {
		check := func(req *proto.CheckRequest) (*proto.CheckResponse, error) {
			return s.abClient.Check(s.ctx, req)
		}

		i := 0
		for {
			resp, err := check(&proto.CheckRequest{
				Ip: "192.168.5.15",
			})
			s.Require().NoError(err)

			if resp.GetOk() == false {
				break
			}

			i++

			if i >= 10000 {
				s.Fail("cannot reach max limit")
			}
		}

		// wrong subnet
		_, err := s.abClient.AddIPToWhiteList(s.ctx, &proto.AddIPToWhiteListRequest{
			Subnet: "192.168.6.0/26",
		})
		s.Require().NoError(err)

		resp, err := check(&proto.CheckRequest{
			Ip: "192.168.5.15",
		})
		s.Require().NoError(err)
		s.Require().False(resp.GetOk())

		// needle subnet
		_, err = s.abClient.AddIPToWhiteList(s.ctx, &proto.AddIPToWhiteListRequest{
			Subnet: "192.168.5.0/26",
		})
		s.Require().NoError(err)

		resp, err = check(&proto.CheckRequest{
			Ip: "192.168.5.15",
		})
		s.Require().NoError(err)
		s.Require().True(resp.GetOk())

	})

	s.Run("no data passed", func() {
		_, err := s.abClient.AddIPToWhiteList(s.ctx, &proto.AddIPToWhiteListRequest{})
		s.Require().Equal(codes.InvalidArgument, status.Code(err))
	})
}
