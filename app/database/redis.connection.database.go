package database

import (
	"errors"
	"fmt"
	"log"

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
		log.Println("ERROR connecting to Redis database:", err.Error())
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
