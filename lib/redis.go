package lib

import "github.com/go-redis/redis/v8"

type Redis struct {
	*redis.Client
}

func NewRedis(env Env) Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     env.RedisCredentials.Address,
		Password: env.RedisCredentials.Password,
		DB:       env.RedisCredentials.Database,
	})

	return Redis{Client: client}
}
