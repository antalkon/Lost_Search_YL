package handler

import (
	"context"

	"github.com/mitchellh/mapstructure"
)

type createUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Data     string `json:"email"`
}

func createUserHandler(ctx context.Context, data map[string]any, as AuthServiceInterface) (map[string]any, error) {
	var cu createUser
	if err := mapstructure.Decode(data, &cu); err != nil {
		return nil, err
	}
	token, err := as.CreateUser(ctx, cu.Login, cu.Password, cu.Data)
	if err != nil {
		return nil, err
	}
	return map[string]any{"token": token}, nil
}

type loginUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func loginUserHandler(ctx context.Context, data map[string]any, as AuthServiceInterface) (map[string]any, error) {
	var lu loginUser
	if err := mapstructure.Decode(data, &lu); err != nil {
		return nil, err
	}
	token, err := as.LoginUser(ctx, lu.Login, lu.Password)
	if err != nil {
		return nil, err
	}
	return map[string]any{"token": token}, nil
}

type validateToken struct {
	Token string `json:"token"`
}

func validateTokenHandler(data map[string]any, as AuthServiceInterface) (map[string]any, error) {
	var vt validateToken
	if err := mapstructure.Decode(data, &vt); err != nil {
		return nil, err
	}
	return map[string]any{"valid": as.IsTokenValid(vt.Token)}, nil
}

type getUserData struct {
	Login string `json:"login"`
}

func getUserDataHandler(ctx context.Context, data map[string]any, as AuthServiceInterface) (map[string]any, error) {
	var gd getUserData
	if err := mapstructure.Decode(data, &gd); err != nil {
		return nil, err
	}
	u_data, err := as.GetUserData(ctx, gd.Login)
	if err != nil {
		return nil, err
	}
	return map[string]any{"email": u_data}, nil
}

type updateUserData struct {
	Token string `json:"token"`
	Data  string `json:"email"`
}

func updateUserDataHandler(ctx context.Context, data map[string]any, as AuthServiceInterface) (map[string]any, error) {
	var ud updateUserData
	if err := mapstructure.Decode(data, &ud); err != nil {
		return nil, err
	}
	u_data, err := as.UpdateUserData(ctx, ud.Token, ud.Data)
	if err != nil {
		return nil, err
	}
	return map[string]any{"email": u_data}, nil
}

type getLoginByToken struct {
	Token string `json:"token"`
}

func getLoginByTokenHandler(ctx context.Context, data map[string]any, as AuthServiceInterface) (map[string]any, error) {
	var gbt getLoginByToken
	if err := mapstructure.Decode(data, &gbt); err != nil {
		return nil, err
	}
	login, err := as.GetLoginByToken(ctx, gbt.Token)
	if err != nil {
		return nil, err
	}
	return map[string]any{"login": login}, nil
}
