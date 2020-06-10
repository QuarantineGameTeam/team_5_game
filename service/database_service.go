package service

import (
	"encoding/json"
	"log"
	"strconv"
	"team_5_game/config"
	"team_5_game/model/database"
	"team_5_game/model/telegram"

	"github.com/go-redis/redis/v8"
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

func SaveUserClan(callbackQuery *telegram.CallbackQuery) {
	log.Println("Start clan saving")

	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
	}

	clanName := string(callbackQuery.Data)

	user.Clan = clanName

	out, err := json.Marshal(user)
	if err != nil {
		log.Println("Could not marshal user", err)
	}

	err = redisClient.Set(context, userPrefix+strconv.FormatInt(callbackQuery.From.ID, 10), string(out), 0).Err()
	if err != nil {
		log.Println("Could not save clan", err)
	}
}

func AppendUserTrack(callbackQuery *telegram.CallbackQuery, position int) {
	log.Println("Start track saving")

	user, err := GetUserFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
	}
	user.Track[position-1] = position

	out, err := json.Marshal(user)
	if err != nil {
		log.Println("Could not marshal user", err)
	}

	err = redisClient.Set(context, userPrefix+strconv.FormatInt(callbackQuery.From.ID, 10), string(out), 0).Err()
	if err != nil {
		log.Println("Could not save track", err)
	}
}
