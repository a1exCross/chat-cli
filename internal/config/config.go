package config

import (
	"encoding/json"
	"github.com/a1exCross/chat-cli/internal/model"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const configFileName = "cfg.txt"

// Load - парсит .env
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

func LoadExecConfig() model.UserInfoConfig {
	data, err := os.ReadFile(configFileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf(color.RedString("%s does not exist", configFileName))
		}

		log.Fatalf("failed to read %s: %v", configFileName, err)
	}

	var userInfo model.UserInfoConfig

	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		log.Fatalf("failed to unmarshal %s: %v", configFileName, err)
	}

	return userInfo
}

func SaveExecConfig(userData model.UserInfoConfig) {
	data, err := json.Marshal(userData)
	if err != nil {
		log.Fatal("failed to marshal user info to JSON")
	}

	file, err := os.Create("cfg.txt")
	if err != nil {
		log.Fatal("failed to create cfg.txt")
	}

	_, err = file.Write(data)
	if err != nil {
		log.Fatal("failed to write user info to file")
	}
}

type AuthConfig interface {
	AuthAddress() string
}

type ChatConfig interface {
	ChatAddress() string
}
