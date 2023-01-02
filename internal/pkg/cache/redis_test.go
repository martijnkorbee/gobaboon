package cache

import (
	"testing"
)

func Test_redis_CreateRedisCache(t *testing.T) {
	_, err := CreateRedisCache(RedisConfig{
		Prefix:   "test-baboon",
		Host:     testRedisClient.Host(),
		Port:     testRedisClient.Port(),
		Password: "",
	})
	if err != nil {
		t.Error("failed to connect to redis:", err)
	}
}
