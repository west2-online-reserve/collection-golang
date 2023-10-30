package myutils

import (
	"bytes"
	"context"
	"fmt"
	"four/config"
	"four/consts"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
)

type PutRet struct {
	Hash string
	key  string
}

var conf *config.QiNiu

func OssInit() {
	con := config.Config.QiNiu
	conf = con
}

func getOSS() (*qbox.Mac, string, string) {
	acK := conf.AccessKey
	seK := conf.SecretKey
	domain := config.Config.QiNiu.Domain
	return qbox.NewMac(acK, seK), conf.Bucket, domain
}

func UploadVideo(ctx context.Context, name string, data []byte) (*PutRet, error) {
	mac, bucket, _ := getOSS()
	key := conf.VideoPath + name
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{Region: &storage.ZoneHuadongZheJiang2}
	resumeUploader := storage.NewResumeUploaderV2(&cfg)

	putExtra := storage.RputV2Extra{PartSize: 2 * consts.MB}
	ret := PutRet{}
	reader := bytes.NewReader(data)
	err := resumeUploader.Put(ctx, &ret, upToken, key, reader, int64(len(data)), &putExtra)
	return &ret, err
}

func DownloadVideo(name string, data chan []byte) error {
	mac, bucket, domain := getOSS()
	key := conf.VideoPath + name
	bm := storage.NewBucketManager(mac, &storage.Config{
		Region: &storage.ZoneHuadongZheJiang2,
	})

	resp, err := bm.Get(bucket, key, &storage.GetObjectInput{
		DownloadDomains: []string{domain},
		PresignUrl:      true,
		Range:           "bytes=0-",
	})
	if err != nil {
		close(data)
		return err
	}
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		close(data)
		return err
	}
	data <- d
	data <- nil
	return nil
}

func UploadAvatar(ctx context.Context, name string, data []byte) (*PutRet, error, string) {
	mac, bucket, _ := getOSS()
	key := conf.AvatarPath + name
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", bucket, key),
	}
	upToken := putPolicy.UploadToken(mac)

	formLoader := storage.NewFormUploader(&storage.Config{Region: &storage.ZoneHuadongZheJiang2})

	ret := &PutRet{}
	putExtra := &storage.PutExtra{}

	length := int64(len(data))

	err := formLoader.Put(ctx, ret, upToken, key, bytes.NewReader(data), length, putExtra)
	return ret, err, key
}
