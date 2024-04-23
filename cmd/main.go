package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uptrace/go-clickhouse/ch"
	"golang.org/x/sync/errgroup"
	"test-stat4market/internal/app"
	"test-stat4market/internal/config"
	"test-stat4market/internal/logger"
	"test-stat4market/internal/repository"
	"test-stat4market/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGQUIT, syscall.SIGTERM)
	defer stop()

	app.InitLogger(ctx)

	conf, err := config.NewConfig(ctx)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{Panic: "Error with reading config"})
	}

	db, err := connectToClickHouse(conf)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{Panic: "Error connecting to DB"})
	}

	dao := repository.NewDAO(db)
	eventsService := service.NewEventsService(dao)

	restSrv := app.ServerInit(ctx, eventsService)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.WarnKV(gCtx, logger.Data{Msg: "Clickhouse connected!"})
		return db.Ping(ctx)
	})
	g.Go(func() error {
		fmt.Println("Server started!")
		addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
		return restSrv.Listen(addr)
	})
	g.Go(func() error {
		<-gCtx.Done()
		fmt.Println("Server is shut down.")
		return restSrv.ShutdownWithTimeout(time.Second * 5)
	})
	g.Go(func() error {
		<-gCtx.Done()
		return db.Close()
	})

	if err = g.Wait(); err != nil {
		logger.ErrorKV(gCtx, logger.Data{
			Error:  err,
			Msg:    "Error group wait got error",
			Detail: fmt.Sprintf("exit reason: %v", err),
		})
	}

	logger.WarnKV(gCtx, logger.Data{Msg: "Server was gracefully shut down."})
}

func connectToClickHouse(config *config.Entity) (*ch.DB, error) {

	opts := []ch.Option{
		ch.WithAddr(fmt.Sprintf("%s:%d", config.DB.Host, config.DB.Port)),
		ch.WithCompression(true),
		ch.WithDatabase(config.DB.Name),
		ch.WithUser(config.DB.User),
		ch.WithPassword(config.DB.Pass),
		ch.WithDialTimeout(time.Duration(config.DB.ConnDialTimeout) * time.Second),
		ch.WithReadTimeout(time.Duration(config.DB.ConnReadTimeout) * time.Second),
		ch.WithWriteTimeout(time.Duration(config.DB.ConnWriteTimeout) * time.Second),
		ch.WithConnMaxLifetime(time.Hour),
		ch.WithPoolSize(config.DB.PoolSize),
		ch.WithInsecure(true),
	}

	db := ch.Connect(opts...)

	return db, nil
}
