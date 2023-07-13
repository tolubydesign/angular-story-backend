package utils

import (
	"context"
	"errors"
	"fmt"

	"time"

	"github.com/redis/go-redis/v9"
)

const (
	redisUser     = ""
	redisHost     = "localhost"
	redisPort     = "6379"
	redisPassword = ""
	redisDatabase = 0
)

var redisClientSingleton *redis.Client

func ConnectToRedisDatabase() (*redis.Client, error) {
	databaseAddress := fmt.Sprintf("redis://%v:%v@%v:%v/%v", redisUser, redisPassword, redisHost, redisPort, redisDatabase)
	opt, err := redis.ParseURL(databaseAddress)
	if err != nil {
		return nil, err
	}

	redisClientSingleton = redis.NewClient(opt)
	return redisClientSingleton, nil
}

func GetRedisClientSingleton() (*redis.Client, error) {
	if redisClientSingleton == nil {
		return nil, errors.New("Redis Client singleton is null")
	}

	return redisClientSingleton, nil
}

// TODO: expand to include "warning" and "error"
// TODO: document function
func ConsoleActionToRedisDatabase(information string) error {
	context := context.Background()
	t := time.Now()
	var oneWeek time.Duration = 7*24*60*60 // 1 week = 7 days = 7×(24 hours) = 7×(24×(60 minutes)) = 7×(24×(60×(60 seconds))).
	var duration time.Duration = 1000000000 * oneWeek // Equals 1000 Milliseconds. Equals 1 second
	if redisClientSingleton == nil {
		return errors.New("Redis Client singleton is null")
	}

	inputInformation := fmt.Sprintf("%v at %v", information, t.String())
	err := redisClientSingleton.Set(context, "logging", inputInformation, duration).Err()
	if err != nil {
		return err
	}

	return nil
}
