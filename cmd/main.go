package main

import (
	"labs/pkg/handler"
	"labs/pkg/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	handler.FilePath = os.Getenv("filePath")
	handler.CommandFilePath = os.Getenv("commandFilePath")

	logFile, err := os.OpenFile("./app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Не удалось открыть файл лога:", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	server.StartServer()
}
