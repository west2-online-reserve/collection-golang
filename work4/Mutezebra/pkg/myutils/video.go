package myutils

import (
	"errors"
	"four/consts"
)

func IsValidVideoSize(size int64) error {
	if size >= consts.MaxVideoSize {
		return errors.New("video size over the limit")
	}
	return nil
}
