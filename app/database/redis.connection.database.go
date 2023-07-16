package database

import (
	"context"
	"errors"
	"fmt"

	"time"

	"github.com/redis/go-redis/v9"
	config "github.com/tolubydesign/angular-story-backend/app/config"
)

var redisClientSingleton *redis.Client

/*
Create and connect to redis database.

Resource: How To Start Logging With Redis - https://betterstack.com/community/guides/logging/how-to-start-logging-with-redis/

Returning redis client and possible error.
*/
func ConnectToRedisDatabase() (*redis.Client, error) {
	var err error
	config, err := config.GetConfiguration()
	if err != nil {
		return nil, err
	}

	var redisConfig = config.Configuration.Redis

	databaseAddress := fmt.Sprintf("redis://%v:%v@%v:%v/%v", redisConfig.User, redisConfig.Password, redisConfig.Host, redisConfig.Port, redisConfig.Database)
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

/*
Log information to redis database.

Resource: How to Get All Keys in Redis - https://chartio.com/resources/tutorials/how-to-get-all-keys-in-redis/

Returning possible error.
*/
func LogEvent(details string) error {
	name := "logging"
	iteration, checkErr := getTotalNumberOfLogs(name)
	if checkErr != nil {
		return checkErr
	}

	key := fmt.Sprintf("%s:%v", name, iteration)
	context := context.Background()
	t := time.Now()
	// var oneWeek time.Duration = 7 * 24 * 60 * 60      // 1 week = 7 days = 7×(24 hours) = 7×(24×(60 minutes)) = 7×(24×(60×(60 seconds))).
	// var duration time.Duration = 1000000000 * oneWeek // Equals 1000 Milliseconds. Equals 1 second
	if redisClientSingleton == nil {
		return errors.New("Redis Client singleton is undefined")
	}

	value := fmt.Sprintf("%v, %v", t.String(), details)
	err := redisClientSingleton.Set(context, key, value, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

/*
Log error details to redis database.

Returning possible error.
*/
func LogError(details string) error {
	var err error
	name := "error"
	iteration, err := getTotalNumberOfLogs(name)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%v", name, iteration)
	context := context.Background()
	t := time.Now()
	if redisClientSingleton == nil {
		return errors.New("Redis Client singleton is undefined")
	}

	value := fmt.Sprintf("%v, %v", t.String(), details)
	err = redisClientSingleton.Set(context, key, value, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

/*
Log warning details to redis database.

Returning possible error.
*/
func LogWarning(details string) error {
	var err error
	name := "warning"
	iteration, err := getTotalNumberOfLogs(name)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%v", name, iteration)
	context := context.Background()
	t := time.Now()
	if redisClientSingleton == nil {
		return errors.New("Redis Client singleton is undefined")
	}

	value := fmt.Sprintf("%v, %v", t.String(), details)
	err = redisClientSingleton.Set(context, key, value, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}
