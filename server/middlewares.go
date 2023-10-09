package server

import (
	"log"
	"net/http"

	"koushoku/cache"

	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

func WithName(name string) Handler {
	return func(c *Context) {
		c.SetData("name", name)
		c.Next()
	}
}

func WithRedirect(relativePath string) Handler {
	return func(c *Context) {
		c.Redirect(http.StatusFound, relativePath)
	}
}

var limiters = make(map[string]Handler)

func WithRateLimit(prefix, formatted string) Handler {
	handler, ok := limiters[prefix]
	if !ok {
		rate, err := limiter.NewRateFromFormatted(formatted)
		if err != nil {
			log.Fatalln(err)
		}

		store, err := redis.NewStoreWithOptions(cache.Redis, limiter.StoreOptions{
			Prefix: prefix,
		})
		if err != nil {
			log.Fatalln(err)
		}

		instance := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
		m := mgin.NewMiddleware(instance)

		handler = func(c *Context) { m(c.Context) }
		limiters[prefix] = handler
	}
	return handler
}
