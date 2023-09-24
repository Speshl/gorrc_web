package config

import (
	"log"
	"os"
	"strconv"

	"github.com/Speshl/gorrc_web/internal/service/server/socketio"
)

const AppEnvBase = "GORRC_"

const DefaultPort = "8181"

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
		Host:     GetStringEnv("DBHOST", "127.0.0.1"),
		Port:     GetIntEnv("DBPORT", 3306),
		User:     GetStringEnv("DBUSER", "speshl"),
		Password: GetStringEnv("DBPASSWORD", "redalert"),
		DBName:   GetStringEnv("DBNAME", "gorrc"),
	}
}

func GetIntEnv(env string, defaultValue int) int {
	envValue, found := os.LookupEnv(AppEnvBase + env)
	if !found {
		return defaultValue
	} else {
		value, err := strconv.ParseInt(envValue, 10, 32)
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
		value, err := strconv.ParseBool(envValue)
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
		return envValue
	}
}
