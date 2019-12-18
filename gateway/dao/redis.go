package dao

import "github.com/go-redis/redis"

func NewRedisClient() (*redis.Client, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
