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
	"strings"
	"time"
)

type Server struct {
	serverPort int
	serv       *echo.Echo
	repo       *kafka.BrokerRepo
	ctx        context.Context
}

func New(ctx context.Context, serverPort int, repo *kafka.BrokerRepo) (*Server, error) {
	e := echo.New()
	return &Server{serverPort: serverPort, serv: e, repo: repo, ctx: ctx}, nil
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
			return ctx.String(http.StatusUnauthorized, "Search Service Unavailable")
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
			var searchResp models.SearchResponse
			var finds []models.SearchKafkaResponse
			d := respData.Data.(map[string]any)["finds"].([]interface{})
			logger.GetLogger(s.ctx).Info(s.ctx, fmt.Sprintf("found %+v", d))
			for _, v := range d {
				var searchData models.SearchKafkaResponse
				searchData.Name = v.(map[string]any)["name"].(string)
				searchData.Type = v.(map[string]any)["type"].(string)
				searchData.Description = v.(map[string]any)["description"].(string)
				searchData.Uuid = v.(map[string]any)["id"].(string)
				var Loc models.Location
				logger.GetLogger(s.ctx).Info(s.ctx, fmt.Sprintf("name type desc and uuid here %v", v.(map[string]any)))
				Loc.City = v.(map[string]any)["location"].(map[string]any)["city"].(string)
				Loc.Country = v.(map[string]any)["location"].(map[string]any)["country"].(string)
				Loc.District = v.(map[string]any)["location"].(map[string]any)["district"].(string)

				searchData.Location = Loc

				finds = append(finds, searchData)
			}
			searchResp.Finds = finds

			_ = ctx.JSON(http.StatusOK, searchResp)
			return nil
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
			return ctx.String(http.StatusUnauthorized, "Search Service Unavailable")
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
			logger.GetLogger(s.ctx).Info(s.ctx, fmt.Sprintf("msg: %v)", string(response)))
			makeData.Uuid = respData.Data.(map[string]any)["find_uuid"].(string)
			if err != nil {
				return ctx.String(http.StatusInternalServerError, err.Error())
			}
			_ = ctx.JSON(http.StatusOK, makeData)
			return nil
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
			return ctx.String(http.StatusUnauthorized, "Search Service Unavailable")
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
			applyData.Login = respData.Data.(map[string]any)["user_login"].(string)

			email := s.GetEmail(applyData.Login)
			err = s.Notify(email)
			applyData.Email = email
			_ = ctx.JSON(http.StatusOK, applyData)
			if err != nil {
				// TODO Log error
			}
			return nil
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
			loginData.Token = respData.Data.(map[string]any)["token"].(string)
			ctx.Response().Header().Add("Authorization", "Bearer "+loginData.Token)
			_ = ctx.JSON(http.StatusOK, models.LoginResponse{Success: true})
			return nil
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
	logger.GetLogger(s.ctx).Info(s.ctx, fmt.Sprintf("send req"))
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
			logger.GetLogger(s.ctx).Info(s.ctx, fmt.Sprintf("got answer %+v", respData))
			if err != nil {
				return ctx.String(http.StatusBadGateway, err.Error())
			}
			if respData.Status != "success" {
				continue
			}
			var regData models.RegisterKafkaResponse
			regData.Token = respData.Data.(map[string]any)["token"].(string)
			ctx.Response().Header().Add("Authorization", "Bearer "+regData.Token)
			_ = ctx.JSON(http.StatusOK, models.RegisterResponse{Success: true})
			return nil
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
	token = strings.TrimPrefix(token, "Bearer ")
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
			tokData.Valid = data.Data.(map[string]any)["valid"].(bool)
			valid := tokData.Valid
			if !valid {
				_ = ctx.String(http.StatusUnauthorized, "Invalid token")
				return false
			}
			return true
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
			return nil
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
			tokData.Email = data.Data.(map[string]any)["email"].(string)
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
		token = strings.TrimPrefix(token, "Bearer ")
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
			tokData.Login = data.Data.(map[string]any)["login"].(string)
			login := tokData.Login
			logger.GetLogger(s.ctx).Info(s.ctx, fmt.Sprintf("got answer %+v", resp))
			logger.GetLogger(s.ctx).Info(s.ctx, fmt.Sprintf("got login %+v", login))

			return login
		case <-time.After(5 * time.Second):
			continue
		}
	}
	return ""
}
