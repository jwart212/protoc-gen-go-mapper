package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/examples/advanced/internal/gen/item_categoriespb"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/examples/advanced/internal/grpcserver/handler"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/examples/advanced/internal/postgres"
	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/examples/advanced/pkg/env"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	GRPCAddr       string `env:"GRPC_LISTEN_ADDRESS" envDefault:":50051"`
	PostgresURL    string `env:"POSTGRES_URL" envDefault:"postgres://postgres:ikt%402025!@localhost:5432/dba_pos?sslmode=disable"`
	ServiceName    string `env:"SERVICE_NAME" envDefault:"pos-service"`
	ServiceVersion string `env:"SERVICE_VERSION" envDefault:"1.0.0"`
	Environment    string `env:"ENVIRONMENT" envDefault:"development"`
}

type CMD struct {
	envPrefix string
	dotenv    bool

	cfg *Config

	grpcServer *grpc.Server

	db *postgres.Postgres

	closers []func(ctx context.Context) error
}

type OptsFunc func(c *CMD) (err error)

func WithPrefix(p string) OptsFunc {
	return func(c *CMD) (err error) {
		c.envPrefix = p
		return
	}
}

func WithoutDotEnv() OptsFunc {
	return func(c *CMD) (err error) {
		c.dotenv = false
		return
	}
}

func New(opts ...OptsFunc) (c *CMD, err error) {
	c = &CMD{
		envPrefix: `POS_`,
		dotenv:    true,
		cfg: &Config{
			PostgresURL: os.Getenv("POSTGRES_URL"),
			GRPCAddr:    os.Getenv("GRPC_LISTEN_ADDRESS"),
		},
	}

	for _, opt := range opts {
		if err = opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	cfg, err := env.MustLoad[Config](
		env.Options{
			Prefix: c.envPrefix,
			DotEnv: c.dotenv,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"load env: %w",
			err,
		)
	}
	c.cfg = cfg

	if err = c.initPostgres(); err != nil {
		return nil, fmt.Errorf(
			"failed init postgres: %w",
			err,
		)
	}
	if err := c.initServer(); err != nil {
		return nil, fmt.Errorf(
			"failed start grpc server: %w",
			err,
		)
	}
	return
}

func (c *CMD) initPostgres() (err error) {
	c.db, err = postgres.New(
		postgres.WithConnUri(c.cfg.PostgresURL),
	)
	if err != nil {
		return
	}
	c.closers = append(c.closers, c.db.Close)
	return
}

func (c *CMD) initServer() (err error) {
	log.Println("SERVER INITIALIZED")
	l, err := net.Listen(
		"tcp",
		c.cfg.GRPCAddr,
	)
	if err != nil {
		return err
	}
	log.Printf(
		"gRPC listening on %s",
		l.Addr().String(),
	)
	server := grpc.NewServer(
		// grpc.ChainUnaryInterceptor(),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	item_categoriespb.RegisterItemCategoryServiceServer(
		server,
		handler.NewItemCategories(
			handler.WithhDatabaseIC(c.db),
		),
	)
	reflection.Register(server)
	c.grpcServer = server
	go func() {
		if serveErr := server.Serve(l); serveErr != nil {
			log.Printf(
				"grpc serve error: %v",
				serveErr,
			)
		}
	}()
	c.closers = append(c.closers, func(ctx context.Context) error {
		server.GracefulStop()
		return nil
	})
	log.Printf("gRPC server started on %s", l.Addr())
	return
}

func (c *CMD) Run(ctx context.Context) error {
	log.Println("WAITING SIGNAL")
	<-ctx.Done()
	log.Println("CONTEXT DONE")
	log.Println("shutting down application")

	return c.close(ctx, nil)
}

func (c *CMD) close(ctx context.Context, err error) error {
	log.Println("CLOSE CALLED")
	for i := len(c.closers) - 1; i >= 0; i-- {
		err = errors.Join(
			err,
			c.closers[i](ctx),
		)
	}
	return err
}
