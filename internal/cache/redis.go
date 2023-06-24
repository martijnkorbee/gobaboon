package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisCache struct {
	Pool   *redis.Pool
	Prefix string
}

type RedisConfig struct {
	Prefix   string
	Host     string
	Port     string
	Password string
}

func CreateRedisCache(c RedisConfig) (*RedisCache, error) {
	var (
		url = fmt.Sprintf("%s:%s", c.Host, c.Port)
	)

	pool := redis.Pool{
		MaxIdle:     50,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", url, redis.DialPassword(c.Password))
		},
	}

	// test connection
	conn := pool.Get()
	defer conn.Close()

	ok, err := conn.Do("PING")
	if err != nil || ok != "PONG" {
		return nil, errors.New(err.Error())
	}

	return &RedisCache{
		Pool:   &pool,
		Prefix: c.Prefix,
	}, nil
}

func (r *RedisCache) GetConnection() any {
	return r.Pool
}

func (r *RedisCache) Close() error {
	return r.Pool.Close()
}

func (c *RedisCache) Has(key string) (bool, error) {
	key = fmt.Sprintf("%s:%s", c.Prefix, key)

	conn := c.Pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (c *RedisCache) Get(key string) (interface{}, error) {
	key = fmt.Sprintf("%s:%s", c.Prefix, key)

	conn := c.Pool.Get()
	defer conn.Close()

	// get cache entry
	cacheEntry, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	// decode entry
	decoded, err := decode(string(cacheEntry))
	if err != nil {
		return nil, err
	}

	return decoded[key], nil
}

func (c *RedisCache) Set(key string, value interface{}, expires ...int) error {
	key = fmt.Sprintf("%s:%s", c.Prefix, key)

	conn := c.Pool.Get()
	defer conn.Close()

	entry := Entry{}
	entry[key] = value
	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(expires) > 0 {
		_, err = conn.Do("SETEX", key, expires[0], string(encoded))
		if err != nil {
			return err
		}
	} else {
		_, err = conn.Do("SET", key, string(encoded))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisCache) Forget(key string) error {
	key = fmt.Sprintf("%s:%s", c.Prefix, key)

	conn := c.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) EmptyByMatch(key string) error {
	key = fmt.Sprintf("%s:%s", c.Prefix, key)

	conn := c.Pool.Get()
	defer conn.Close()

	// get keys
	keys, err := c.getKeys(key)
	if err != nil {
		return err
	}

	for _, x := range keys {
		_, err := conn.Do("DEL", x)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisCache) Empty() error {
	key := fmt.Sprintf("%s:", c.Prefix)

	conn := c.Pool.Get()
	defer conn.Close()

	// get keys
	keys, err := c.getKeys(key)
	if err != nil {
		return err
	}

	for _, x := range keys {
		_, err := conn.Do("DEL", x)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisCache) getKeys(pattern string) ([]string, error) {
	conn := c.Pool.Get()
	defer conn.Close()

	i := 0
	keys := []string{}

	for {
		arr, err := redis.Values(conn.Do("SCAN", i, "MATCH", fmt.Sprintf("%s*", pattern)))
		if err != nil {
			return keys, err
		}

		i, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if i == 0 {
			break
		}
	}

	return keys, nil
}
