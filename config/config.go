package config

import (
	"context"
	"fmt"
	"log"

	"github.com/caarlos0/env/v9"
)

type ContextKey string

const CfgKey ContextKey = "_cfg"

type Config struct {
	DB_URI        string `env:"DB_URI,required"`
	DB_DriverName string `env:"DB_DRIVER_NAME" envDefault:"sqlite3"`

	AUTH_JWTSecret string `env:"AUTH_JWT_SECRET,required"`

	HTTP_StaticDir      string `env:"HTTP_STATIC_DIR" envDefault:"_ui/public"`
	HTTP_IdleShutdownMS int64  `env:"HTTP_IDLE_SHUTDOWN_MS"`

	LOG_MinLevel     int `env:"LOG_LEVEL" envDefault:"-4"` // debug
	LOG_OutputFormat string `env:"LOG_OUTPUT_FORMAT"`

	RPC_Host       string `env:"RPC_HOST"`
	RPC_PathPrefix string `env:"RPC_PATH_PREFIX" envDefault:"/rpc"`
	RPC_Port       string `env:"RPC_PORT" envDefault:"4001"`
}

// TODO: Eventually use viper for this...

func InitConfig(ctx context.Context) (context.Context, error) {
	c := &Config{}

	if err := env.Parse(c); err != nil {
		return nil, err
	}
	fmt.Printf("Parsed config. E.g. DB_DRIVER_NAME = %s", c.DB_DriverName)

	withContext := context.WithValue(ctx, CfgKey, c)
	return withContext, nil
}

func MustContext(ctx context.Context) *Config {
	cfg, ok := ctx.Value(CfgKey).(*Config)
	if !ok || cfg == nil {
		log.Fatalln("no config found in context")
	}
	return cfg
}
