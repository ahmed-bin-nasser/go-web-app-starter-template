package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"example.com/internal/database"

	"github.com/caarlos0/env"
	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	gob.Register(time.Time{})

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	BaseURL  string `env:"BASE_URL"`
	HttpPort int    `env:"HTTP_PORT"`
	AppEnv   string `env:"APP_ENV"`

	Dsn string `env:"DSN"`

	SecretKey     string `env:"SECRET_KEY"`
	OldSecretKey  string `env:"OLD_SECRET_KEY"`
	CookieMaxAge  int    `env:"COOKIE_MAX_AGE"`
	CsrfSecretKey string `env:"CSRF_SECRET_KEY"`

	Version string `env:"VERSION"`
}

type application struct {
	config        config
	wg            sync.WaitGroup
	db            *database.DB
	logger        *slog.Logger
	sessionStore  *sessions.CookieStore
	templateCache map[string]*template.Template
}

func run(logger *slog.Logger) error {
	var cfg config

	if err := env.Parse(&cfg); err != nil {
		return err
	}

	showVersion := flag.Bool("version", false, "display version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("version: %s\n", cfg.Version)
		return nil
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		return err
	}

	db, err := database.New(cfg.Dsn)
	if err != nil {
		return err
	}

	keyPairs := [][]byte{[]byte(cfg.SecretKey), nil}
	if cfg.OldSecretKey != "" {
		keyPairs = append(keyPairs, []byte(cfg.OldSecretKey), nil)
	}

	sessionStore := sessions.NewCookieStore(keyPairs...)
	sessionStore.Options = &sessions.Options{
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * cfg.CookieMaxAge,
	}

	app := &application{
		wg:            sync.WaitGroup{},
		db:            db,
		config:        cfg,
		logger:        logger,
		sessionStore:  sessionStore,
		templateCache: templateCache,
	}

	return app.serveHTTP()
}
