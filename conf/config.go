package conf

import (
	"log"

	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func init() {
	var err error
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath("config/")

	err = v.ReadInConfig()

	if err != nil {
		log.Fatal("error on parsing configuration file ", err.Error())
	}
	config = v
}

func GetConfig() *viper.Viper {
	return config
}
