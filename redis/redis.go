package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type RedisConnection struct {
	redis.Conn
}

func NewRedisConnection(uri string) *RedisConnection {
	c, err := redis.DialURL(uri)
	if err != nil {
		panic(err)
	}

	return &RedisConnection{c}
}

// FlushAll flushes all the key in the redis db
func (c *RedisConnection) FlushAll() {
	c.Do("FLUSHALL")
}

func (c *RedisConnection) StartTx() {
	c.Do("MULTI")
}

func (c *RedisConnection) ExecuteTx() {
	c.Do("EXEC")
}

// GetValue gets the value saved at a given key
func (c *RedisConnection) GetValue(key string) (string, error) {
	value, err := redis.String(c.Do("GET", key))
	if err != nil {
		return "", err
	}

	return value, nil
}

// GetSortedSet fetches complete sorted set using ZRANGE query.
// Returns map with value and its rank in set.
func (c *RedisConnection) GetSortedSet(key string) (map[string]float64, error) {
	resMap := make(map[string]float64)
	res, err := redis.Strings(c.Do("ZRANGE", key, "0", "-1", "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(res); i = i + 2 {
		resMap[res[i]], _ = strconv.ParseFloat(res[i+1], 64)
	}

	return resMap, nil
}

// Exists checks if a key exists in redis
func (c *RedisConnection) Exists(key string) bool {
	exists, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// ZAdd inserts value in sorted set.
// Cmd Returns: number of insertions and error
// Returns: error
func (c *RedisConnection) ZAdd(key string, rank int64, member string) error {
	_, err := redis.Int64(c.Do("ZADD", key, "NX", rank, member))
	return err
}

// ZRem removes value in sorted set.
// Cmd Returns: number of deletions and error
// Returns: error
func (c *RedisConnection) ZRem(key string, member string) error {
	_, err := redis.Int64(c.Do("ZREM", key, member))
	return err
}

func (c *RedisConnection) ZCount(key string) (int64, error) {
	count, err := redis.Int64(c.Do("ZCOUNT", key, "-inf", "+inf"))
	if err != nil {
		return 0, err
	}

	return count, nil
}

// IncrBy increment value of a key by passed amount. Returns: currentValue of key
func (c *RedisConnection) IncrBy(key string, value int64) (int64, error) {
	return redis.Int64(c.Do("INCRBY", key, value))
}

// Set sets the value of a key to passed key.
// Cmd Returns: "OK" if successfull and error
//Returns error if error occured
func (c *RedisConnection) Set(key string, value uint64) error {
	ok, err := redis.String(c.Do("SET", key, value))
	if err != nil {
		return err
	} else if ok != "OK" {
		return fmt.Errorf("Some error occured while running SET command on key: %v", key)
	}
	return nil
}

func (c *RedisConnection) Get(key string) (uint64, error) {
	return redis.Uint64(c.Do("GET", key))
}

// ZRangeByLex executes ZRANGEBYLEX expecting []string as return
func (c *RedisConnection) GetKeys(key string) ([]string, error) {
	return redis.Strings(c.Do("KEYS", key))
}

// Del removes given key from redis
// Cmd Returns: number of deletions and error
// Returns: error
func (c *RedisConnection) Del(key string) error {
	_, err := redis.Int64(c.Do("DEL", key))
	return err
}

// ZRangeByLex executes ZRANGEBYLEX expecting []string as return
func (c *RedisConnection) ZRangeByLex(key, min, max string) ([]string, error) {
	return redis.Strings(c.Do("ZRANGEBYLEX", key, min, max))
}

// ZRangeByLexInt executes ZRANGEBYLEX expecting []int64 as return
func (c *RedisConnection) ZRangeByLexInt(key, min, max string) ([]int64, error) {
	return redis.Int64s(c.Do("ZRANGEBYLEX", key, min, max))
}

// ZRevRangeByLex executes ZREVRANGEBYLEX expecting []int64 as return
func (c *RedisConnection) ZRevRangeByLex(key, min, max string) ([]string, error) {
	return redis.Strings(c.Do("ZREVRANGEBYLEX", key, min, max))
}

// ZRevRangeByLexInt executes ZREVRANGEBYLEX expecting []int64 as return
func (c *RedisConnection) ZRevRangeByLexInt(key, min, max string) ([]int64, error) {
	return redis.Int64s(c.Do("ZREVRANGEBYLEX", key, min, max))
}

// Sort executes SORT command. Returns byteslices [][]byte and error
func (c *RedisConnection) Sort(key, by string, alpha, desc bool, get ...string) ([][]byte, error) {
	args := []interface{}{key}
	if by != "" {
		args = append(args, "BY", by)
	}
	for _, g := range get {
		args = append(args, "GET", g)
	}
	if alpha {
		args = append(args, "ALPHA")
	}
	if desc {
		args = append(args, "DESC")
	} else {
		args = append(args, "ASC")
	}
	return redis.ByteSlices(c.Do("SORT", args...))
}

// Keys returns the keys stored in redis with specified pattern
func (c *RedisConnection) Keys(pattern string) (res []string, err error) {
	return redis.Strings(c.Do("KEYS", pattern))
}

// MGet returns the value for keys passed
func (c *RedisConnection) MGet(keys ...string) (res []string, err error) {
	args := make([]interface{}, len(keys))
	for _, key := range keys {
		args = append(args, key)
	}
	return redis.Strings(c.Do("MGET", args...))
}
