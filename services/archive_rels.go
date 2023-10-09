package services

import (
	"sort"
	"strings"
	"sync"

	"koushoku/models"
	"koushoku/modext"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

var relsCache struct {
	Artists   map[string]*models.Artist
	Circles   map[string]*models.Circle
	Magazines map[string]*models.Magazine
	Parodies  map[string]*models.Parody
	Tags      map[string]*models.Tag

	sync.RWMutex
	sync.Once
}

func PopulateArchiveRels(e boil.Executor, model *models.Archive, archive *modext.Archive) error {
	relsCache.Do(func() {
		relsCache.Lock()
		defer relsCache.Unlock()

		relsCache.Artists = make(map[string]*models.Artist)
		relsCache.Circles = make(map[string]*models.Circle)
		relsCache.Magazines = make(map[string]*models.Magazine)
		relsCache.Parodies = make(map[string]*models.Parody)
		relsCache.Tags = make(map[string]*models.Tag)
	})

	var err error
	if len(archive.Artists) > 0 {
		var artists []*models.Artist
		for _, artist := range archive.Artists {
			relsCache.RLock()
			artistModel, ok := relsCache.Artists[artist.Name]
			relsCache.RUnlock()

			if ok {
				artists = append(artists, artistModel)
				continue
			}

			relsCache.Lock()
			artist, err = CreateArtist(artist.Name)
			if err != nil {
				relsCache.Unlock()
				return err
			}

			artistModel = &models.Artist{ID: artist.ID, Slug: artist.Slug, Name: artist.Name}
			relsCache.Artists[artist.Name] = artistModel
			relsCache.Unlock()

			artists = append(artists, artistModel)
		}
		if err := model.SetArtists(e, false, artists...); err != nil {
			return err
		}
	}

	if len(archive.Circles) > 0 {
		var circles []*models.Circle
		for _, circle := range archive.Circles {
			relsCache.RLock()
			circleModel, ok := relsCache.Circles[circle.Name]
			relsCache.RUnlock()

			if ok {
				circles = append(circles, circleModel)
				continue
			}

			relsCache.Lock()
			circle, err := CreateCircle(circle.Name)
			if err != nil {
				relsCache.Unlock()
				return err
			}

			circleModel = &models.Circle{ID: circle.ID, Slug: circle.Slug, Name: circle.Name}
			relsCache.Circles[circle.Name] = circleModel
			relsCache.Unlock()

			circles = append(circles, circleModel)
		}
		if err := model.SetCircles(e, false, circles...); err != nil {
			return err
		}
	}

	if len(archive.Magazines) > 0 {
		var magazines []*models.Magazine
		for _, magazine := range archive.Magazines {
			relsCache.RLock()
			magazineModel, ok := relsCache.Magazines[magazine.Name]
			relsCache.RUnlock()

			if ok {
				magazines = append(magazines, magazineModel)
				continue
			}

			relsCache.Lock()
			magazine, err := CreateMagazine(magazine.Name)
			if err != nil {
				relsCache.Unlock()
				return err
			}

			magazineModel = &models.Magazine{ID: magazine.ID, Slug: magazine.Slug, Name: magazine.Name}
			relsCache.Magazines[magazine.Name] = magazineModel
			relsCache.Unlock()

			magazines = append(magazines, magazineModel)
		}
		if err := model.SetMagazines(e, false, magazines...); err != nil {
			return err
		}
	}

	if len(archive.Parodies) > 0 {
		var parodies []*models.Parody
		for _, parody := range archive.Parodies {
			relsCache.RLock()
			parodyModel, ok := relsCache.Parodies[parody.Name]
			relsCache.RUnlock()

			if ok {
				parodies = append(parodies, parodyModel)
				continue
			}

			relsCache.Lock()
			parody, err := CreateParody(parody.Name)
			if err != nil {
				relsCache.Unlock()
				return err
			}

			parodyModel = &models.Parody{ID: parody.ID, Slug: parody.Slug, Name: parody.Name}
			relsCache.Parodies[parody.Name] = parodyModel
			relsCache.Unlock()

			parodies = append(parodies, parodyModel)
		}
		if err := model.SetParodies(e, false, parodies...); err != nil {
			return err
		}
	}

	if len(archive.Tags) > 0 {
		var tags []*models.Tag
		for _, tag := range archive.Tags {
			relsCache.RLock()
			tagModel, ok := relsCache.Tags[tag.Name]
			relsCache.RUnlock()

			if ok {
				tags = append(tags, tagModel)
				continue
			}

			relsCache.Lock()
			tag, err := CreateTag(tag.Name)
			if err != nil {
				relsCache.Unlock()
				return err
			}

			tagModel = &models.Tag{ID: tag.ID, Slug: tag.Slug, Name: tag.Name}
			relsCache.Tags[tag.Name] = tagModel
			relsCache.Unlock()

			tags = append(tags, tagModel)
		}
		if err := model.SetTags(e, false, tags...); err != nil {
			return err
		}
	}
	return nil
}

func validateArchiveRels(rels []string) (result []string) {
	for _, v := range rels {
		if strings.EqualFold(v, ArchiveRels.Artists) {
			result = append(result, ArchiveRels.Artists)
		} else if strings.EqualFold(v, ArchiveRels.Circles) {
			result = append(result, ArchiveRels.Circles)
		} else if strings.EqualFold(v, ArchiveRels.Magazines) {
			result = append(result, ArchiveRels.Magazines)
		} else if strings.EqualFold(v, ArchiveRels.Parodies) {
			result = append(result, ArchiveRels.Parodies)
		} else if strings.EqualFold(v, ArchiveRels.Tags) {
			result = append(result, ArchiveRels.Tags)
		} else if strings.EqualFold(v, ArchiveRels.Submission) {
			result = append(result, ArchiveRels.Submission)
		}
	}
	sort.Strings(result)
	return
}
