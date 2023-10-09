package services

import (
	"database/sql"
	"log"
	"regexp"
	"strings"

	"koushoku/cache"
	"koushoku/errs"
	"koushoku/models"
	"koushoku/modext"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/crypto/bcrypt"
)

var emailRgx = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func isEmail(e string) bool {
	return emailRgx.MatchString(e)
}

func hashPassword(rawPassword string) (string, error) {
	buf, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	return string(buf), err
}

type CreateUserOptions struct {
	Name     string
	Email    string
	Password string
}

func CreateUser(opts CreateUserOptions) (*modext.User, error) {
	opts.Name = strings.TrimSpace(opts.Name)
	opts.Email = strings.TrimSpace(opts.Email)

	switch {
	case len(opts.Name) < 3:
		return nil, errs.UserNameTooShort
	case len(opts.Name) > 32:
		return nil, errs.UserNameTooLong
	case len(opts.Email) == 0:
		return nil, errs.EmailRequired
	case len(opts.Email) > 255:
		return nil, errs.EmailTooLong
	case !isEmail(opts.Email):
		return nil, errs.EmailInvalid
	case len(opts.Password) < 6:
		return nil, errs.PasswordTooShort
	}

	hashedPassword, err := hashPassword(opts.Password)
	if err != nil {
		log.Println(err)
		return nil, errs.Unknown
	}

	user := &models.User{
		Name:     opts.Name,
		Email:    opts.Email,
		Password: hashedPassword,
	}

	if err := user.InsertG(boil.Infer()); err != nil {
		log.Println(err)
		return nil, errs.Unknown
	}
	return modext.NewUser(user), nil
}

type GetUserResult struct {
	User *modext.User
	Err  error
}

func GetUser(id int64) (result *GetUserResult) {
	if c, err := cache.Users.GetWithInt64(id); err == nil {
		return c.(*GetUserResult)
	}

	result = &GetUserResult{}
	defer func() {
		if result.User != nil || result.Err != nil {
			cache.Users.RemoveWithInt64(id)
			cache.Users.SetWithInt64(id, result, 0)
		}
	}()

	user, err := models.FindUserG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			result.Err = errs.UserNotFound
		} else {
			log.Println(err)
			result.Err = errs.Unknown
		}
		return
	}

	result.User = modext.NewUser(user)
	return
}

func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func UpdatePassword(id int64, password, newPassword string) error {
	user, err := models.FindUserG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.UserNotFound
		}
		log.Println(err)
		return errs.Unknown
	}

	if err := checkPassword(user.Password, password); err != nil {
		log.Println(err)
		return errs.InvalidCredentials
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		log.Println(err)
		return errs.Unknown
	}

	user.Password = hashedPassword
	if err := user.UpdateG(boil.Infer()); err != nil {
		log.Println(err)
		return errs.Unknown
	}
	return nil
}

func DeleteUser(id int64, password string) error {
	user, err := models.FindUserG(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.UserNotFound
		}
		log.Println(err)
		return errs.Unknown
	}

	if err := checkPassword(user.Password, password); err != nil {
		log.Println(err)
		return errs.InvalidCredentials
	}

	if err := user.DeleteG(); err != nil {
		log.Println(err)
		return errs.Unknown
	}
	return nil
}
