package service

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"team_5_game/config"
	"team_5_game/model/database"
)

const userPrefix = "USER_"

var (
	context     = redisClient.Context()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress(),
		Password: config.RedisPassword(),
		DB:       config.RedisDB(),
	})
)

func GetUserFromDB(id int64) (*database.User, error) {
	log.Println("Get user from DB, user ID", id)

	result, err := redisClient.Get(context, userPrefix+strconv.FormatInt(id, 10)).Result()
	if err == redis.Nil {
		log.Println("User not found", err)
		return nil, err
	} else if err != nil {
		log.Println("Could not read user from DB", err)
		return nil, err
	} else {
		user := &database.User{}

		err := json.Unmarshal([]byte(result), user)
		if err != nil {
			log.Println("Could not unmarshal user", err)
			return nil, err
		}

		log.Println("User successfully received from DB, ID", user.ID)
		return user, nil
	}
}

func SaveUserToDB(user *database.User) error {
	log.Println("Save user to the DB, ID", user.ID)

	out, err := json.Marshal(user)
	if err != nil {
		log.Println("Could not marshal user", err)
		return err
	}

	err = redisClient.Set(context, userPrefix+strconv.FormatInt(user.ID, 10), string(out), 0).Err()
	if err != nil {
		log.Println("Could not save user", err)
		return err
	}

	log.Println("User successfully saved to DB, ID", user.ID)
	return nil
}
