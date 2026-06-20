package redis

import (
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

func NewClient(host, port string) *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
}
