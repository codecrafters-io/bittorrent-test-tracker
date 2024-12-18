package redis

import (
	"crypto/tls"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	redigolib "github.com/gomodule/redigo/redis"
)

// redisBackend represents a redis handler.
type redisBackend struct {
	pool    *redigolib.Pool
	redsync *redsync.Redsync
}

// newRedisBackend creates a redisBackend instance.
func newRedisBackend(cfg *Config, u *redisURL, socketPath string) *redisBackend {
	rc := &redisConnector{
		URL:            u,
		SocketPath:     socketPath,
		ReadTimeout:    cfg.RedisReadTimeout,
		WriteTimeout:   cfg.RedisWriteTimeout,
		ConnectTimeout: cfg.RedisConnectTimeout,
	}
	pool := rc.NewPool()
	redsync := redsync.New(redigo.NewPool(pool))
	return &redisBackend{
		pool:    pool,
		redsync: redsync,
	}
}

// open returns or creates instance of Redis connection.
func (rb *redisBackend) open() redigolib.Conn {
	return rb.pool.Get()
}

type redisConnector struct {
	URL            *redisURL
	SocketPath     string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	ConnectTimeout time.Duration
}

// NewPool returns a new pool of Redis connections
func (rc *redisConnector) NewPool() *redigolib.Pool {
	return &redigolib.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigolib.Conn, error) {
			c, err := rc.open()
			if err != nil {
				return nil, err
			}

			if rc.URL.DB != 0 {
				_, err = c.Do("SELECT", rc.URL.DB)
				if err != nil {
					return nil, err
				}
			}

			return c, err
		},
		// PINGs connections that have been idle more than 10 seconds
		TestOnBorrow: func(c redigolib.Conn, t time.Time) error {
			if time.Since(t) < 10*time.Second {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// Open a new Redis connection
func (rc *redisConnector) open() (redigolib.Conn, error) {
	opts := []redigolib.DialOption{
		redigolib.DialDatabase(rc.URL.DB),
		redigolib.DialReadTimeout(rc.ReadTimeout),
		redigolib.DialWriteTimeout(rc.WriteTimeout),
		redigolib.DialConnectTimeout(rc.ConnectTimeout),
	}

	if rc.URL.Password != "" {
		opts = append(opts, redigolib.DialPassword(rc.URL.Password))
	}

	if rc.URL.UseTLS {
		opts = append(opts, redigolib.DialUseTLS(true))
		opts = append(opts, redigolib.DialTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		}))
	}

	if rc.SocketPath != "" {
		return redigolib.Dial("unix", rc.SocketPath, opts...)
	}

	return redigolib.Dial("tcp", rc.URL.Host, opts...)
}

// A redisURL represents a parsed redisURL
// The general form represented is:
//
//	redis://[password@]host][/][db]
type redisURL struct {
	Host     string
	Password string
	DB       int
	UseTLS   bool
}

// parseRedisURL parse rawurl into redisURL
func parseRedisURL(target string) (*redisURL, error) {
	var u *url.URL
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "redis" && u.Scheme != "rediss" {
		return nil, errors.New("no redis/rediss scheme found")
	}

	db := 0 // default redis db
	parts := strings.Split(u.Path, "/")
	if len(parts) != 1 {
		db, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
	}

	return &redisURL{
		Host:     u.Host,
		Password: u.User.String(),
		DB:       db,
		UseTLS:   u.Scheme == "rediss", // (redis:// is unencrypted, rediss:// is encrypted)
	}, nil
}
