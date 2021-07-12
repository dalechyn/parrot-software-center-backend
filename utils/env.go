package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)


type Keys struct {
	MainKey string `yaml:"main"`
	EmailKey string `yaml:"email"`
	ModeratorKey string `yaml:"moderator"`
}

type RedisPasswords struct {
	MasterPassword string `yaml:"master"`
	SentinelPassword string `yaml:"sentinel"`
}

type EmailCredentials struct {
	Mail string `yaml:"mail"`
	Password string `yaml:"password"`
}

type Settings struct {
	Moderators []string `yaml:"moderators"`
	IsolatedPackages []string `yaml:"isolatedPackages"`
	Keys Keys `yaml:"keys"`
	RedisPasswords RedisPasswords `yaml:"redis"`
	EmailCredentials EmailCredentials `yaml:"email"`
}

var Root struct {
	Settings Settings `yaml:"settings"`
}

func Load() {
	settingsRaw, err := os.ReadFile("settings.yaml")
	if err != nil {
		log.Fatalf(`Oops! Looks like you don't have settings.yaml.
Check settings.example.yaml for an example and create look-a-like Settings.yaml`)
	}

	err = yaml.Unmarshal(settingsRaw, &Root)
	if err != nil {
		log.Fatalf("parsing error: %v", err)
	}

	// Connecting to Redis
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs:    []string{"sentinel:26379", "sentinel:26380", "sentinel:26381"},
		MasterName:       "mymaster",
		SentinelPassword: Root.Settings.RedisPasswords.SentinelPassword,
		Password:         Root.Settings.RedisPasswords.MasterPassword,
	})

	{
		var args []interface{}
		for _, m := range Root.Settings.Moderators {
			args = append(args, m)
		}

		if _, err := rdb.SAdd(context.Background(), "moderators", args...).Result(); err != nil {
			log.Fatal(err)
		}
	}
	{
		var args []interface{}
		for _, m := range Root.Settings.IsolatedPackages {
			args = append(args, m)
		}

		if _, err := rdb.SAdd(context.Background(), "isolated", args...).Result(); err != nil {
			log.Fatal(err)
		}
	}

}

func GetSecret() string {
	return Root.Settings.Keys.MainKey
}

func GetRedisPassword() string {
	return Root.Settings.RedisPasswords.MasterPassword
}

func GetSentinelPassword() string {
	return Root.Settings.RedisPasswords.SentinelPassword
}

func GetEmail() string {
	return Root.Settings.EmailCredentials.Mail
}

func GetEmailPassword() string {
	return Root.Settings.EmailCredentials.Password
}