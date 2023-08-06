package config

import (
	"log"
	"net/url"
	"os"
	"scheduleme/util"
	"strings"
)

type ConfigStruct struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GoogleRedirectPath string
	Port               string
	ENV                ENV
	Dsn                string
	Secret             string
}
type ENV string

func (e ENV) IsDev() bool {
	return strings.ToLower(string(e)) == "dev"
}
func (e ENV) IsProd() bool {
	return strings.ToLower(string(e)) == "prod"
}

func ConfigFromEnv() *ConfigStruct {

	c := &ConfigStruct{
		GoogleClientSecret: mustGetEnv("GOOGLE_CLIENT_SECRET"),
		GoogleClientID:     mustGetEnv("GOOGLE_CLIENT_ID"),
		GoogleRedirectURL:  mustGetEnv("GOOGLE_REDIRECT_URL"),
		Port:               getEnv("PORT", "8080"),
		ENV:                ENV(getEnv("ENV", "prod")),
		Dsn:                mustGetEnv("DSN"),
		Secret:             getEnv("SECRET", util.RandomStr(62)),
	}
	c.GoogleRedirectPath = getPath(c.GoogleRedirectURL)
	return c
}
func getPath(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Fatal(err)
	}
	return u.Path
}
func mustGetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic("env var " + key + " not found")
}
func getEnv(key string, defV string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defV
}

func InitConfig() *ConfigStruct {
	println("Initializing config from env vars...")
	cfg := ConfigFromEnv()
	return cfg
}
