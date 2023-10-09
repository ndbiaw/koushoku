package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	. "koushoku/config"
)

type ApiOptions struct {
	ApiKey string `json:"key"`
}

type PurgeCacheOptions struct {
	ApiOptions
	Archives    bool `json:"archives,omitempty"`
	Taxonomies  bool `json:"taxonomies,omitempty"`
	Templates   bool `json:"templates,omitempty"`
	Submissions bool `json:"submissions,omitempty"`
}

var ports []int

func scanPorts(startPort, endPort int) {
	if len(ports) > 0 {
		return
	}

	if startPort == 0 {
		startPort = 42073
	}

	if endPort == 0 {
		endPort = 42074
	}

	for port := startPort; port <= endPort && len(ports) < 4; port++ {
		conn, err := net.Dial("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)))
		if err != nil {
			continue
		}
		conn.Close()
		ports = append(ports, port)
	}
}

func purgeCaches(startPort, endPort int, opts PurgeCacheOptions) {
	scanPorts(startPort, endPort)
	opts.ApiKey = Config.HTTP.ApiKey

	buf, err := json.Marshal(opts)
	if err != nil {
		log.Fatalln(err)
	}

	for _, port := range ports {
		req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/api/purge-cache", port), bytes.NewBuffer(buf))
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode == 200 {
			log.Printf("Purged caches on port %d\n", port)
		} else {
			log.Fatalf("Failed to purge archives cache: %s", res.Status)
		}
	}
}

func reloadTemplates(startPort, endPort int) {
	scanPorts(startPort, endPort)
	opts := ApiOptions{ApiKey: Config.HTTP.ApiKey}

	buf, err := json.Marshal(opts)
	if err != nil {
		log.Fatalln(err)
	}

	for _, port := range ports {
		req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/api/reload-templates", port), bytes.NewBuffer(buf))
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode == 200 {
			log.Printf("Reloaded templates on port %d\n", port)
		} else {
			log.Fatalf("Failed to reload templates: %s", res.Status)
		}
	}
}
