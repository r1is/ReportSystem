package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("读取配置文件错误")
	}
}
