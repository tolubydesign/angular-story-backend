package database

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

// Create and connect to redis database
// How To Start Logging With Redis - https://betterstack.com/community/guides/logging/how-to-start-logging-with-redis/
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

func getTotalNumberOfLogs(eventType string) (int, error) {
	context := context.Background()
	var cursor uint64
	var numeration int
	key := fmt.Sprintf("%s*", eventType)

	for {
		var keys []string
		var err error
		keys, cursor, err = redisClientSingleton.Scan(context, cursor, key, 10).Result()
		if err != nil {
			return numeration, err
		}

		numeration += len(keys)
		if cursor == 0 {
			break
		}
	}

	return numeration, nil
}

// TODO: expand to include "warning" and "error"
/*
Log information to redis database.

How to Get All Keys in Redis - https://chartio.com/resources/tutorials/how-to-get-all-keys-in-redis/

Returns and error if found.
*/
func LogEvent(information string) error {
	eventType := "logging"
	iteration, checkErr := getTotalNumberOfLogs(eventType)
	if checkErr != nil {
		return checkErr
	}

	fmt.Printf("total  = %v\n", iteration)
	key := fmt.Sprintf("%s:%v", eventType, iteration)
	context := context.Background()
	t := time.Now()
	// var oneWeek time.Duration = 7 * 24 * 60 * 60      // 1 week = 7 days = 7×(24 hours) = 7×(24×(60 minutes)) = 7×(24×(60×(60 seconds))).
	// var duration time.Duration = 1000000000 * oneWeek // Equals 1000 Milliseconds. Equals 1 second
	if redisClientSingleton == nil {
		return errors.New("Redis Client singleton is null")
	}

	value := fmt.Sprintf("%v, %v", information, t.String())
	// err := redisClientSingleton.RPush(context, "logging", inputInformation).Err()
	err := redisClientSingleton.Set(context, key, value, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func LogError() {}

func LogWarning() {}
