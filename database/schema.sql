CREATE TABLE IF NOT EXISTS artist (
  id   BIGSERIAL PRIMARY KEY,
  slug VARCHAR(128) NOT NULL DEFAULT NULL,
  name VARCHAR(128) NOT NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS artist_slug_uindex ON artist(slug);
CREATE UNIQUE INDEX IF NOT EXISTS artist_name_uindex ON artist(name);

CREATE TABLE IF NOT EXISTS circle (
  id BIGSERIAL PRIMARY KEY,
  slug VARCHAR(128) NOT NULL DEFAULT NULL,
  name VARCHAR(128) NOT NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS circle_slug_uindex ON circle(slug);
CREATE UNIQUE INDEX IF NOT EXISTS circle_name_uindex ON circle(name);

CREATE TABLE IF NOT EXISTS tag (
  id   BIGSERIAL PRIMARY KEY,
  slug VARCHAR(32) NOT NULL DEFAULT NULL,
  name VARCHAR(32) NOT NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS tag_slug_uindex ON tag(slug);
CREATE UNIQUE INDEX IF NOT EXISTS tag_name_uindex ON tag(name);

CREATE TABLE IF NOT EXISTS magazine (
  id BIGSERIAL PRIMARY KEY,
  slug VARCHAR(128) NOT NULL DEFAULT NULL,
  name VARCHAR(128) NOT NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS magazine_slug_uindex ON magazine(slug);
CREATE UNIQUE INDEX IF NOT EXISTS magazine_name_uindex ON magazine(name);

CREATE TABLE IF NOT EXISTS parody (
  id BIGSERIAL PRIMARY KEY,
  slug VARCHAR(128) NOT NULL DEFAULT NULL,
  name VARCHAR(128) NOT NULL DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS parody_slug_uindex ON parody(slug);
CREATE UNIQUE INDEX IF NOT EXISTS parody_name_uindex ON parody(name);

CREATE TABLE IF NOT EXISTS archive (
  id                BIGSERIAL PRIMARY KEY,
  path              TEXT NOT NULL DEFAULT NULL,
  
  created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMP NOT NULL DEFAULT NOW(),
  published_at      TIMESTAMP,

  title             VARCHAR(1024) NOT NULL DEFAULT NULL,
  slug              VARCHAR(1024) NOT NULL DEFAULT NULL,
  pages             SMALLINT NOT NULL DEFAULT NULL,
  size              BIGINT NOT NULL DEFAULT NULL
);

ALTER TABLE archive
  ADD COLUMN IF NOT EXISTS expunged BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS source VARCHAR(1024) DEFAULT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS archive_path_uindex ON archive(path);
CREATE INDEX IF NOT EXISTS archive_title_index ON archive(title);
CREATE INDEX IF NOT EXISTS archive_slug_index ON archive(slug);
CREATE INDEX IF NOT EXISTS archive_pages_index ON archive(pages);
CREATE INDEX IF NOT EXISTS archive_created_at_index ON archive(created_at);
CREATE INDEX IF NOT EXISTS archive_updated_at_index ON archive(updated_at);
CREATE INDEX IF NOT EXISTS archive_published_at_index ON archive(published_at);
CREATE INDEX IF NOT EXISTS archive_expunged_index ON archive(expunged);

CREATE TABLE IF NOT EXISTS archive_artists (
  archive_id BIGINT NOT NULL DEFAULT NULL REFERENCES archive(id) ON DELETE CASCADE,
  artist_id BIGINT NOT NULL DEFAULT NULL REFERENCES artist(id) ON DELETE CASCADE,
  PRIMARY KEY(archive_id, artist_id)
);

CREATE INDEX IF NOT EXISTS archive_artists_archive_id_index ON archive_artists(archive_id);
CREATE INDEX IF NOT EXISTS archive_artists_artist_id_index ON archive_artists(artist_id);

CREATE TABLE IF NOT EXISTS archive_circles (
  archive_id BIGINT NOT NULL DEFAULT NULL REFERENCES archive(id) ON DELETE CASCADE,
  circle_id BIGINT NOT NULL DEFAULT NULL REFERENCES circle(id) ON DELETE CASCADE,
  PRIMARY KEY(archive_id, circle_id)
);

CREATE INDEX IF NOT EXISTS archive_circles_archive_id_index ON archive_circles(archive_id);
CREATE INDEX IF NOT EXISTS archive_circles_circle_id_index ON archive_circles(circle_id);

CREATE TABLE IF NOT EXISTS archive_magazines (
  archive_id BIGINT NOT NULL DEFAULT NULL REFERENCES archive(id) ON DELETE CASCADE,
  magazine_id BIGINT NOT NULL DEFAULT NULL REFERENCES magazine(id) ON DELETE CASCADE,
  PRIMARY KEY(archive_id, magazine_id)
);

CREATE INDEX IF NOT EXISTS archive_magazines_archive_id_index ON archive_magazines(archive_id);
CREATE INDEX IF NOT EXISTS archive_magazines_magazine_id_index ON archive_magazines(magazine_id);

CREATE TABLE IF NOT EXISTS archive_parodies (
  archive_id BIGINT NOT NULL DEFAULT NULL REFERENCES archive(id) ON DELETE CASCADE,
  parody_id BIGINT NOT NULL DEFAULT NULL REFERENCES parody(id) ON DELETE CASCADE,
  PRIMARY KEY(archive_id, parody_id)
);

CREATE INDEX IF NOT EXISTS archive_parodies_archive_id_index ON archive_parodies(archive_id);
CREATE INDEX IF NOT EXISTS archive_parodies_parody_id_index ON archive_parodies(parody_id);

CREATE TABLE IF NOT EXISTS archive_tags (
  archive_id BIGINT NOT NULL DEFAULT NULL REFERENCES archive(id) ON DELETE CASCADE,
  tag_id    BIGINT NOT NULL DEFAULT NULL REFERENCES tag(id) ON DELETE CASCADE,
  PRIMARY KEY(archive_id, tag_id)
);

CREATE INDEX IF NOT EXISTS archive_tags_archive_id_index ON archive_tags(archive_id);
CREATE INDEX IF NOT EXISTS archive_tags_tag_id_index ON archive_tags(tag_id);

CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,

  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  
  email VARCHAR(256) NOT NULL DEFAULT NULL,
  password varchar(2048) NOT NULL DEFAULT NULL,
  name VARCHAR(32) NOT NULL DEFAULT NULL,
  
  is_banned BOOLEAN NOT NULL DEFAULT FALSE,
  is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS is_banned BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS is_admin BOOLEAN NOT NULL DEFAULT FALSE;

CREATE UNIQUE INDEX IF NOT EXISTS users_email_uindex ON users(email);
CREATE INDEX IF NOT EXISTS users_created_at_index ON users(created_at);
CREATE INDEX IF NOT EXISTS users_updated_at_index ON users(updated_at);
CREATE INDEX IF NOT EXISTS users_is_banned_index ON users(is_banned);
CREATE INDEX IF NOT EXISTS users_is_admin_index ON users(is_admin);

CREATE TABLE IF NOT EXISTS user_favorites (
  user_id BIGINT NOT NULL DEFAULT NULL REFERENCES users(id) ON DELETE CASCADE,
  archive_id BIGINT NOT NULL DEFAULT NULL REFERENCES archive(id) ON DELETE CASCADE,
  PRIMARY KEY(user_id, archive_id)
);

CREATE INDEX IF NOT EXISTS user_favorites_user_id_index ON user_favorites(user_id);
CREATE INDEX IF NOT EXISTS user_favorites_archive_id_index ON user_favorites(archive_id);

CREATE TABLE IF NOT EXISTS submission (
  id  BIGSERIAL PRIMARY KEY,
  
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  
  name VARCHAR(1024) NOT NULL DEFAULT NULL,
  submitter VARCHAR(128) DEFAULT NULL,
  content VARCHAR(4096) NOT NULL DEFAULT NULL,
  notes VARCHAR(4096) DEFAULT NULL,

  accepted_at TIMESTAMP,
  rejected_at TIMESTAMP,

  accepted BOOLEAN NOT NULL DEFAULT FALSE,
  rejected BOOLEAN NOT NULL DEFAULT FALSE
);

ALTER TABLE submission
  ALTER COLUMN content TYPE VARCHAR(10240);

CREATE INDEX IF NOT EXISTS submission_created_at_index ON submission(created_at);
CREATE INDEX IF NOT EXISTS submission_updated_at_index ON submission(updated_at);
CREATE INDEX IF NOT EXISTS submission_submitter_index ON submission(submitter);
CREATE INDEX IF NOT EXISTS submission_accepted_at_index ON submission(accepted_at);
CREATE INDEX IF NOT EXISTS submission_accepted_index ON submission(accepted);
CREATE INDEX IF NOT EXISTS submission_rejected_at_index ON submission(rejected_at);
CREATE INDEX IF NOT EXISTS submission_rejected_index ON submission(rejected);

DO $$
  DECLARE r record;
BEGIN
  FOR r IN
    SELECT conname FROM pg_constraint
      JOIN pg_class ON pg_constraint.conrelid = pg_class.oid
      WHERE pg_class.relname = 'archive'
      AND conname LIKE ANY (ARRAY['archive_redirect_id_fkey%', 'archive_submission_id_fkey%'])
      AND conname != ALL (ARRAY['archive_redirect_id_fkey', 'archive_submission_id_fkey'])
  LOOP
    EXECUTE 'ALTER TABLE archive DROP CONSTRAINT ' || r.conname;
  END LOOP;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'archive_submission_id_fkey') THEN
    ALTER TABLE archive
      ADD COLUMN submission_id BIGINT DEFAULT NULL REFERENCES submission(id) ON DELETE CASCADE;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'archive_redirect_id_fkey') THEN
    ALTER TABLE archive
      ADD COLUMN redirect_id BIGINT DEFAULT NULL REFERENCES archive(id) ON DELETE CASCADE;
  END IF;
END;
$$;