package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/hardyinalexander/make-auth-simple/authentication"
	"github.com/hardyinalexander/make-auth-simple/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	cfg := config.Get()
	ctx := context.Background()

	logger := initLogger()
	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	db, err := initDB(cfg)
	if err != nil {
		level.Error(logger).Log("exit", err)
		os.Exit(-1)
	}
	defer db.Close()

	repo := authentication.InitRepository(db)

	googleOauthCfg := initGoogleOauthConfig(cfg)
	service := authentication.InitService(repo, logger, googleOauthCfg, cfg.OauthStateString, cfg.SecretKey)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	clientHandler := authentication.InitHandler(googleOauthCfg, cfg.OauthStateString)
	endpoints := authentication.InitEndpoints(service)

	go func() {
		fmt.Println("listening on port:", cfg.Port)
		server := authentication.NewHTTPServer(ctx, endpoints, clientHandler)
		errs <- http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), server)
	}()

	level.Error(logger).Log("exit", <-errs)

}

func initLogger() log.Logger {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "authentication",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	return logger
}

func initDB(cfg config.Config) (*gorm.DB, error) {
	connInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword)
	db, err := gorm.Open("postgres", connInfo)

	return db, err
}

func initGoogleOauthConfig(cfg config.Config) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}
