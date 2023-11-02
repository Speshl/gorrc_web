package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Speshl/gorrc_web/internal/service/server/socketio"
)

const (
	AppEnvBase        = "GORRC_"
	DefaultPort       = "8181"
	DefaultDBHost     = "192.168.1.22"
	DefaultDBPort     = 3306
	DefaultDBUser     = "test"
	DefaultDBPassword = "testpassword"
	DefaultDBName     = "gorrc"
)

type Config struct {
	ServerCfg ServerCfg
	DBCfg     DBCfg
}

type ServerCfg struct {
	SocketIOCfg socketio.SocketIOServerCfg
	Port        string
}

type DBCfg struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func GetConfig() Config {
	cfg := Config{
		DBCfg:     GetDBCfg(),
		ServerCfg: GetServerCfg(),
	}

	log.Printf("Server Config: \n%+v\n", cfg)
	return cfg
}

func GetServerCfg() ServerCfg {
	return ServerCfg{
		Port:        GetStringEnv("PORT", DefaultPort),
		SocketIOCfg: GetSocketIOCfg(),
	}
}

func GetSocketIOCfg() socketio.SocketIOServerCfg {
	return socketio.SocketIOServerCfg{}
}

func GetDBCfg() DBCfg {
	return DBCfg{
		Host:     GetStringEnv("DBHOST", DefaultDBHost),
		Port:     GetIntEnv("DBPORT", DefaultDBPort),
		User:     GetStringEnv("DBUSER", DefaultDBUser),
		Password: GetStringEnv("DBPASSWORD", DefaultDBPassword),
		DBName:   GetStringEnv("DBNAME", DefaultDBName),
	}
}

func GetIntEnv(env string, defaultValue int) int {
	envValue, found := os.LookupEnv(AppEnvBase + env)
	if !found {
		return defaultValue
	} else {
		value, err := strconv.ParseInt(strings.Trim(envValue, "\r"), 10, 32)
		if err != nil {
			log.Printf("warning:%s not parsed - error: %s\n", env, err)
			return defaultValue
		} else {
			return int(value)
		}
	}
}

func GetBoolEnv(env string, defaultValue bool) bool {
	envValue, found := os.LookupEnv(AppEnvBase + env)
	if !found {
		return defaultValue
	} else {
		value, err := strconv.ParseBool(strings.Trim(envValue, "\r"))
		if err != nil {
			log.Printf("warning:%s not parsed - error: %s\n", env, err)
			return defaultValue
		} else {
			return value
		}
	}
}

func GetStringEnv(env string, defaultValue string) string {
	envValue, found := os.LookupEnv(AppEnvBase + env)
	if !found {
		return defaultValue
	} else {
		return strings.Trim(envValue, "\r")
	}
}

func GetFloatEnv(env string, defaultValue float64) float64 {
	envValue, found := os.LookupEnv(AppEnvBase + env)
	if !found {
		return defaultValue
	} else {
		value, err := strconv.ParseFloat(envValue, 64)
		if err != nil {
			return defaultValue
		}
		return value
	}
}
