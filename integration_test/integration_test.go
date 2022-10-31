//go:build integration

package integration_test

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"github.com/subosito/gotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "github.com/romsar/antibrut/proto/antibrut/v1"
)

var grpcConn *grpc.ClientConn

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	network, err := pool.CreateNetwork("antibrut_tests")
	if err != nil {
		log.Fatalf("Cannot create Docker network: %s", err)
	}

	envMap, err := gotenv.Read("../.env.testing")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Cannot load .env.testing file: %s", err)
	}
	envSlice := make([]string, 0, len(envMap))
	for key, val := range envMap {
		envSlice = append(envSlice, key+"="+val)
	}

	resource, err := pool.BuildAndRunWithBuildOptions(
		&dockertest.BuildOptions{
			ContextDir: "../",
			Dockerfile: "Dockerfile",
		},
		&dockertest.RunOptions{
			Hostname: "antibrut_tests",
			Name:     "antibrut_tests",
			Networks: []*dockertest.Network{network},
			Env:      envSlice,
		},
	)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		grpcConn, err = grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
		return err
	}); err != nil {
		log.Fatalf("Could not connect to GRPC: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestSomething(t *testing.T) {
	ctx := context.Background()

	service := proto.NewAntiBrutServiceClient(grpcConn)

	resp, err := service.Check(ctx, &proto.CheckRequest{
		Login:    "rasarvarov",
		Password: "foobar",
		Ip:       "93.80.254.151",
	})
	require.NoError(t, err)
	require.True(t, resp.GetOk())
}
