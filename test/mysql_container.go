package test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MysqlContainer struct {
	testcontainers.Container
	URI string
}

var (
	image        = "mysql:8"
	port         = "3306"
	dbName       = "test"
	rootPassword = "password"
)

func SetupMysql(ctx context.Context) (*MysqlContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{nat.Port(port).Port()},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": rootPassword,
			"MYSQL_DATABASE":      dbName,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog(fmt.Sprintf("port: %s  MySQL Community Server - GPL", port)),
			wait.ForListeningPort(nat.Port(port)),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("root:%s@tcp(%s:%s)/%s", rootPassword, hostIP, mappedPort.Port(), dbName)

	return &MysqlContainer{Container: container, URI: uri}, nil
}
