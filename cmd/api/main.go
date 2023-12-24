package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/movapi/internal/data"
	"github.com/movapi/internal/jsonlog"
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
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cnfg config

	flag.IntVar(&cnfg.port, "port", 3000, "API server port")
	flag.StringVar(&cnfg.env, "env", "development", "Envirnoment -> Dev|Staging|Prod")
	flag.StringVar(&cnfg.db.dsn, "db-dsn", os.Getenv("MOVAPI_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cnfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open conns")
	flag.IntVar(&cnfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idel conns")
	flag.StringVar((&cnfg.db.maxIdleTime), "db-max-idle-time", "15m", "PostgreSQL max idel time")

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
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cnfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cnfg.env,
	})
	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
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
