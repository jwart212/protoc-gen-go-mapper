package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jwart212/protoc-gen-go-mapper/examples/advanced/internal/postgres/sqlc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Postgres struct {
	dbUri string
	D     *pgxpool.Pool
	Q     *sqlc.Queries

	tracer trace.Tracer

	closers []func(ctx context.Context) (err error)
}

type OptsFunc func(pg *Postgres) (err error)

func WithTracer(name string) OptsFunc {
	return func(p *Postgres) (err error) {
		p.tracer = otel.Tracer(name)
		return
	}
}

func WithConnUri(uri string) OptsFunc {
	return func(p *Postgres) (err error) {
		p.dbUri = uri
		cfg, err := pgxpool.ParseConfig(uri)
		if err != nil {
			return err
		}
		cfg.ConnConfig.Tracer = otelpgx.NewTracer()
		cfg.MaxConns = 50
		cfg.MinConns = 10

		cfg.MaxConnLifetime = time.Hour
		cfg.MaxConnIdleTime = 30 * time.Minute

		cfg.HealthCheckPeriod = time.Minute

		p.D, err = pgxpool.NewWithConfig(
			context.Background(),
			cfg,
		)
		if err != nil {
			return
		}
		if err := p.D.Ping(context.Background()); err != nil {
			p.D.Close()
			return err
		}
		return
	}
}

func New(opts ...OptsFunc) (pg *Postgres, err error) {
	pg = &Postgres{}

	for _, opt := range opts {
		if err = opt(pg); err != nil {
			return
		}
	}
	if pg.D == nil {
		return nil, fmt.Errorf(
			"missing db connection",
		)
	}

	pg.Q = sqlc.New(pg.D)

	return pg, nil
}

func (p *Postgres) Close(ctx context.Context) (err error) {
	for _, closer := range p.closers {
		if err = closer(ctx); err != nil {
			return
		}
	}
	return
}
