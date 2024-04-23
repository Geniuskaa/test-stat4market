package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uptrace/go-clickhouse/ch"
	"test-stat4market/internal/config"
	"test-stat4market/internal/logger"
	"test-stat4market/internal/repository"
	"test-stat4market/internal/repository/entity"
)

const dataSize = 1100

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGQUIT, syscall.SIGTERM)
	defer stop()
	conf, err := config.NewConfig(ctx)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{Panic: "Error with reading config"})
	}

	db, err := connectToClickHouse(conf)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{Panic: "Error connecting to DB"})
	}
	err = db.Ping(ctx)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{Panic: "Error connecting to DB"})
	}
	defer db.Close()

	dao := repository.NewDAO(db)
	err = generateRandomData(ctx, dao)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{Panic: err})
	}
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

func generateRandomData(ctx context.Context, dao repository.DAO) error {
	q := dao.NewEventsQuery(ctx)
	for i := 1; i < dataSize; i++ {
		data := entity.Event{
			EventType: fmt.Sprintf("event_%v", i),
			UserID:    i,
			EventTime: time.Now(),
			Payload:   fmt.Sprintf("payload_%v", i),
		}
		err := q.Insert(&data)
		if err != nil {
			return err
		}
	}
	return nil
}
