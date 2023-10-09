package main

import (
	"log"
	"net/http"

	"koushoku/cache"
	. "koushoku/config"
	"koushoku/server"
)

type ApiPayload struct {
	ApiKey string `json:"key"`
}

type PurgeCachePayload struct {
	ApiPayload
	Archives    bool `json:"archives"`
	Taxonomies  bool `json:"taxonomies"`
	Templates   bool `json:"templates"`
	Submissions bool `json:"submissions"`
}

func purgeCache(c *server.Context) {
	payload := &PurgeCachePayload{}
	if err := c.BindJSON(payload); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if payload.ApiKey != Config.HTTP.ApiKey {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if payload.Archives {
		log.Println("Purging archives cache...")
		cache.Archives.Purge()
	}

	if payload.Taxonomies {
		log.Println("Purging taxonomies cache...")
		cache.Taxonomies.Purge()
	}

	if payload.Templates {
		log.Println("Purging templates cache...")
		cache.Templates.Purge()
	}

	if payload.Submissions {
		log.Println("Purging submissions cache...")
		cache.Submissions.Purge()
	}
}

func reloadTemplates(c *server.Context) {
	payload := &ApiPayload{}
	if err := c.BindJSON(payload); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if payload.ApiKey != Config.HTTP.ApiKey {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	log.Println("Reloading templates...")
	server.LoadTemplates()
	cache.Templates.Purge()
}
