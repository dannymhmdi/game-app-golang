package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"mymodule/adaptor/redis"
	"mymodule/repository/mysql"
	"mymodule/repository/redis/redisPresence"
	"mymodule/service/authService"
	"mymodule/service/matchmakingService"
	"strings"
)

type HttpServer struct {
	Port string `koanf:"port"`
}

type Config struct {
	HttpConfig        HttpServer                `koanf:"http_config"`
	AuthConfig        authService.Config        `koanf:"auth"`
	DbConfig          mysql.Config              `koanf:"db_config"`
	RedisConfig       redis.Config              `koanf:"redis_config"`
	RedisPresence     redisPresence.Config      `koanf:"redis_presence"`
	MatchMakingConfig matchmakingService.Config `koanf:"matchmaking_config"`
}

func Load() Config {
	k := koanf.New(".")

	lErr := k.Load(confmap.Provider(map[string]interface{}{
		"auth.access_subject":                         accessSubject,
		"auth.refresh_subject":                        refreshSubject,
		"db_config.password":                          mysqlPassword,
		"redis_presence.presence_key_expiration_time": redisPresenceKeyExpirationTime,
	}, "."), nil)

	if lErr != nil {
		panic(lErr)
	}

	err := k.Load(file.Provider("config.yml"), yaml.Parser())
	if err != nil {
		panic(err)
	}

	k.Load(env.Provider("GAMEAPP_", ".", func(s string) string {
		s = strings.TrimPrefix(s, "GAMEAPP_")
		s = strings.ToLower(s)
		s = strings.Replace(s, "_", ".", -1)
		return s
	}), nil)

	var config Config

	uErr := k.Unmarshal("", &config)
	if uErr != nil {
		panic(uErr)
	}
	fmt.Printf("config: %+v\n", config)
	return config
}
