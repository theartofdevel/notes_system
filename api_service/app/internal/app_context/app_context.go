package app_context

import (
	"github.com/theartofdevel/notes_system/api_service/internal/config"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"sync"
)

type AppContext struct {
	Config *config.Config
}

var instance *AppContext
var once sync.Once

func GetInstance() *AppContext {
	once.Do(func() {
		logging.GetLogger().Println("initializing application context")
		instance = &AppContext{
			Config: config.GetConfig(),
		}
	})

	return instance
}
