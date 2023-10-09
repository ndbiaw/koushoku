package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"

	. "koushoku/config"
	. "koushoku/services"

	"koushoku/database"
	"koushoku/models"
	"koushoku/modext"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"

	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	fBaseURL  = "https://www.fakku.net"
	iBaseURL  = "https://irodoricomics.com"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36"
)

var httpClient struct {
	*http.Client
	once sync.Once
}

func initHttpClient() {
	httpClient.once.Do(func() {
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatalln(err)
		}
		u, err := url.Parse("https://fakku.net")
		if err != nil {
			log.Fatalln(err)
		}
		httpClient.Client = &http.Client{Jar: jar}
		jar.SetCookies(u, []*http.Cookie{{
			Name:     "fakku_sid",
			Value:    Config.HTTP.Cookie,
			Domain:   "fakku.net",
			HttpOnly: true,
		}})
	})
}

func sendRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	if strings.Contains(url, "irodoricomics.com") {
		req.AddCookie(&http.Cookie{
			Name:   "irodori_splash",
			Value:  "1",
			Domain: ".irodoricomics.com",
		})
	}
	return httpClient.Do(req)
}

func searchF(model *models.Archive) (path string, err error) {
	var (
		res      *http.Response
		document *goquery.Document
	)

	res, err = sendRequest(fmt.Sprintf("%s/search/%s", fBaseURL, model.Slug))
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return
	}

	document, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	document.Find("body > div .grid > div[id^='content-']").
		EachWithBreak(func(i int, s *goquery.Selection) bool {
			titleElement := s.Find("a.text-lg").First()
			if Slugify(titleElement.Text()) != model.Slug {
				return true
			}

			artistSlug := Slugify(s.Find("a.text-sm").First().Text())
			if len(artistSlug) == 0 {
				return true
			}

			if v, ok := Aliases.ArtistMatches[artistSlug]; ok {
				artistSlug = Slugify(v)
			}

			for _, artist := range model.R.Artists {
				if artistSlug == artist.Slug {
					path, _ = titleElement.Attr("href")
					break
				}
			}
			return len(path) == 0
		})
	return
}

func scrapeF(fn, fnSlug, path string, model *models.Archive) (ok bool) {
	if len(path) == 0 {
		var err error
		path, err = searchF(model)
		if err != nil {
			log.Fatalln(err)
		}

		if len(path) == 0 {
			path = fmt.Sprintf("/hentai/%s-english", model.Slug)
		}
	}

	res, err := sendRequest(fBaseURL + path)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("[F] metadata not available:", fn)
		return
	}

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("[F] metadata found:", fn)
	metadata, ok := Metadatas.Map[fnSlug]
	if !ok {
		metadata = &Metadata{}
		Metadatas.Map[fnSlug] = metadata
	}

	metadata.Title = strings.TrimSpace(document.Find("body > div > div.grid > div > div > div[class*='table-cell'] > h1").Text())
	fields := document.Find("body > div > div.grid > div > div > div[class*='table-cell'] > .text-sm")
	fields.Each(func(i int, s *goquery.Selection) {
		if s.Children().Length() == 1 {
			return
		}

		section := strings.ToLower(s.Children().First().Text())
		if strings.Contains(section, "artist") {
			artists := strings.Split(strings.TrimSpace(s.Children().Last().Text()), ",")
			for _, artist := range artists {
				artist = strings.TrimSpace(artist)
				if v, ok := Aliases.ArtistMatches[Slugify(artist)]; ok {
					artist = v
				}

				duplicate := false
				for _, v := range metadata.Artists {
					if v == artist {
						duplicate = true
						break
					}
				}
				if !duplicate {
					metadata.Artists = append(metadata.Artists, artist)
				}
			}
		} else if strings.Contains(section, "circle") {
			circles := strings.Split(strings.TrimSpace(s.Children().Last().Text()), ",")
			for _, circle := range circles {
				circle = strings.TrimSpace(circle)
				if v, ok := Aliases.CircleMatches[Slugify(circle)]; ok {
					circle = v
				}

				duplicate := false
				for _, v := range metadata.Circles {
					if v == circle {
						duplicate = true
						break
					}
				}
				if !duplicate {
					metadata.Circles = append(metadata.Circles, circle)
				}
			}
		} else if strings.Contains(section, "parody") {
			parodies := strings.Split(strings.TrimSpace(s.Children().Last().Text()), ",")
			for _, parody := range parodies {
				parody = strings.TrimSpace(parody)
				if v, ok := Aliases.ParodyMatches[Slugify(parody)]; ok {
					parody = v
				}

				duplicate := false
				for _, v := range metadata.Parodies {
					if v == parody {
						duplicate = true
						break
					}
				}
				if !duplicate {
					metadata.Parodies = append(metadata.Parodies, parody)
				}
			}
		} else if strings.Contains(section, "magazine") {
			magazines := strings.Split(strings.TrimSpace(s.Children().Last().Text()), ",")
			for _, magazine := range magazines {
				magazine = strings.TrimSpace(magazine)
				if v, ok := Aliases.MagazineMatches[Slugify(magazine)]; ok {
					magazine = v
				}

				duplicate := false
				for _, v := range metadata.Magazines {
					if v == magazine {
						duplicate = true
						break
					}
				}
				if !duplicate {
					metadata.Magazines = append(metadata.Magazines, magazine)
				}
			}
		}
	})

	// Parse tags
	fields.Last().Children().First().Children().Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if len(href) > 0 {
			tag := strings.TrimSpace(s.Text())
			if v, ok := Aliases.TagMatches[Slugify(tag)]; ok {
				tag = v
			}

			duplicate := false
			for _, v := range metadata.Tags {
				if v == tag {
					duplicate = true
					break
				}
			}
			if !duplicate {
				metadata.Tags = append(metadata.Tags, tag)
			}
		}
	})
	return true
}

func searchI(model *models.Archive) (path string, err error) {
	var (
		res      *http.Response
		document *goquery.Document
	)

	res, err = sendRequest(fmt.Sprintf("%s/index.php?route=product/search&search=%s", iBaseURL, model.Title))
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return
	}

	document, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
	}

	entries := document.Find(".main-products > .product-layout")
	entries.EachWithBreak(func(i int, s *goquery.Selection) bool {
		titleElement := s.Find(".caption > .name a")
		if Slugify(titleElement.Text()) == model.Slug {
			return true
		}

		artistSlug := Slugify(s.Find(".caption > .stats span a").Text())
		if len(artistSlug) == 0 {
			return true
		}

		for _, artist := range model.R.Artists {
			if artistSlug == artist.Slug {
				path, _ = titleElement.Attr("href")
				break
			}
		}
		return len(path) == 0
	})

	if len(path) == 0 {
		document.EachWithBreak(func(i int, s *goquery.Selection) bool {
			titleElement := s.Find(".caption > .name a")
			if strings.Contains(Slugify(titleElement.Text()), model.Slug) {
				return true
			}

			artistSlug := Slugify(s.Find(".caption > .stats span a").Text())
			if len(artistSlug) == 0 {
				return true
			}

			for _, artist := range model.R.Artists {
				if artistSlug == artist.Slug {
					path, _ = titleElement.Attr("href")
					break
				}
			}
			return len(path) == 0
		})
	}
	return
}

func scrapeI(fn, fnSlug, path string, model *models.Archive) (ok bool) {
	if len(path) == 0 {
		var err error
		path, err = searchI(model)
		if err != nil {
			log.Fatalln(err)
		}

		if len(path) == 0 && len(model.R.Artists) == 1 {
			path = fmt.Sprintf("/%s/%s", model.R.Artists[0].Slug, model.Slug)
		}
	}

	if !strings.HasPrefix(path, "http") {
		path = iBaseURL + path
	}

	res, err := sendRequest(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("[I] metadata not available:", fn)
		return
	}

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Print("[I] metadata found:", fn)
	metadata, ok := Metadatas.Map[fnSlug]
	if !ok {
		metadata = &Metadata{}
		Metadatas.Map[fnSlug] = metadata
	}

	metadata.Title = strings.TrimSpace(document.Find("h1.title.page-title").Text())
	artists := document.Find(".product-manufacturer a")
	if artists.Length() > 0 {
		artists.Each(func(i int, s *goquery.Selection) {
			artist := strings.TrimSpace(s.Text())
			if v, ok := Aliases.ArtistMatches[Slugify(artist)]; ok {
				artist = v
			}

			duplicate := false
			for _, v := range metadata.Artists {
				if v == artist {
					duplicate = true
					break
				}
			}
			if !duplicate {
				metadata.Artists = append(metadata.Artists, artist)
			}
		})
	}

	productId, _ := document.Find("#product-id").Attr("value")
	res, err = sendRequest(fmt.Sprintf("%s/index.php?route=product/product/cattags&product_id=%s", iBaseURL, productId))
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("[I] failed to get tags:", fn)
	}

	document, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	tags := document.Find(".ctags")
	if tags.Length() > 0 {
		tags.Each(func(i int, s *goquery.Selection) {
			tag := strings.TrimSpace(s.Text())
			if v, ok := Aliases.TagMatches[Slugify(tag)]; ok {
				tag = v
			}

			duplicate := false
			for _, v := range metadata.Tags {
				if v == tag {
					duplicate = true
					break
				}
			}
			if !duplicate {
				metadata.Tags = append(metadata.Tags, tag)
			}
		})
	}
	return true
}

func scrapeMetadata() {
	initHttpClient()
	InitAliases()
	InitMetadatas()

	archives, err := models.Archives(
		Load(ArchiveRels.Artists),
		Load(ArchiveRels.Parodies),
		Load(ArchiveRels.Tags),
	).AllG()
	if err != nil {
		log.Fatalln(err)
	}

	total := len(archives)
	log.Println(fmt.Sprintf("%d archives found", total))

	c := make(chan bool, 10)
	defer close(c)

	var wg sync.WaitGroup
	wg.Add(total)

	for i, model := range archives {
		c <- true
		go func(i int, model *models.Archive) {
			defer func() {
				wg.Done()
				<-c
			}()

			fn := FileName(model.Path)
			fnSlug := Slugify(fn)

			if _, ok := Metadatas.Map[fnSlug]; ok {
				return
			}

			if !scrapeF(fn, fnSlug, "", model) {
				scrapeI(fn, fnSlug, "", model)
			}
		}(i, model)
	}
	wg.Wait()

	buf, err := json.Marshal(Metadatas.Map)
	if err == nil {
		err = os.WriteFile("metadata.json", buf, 755)
	}

	if err != nil {
		log.Fatalln(errors.WithStack(err))
	}
}

func scrapeMetadataById(id int64, fpath, ipath string) {
	initHttpClient()
	InitAliases()
	InitMetadatas()

	model, err := models.Archives(
		Where("id = ?", id),
		Load(ArchiveRels.Artists),
		Load(ArchiveRels.Parodies),
		Load(ArchiveRels.Tags),
	).OneG()
	if err != nil {
		log.Fatalln(err)
	}

	fn := FileName(model.Path)
	fnSlug := Slugify(fn)

	if !scrapeF(fn, fnSlug, fpath, model) {
		scrapeI(fn, fnSlug, ipath, model)
	}

	buf, err := json.Marshal(Metadatas.Map)
	if err == nil {
		err = os.WriteFile("metadata.json", buf, 755)
	}

	if err != nil {
		log.Fatalln(errors.WithStack(err))
	}
}

func importMetadata() {
	InitMetadatas()

	archives, err := models.Archives(
		Load(ArchiveRels.Artists),
		Load(ArchiveRels.Circles),
		Load(ArchiveRels.Magazines),
		Load(ArchiveRels.Parodies),
		Load(ArchiveRels.Tags),
	).AllG()
	if err != nil {
		log.Fatalln(err)
	}

	tx, err := database.Conn.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(archives))

	c := make(chan bool, 20)
	defer close(c)

	for _, model := range archives {
		c <- true
		go func(model *models.Archive) {
			defer func() {
				wg.Done()
				<-c
			}()

			fn := FileName(model.Path)
			fnSlug := Slugify(fn)

			metadata, ok := Metadatas.Map[fnSlug]
			if !ok {
				return
			}

			log.Println("Importing metadata of", fn)
			archive := modext.NewArchive(model).LoadRels(model)

			for _, artist := range metadata.Artists {
				slug := Slugify(artist)
				if v, ok := Aliases.ArtistMatches[slug]; ok {
					slug = Slugify(v)
					artist = v
				}

				isDuplicate := false
				for _, a := range archive.Artists {
					if a.Slug == slug {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					archive.Artists = append(archive.Artists,
						&modext.Artist{Name: artist})
				}
			}

			for _, circle := range metadata.Circles {
				slug := Slugify(circle)
				if v, ok := Aliases.CircleMatches[slug]; ok {
					slug = Slugify(v)
					circle = v
				}

				isDuplicate := false
				for _, c := range archive.Circles {
					if c.Slug == slug {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					archive.Circles = append(archive.Circles,
						&modext.Circle{Name: circle})
				}
			}

			for _, magazine := range metadata.Magazines {
				slug := Slugify(magazine)
				if v, ok := Aliases.MagazineMatches[slug]; ok {
					slug = Slugify(v)
					magazine = v
				}

				isDuplicate := false
				for _, m := range archive.Magazines {
					if m.Slug == slug {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					archive.Magazines = append(archive.Magazines,
						&modext.Magazine{Name: magazine})
				}
			}

			for _, parody := range metadata.Parodies {
				slug := Slugify(parody)
				if v, ok := Aliases.ParodyMatches[slug]; ok {
					slug = Slugify(v)
					parody = v
				}

				isDuplicate := false
				for _, p := range archive.Parodies {
					if p.Slug == slug {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					archive.Parodies = append(archive.Parodies,
						&modext.Parody{Name: parody})
				}
			}

			for _, tag := range metadata.Tags {
				slug := Slugify(tag)
				if v, ok := Aliases.TagMatches[slug]; ok {
					slug = Slugify(v)
					tag = v
				}

				isDuplicate := false
				for _, t := range archive.Tags {
					if t.Slug == slug {
						isDuplicate = true
						break
					}
				}

				if !isDuplicate {
					archive.Tags = append(archive.Tags,
						&modext.Tag{Name: tag})
				}
			}

			if len(metadata.Title) > 0 && metadata.Title != archive.Title {
				model.Title = metadata.Title
				model.Slug = Slugify(model.Title)

				if v, ok := Aliases.ArchiveMatches[model.Slug]; ok {
					model.Slug = Slugify(v)
					model.Title = v
				}

				if err := model.Update(tx, boil.Whitelist(ArchiveCols.Title, ArchiveCols.Slug)); err != nil {
					log.Fatalln(err)
				}
			}

			if err := PopulateArchiveRels(tx, model, archive); err != nil {
				log.Fatalln(err)
			}
		}(model)
	}
	wg.Wait()

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}
}
