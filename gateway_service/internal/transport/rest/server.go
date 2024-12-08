package rest

import (
	"context"
	"fmt"
	"gateway_service/pkg/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Server struct {
	serverPort int
	serv       *echo.Echo
}

func New(ctx context.Context, serverPort int) (*Server, error) {
	e := echo.New()
	e.GET("/v1/api/SearchAds", func(ctx echo.Context) error {
		ctx.String(http.StatusOK, "hello")
		return nil
	})
	e.POST("/v1/api/MakeAds", func(ctx echo.Context) error { return nil })
	e.POST("/v1/api/ApplyAds", func(ctx echo.Context) error { return nil })
	e.POST("/v1/api/Login", func(ctx echo.Context) error { return nil })
	e.POST("/v1/api/Signin", func(ctx echo.Context) error { return nil })
	return &Server{serverPort: serverPort, serv: e}, nil
}

func (s *Server) Start(ctx context.Context) error {
	log := logger.GetLogger(ctx)
	err := s.serv.Start(fmt.Sprintf(":%d", s.serverPort))
	if err != nil {
		return err
	}
	log.Info(ctx, "Server started")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.serv.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
