package config

import (
	"os"
	"fmt"
	"io"
	"encoding/json"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_URL				string	`json:"db_url"`				// key will be "db_url"
	CurrentUserName		string	`json:"current_user_name"`	// key will be "current_user_name"
}

func ReadConfigFile() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Could not get path to config file: ", configFileName)
		return Config{}, err
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return Config{}, err
	}
	defer file.Close()

	// Read the file into memory
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return Config{}, err
	}

	// Unmarshal into the struct
	cfg := Config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("Error parsing JSON into struct: %v\n", err)
		return Config{}, err
	}
	// fmt.Printf("Struct result: %+v\n", cfg)

	return cfg, nil
}

func SetUserName(cfg Config, currentUserName string) error {
	cfg.CurrentUserName = currentUserName

	jsonData, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		fmt.Println(err)
		return err
	}
	// fmt.Println(string(jsonData))

	filePath, err := getConfigFilePath()
	if err != nil {
		println("Could not get user's home directory...")
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(string(jsonData))

	return nil
}


func getConfigFilePath() (string, error) {
	filePath, err := os.UserHomeDir()
	if err != nil {
		println("Could not get user's home directory...")
		return "", err
	}

	filePath += "/" + configFileName
	return filePath, nil
}