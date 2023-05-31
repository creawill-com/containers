package driver

import (
	"context"
	"fmt"

	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Mysql struct {
	dbName, dbPassword string
}

func NewMysql() *Postgres {
	return &Postgres{
		dbName:     "test",
		dbPassword: "my-secret-pwd",
	}
}

func (m *Mysql) Type() Type {
	return TypeMysql
}

func (m *Mysql) GenericContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "mysql:8.0.27",
		Env: map[string]string{
			"MYSQL_DATABASE": m.dbName, "MYSQL_ROOT_PASSWORD": m.dbPassword},
		ExposedPorts: []string{"3306/tcp"},
		WaitingFor:   wait.ForListeningPort("3306/tcp"),
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func (m *Mysql) Dsn(ctx context.Context, c testcontainers.Container) string {
	port, err := c.MappedPort(ctx, "3306")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(
		"root:%s@(localhost:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.dbPassword,
		port.Port(),
		m.dbName,
	)
}
