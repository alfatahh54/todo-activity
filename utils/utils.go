package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GoDotEnvVariable(key string) string {
	err := godotenv.Load(Dir(".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(key)
}
func Dir(path string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, path)
}
func TestMode() bool {
	tesMode := GoDotEnvVariable("TEST_MODE")
	if tesMode == "TRUE" {
		return true
	}
	return false
}
