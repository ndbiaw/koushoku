package services

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"

	. "koushoku/config"
)

var Blacklists struct {
	ArchiveMatches   map[string]bool
	ArchiveWildcards []string
	ArtistMatches    map[string]bool
	CircleMatches    map[string]bool
	MagazineMatches  map[string]bool
	TagMatches       map[string]bool

	once sync.Once
}

func InitBlacklists() {
	Blacklists.once.Do(func() {
		Blacklists.ArchiveMatches = make(map[string]bool)
		Blacklists.ArtistMatches = make(map[string]bool)
		Blacklists.CircleMatches = make(map[string]bool)
		Blacklists.MagazineMatches = make(map[string]bool)
		Blacklists.TagMatches = make(map[string]bool)

		stat, err := os.Stat(Config.Paths.Blacklist)
		if os.IsNotExist(err) || stat.IsDir() {
			return
		}

		f, err := os.Open(Config.Paths.Blacklist)
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
			if len(strs) < 2 {
				continue
			}

			v := Slugify(strings.Join(strs[1:], ":"))

			switch strings.TrimSpace(strs[0]) {
			case "title":
				Blacklists.ArchiveMatches[v] = true
			case "title*":
				Blacklists.ArchiveWildcards = append(Blacklists.ArchiveWildcards, v)
			case "artist":
				Blacklists.ArtistMatches[v] = true
			case "circle":
				Blacklists.CircleMatches[v] = true
			case "magazine":
				Blacklists.MagazineMatches[v] = true
			case "tag":
				Blacklists.TagMatches[v] = true
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	})
}
