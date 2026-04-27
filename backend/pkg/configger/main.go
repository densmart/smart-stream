package configger

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const DefaultCfgPath = "backend/config"

// InitConfig initializes the configuration by reading from the specified file
// and replacing any environment variable placeholders with their actual values.
func InitConfig(cfgPath string, name string, ctype string) {
	viper.SetConfigName(name)
	viper.SetConfigType(ctype)
	viper.AddConfigPath(cfgPath)

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("error loading config file: %s", err.Error())
		log.Fatalf("couldn't load configuration, cannot start")
	}
	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")))
		}
	}
}

// getEnvOrPanic retrieves the value of the environment variable named by the key.
// If the variable is not present, it returns the default value specified after the '|' character.
// If the format is incorrect (missing '|'), it panics.
func getEnvOrPanic(env string) string {
	if !strings.Contains(env, "|") {
		log.Fatalf("Log format variable %s is incorrect. '|' missing", env)
	}
	varSplit := strings.Split(env, "|")
	envVar := varSplit[0]
	defaultVar := varSplit[1]

	res := os.Getenv(envVar)
	if len(res) == 0 {
		res = defaultVar
	}
	return res
}
