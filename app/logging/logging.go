package logging

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tolubydesign/angular-story-backend/app/database"
)

func getTotalNumberOfLogs(eventType string, client *redis.Client) (int, error) {
	context := context.Background()
	var cursor uint64
	var numeration int
	key := fmt.Sprintf("%s*", eventType)

	for {
		var keys []string
		var err error
		keys, cursor, err = client.Scan(context, cursor, key, 10).Result()
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
*/
func Event(d ...any) {
	var err error
	name := "logging"
	s := fmt.Sprint(d...)

	// Connect to Redis
	singleton, err := database.GetRedisClientSingleton()
	if err != nil {
		logBasicEvent(s)
		return
	}

	if singleton != nil {
		iteration, err := getTotalNumberOfLogs(name, singleton)
		if err != nil {
			logIterationError(err)
			return
		}

		key := fmt.Sprintf("%s:%v", name, iteration)
		context := context.Background()
		t := time.Now()

		value := fmt.Sprintf("%v, %v", t.String(), d)
		err = singleton.Set(context, key, value, time.Hour).Err()
		if err != nil {
			logRedisSingletonSetError(err)
			return
		}
	}
}

// Log error details to redis database.
func Error(details ...any) {
	var err error
	name := "error"
	s := fmt.Sprint(details...)

	// Connect to Redis
	singleton, err := database.GetRedisClientSingleton()
	if err != nil {
		// log.Println("TIMESTAMP:", time.Now().UTC(), " | LOG:ERROR | Redis Connection Error:", err.Error())
		logBasicEvent(s)
		return
	}

	if singleton != nil {
		iteration, err := getTotalNumberOfLogs(name, singleton)
		if err != nil {
			logIterationError(err)
			return
		}

		key := fmt.Sprintf("%s:%v", name, iteration)
		context := context.Background()
		t := time.Now()

		value := fmt.Sprintf("%v, %v", t.String(), details)
		err = singleton.Set(context, key, value, time.Hour).Err()
		if err != nil {
			logRedisSingletonSetError(err)
			return
		}
	}
}

// Log warning details to redis database.
func Warning(details ...any) {
	var err error
	name := "warning"
	s := fmt.Sprint(details...)

	// Connect to Redis
	singleton, err := database.GetRedisClientSingleton()
	if err != nil {
		// log.Println("TIMESTAMP:", time.Now().UTC(), " | LOG:ERROR | Redis Connection Error:", err.Error())
		logBasicEvent(s)
		return
	}

	if singleton != nil {
		iteration, err := getTotalNumberOfLogs(name, singleton)
		if err != nil {
			logIterationError(err)
			return
		}

		key := fmt.Sprintf("%s:%v", name, iteration)
		context := context.Background()
		t := time.Now()

		value := fmt.Sprintf("%v, %v", t.String(), details)
		err = singleton.Set(context, key, value, time.Hour).Err()
		if err != nil {
			logRedisSingletonSetError(err)
			return
		}
	}
}

func logIterationError(err error) {
	log.Println("| LOG:ERROR  | Failed to get Iteration:", err.Error())
}

func logRedisSingletonSetError(err error) {
	log.Println("| LOG:ERROR | Redis Client Set Error:", err.Error())
}

func logBasicEvent(message string) {
	log.Println("| LOG:EVENT | Message:", message)
}
