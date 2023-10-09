package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	. "koushoku/config"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client
var Archives *LRU
var Taxonomies *LRU
var Templates *LRU
var Submissions *LRU

var Users *LRU
var Favorites *LRU
var Cache *LRU

const defaultExpr = time.Hour
const templateExpr = 5 * time.Minute

func Init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Config.Redis.Host, Config.Redis.Port),
		DB:       Config.Redis.DB,
		Password: Config.Redis.Passwd,
	})
	if result := Redis.Ping(context.Background()); result.Err() != nil {
		log.Fatalln(result.Err())
	}

	Archives = New(4096, defaultExpr)
	Taxonomies = New(4096, defaultExpr)
	Templates = New(4096, templateExpr)
	Submissions = New(4096, defaultExpr)
	Users = New(2048, defaultExpr)
	Favorites = New(2048, defaultExpr)
	Cache = New(512, defaultExpr)
}
