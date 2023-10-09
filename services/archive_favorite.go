package services

import (
	"database/sql"
	"log"

	"koushoku/cache"
	"koushoku/errs"
	"koushoku/models"
	"koushoku/modext"

	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func AddFavorite(id int64, user *modext.User) error {
	archive, err := models.FindArchiveG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.ArchiveNotFound
		}
		log.Println(err)
		return errs.Unknown
	}

	if exists, _ := archive.Users(Where("id = ?", user.ID)).ExistsG(); exists {
		return nil
	}

	if err := archive.AddUsersG(false, &models.User{ID: user.ID}); err != nil {
		log.Println(err)
		return errs.Unknown
	}
	return nil
}

type GetFavoritesOptions = GetArchivesOptions
type GetFavoritesResult = GetArchivesResult

func GetFavorites(user *modext.User, opts GetFavoritesOptions) (result *GetFavoritesResult) {
	opts.Validate()

	cacheKey := makeCacheKey(opts)
	if c, err := cache.Favorites.GetWithPrefix(user.ID, cacheKey); err == nil {
		return c.(*GetFavoritesResult)
	}

	result = &GetFavoritesResult{Archives: []*modext.Archive{}}
	defer func() {
		if len(result.Archives) > 0 || result.Total > 0 || result.Err != nil {
			cache.Favorites.RemoveWithPrefix(user.ID, cacheKey)
			cache.Favorites.SetWithPrefix(user.ID, cacheKey, result, 0)
		}
	}()

	selectMods, countMods := opts.ToQueries()
	model := &models.User{ID: user.ID}
	archives, err := model.Archives(selectMods...).AllG()
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	count, err := model.Archives(countMods...).CountG()
	if err != nil {
		log.Println(err)
		result.Err = errs.Unknown
		return
	}

	result.Archives = make([]*modext.Archive, len(archives))
	result.Total = int(count)

	for i, archive := range archives {
		result.Archives[i] = modext.NewArchive(archive).LoadRels(archive)
	}
	return
}

func DeleteFavorite(id int64, user *modext.User) error {
	archive, err := models.FindArchiveG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.ArchiveNotFound
		}
		log.Println(err)
		return errs.Unknown
	}

	if err := archive.RemoveUsersG(&models.User{ID: user.ID}); err != nil {
		log.Println(err)
		return errs.Unknown
	}
	// TODO: purge cache
	return nil
}
