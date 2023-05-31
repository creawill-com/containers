package driver

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Postgres struct {
	dbName, dbUser, dbPassword string
}

func NewPostgres() *Postgres {
	return &Postgres{
		dbName:     "test",
		dbUser:     "my-secret-user",
		dbPassword: "my-secret-pwd",
	}
}

func (p *Postgres) Type() Type {
	return TypePostgres
}

func (p *Postgres) GenericContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "postgres:15-alpine",
		Env: map[string]string{
			"POSTGRES_DB":       p.dbName,
			"POSTGRES_USER":     p.dbUser,
			"POSTGRES_PASSWORD": p.dbPassword,
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func (p *Postgres) Dsn(ctx context.Context, c testcontainers.Container) string {
	port, err := c.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		p.dbUser,
		p.dbPassword,
		p.dbName,
		port.Port(),
	)
}
