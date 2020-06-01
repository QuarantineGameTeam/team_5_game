package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func init() {
	// loads values from config.env into the system
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalln("No config.env file found")
	}
}

func BotToken() string {
	return os.Getenv("BOT_TOKEN")
}

func RedisAddress() string {
	return os.Getenv("REDIS_ADDRESS")
}

func RedisPassword() string {
	return os.Getenv("REDIS_PASSWORD")
}

func RedisDB() int {
	result, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	return result
}

func RegisterWebhook() bool {
	registerWebhook, _ := strconv.ParseBool(os.Getenv("REGISTER_WEBHOOK"))
	return registerWebhook
}

func ServerPort() string {
	return os.Getenv("SERVER_PORT")
}

func ServerURL() string {
	return os.Getenv("SERVER_URL")
}
