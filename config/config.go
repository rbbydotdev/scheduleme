package config

import (
	"fmt"
	"net/url"
	"os"
	"scheduleme/util"
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

const (
	EnvDev  ENV = "dev"
	EnvProd ENV = "prod"
)

func (e ENV) IsDev() bool {
	return e == EnvDev
}

func (e ENV) IsProd() bool {
	return e == EnvProd
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
		panic(fmt.Sprintf("bad url=%s err=%v", rawurl, err))
	}
	return u.Path
}
func mustGetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf("env var=%s not found", key))
}
func getEnv(key string, defV string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defV
}

func InitConfig() *ConfigStruct {
	cfg := ConfigFromEnv()
	return cfg
}
