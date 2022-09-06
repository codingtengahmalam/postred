package src

import (
	"fmt"
	validatorEngine "github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"postred/config"
	"postred/src/delivery"
	"postred/src/helper/validator"
	"postred/src/repository"
	"postred/src/usecase"
)

type (
	server struct {
		httpServer *echo.Echo
		cfg        config.Config
	}

	Server interface {
		Run()
	}
)

func InitServer(cfg config.Config) Server {
	e := echo.New()
	e.HideBanner = true
	e.Validator = &validator.GoPlaygroundValidator{
		Validator: validatorEngine.New(),
	}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return &server{
		httpServer: e,
		cfg:        cfg,
	}
}

func (s *server) Run() {

	s.httpServer.GET("", func(e echo.Context) error {

		return e.JSON(http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "Hello, World!" + s.cfg.ServiceName() + " " + s.cfg.ServiceEnvironment(),
		})
	})

	postRepo := repository.NewPostRepository(s.cfg)
	postUsecase := usecase.NewPostUsecase(postRepo)
	postDelivery := delivery.NewPostDelivery(postUsecase)
	postGroup := s.httpServer.Group("/posts")
	postDelivery.Mount(postGroup)

	if err := s.httpServer.Start(fmt.Sprintf(":%d", s.cfg.ServicePort())); err != nil {
		log.Panic(err)
	}
}
