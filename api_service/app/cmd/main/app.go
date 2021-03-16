package main

import (
	"github.com/theartofdevel/notes_system/api_service/internal/router"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	defer router.Init()

	logger.Print("application initialized and started")
}
