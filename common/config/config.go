package config

import (
	"encoding/json"
	"log"
	"os"
)

var configInfo map[string]interface{}

func init() {
	configInfo = make(map[string]interface{})
	pf, err := os.OpenFile("config/config.json", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("读取配置文件失败,[%s]\n", err.Error())
	}
	defer pf.Close()

	decoder := json.NewDecoder(pf)
	if err = decoder.Decode(&configInfo); err != nil {
		log.Fatalf("配置文件序列化失败,[%s]\n", err.Error())
	}
}

func GetChanelMax(key string) int {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("获取[%s]管道缓冲大小失败,请根据关键字检查配置文件[%s]\n", key, err)
		}
	}()
	chanel := configInfo["Chanel"].(map[string]interface{})
	entity := chanel[key].(map[string]interface{})
	return int(entity["Max"].(float64)) // MMP golang默认为float64
}

func GetDatabase(key string) string {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("获取[%s]管道缓冲大小失败,请根据关键字检查配置文件[%s]\n", key, err)
		}
	}()
	database := configInfo["Database"].(map[string]interface{})
	return database[key].(string)
}
