package cache

import (
	"encoding/json"
	"fmt"
	"review/internal/domain/model"
	"time"

	"github.com/go-redis/redis"
)

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-required:"true"`
	Port     string `env:"REDIS_PORT" env-required:"true"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
}

func New(cfg RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
	})
	if err := client.Ping().Err(); err != nil {
		panic(fmt.Errorf("failed to connect Redis:%w", err))
	}
	return client
}

type RedisRepository struct {
	Client *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	return &RedisRepository{Client: redisClient}
}

func (r *RedisRepository) GetReviews(email string) ([]model.Review, error, bool) {
	const op = "storage.cache.GetReviews"
	key := "reviews:" + email

	cacheData, err := r.Client.Get(key).Result()
	if err == redis.Nil {
		return nil, nil, false
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err), false
	}

	var reviews []model.Review
	if err := json.Unmarshal([]byte(cacheData), &reviews); err != nil {
		_ = r.Client.Del(key)
		return nil, fmt.Errorf("%s: invalid cache data: %w", op, err), false
	}

	return reviews, nil, true

}

func (r *RedisRepository) SaveReviews(email string, reviews []model.Review) error {
	const op = "storage.cache.SaveReviews"
	key := "reviews:" + email

	data, err := json.Marshal(reviews)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := r.Client.Set(key, data, 1*time.Minute).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}


