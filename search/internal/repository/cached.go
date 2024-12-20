package repository

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"gitlab.crja72.ru/gospec/go21/go_final_project/pkg/db/postgres"
	cache "gitlab.crja72.ru/gospec/go21/go_final_project/pkg/db/redis"
)

type CachedPostgres struct {
	postgres *PostgresRepo
	redis    *redis.Client
}

func NewCachedRepo(postgresCfg postgres.PostgresConfig, redisCfg cache.RedisConfig) (*CachedPostgres, error) {
	postgres, err := NewPostgresRepo(postgresCfg)
	if err != nil {
		return nil, err
	}

	redis, err := cache.New(redisCfg)
	if err != nil {
		return nil, err
	}

	return &CachedPostgres{
		postgres: postgres,
		redis:    redis,
	}, nil
}

func (c *CachedPostgres) AddFind(req AddReq) (AddResp, error) {
	return c.postgres.AddFind(req)
}

func (c *CachedPostgres) RespondToFind(req RespondReq) (RespondResp, error) {
	return c.postgres.RespondToFind(req)
}

func (c *CachedPostgres) GetFind(req GetReq) (GetResp, error) {
	var (
		hash = getHash(fmt.Sprintf("%+v", req))
		r    GetResp
	)

	val, err := c.redis.Get(context.Background(), hash).Result()
	if err != nil {
		resp, err := c.postgres.GetFind(req)
		{
			data, err := json.Marshal(resp)
			if err != nil {
				return GetResp{}, fmt.Errorf("CachedPostgres GetFind: marshal data: %w", err)
			}
			if err := c.redis.Set(context.Background(), hash, data, cache.DefaultExpiration).Err(); err != nil {
				log.Printf("redis set error: %s\n", err.Error())
				return GetResp{}, fmt.Errorf("CachedPostgres GetFind: set data to redis: %w", err)
			}

		}
		if err == nil {
			return resp, err
		}
	}
	if err := json.Unmarshal([]byte(val), &r); err != nil {
		return GetResp{}, fmt.Errorf("CachedPostgres GetFind: can't unmarshal data from redis: %w", err)
	}
	log.Println("use cached value")
	return r, nil
}

func getHash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
