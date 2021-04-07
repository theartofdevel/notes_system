package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/theartofdevel/notes_system/api_service/internal/client/category_service"
	"github.com/theartofdevel/notes_system/api_service/internal/client/note_service"
	"github.com/theartofdevel/notes_system/api_service/internal/client/tag_service"
	"github.com/theartofdevel/notes_system/api_service/internal/client/user_service"
	"github.com/theartofdevel/notes_system/api_service/internal/config"
	"github.com/theartofdevel/notes_system/api_service/internal/handlers/auth"
	"github.com/theartofdevel/notes_system/api_service/internal/handlers/categories"
	"github.com/theartofdevel/notes_system/api_service/internal/handlers/notes"
	"github.com/theartofdevel/notes_system/api_service/internal/handlers/tags"
	"github.com/theartofdevel/notes_system/api_service/pkg/cache/freecache"
	"github.com/theartofdevel/notes_system/api_service/pkg/handlers/metric"
	"github.com/theartofdevel/notes_system/api_service/pkg/jwt"
	"github.com/theartofdevel/notes_system/api_service/pkg/logging"
	"github.com/theartofdevel/notes_system/api_service/pkg/shutdown"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	logger.Println("config initializing")
	cfg := config.GetConfig()

	logger.Println("router initializing")
	router := httprouter.New()

	logger.Println("cache initializing")
	refreshTokenCache := freecache.NewCacheRepo(104857600) // 100MB

	logger.Println("helpers initializing")
	jwtHelper := jwt.NewHelper(refreshTokenCache, logger)

	logger.Println("create and register handlers")

	metricHandler := metric.Handler{Logger: logger}
	metricHandler.Register(router)

	userService := user_service.NewService(cfg.UserService.URL, "/users", logger)
	authHandler := auth.Handler{JWTHelper: jwtHelper, UserService: userService, Logger: logger}
	authHandler.Register(router)

	categoryService := category_service.NewService(cfg.CategoryService.URL, "/categories", logger)
	categoriesHandler := categories.Handler{CategoryService: categoryService, Logger: logger}
	categoriesHandler.Register(router)

	noteService := note_service.NewService(cfg.NoteService.URL, "/notes", logger)
	notesHandler := notes.Handler{NoteService: noteService, Logger: logger}
	notesHandler.Register(router)

	tagService := tag_service.NewService(cfg.TagService.URL, "/tags", logger)
	tagsHandler := tags.Handler{TagService: tagService, Logger: logger}
	tagsHandler.Register(router)

	logger.Println("start application")
	start(router, logger, cfg)
}

func start(router *httprouter.Router, logger logging.Logger, cfg *config.Config) {
	var server *http.Server
	var listener net.Listener

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		socketPath := path.Join(appDir, "app.sock")
		logger.Infof("socket path: %s", socketPath)

		logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		logger.Infof("bind application to host: %s and port: %s", cfg.Listen.BindIP, cfg.Listen.Port)

		var err error

		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		if err != nil {
			logger.Fatal(err)
		}
	}

	server = &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		server)

	logger.Println("application initialized and started")

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
