package errs

import "errors"

var Unknown = errors.New("Unknown error")

var (
	ArchiveNotFound    = errors.New("Archive does not exist")
	ArtistNotFound     = errors.New("Artist does not exist")
	CircleNotFound     = errors.New("Circle does not exist")
	MagazineNotFound   = errors.New("Magazine does not exist")
	TagNotFound        = errors.New("Tag does not exist")
	ParodyNotFound     = errors.New("Parody does not exist")
	UserNotFound       = errors.New("User does not exist")
	SubmissionNotFound = errors.New("Submission does not exist")
)

var (
	ArchivePathRequired        = errors.New("Archive path is required")
	ArtistNameRequired         = errors.New("Artist name is required")
	ArtistNameTooLong          = errors.New("Artist name must be at most 128 characters")
	CircleNameRequired         = errors.New("Circle name is required")
	CircleNameTooLong          = errors.New("CIrcle name must be at most 128 characters")
	MagazineNameRequired       = errors.New("Magazine name is required")
	MagazineNameTooLong        = errors.New("Magazine name must be at most 128 characters")
	ParodyNameRequired         = errors.New("Parody name is required")
	ParodyNameTooLong          = errors.New("Parody name must be at most 128 characters")
	TagNameRequired            = errors.New("Tag name is required")
	TagNameTooLong             = errors.New("Tag name must be at most 128 characters")
	SubmissionNameRequired     = errors.New("Submission name is required")
	SubmissionNameTooLong      = errors.New("Submission name must be at most 1024 characters")
	SubmissionSubmitterTooLong = errors.New("Submission submitter must be at most 128 characters")
	SubmissionContentRequired  = errors.New("Submission content is required")
	SubmissionContentTooLong   = errors.New("Submission content must be at most 10240 characters")
)

var (
	UserNameTooShort   = errors.New("User name must be at least 3 characters")
	UserNameTooLong    = errors.New("User name must be at most 32 characters")
	EmailRequired      = errors.New("Email is required")
	EmailTooLong       = errors.New("Email must be at most 255 characters")
	EmailInvalid       = errors.New("Email is invalid")
	PasswordTooShort   = errors.New("Password must be at least 6 characters")
	InvalidCredentials = errors.New("Invalid credentials")
)
