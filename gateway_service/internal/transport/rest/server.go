package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gateway_service/internal/models"
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

func New(ctx context.Context, serverPort int, repo *kafka.BrokerRepo) (*Server, error) {
	e := echo.New()
	return &Server{serverPort: serverPort, serv: e, repo: repo}, nil
}

func (s *Server) Start(ctx context.Context) error {
	log := logger.GetLogger(ctx)
	s.serv.GET("/v1/api/SearchAds", s.SearchAds)
	s.serv.POST("/v1/api/MakeAds", s.MakeAds)
	s.serv.POST("/v1/api/ApplyAds", s.ApplyAds)
	s.serv.POST("/v1/api/Login", s.Login)
	s.serv.POST("/v1/api/Register", s.Register)
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
	s.repo.Stop()
	return nil
}

func (s *Server) SearchAds(ctx echo.Context) error {
	valid := s.VerifyToken(ctx)
	if !valid {
		return nil
	}

	var data models.SearchRequest
	if err := ctx.Bind(&data); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	for i := range 6 {
		if i == 5 {
			return ctx.String(http.StatusUnauthorized, "Auth Service Unavailable")
		}
		id := uuid.New().String()
		ch, err := s.repo.SearchAds(id, data.Name, data.TypeOfFinding, data.Location)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		select {
		case response := <-ch:
			var respData models.KafkaResponse
			err = json.Unmarshal(response, &respData)
			if err != nil {
				return ctx.String(http.StatusBadGateway, err.Error())
			}
			if respData.Status != "success" {
				continue
			}
			searchData := models.SearchResponse{
				Finds: respData.Data.([]models.Finding),
			}

			if err != nil {
				return ctx.String(http.StatusInternalServerError, err.Error())
			}
			_ = ctx.JSON(http.StatusOK, searchData)
			break
		case <-time.After(5 * time.Second):
			continue
		}
	}
	return nil
}

func (s *Server) MakeAds(ctx echo.Context) error {
	valid := s.VerifyToken(ctx)
	if !valid {
		return nil
	}
	login := s.GetLogin(ctx.Request().Header.Get("Authorization"))
	var data models.MakeAdsRequest
	if err := ctx.Bind(&data); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	for i := range 6 {
		if i == 5 {
			return ctx.String(http.StatusUnauthorized, "Auth Service Unavailable")
		}
		id := uuid.New().String()
		ch, err := s.repo.MakeAds(login, id, data.Name, data.Description, data.TypeOfFinding, data.Location)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		select {
		case response := <-ch:
			var respData models.KafkaResponse
			err = json.Unmarshal(response, &respData)
			if err != nil {
				return ctx.String(http.StatusBadGateway, err.Error())
			}
			if respData.Status != "success" {
				continue
			}
			var makeData models.MakeAdsResponse
			err = json.Unmarshal([]byte(respData.Data), &makeData)
			if err != nil {
				return ctx.String(http.StatusInternalServerError, err.Error())
			}
			_ = ctx.JSON(http.StatusOK, makeData)
			break
		case <-time.After(5 * time.Second):
			continue
		}
	}
	return nil
}

func (s *Server) ApplyAds(ctx echo.Context) error {
	valid := s.VerifyToken(ctx)
	if !valid {
		return nil
	}

	var data models.ApplyRequest
	if err := ctx.Bind(&data); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	for i := range 6 {
		if i == 5 {
			return ctx.String(http.StatusUnauthorized, "Auth Service Unavailable")
		}
		id := uuid.New().String()
		ch, err := s.repo.ApplyAds(id, data.Uuid)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		select {
		case response := <-ch:
			var respData models.KafkaResponse
			err = json.Unmarshal(response, &respData)
			if err != nil {
				return ctx.String(http.StatusBadGateway, err.Error())
			}
			if respData.Status != "success" {
				continue
			}
			var applyData models.ApplyResponse
			err = json.Unmarshal([]byte(respData.Data), &applyData)
			if err != nil {
				return ctx.String(http.StatusInternalServerError, err.Error())
			}
			_ = ctx.JSON(http.StatusOK, applyData)
			email := s.GetEmail(applyData.Login)
			err = s.Notify(email)
			if err != nil {
				// TODO Log error
			}
			break
		case <-time.After(5 * time.Second):
			continue
		}
	}
	return nil
}

func (s *Server) Login(ctx echo.Context) error {
	var data models.LoginRequest
	if err := ctx.Bind(&data); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	for i := range 6 {
		if i == 5 {
			return ctx.String(http.StatusUnauthorized, "Auth Service Unavailable")
		}
		id := uuid.New().String()
		ch, err := s.repo.Login(id, data.Login, data.Password)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		select {
		case response := <-ch:
			var respData models.KafkaResponse
			err = json.Unmarshal(response, &respData)
			if err != nil {
				return ctx.String(http.StatusBadGateway, err.Error())
			}
			if respData.Status != "success" {
				continue
			}
			var loginData models.LoginKafkaResponse
			err = json.Unmarshal([]byte(respData.Data), &loginData)
			if err != nil {
				return ctx.String(http.StatusInternalServerError, err.Error())
			}
			ctx.Response().Header().Add("Authorization", "Bearer "+loginData.Token)
			_ = ctx.JSON(http.StatusOK, models.LoginResponse{Success: true})
			break

		case <-time.After(5 * time.Second):
			continue
		}
	}
	return nil
}

func (s *Server) Register(ctx echo.Context) error {
	var data models.RegisterRequest
	if err := ctx.Bind(&data); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	for i := range 6 {
		if i == 5 {
			return ctx.String(http.StatusUnauthorized, "Auth Service Unavailable")
		}
		id := uuid.New().String()
		ch, err := s.repo.Register(id, data.Login, data.Password, data.Email)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		select {
		case response := <-ch:
			var respData models.KafkaResponse
			err = json.Unmarshal(response, &respData)
			if err != nil {
				return ctx.String(http.StatusBadGateway, err.Error())
			}
			if respData.Status != "success" {
				continue
			}
			var regData models.RegisterKafkaResponse
			err = json.Unmarshal([]byte(respData.Data), &regData)
			if err != nil {
				return ctx.String(http.StatusInternalServerError, err.Error())
			}
			ctx.Response().Header().Add("Authorization", "Bearer "+regData.Token)
			_ = ctx.JSON(http.StatusOK, models.RegisterResponse{Success: true})
			break
		case <-time.After(5 * time.Second):
			continue
		}
	}
	return nil
}

func (s *Server) VerifyToken(ctx echo.Context) bool {
	token := ctx.Request().Header.Get("Authorization")
	if token == "" {
		_ = ctx.String(http.StatusUnauthorized, "Authorization header required")
		return false
	}
	// Validate Token
	for i := range 6 {
		if i == 5 {
			_ = ctx.String(http.StatusUnauthorized, "Auth Service Unavailable")
			return false
		}
		id := uuid.New().String()
		ch, err := s.repo.ValidateToken(id, token)
		if err != nil {
			_ = ctx.String(http.StatusBadGateway, err.Error())
			return false
		}
		select {
		case resp := <-ch:
			var data models.KafkaResponse
			_ = json.Unmarshal(resp, &data)
			if data.Status != "success" {
				continue
			}
			var tokData models.ValidateTokenResponse
			_ = json.Unmarshal([]byte(data.Data), &tokData)
			valid := tokData.Valid
			if !valid {
				_ = ctx.String(http.StatusUnauthorized, "Invalid token")
				return false
			}
			break
		case <-time.After(5 * time.Second):
			continue
		}
	}
	return true
}

func (s *Server) Notify(email string) error {
	for i := range 6 {
		if i == 5 {
			return errors.New("notify service failed")
		}
		id := uuid.New().String()
		subject := ""
		body := ""
		ch, err := s.repo.NotifyUser(id, email, subject, body)
		if err != nil {
			return err
		}
		select {
		case resp := <-ch:
			var data models.KafkaResponse
			_ = json.Unmarshal(resp, &data)
			if data.Status != "success" {
				continue
			}
			var tokData = data.Data
			fmt.Println(tokData)
			break
		case <-time.After(5 * time.Second):
			continue
		}
	}
	return nil
}

func (s *Server) GetEmail(login string) string {
	for i := range 6 {
		if i == 5 {
			return ""
		}
		id := uuid.New().String()
		ch, err := s.repo.GetEmail(id, login)
		if err != nil {
			return ""
		}
		select {
		case resp := <-ch:
			var data models.KafkaResponse
			_ = json.Unmarshal(resp, &data)
			if data.Status != "success" {
				continue
			}
			var tokData models.GetEmailResponse
			_ = json.Unmarshal([]byte(data.Data), &tokData)
			email := tokData.Email
			return email
		case <-time.After(5 * time.Second):
			continue
		}
	}
	return ""
}

func (s *Server) GetLogin(token string) string {
	for i := range 6 {
		if i == 5 {
			return ""
		}
		id := uuid.New().String()
		ch, err := s.repo.GetLogin(id, token)
		if err != nil {
			return ""
		}
		select {
		case resp := <-ch:
			var data models.KafkaResponse
			_ = json.Unmarshal(resp, &data)
			if data.Status != "success" {
				continue
			}
			var tokData models.GetLoginResponse
			_ = json.Unmarshal([]byte(data.Data), &tokData)
			login := tokData.Login
			return login
		case <-time.After(5 * time.Second):
			continue
		}
	}
	return ""
}
