package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"tunnerse/logger"
)

type Config struct {
	HTTPPort string

	SUBDOMAIN     bool
	WARNS_ON_HTML bool

	TUNNEL_LIFE_TIME            int
	TUNNEL_INACTIVITY_LIFE_TIME int

	DBHost    string
	DBPort    string
	DBUser    string
	DBName    string
	DBPwd     string
	DBSSLMode string
}

var AppConfig Config

func LoadAppConfig() error {

	err := godotenv.Load()
	if err != nil {
		logger.Log("DEBUG", "Error on read .env file", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		})
	}

	AppConfig = Config{
		HTTPPort: getEnvStr("HTTPPort", "8080"),

		SUBDOMAIN:     getEnvBool("SUBDOMAIN", false),
		WARNS_ON_HTML: getEnvBool("WARNS_ON_HTML", true),

		TUNNEL_LIFE_TIME:            getEnvInt("TUNNEL_LIFE_TIME", 86400),
		TUNNEL_INACTIVITY_LIFE_TIME: getEnvInt("TUNNEL_INACTIVITY_LIFE_TIME", 10),

		DBHost:    getEnvStr("DBHost", "localhost"),
		DBPort:    getEnvStr("DBPort", "27017"),
		DBUser:    getEnvStr("DBUser", "tunnerse_admin"),
		DBName:    getEnvStr("DBName", "tunnerse"),
		DBPwd:     getEnvStr("DBPwd", "s1ea021274"),
		DBSSLMode: getEnvStr("DBSSLMode", "disable"),
	}

	logger.Log("ENV", "Defined environment variables", []logger.LogDetail{
		{Key: "HTTPPort", Value: AppConfig.HTTPPort},
		{Key: "SUBDOMAIN", Value: AppConfig.SUBDOMAIN},
	})
	return nil
}

func getEnvStr(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
