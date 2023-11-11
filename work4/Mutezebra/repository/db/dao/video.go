package dao

import (
	"context"
	"fmt"
	"four/consts"
	"four/repository/db/model"
	"gorm.io/gorm"
)

type VideoDao struct {
	*gorm.DB
}

func NewVideoDao(ctx context.Context) *VideoDao {
	return &VideoDao{NewDBClient(ctx)}
}

func (dao *VideoDao) FindVideoByVideoID(vid uint) (video *model.Video, err error) {
	err = dao.DB.Model(&model.Video{}).Where("id=?", vid).First(&video).Error
	return video, err
}

func (dao *VideoDao) UpdateVideoViews(vid uint, view int) error {
	err := dao.DB.Model(&model.Video{}).Where("id=?", vid).UpdateColumn("view", view).Error
	return err
}

func (dao *VideoDao) Create(ch chan *model.Video, chErr chan error) {
	var err error
	v := &model.Video{}
	err = dao.DB.Model(&model.Video{}).Create(&v).Error
	if err != nil {
		chErr <- err
		close(chErr)
		close(ch)
		return
	}
	ch <- v
	v = <-ch
	err = dao.Update(v)
	chErr <- err

	close(ch)
	close(chErr)
	return
}

func (dao *VideoDao) Update(v *model.Video) (err error) {
	err = dao.DB.Model(&v).Save(v).Error
	return err
}

func (dao *VideoDao) FindCommentRoot(c *model.Comment) (root uint, err error) {
	tableName := fmt.Sprintf("comment%d", c.VideoID/consts.EachVideoRecordACommentTable)
	err = dao.Raw(c.FindCommentRootSQL(tableName)).Scan(&root).Error
	if err != nil {
		fmt.Println("find comment root failed")
	}
	return
}

func (dao *VideoDao) CreateComment(c *model.Comment) (err error) {
	index := c.VideoID / consts.EachVideoRecordACommentTable
	err = dao.DB.Exec(consts.CreateNewCommentTable(index)).Error
	if err != nil {
		return err
	}

	tableName := fmt.Sprintf("comment%d", index)

	err = dao.DB.Exec(c.InsertNewCommentSQL(tableName)).Error
	return err
}

func (dao *VideoDao) FindVideos(vids []uint) (videos []model.Video, err error) {
	err = dao.DB.Model(&model.Video{}).Find(&videos, "id in (?)", vids).Error
	return
}

func (dao *VideoDao) Delete(vid uint) (err error) {
	err = dao.DB.Delete(&model.Video{}, "id=?", vid).Error
	return err
}
