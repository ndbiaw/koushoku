package services

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	. "koushoku/config"
)

type Metadata struct {
	Title     string
	Artists   []string
	Circles   []string
	Magazines []string
	Parodies  []string
	Tags      []string
}

var Metadatas struct {
	Map  map[string]*Metadata
	once sync.Once
}

func InitMetadatas() {
	Metadatas.once.Do(func() {
		Metadatas.Map = make(map[string]*Metadata)
		path := filepath.Join(Config.Paths.Metadata)

		stat, err := os.Stat(path)
		if os.IsNotExist(err) || stat.IsDir() {
			return
		}

		buf, err := os.ReadFile(path)
		if err != nil {
			log.Println(err)
			return
		}

		if err := json.Unmarshal(buf, &Metadatas.Map); err != nil {
			log.Println(err)
		}
	})
}
