package app

import (
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/rs/zerolog"
	"time"
)

type Config struct {
	App        app
	DB         db
	Logger     logger
	Middleware middleware
}

type app = struct {
	Name    string
	Port    string
	Timeout time.Duration
	TLS     struct {
		Enable   bool
		CertFile string
		KeyFile  string
	}
}

type db = struct {
	Postgres struct {
		DSN string
	}
}

type logger = struct {
	TimeFormat string
	Level      zerolog.Level
	Prettier   bool
}

type middleware = struct {
	Compress struct {
		Enable bool
		Level  compress.Level
	}

	Recover struct {
		Enable bool
	}

	Monitor struct {
		Enable bool
		Path   string
	}

	Pprof struct {
		Enable bool
	}

	Limiter struct {
		Enable  bool
		Max     int
		ExpSecs time.Duration
	}

	Filesystem struct {
		Enable bool
		Browse bool
		MaxAge int
		Index  string
		Root   string
	}
}

func NewConfig() *Config {
	config := Config{
		App: app{
			Name: "a",
			Port: "5000",
		},
		DB:         db{},
		Logger:     logger{},
		Middleware: middleware{},
	}
	return &config
}
