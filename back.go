package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

type back struct {
	db         *sql.DB
	bot        *tgbotapi.BotAPI
	httpServer *http.Server
	wg         sync.WaitGroup
}

type models struct {
}

type config struct {
	Database databaseConfig
	Server   serverConfig
	Redis    redisConfig
	Telegram telegramConfig
}

type databaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type telegramConfig struct {
	Token string
}

type serverConfig struct {
	Bind              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

type redisConfig struct {
	Addr string
	Port int
}

func newBack() (*back, error) {
	// db, err := sql.Open("postgres", cfg.Database.GetConn())
	//
	// if err != nil {
	// 	return nil, err
	// }
	//
	// err = db.Ping()
	//
	// if err != nil {
	// 	return nil, err
	// }

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr: fmt.Sprintf("%s:%d", cfg.Redis.Addr, cfg.Redis.Port),
	// })
	//
	// _, err := rdb.Ping(context.Background()).Result()
	//
	// if err != nil {
	// 	return nil, err
	// }

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)

	if err != nil {
		log.Panic(err)
	}

	a := back{
		// db: db,
		bot: bot,
	}

	err = a.setupHTTPServer(cfg.Server)

	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (b *back) Run() {
	b.runHTTPServer()

	log.Println(fmt.Sprintf("Starting server at port %s", b.httpServer.Addr))
}

func (b *back) runHTTPServer() {
	b.wg.Add(1)

	go func() {
		defer b.wg.Done()

		err := b.httpServer.ListenAndServe()

		if err != http.ErrServerClosed {
			b.Stop()
		}
	}()
}

func (b *back) Stop() {
	err := b.httpServer.Shutdown(context.Background())

	if err != nil {
		log.Println(err)
	}

	b.wg.Wait()
}

func (d *databaseConfig) GetConn() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Database,
	)
}
