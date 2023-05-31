package containers

import (
	"context"

	"github.com/creawill-com/containers/driver"
	testcontainers "github.com/testcontainers/testcontainers-go"
)

type Driver interface {
	Type() driver.Type
	GenericContainer(context.Context) (testcontainers.Container, error)
	Dsn(context.Context, testcontainers.Container) string
}
