package services

const rawSqlArtistsMatch = `
(
	SELECT COUNT(*) FROM archive_artists
	LEFT JOIN artist ON artist.id = archive_artists.artist_id
		AND archive_artists.archive_id = archive.id
	WHERE artist.slug = ?
) > 0`

const rawSqlArtistsWildcard = `
(
	SELECT COUNT(*) FROM archive_artists
	LEFT JOIN artist ON artist.id = archive_artists.artist_id
		AND archive_artists.archive_id = archive.id
	WHERE artist.slug ILIKE '%' || ? || '%'
) > 0`

const rawSqlExcludeArtistsMatch = `(
	SELECT COUNT(*) FROM archive_artists
	LEFT JOIN artist ON artist.id = archive_artists.artist_id
		AND archive_artists.archive_id = archive.id
	WHERE artist.slug = ?
) = 0`

const rawSqlExcludeArtistsWildcard = `(
	SELECT COUNT(*) FROM archive_artists
	LEFT JOIN artist ON artist.id = archive_artists.artist_id
		AND archive_artists.archive_id = archive.id
	WHERE artist.slug ILIKE '%' || ? || '%'
) = 0`

const rawSqlCirclesMatch = `(
	SELECT COUNT(*) FROM archive_circles
	LEFT JOIN circle ON circle.id = archive_circles.circle_id
		AND archive_circles.archive_id = archive.id
	WHERE circle.slug = ?
) > 0`

const rawSqlCirclesWildcard = `(
	SELECT COUNT(*) FROM archive_circles
	LEFT JOIN circle ON circle.id = archive_circles.circle_id
		AND archive_circles.archive_id = archive.id
	WHERE circle.slug ILIKE '%' || ? || '%'
) > 0`

const rawSqlExcludeCirclesMatch = `(
	SELECT COUNT(*) FROM archive_circles
	LEFT JOIN circle ON circle.id = archive_circles.circle_id
		AND archive_circles.archive_id = archive.id
	WHERE circle.slug = ?
) = 0`

const rawSqlExcludeCirclesWildcard = `(
	SELECT COUNT(*) FROM archive_circles
	LEFT JOIN circle ON circle.id = archive_circles.circle_id	
		AND archive_circles.archive_id = archive.id
	WHERE circle.slug ILIKE '%' || ? || '%'
) = 0`

const rawSqlMagazinesMatch = `(
	SELECT COUNT(*) FROM archive_magazines
	LEFT JOIN magazine ON magazine.id = archive_magazines.magazine_id
		AND archive_magazines.archive_id = archive.id
	WHERE magazine.slug = ?
) > 0`

const rawSqlMagazinesWildcard = `(
	SELECT COUNT(*) FROM archive_magazines
	LEFT JOIN magazine ON magazine.id = archive_magazines.magazine_id
		AND archive_magazines.archive_id = archive.id
	WHERE magazine.slug ILIKE '%' || ? || '%'
) > 0`

const rawSqlExcludeMagazinesMatch = `(
	SELECT COUNT(*) FROM archive_magazines
	LEFT JOIN magazine ON magazine.id = archive_magazines.magazine_id
		AND archive_magazines.archive_id = archive.id
	WHERE magazine.slug = ?
) = 0`

const rawSqlExcludeMagazinesWildcard = `(
	SELECT COUNT(*) FROM archive_magazines
	LEFT JOIN magazine ON magazine.id = archive_magazines.magazine_id
		AND archive_magazines.archive_id = archive.id
	WHERE magazine.slug ILIKE '%' || ? || '%'
) = 0`

const rawSqlParodiesMatch = `(
	SELECT COUNT(*) FROM archive_parodies
	LEFT JOIN parody ON parody.id = archive_parodies.parody_id
		AND archive_parodies.archive_id = archive.id
	WHERE parody.slug = ?
) > 0`

const rawSqlParodiesWildcard = `(
	SELECT COUNT(*) FROM archive_parodies
	LEFT JOIN parody ON parody.id = archive_parodies.parody_id
		AND archive_parodies.archive_id = archive.id
	WHERE parody.slug ILIKE '%' || ? || '%'
) > 0`

const rawSqlExcludeParodiesMatch = `(	
	SELECT COUNT(*) FROM archive_parodies
	LEFT JOIN parody ON parody.id = archive_parodies.parody_id
		AND archive_parodies.archive_id = archive.id
	WHERE parody.slug ILIKE ?
) = 0`

const rawSqlExcludeParodiesWildcard = `(	
	SELECT COUNT(*) FROM archive_parodies
	LEFT JOIN parody ON parody.id = archive_parodies.parody_id
		AND archive_parodies.archive_id = archive.id
	WHERE parody.slug ILIKE '%' || ? || '%'
) = 0`

const rawSqlTagsMatch = `
(
	SELECT COUNT(*) FROM archive_tags
	LEFT JOIN tag ON tag.id = archive_tags.tag_id
		AND archive_tags.archive_id = archive.id
	WHERE tag.slug = ?
) > 0`

const rawSqlTagsWildcard = `(
	SELECT COUNT(*) FROM archive_tags
	LEFT JOIN tag ON tag.id = archive_tags.tag_id	
		AND archive_tags.archive_id = archive.id
	WHERE tag.slug ILIKE '%' || ? || '%'
) > 0`

const rawSqlExcludeTagsMatch = `(
	SELECT COUNT(*) FROM archive_tags
	LEFT JOIN tag ON tag.id = archive_tags.tag_id
		AND archive_tags.archive_id = archive.id
	WHERE tag.slug = ?
) = 0`

const rawSqlExcludeTagsWildcard = `(
	SELECT COUNT(*) FROM archive_tags
	LEFT JOIN tag ON tag.id = archive_tags.tag_id
		AND archive_tags.archive_id = archive.id
	WHERE tag.slug ILIKE '%' || ? || '%'
) = 0`
