package authservice

import (
	"auth/internal/userrepo"
	"context"
	"testing"
	"time"
)

func TestAuthService(t *testing.T) {
	repo, err := userrepo.NewUserRepo(userrepo.AuthRepoConfig{
		UserName: "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     "5432",
		DbName:   "db",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	as := NewAuthService(AuthServiceConfig{
		TokenExpiration: 5,
	}, repo)

	token, err := as.CreateUser(context.Background(), "test", "password", "data")
	if err != nil {
		t.Fatal(err)
	}

	login, err := as.GetLoginByToken(context.Background(), token)
	if err != nil {
		t.Fatal(err)
	}

	if login != "test" {
		t.Fatal("login != test")
	}

	data, err := as.GetUserData(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}

	if data != "data" {
		t.Fatal("data != data")
	}

	token, err = as.LoginUser(context.Background(), "test", "password")
	if err != nil {
		t.Fatal(err)
	}

	login, err = as.GetLoginByToken(context.Background(), token)
	if err != nil {
		t.Fatal(err)
	}

	if login != "test" {
		t.Fatal("login != test")
	}

	ok := as.IsTokenValid(token)
	if !ok {
		t.Fatal("token is not valid")
	}

	u_data, err := as.UpdateUserData(context.Background(), token, "new_data")
	if err != nil {
		t.Fatal(err)
	}

	if u_data != "new_data" {
		t.Fatal("u_data != new_data")
	}

	u_data, err = as.GetUserData(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}

	if u_data != "new_data" {
		t.Fatal("u_data != new_data")
	}

	time.Sleep(time.Second * 6)

	ok = as.IsTokenValid(token)
	if ok {
		t.Fatal("token is valid")
	}
}
