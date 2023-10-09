package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	. "koushoku/config"
	. "koushoku/services"

	"koushoku/models"
	"koushoku/modext"

	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	archiveRgx = regexp.MustCompile(`(\(|\[|\{)?[^\(\[\{\}\]\)]+(\}\)|\])?`)
	miscRgx    = regexp.MustCompile(`(?i)(fakku|irodori comics|x?\d+00x?)`)
)

func getArchivePaths() (paths []string, err error) {
	walkFn := func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() ||
			!(strings.HasSuffix(path, ".zip") ||
				strings.HasSuffix(path, ".cbz") ||
				strings.HasSuffix(path, ".rar")) {
			return err
		}
		paths = append(paths, path)
		return nil
	}
	return paths, filepath.Walk(Config.Directories.Data, walkFn)
}

func populateArchive(archive *modext.Archive) error {
	fileName := FileName(archive.Path)
	if stat, err := os.Stat(archive.Path); err == nil {
		archive.Size = stat.Size()
	} else {
		return err
	}

	var (
		artists   = make(map[string]string)
		circles   = make(map[string]string)
		magazines = make(map[string]string)
		parodies  = make(map[string]string)
		tags      = make(map[string]string)
	)

	if metadata, ok := Metadatas.Map[Slugify(fileName)]; ok {
		for _, parody := range metadata.Parodies {
			slug := Slugify(parody)
			if v, ok := Aliases.ParodyMatches[slug]; ok {
				slug = Slugify(v)
				parody = v
			}
			parodies[slug] = parody
		}

		for _, tag := range metadata.Tags {
			slug := Slugify(tag)
			if v, ok := Aliases.TagMatches[slug]; ok {
				tag = v
				slug = Slugify(v)
			}
			if _, ok := Blacklists.TagMatches[slug]; ok {
				return nil
			}
			tags[slug] = tag
		}
	}

	matches := archiveRgx.FindAllString(fileName, -1)
	if len(matches) == 0 {
		return nil
	}

	var title string
	for i, match := range matches {
		match = strings.TrimSpace(match)
		if len(match) == 0 {
			continue
		}

		if strings.HasPrefix(match, "[") {
			if i == 0 {
				match = strings.TrimSuffix(strings.TrimPrefix(match, "["), "]")
				if match = strings.TrimSpace(match); len(match) == 0 {
					continue
				}

				names := strings.Split(match, ",")
				for _, name := range names {
					if name = strings.TrimSpace(name); len(name) > 0 {
						artists[Slugify(name)] = name
					}
				}
			}
		} else if strings.HasPrefix(match, "(") {
			if i == 1 {
				match = strings.TrimSuffix(strings.TrimPrefix(match, "("), ")")
				if match = strings.TrimSpace(match); len(match) == 0 {
					continue
				}

				if len(artists) > 0 {
					for k, v := range artists {
						circles[k] = v
						delete(artists, k)
					}
				}

				names := strings.Split(match, ",")
				for _, name := range names {
					if name = strings.TrimSpace(name); len(name) > 0 {
						artists[Slugify(name)] = name
					}
				}
			} else if i == 2 || i == 3 {
				match = strings.TrimSuffix(strings.TrimPrefix(match, "("), ")")
				if match = strings.TrimSpace(match); len(match) == 0 {
					continue
				}

				if i < len(matches)-1 {
					next := matches[i+1]
					if len(next) > 0 &&
						!(strings.HasPrefix(match, "[") ||
							strings.HasPrefix(match, "(") ||
							strings.HasPrefix(next, "{")) {
						continue
					}
				}

				if strings.HasPrefix(match, "x") ||
					strings.EqualFold(match, "temp") ||
					strings.EqualFold(match, "strong") ||
					strings.EqualFold(match, "complete") {
					continue
				}

				names := strings.Split(match, ", ")
				for _, name := range names {
					if name = strings.TrimSpace(name); len(name) > 0 {
						magazines[Slugify(name)] = name
					}
				}
			}
		} else if strings.HasPrefix(match, "{") {
			match = strings.TrimSuffix(strings.TrimPrefix(match, "{"), "}")
			match = strings.TrimSpace(match)

			if len(match) == 0 ||
				strings.Contains(Slugify(match), "comic") ||
				strings.Contains(Slugify(match), "2d-market") {
				continue
			}

			// Fix inconsistent tag names
			match = strings.ReplaceAll(match, "zero gravity", "zero-gravity")
			match = strings.ReplaceAll(match, "dark skin", "dark-skin")
			match = strings.ReplaceAll(match, "heart pupil", "heart-pupil")

			names := strings.Split(match, " ")
			for _, name := range names {
				name = strings.TrimSpace(name)
				if len(name) > 0 {
					tags[Slugify(name)] = strings.ReplaceAll(name, "-", " ")
				}
			}
		} else if i == 1 || i == 2 {
			match = strings.TrimSpace(miscRgx.ReplaceAllString(match, ""))
			if len(match) > 0 {
				title = match
			}
		}
	}

	if len(title) == 0 {
		return nil
	}

	titleSlug := Slugify(title)
	if v, ok := Aliases.ArchiveMatches[titleSlug]; ok {
		titleSlug = Slugify(title)
		title = v
	}

	if _, ok := Blacklists.ArchiveMatches[titleSlug]; ok {
		return nil
	}

	for _, v := range Blacklists.ArchiveWildcards {
		if strings.Contains(titleSlug, v) {
			return nil
		}
	}

	for slug, artist := range artists {
		if v, ok := Aliases.ArtistMatches[slug]; ok {
			slug = Slugify(v)
			artist = v
		}
		if _, ok := Blacklists.ArtistMatches[slug]; ok {
			return nil
		}
		archive.Artists = append(archive.Artists,
			&modext.Artist{Slug: slug, Name: artist})
	}

	for slug, circle := range circles {
		if v, ok := Aliases.CircleMatches[slug]; ok {
			slug = Slugify(v)
			circle = v
		}
		if _, ok := Blacklists.CircleMatches[slug]; ok {
			return nil
		}
		archive.Circles = append(archive.Circles,
			&modext.Circle{Slug: slug, Name: circle})
	}

	for slug, magazine := range magazines {
		if v, ok := Aliases.MagazineMatches[slug]; ok {
			slug = Slugify(v)
			magazine = v
		}
		if _, ok := Blacklists.MagazineMatches[slug]; ok {
			return nil
		}
		archive.Magazines = append(archive.Magazines,
			&modext.Magazine{Slug: slug, Name: magazine})
	}

	for slug, parody := range parodies {
		if v, ok := Aliases.ParodyMatches[slug]; ok {
			slug = Slugify(v)
			parody = v
		}
		archive.Parodies = append(archive.Parodies,
			&modext.Parody{Slug: slug, Name: parody})
	}

	for slug, tag := range tags {
		if v, ok := Aliases.TagMatches[slug]; ok {
			slug = Slugify(v)
			tag = v
		}
		if _, ok := Blacklists.TagMatches[slug]; ok {
			return nil
		}

		isDuplicate := false
		for _, t := range archive.Tags {
			if slug == Slugify(t.Name) {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			archive.Tags = append(archive.Tags,
				&modext.Tag{Slug: slug, Name: tag})
		}
	}

	zf, err := zip.OpenReader(archive.Path)
	if err != nil {
		if err == zip.ErrFormat {
			log.Println(err, archive.Path)
			return nil
		}
		return err
	}
	defer zf.Close()

	for _, f := range zf.File {
		stat := f.FileInfo()
		name := stat.Name()

		if stat.IsDir() || !IsImage(name) {
			continue
		}

		archive.Pages++
		if archive.CreatedAt == 0 {
			archive.CreatedAt = stat.ModTime().Unix()
		}
	}

	if archive.Pages == 0 {
		return nil
	}

	archive.Title = title
	archive.Slug = titleSlug

	return nil
}

func moderateArchives() {
	InitBlacklists()

	archives, err := models.Archives(Load(ArchiveRels.Artists), Load(ArchiveRels.Tags)).AllG()
	if err != nil {
		log.Fatalln(err)
	}

	for _, archive := range archives {
		titleSlug := Slugify(archive.Title)
		_, isRemove := Blacklists.ArchiveMatches[titleSlug]

		if archive.R != nil && len(archive.R.Artists) > 0 {
			for _, artist := range archive.R.Artists {
				if _, ok := Blacklists.ArtistMatches[Slugify(artist.Name)]; ok {
					artist.DeleteG()
					isRemove = true
				}
			}
		}

		if !isRemove {
			for _, slug := range Blacklists.ArchiveWildcards {
				if strings.Contains(titleSlug, slug) {
					isRemove = true
					break
				}
			}
		}

		if !isRemove && archive.R != nil && len(archive.R.Tags) > 0 {
			for _, tag := range archive.R.Tags {
				if _, ok := Blacklists.TagMatches[Slugify(tag.Name)]; ok {
					tag.DeleteG()
					isRemove = true
				}
			}
		}

		if isRemove {
			log.Println("Removing archive", archive.Path)
			DeleteArchive(archive.ID)
		}
	}
}

func indexArchive(path string, reindex bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}

	archive := &modext.Archive{Path: path}
	log.Println("Populating archive", filepath.Base(path))

	if err := populateArchive(archive); err != nil {
		log.Fatalln(err)
	}

	if len(archive.Title) == 0 {
		return
	}

	log.Println("Indexing archive", filepath.Base(path))

	var model *modext.Archive
	var err error

	if reindex {
		model, err = UpdateArchive(archive)
	} else {
		model, err = CreateArchive(archive)
	}

	if model != nil && err == nil {
		log.Println("Creating symlink")
		CreateArchiveSymlink(model)
	}
	if err != nil {
		log.Println(err)
	}
}

func indexArchives(reindex bool) {
	InitAliases()
	InitBlacklists()
	InitMetadatas()

	paths, err := getArchivePaths()
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(paths))

	c := make(chan bool, 20)
	defer close(c)

	var archives []*modext.Archive
	var mutex sync.Mutex

	for _, path := range paths {
		c <- true

		go func(path string) {
			defer func() {
				wg.Done()
				<-c
			}()

			archive := &modext.Archive{Path: path}
			log.Println("Populating archive", filepath.Base(path))
			if err := populateArchive(archive); err != nil {
				log.Fatalln(err)
			}

			if len(archive.Title) > 0 {
				mutex.Lock()
				archives = append(archives, archive)
				mutex.Unlock()
			}
		}(path)
	}
	wg.Wait()

	sort.SliceStable(archives, func(i, j int) bool {
		return archives[i].CreatedAt < archives[j].CreatedAt
	})

	wg.Add(len(archives))
	for _, archive := range archives {
		c <- true
		go func(archive *modext.Archive) {
			defer func() {
				wg.Done()
				<-c
			}()

			log.Println("Indexing archive", filepath.Base(archive.Path))

			var model *modext.Archive
			var err error

			if reindex {
				model, err = UpdateArchive(archive)
			} else {
				model, err = CreateArchive(archive)
			}

			if model != nil && err == nil {
				CreateArchiveSymlink(model)
			}
		}(archive)
	}
	wg.Wait()
}

func updateSlugs() {
	archives, err := models.Archives().AllG()
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(archives))

	c := make(chan bool, 20)
	defer close(c)

	for _, archive := range archives {
		c <- true
		go func(archive *models.Archive) {
			defer func() {
				wg.Done()
				<-c
			}()
			archive.Slug = Slugify(archive.Title)
			if err := archive.UpdateG(boil.Whitelist("slug")); err != nil {
				log.Fatalln(err)
			}
		}(archive)
	}
	wg.Wait()

	artists, err := models.Artists().AllG()
	if err != nil {
		log.Fatalln(err)
	}

	wg.Add(len(artists))
	for _, artist := range artists {
		c <- true
		go func(artist *models.Artist) {
			defer func() {
				wg.Done()
				<-c
			}()
			artist.Slug = Slugify(artist.Name)
			if err := artist.UpdateG(boil.Whitelist("slug")); err != nil {
				log.Fatalln(err)
			}
		}(artist)
	}
	wg.Wait()

	circles, err := models.Circles().AllG()
	if err != nil {
		log.Fatalln(err)
	}

	wg.Add(len(circles))
	for _, circle := range circles {
		c <- true
		go func(circle *models.Circle) {
			defer func() {
				wg.Done()
				<-c
			}()
			circle.Slug = Slugify(circle.Name)
			if err := circle.UpdateG(boil.Whitelist("slug")); err != nil {
				log.Fatalln(err)
			}
		}(circle)
	}
	wg.Wait()

	magazines, err := models.Magazines().AllG()
	if err != nil {
		log.Fatalln(err)
	}

	wg.Add(len(magazines))
	for _, magazine := range magazines {
		c <- true
		go func(magazine *models.Magazine) {
			defer func() {
				wg.Done()
				<-c
			}()
			magazine.Slug = Slugify(magazine.Name)
			if err := magazine.UpdateG(boil.Whitelist("slug")); err != nil {
				log.Fatalln(err)
			}
		}(magazine)
	}
	wg.Wait()

	parody, err := models.Parodies().AllG()
	if err != nil {
		log.Fatalln(err)
	}

	wg.Add(len(parody))
	for _, p := range parody {
		c <- true
		go func(p *models.Parody) {
			defer func() {
				wg.Done()
				<-c
			}()
			p.Slug = Slugify(p.Name)
			if err := p.UpdateG(boil.Whitelist("slug")); err != nil {
				log.Fatalln(err)
			}
		}(p)
	}
	wg.Wait()

	tag, err := models.Tags().AllG()
	if err != nil {
		log.Fatalln(err)
	}

	wg.Add(len(tag))
	for _, t := range tag {
		c <- true
		go func(t *models.Tag) {
			defer func() {
				wg.Done()
				<-c
			}()
			t.Slug = Slugify(t.Name)
			if err := t.UpdateG(boil.Whitelist("slug")); err != nil {
				log.Fatalln(err)
			}
		}(t)
	}
	wg.Wait()
}

// remapArchives regenerates symlinks for all archives
func remapArchives() {
	archives, err := models.Archives().AllG()
	if err != nil {
		log.Fatalln(err)
	}

	os.RemoveAll(Config.Directories.Symlinks)
	for _, archive := range archives {
		CreateArchiveSymlink(modext.NewArchive(archive))
	}
}

func generateThumbnails() {
	archives, err := models.Archives(Where("published_at IS NOT NULL AND expunged IS FALSE"), OrderBy("id ASC")).AllG()
	if err != nil {
		log.Fatalln(err)
	}

	wg := &sync.WaitGroup{}
	c := make(chan bool, 5)
	defer close(c)

	for _, archive := range archives {
		log.Println("Generating thumbnails for", archive.ID, "-", archive.Title)

		zf, err := zip.OpenReader(archive.Path)
		if err != nil {
			log.Fatalln(err)
		}

		var files []*zip.File
		for _, f := range zf.File {
			stat := f.FileInfo()
			name := stat.Name()

			if stat.IsDir() || !IsImage(name) {
				continue
			}

			files = append(files, f)
		}

		sort.SliceStable(files, func(i, j int) bool {
			return GetPageNum(filepath.Base(files[i].Name)) < GetPageNum(filepath.Base(files[j].Name))
		})

		wg.Add(len(files))
		for i, f := range files {
			c <- true
			go func(n int, f *zip.File) {
				defer func() {
					wg.Done()
					<-c
				}()

				log.Println("Generating thumbnail of page", n)
				width := 288
				if n > 1 {
					width = 320
				}

				fp := filepath.Join(Config.Directories.Thumbnails,
					fmt.Sprintf("%d-%d.%d.webp", archive.ID, n, width))

				reader, err := f.Open()
				if err != nil {
					log.Fatalln(err)
				}
				defer reader.Close()

				tmp, err := os.CreateTemp("", "tmp-")
				if err != nil {
					log.Fatalln(err)
				}
				defer func() {
					tmp.Close()
					os.Remove(tmp.Name())
				}()

				if _, err := io.Copy(tmp, reader); err != nil {
					log.Fatalln(err)
				}

			Resize:
				if _, err := os.Stat(fp); os.IsNotExist(err) {
					opts := ResizeOptions{Width: width, Height: width * 3 / 2}
					opts.PNG = strings.HasSuffix(strings.ToLower(f.FileHeader.Name), ".png")
					if err := ResizeImage(tmp.Name(), fp, opts); err != nil {
						log.Fatalln(err)
					}
					time.Sleep(time.Second)
				}

				if width == 288 {
					width = 896
					fp = filepath.Join(Config.Directories.Thumbnails,
						fmt.Sprintf("%d-%d.%d.webp", archive.ID, n, width))
					goto Resize
				}
			}(i+1, f)
		}
		wg.Wait()
	}
}

func purgeThumbnails() {
	if err := os.RemoveAll(Config.Directories.Thumbnails); err != nil {
		log.Fatalln(err)
	} else if err := os.MkdirAll(Config.Directories.Thumbnails, 0755); err != nil {
		log.Fatalln(err)
	}
}

func purgeSymlinks() {
	if err := os.RemoveAll(Config.Directories.Symlinks); err != nil {
		log.Fatalln(err)
	}
}
