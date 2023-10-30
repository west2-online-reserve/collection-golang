package myutils

import (
	"errors"
	"four/consts"
	"io"
	"mime/multipart"
	"os"
	"path"
	"regexp"
	"strings"
)

func ValidAvatar(avatar *multipart.FileHeader) (ext string, err error) {
	ext, ok := IsImg(avatar.Filename)
	if !ok {
		err = errors.New("avatar`s type is invalid")
		return
	}
	ok = IsValidAvatarSize(avatar.Size)
	if !ok {
		err = errors.New("avatar`s size it too big")
		return
	}
	return
}

func SavedAvatarFile(avatar *multipart.FileHeader, avatarPath string) (err error) {
	out, err := os.OpenFile(avatarPath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	src, err := avatar.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	_, err = io.Copy(out, src)
	return err
}

func IsEmail(email string) bool {
	result, _ := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, email)
	return result
}

func IsImg(fileName string) (string, bool) {
	ext := GetFileSuffix(fileName)
	if ext == ".jpg" || ext == ".png" {
		return ext, true
	}
	return "", false
}

func IsValidAvatarSize(size int64) bool {
	if size >= consts.MaxAvatarSize {
		return false
	}
	return true
}

func GetFileSuffix(fileName string) string {
	ext := strings.ToLower(path.Ext(fileName))
	return ext
}
