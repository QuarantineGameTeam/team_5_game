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

const (
	userPrefix   = "USER_"
	battlePrefix = "BATTLE_"
)

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

func AppendUserTrack(callbackQuery *telegram.CallbackQuery, position int) {
	log.Println("Start track saving")

	//user, err := GetUserFromDB(callbackQuery.From.ID)
	battle, err := GetBattleFromDB(callbackQuery.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
	}
	battle.Sector[position-1].OwnedBy[0] = callbackQuery.From.ID
	battle.Sector[position-1].IsCaptured = true

	out, err := json.Marshal(battle)
	if err != nil {
		log.Println("Could not marshal battle", err)
	}

	err = redisClient.Set(context, battlePrefix+strconv.FormatInt(callbackQuery.From.ID, 10), string(out), 0).Err()
	if err != nil {
		log.Println("Could not save track", err)
	}
}

func GetBattleFromDB(id int64) (*database.Battle, error) {
	log.Println("Get battle from DB, battle ID", id)

	result, err := redisClient.Get(context, battlePrefix+strconv.FormatInt(id, 10)).Result()
	if err == redis.Nil {
		log.Println("Battle not found", err)
		return nil, err
	} else if err != nil {
		log.Println("Could not read battle from DB", err)
		return nil, err
	} else {
		battle := &database.Battle{}

		err := json.Unmarshal([]byte(result), battle)
		if err != nil {
			log.Println("Could not unmarshal battle", err)
			return nil, err
		}

		log.Println("Battle successfully received from DB, ID", battle.ID)
		return battle, nil
	}
}

func SaveBattleToDB(battle *database.Battle) error {
	log.Println("Save battle to the DB, ID", battle.ID)

	out, err := json.Marshal(battle)
	if err != nil {
		log.Println("Could not marshal battle", err)
		return err
	}

	err = redisClient.Set(context, battlePrefix+strconv.FormatInt(battle.ID, 10), string(out), 0).Err()
	if err != nil {
		log.Println("Could not save battle", err)
		return err
	}

	log.Println("Battle successfully saved to DB, ID", battle.ID)
	return nil
}
