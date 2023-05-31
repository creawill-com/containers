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

type Gorm struct {
	db *gorm.DB
	c  testcontainers.Container
}

func NewGorm(ctx context.Context, d Driver) *Gorm {
	c, err := d.GenericContainer(ctx)
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

	db, err := gorm.Open(resolveGormDialector(ctx, d, c), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	return &Gorm{
		db: db,
		c:  c,
	}
}

func (g *Gorm) DB() *gorm.DB {
	if g.db == nil {
		panic(errors.New("db is nil"))
	}

	return g.db
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

func resolveGormDialector(ctx context.Context, d Driver, c testcontainers.Container) gorm.Dialector {
	switch d.Type() {
	case driver.TypePostgres:
		return postgres.Open(d.Dsn(ctx, c))
	}

	panic(errors.New("can`t resolve dialector"))
}
