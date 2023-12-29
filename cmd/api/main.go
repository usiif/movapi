package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/movapi/internal/data"
	"github.com/movapi/internal/jsonlog"
	"github.com/movapi/internal/mailer"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		corsSafeList []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cnfg config

	flag.IntVar(&cnfg.port, "port", 3000, "API server port")
	flag.StringVar(&cnfg.env, "env", "development", "Envirnoment -> Dev|Staging|Prod")
	flag.StringVar(&cnfg.db.dsn, "db-dsn", os.Getenv("MOVAPI_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cnfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open conns")
	flag.IntVar(&cnfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idel conns")
	flag.StringVar((&cnfg.db.maxIdleTime), "db-max-idle-time", "15m", "PostgreSQL max idel time")
	flag.Float64Var(&cnfg.limiter.rps, "limiter-rps", 2, "Rate limiter max requests per second")
	flag.IntVar(&cnfg.limiter.burst, "limiter-burst", 4, "Rate limiter burst")
	flag.BoolVar(&cnfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cnfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cnfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cnfg.smtp.username, "smtp-username", "4cfc466eea6664", "SMTP username")
	flag.StringVar(&cnfg.smtp.password, "smtp-password", "407823874fac4d", "SMTP password")
	flag.StringVar(&cnfg.smtp.sender, "smtp-sender", "Movapi <no-reply@Movapi.com>", "SMTP sender")
	flag.Func("cors-trusted-origins", "Trust CORS origins - Space seperated", func(s string) error {
		cnfg.cors.corsSafeList = strings.Split(s, " ")
		return nil
	})

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LeveInfo)

	db, err := openDB(cnfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("Database connection bool established", nil)

	app := &application{
		config: cnfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cnfg.smtp.host, cnfg.smtp.port, cnfg.smtp.username, cnfg.smtp.password, cnfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxIdleConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	durtation, err := time.ParseDuration(cfg.db.maxIdleTime)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(durtation)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
