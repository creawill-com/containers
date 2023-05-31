package containers

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/creawill-com/containers/driver"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGorm(d Driver) *Gorm {
	return &Gorm{
		driver: d,
	}
}

type Gorm struct {
	driver Driver

	c testcontainers.Container
}

func (g *Gorm) DB(ctx context.Context) *gorm.DB {
	var err error

	g.c, err = g.driver.GenericContainer(ctx)
	if err != nil {
		panic(err)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
		},
	)

	db, err := gorm.Open(g.resolveDialector(ctx, g.c), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	return db
}

func (g *Gorm) Terminate(ctx context.Context) {
	if g.c == nil {
		panic(errors.New("container is nil"))
	}

	err := g.c.Terminate(ctx)
	if err != nil {
		panic(err)
	}
}

func (g *Gorm) resolveDialector(ctx context.Context, c testcontainers.Container) gorm.Dialector {
	switch g.driver.Type() {
	case driver.TypePostgres:
		return postgres.Open(g.driver.Dsn(ctx, c))
	}

	panic(errors.New("can`t resolve dialector"))
}
