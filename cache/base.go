package cache

import (
	"eshop_cart/log"
	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	url    string
)

func init() {
	url = "redis://root:@117.72.72.114:16379/0"
	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Errorf("err: %v", err)
		panic(err)
	}

	Client = redis.NewClient(opt)
}
