package rest

import (
	"context"
	"fmt"
	"gateway_service/pkg/kafka"
	"gateway_service/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Server struct {
	serverPort int
	serv       *echo.Echo
	repo       *kafka.BrokerRepo
}

func New(ctx context.Context, serverPort int) (*Server, error) {
	e := echo.New()
	return &Server{serverPort: serverPort, serv: e}, nil
}

func (s *Server) Start(ctx context.Context) error {
	log := logger.GetLogger(ctx)
	s.serv.GET("/v1/api/SearchAds", s.SearchAds)
	s.serv.POST("/v1/api/MakeAds", s.MakeAds)
	s.serv.POST("/v1/api/ApplyAds", s.ApplyAds)
	s.serv.POST("/v1/api/Login", s.Login)
	s.serv.POST("/v1/api/Signin", s.Signin)
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

func (s *Server) SearchAds(ctx echo.Context) error {
	//TODO jwt check
	id := uuid.New().String()
	var data []byte
	_, err := ctx.Request().Body.Read(data)
	if err != nil {
		return err
	}
	err = ctx.Request().Body.Close()
	if err != nil {
		return err
	}
	respose_ch, err := s.repo.SearchAds(id, string(data))
	if err != nil {
		return err
	}

	select {
	case respose := <-respose_ch:
		_, err = ctx.Response().Write([]byte(respose))
		if err != nil {
			return err
		}
	case <-time.After(5 * time.Second):
		ctx.Response().Status = http.StatusGatewayTimeout
		_, err = ctx.Response().Write([]byte("timeout"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) MakeAds(ctx echo.Context) error {
	return nil
}

func (s *Server) ApplyAds(ctx echo.Context) error {
	return nil
}

func (s *Server) Login(ctx echo.Context) error {
	return nil
}

func (s *Server) Signin(ctx echo.Context) error {
	return nil
}
