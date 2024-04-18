package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/go-pg/pg/v10"
)

const _defaultMaxPoolSize = 2

var ErrNoRows = pg.ErrNoRows

type Postgres struct {
	maxPoolSize int
	DB          *pg.DB
}

func New(url string, opts ...Option) (*Postgres, error) {
	res := &Postgres{
		maxPoolSize: _defaultMaxPoolSize,
	}

	for _, opt := range opts {
		opt(res)
	}
	// To connect to a database
	opt, err := pg.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pg.ParseURL: %w", err)
	}

	opt.PoolSize = res.maxPoolSize

	db := pg.Connect(opt)
	// To check if database is up and running
	ctx := context.Background()
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgres - New - db.Ping: %w", err)
	}

	res.DB = db

	return res, nil
}

func (pg *Postgres) Migrate(filePath string) error {
	c, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("postgres - Migrate - os.ReadFile: %w", err)
	}

	_, err = pg.DB.Exec(string(c))
	if err != nil {
		return fmt.Errorf("postgres - Migrate - pg.DB.Exec: %w", err)
	}

	return nil
}

func (pg *Postgres) Close() error {
	return fmt.Errorf("postgres - Close - pg.DB.Close: %w", pg.DB.Close())
}
