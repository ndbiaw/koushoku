package services

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"

	. "koushoku/config"
)

var Aliases struct {
	ArchiveMatches  map[string]string
	ArtistMatches   map[string]string
	CircleMatches   map[string]string
	MagazineMatches map[string]string
	ParodyMatches   map[string]string
	TagMatches      map[string]string

	once sync.Once
}

func InitAliases() {
	Aliases.once.Do(func() {
		Aliases.ArchiveMatches = make(map[string]string)
		Aliases.ArtistMatches = make(map[string]string)
		Aliases.CircleMatches = make(map[string]string)
		Aliases.MagazineMatches = make(map[string]string)
		Aliases.ParodyMatches = make(map[string]string)
		Aliases.TagMatches = make(map[string]string)

		stat, err := os.Stat(Config.Paths.Alias)
		if os.IsNotExist(err) || stat.IsDir() {
			return
		}

		f, err := os.Open(Config.Paths.Alias)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(line) == 0 {
				continue
			}

			strs := strings.Split(strings.ToLower(line), ":")
			if len(strs) < 3 {
				continue
			}

			k := Slugify(strs[1])
			v := strings.TrimSpace(strings.Join(strs[2:], ":"))

			switch strings.TrimSpace(strs[0]) {
			case "title":
				Aliases.ArchiveMatches[k] = v
			case "artist":
				Aliases.ArtistMatches[k] = v
			case "circle":
				Aliases.CircleMatches[k] = v
			case "magazine":
				Aliases.MagazineMatches[k] = v
			case "parody":
				Aliases.ParodyMatches[k] = v
			case "tag":
				Aliases.TagMatches[k] = v
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	})
}
