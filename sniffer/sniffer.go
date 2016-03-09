package sniffer

import (
	"fmt"
	"sync"

	"github.com/garyburd/redigo/redis"
)

// Result has a Key and TTL
type Result struct {
	Key string
	TTL int64
}

func (result *Result) String() string {
	return fmt.Sprintf("Key: %v, TTL: %d", result.Key, result.TTL)
}

// Sniffer gets the TTLs of random keys
type Sniffer struct {
	redisURI      string
	redisConn     redis.Conn
	redisConnOnce sync.Once
}

// New constructs a new Sniffer
func New(redisURI string) *Sniffer {
	return &Sniffer{redisURI: redisURI}
}

// Sniff returns a result
func (sniffer *Sniffer) Sniff() (*Result, error) {
	conn, err := redis.DialURL(sniffer.redisURI)
	if err != nil {
		return nil, err
	}

	keyResult, err := conn.Do("RANDOMKEY")
	if err != nil {
		return nil, err
	}

	Key := string(keyResult.([]byte))

	ttlResult, err := conn.Do("TTL", Key)
	if err != nil {
		return nil, err
	}

	TTL := ttlResult.(int64)

	return &Result{Key, TTL}, nil
}

func (sniffer *Sniffer) redis() (redis.Conn, error) {
	var err error

	sniffer.redisConnOnce.Do(func() {
		var conn redis.Conn
		conn, err = redis.DialURL(sniffer.redisURI)
		sniffer.redisConn = conn
	})

	if err != nil {
		return nil, err
	}

	return sniffer.redisConn, nil
}
