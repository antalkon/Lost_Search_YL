package userrepo

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// названия столбцов в БД
const (
	_login    = "login"
	_password = "password"
	_data     = "data"
)

// названия таблиц в БД
const (
	_users = "users"
)

// конфигурация репозитория
type AuthRepoConfig struct {
	UserName string `env:"POSTGRES_USER" env-default:"root"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"123"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	DbName   string `env:"POSTGRES_DB" env-default:"db"`
}

// репозиторий
type UserRepo struct {
	db    *sqlx.DB
	cache *redis.Client //TODO: чтоб кэш работал
}

// пользователь
type User struct {
	Login    string `json:"login"` //без id, ибо логины уникальные
	Password string `json:"password"`
	Data     string `json:"data"`
}

// создание репозитория
func NewUserRepo(c AuthRepoConfig, cache *redis.Client) (*UserRepo, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		c.UserName, c.Password, c.DbName, c.Host, c.Port)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := db.Conn(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &UserRepo{db: db, cache: cache}, nil
}

// создание пользователя (добавление в БД)
func (ur *UserRepo) CreateUser(ctx context.Context, u User) error {
	_, err := sq.Insert(_users).Columns(_login, _password, _data).
		PlaceholderFormat(sq.Dollar).
		Values(u.Login, u.Password, u.Data).
		RunWith(ur.db).QueryContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

// получение пользователя из БД
var ErrUserNotFound = fmt.Errorf("user not found")

func (ur *UserRepo) GetUser(ctx context.Context, login string) (User, error) {
	res := User{}
	err := sq.Select(_login, _password, _data).From(_users).
		Where(sq.Eq{_login: login}).
		PlaceholderFormat(sq.Dollar).
		RunWith(ur.db).
		QueryRowContext(ctx).Scan(&res.Login, &res.Password, &res.Data)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	return res, nil
}

// обновление данных пользователя в БД
func (ur *UserRepo) UpdateUser(ctx context.Context, u User) (User, error) {
	res := User{}
	err := sq.Update(_users).
		Set(_login, u.Login).
		Set(_password, u.Password).
		Set(_data, u.Data).
		Where(sq.Eq{_login: u.Login}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		RunWith(ur.db).QueryRowContext(ctx).
		Scan(&res.Login, &res.Password, &res.Data)
	if err != nil {
		return User{}, err
	}
	return res, nil
}
