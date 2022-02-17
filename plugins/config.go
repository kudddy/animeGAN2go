package plugins

import (
	"github.com/joho/godotenv"
	"os"
)

var e = godotenv.Load()

var Username = os.Getenv("db_user")
var Password = os.Getenv("db_pass")
var DbName = os.Getenv("db_name")
var DbHost = os.Getenv("db_host")
var DbPort = os.Getenv("db_port")

var Token = os.Getenv("bot_token")

var RedisHost = os.Getenv("redis_host")
