package dao

import (
	"context"
	"errors"
	"fmt"
	"four/consts"
	"four/repository/db/model"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

func (dao *UserDao) FindUserByUserName(userName string) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).First(&user, "user_name=?", userName).Error
	if err != nil {
		return nil, err
	}
	return
}

func (dao *UserDao) FindUserByUserEmail(email string) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).First(&user, "email=?", email).Error
	if err != nil {
		return nil, err
	}
	return
}

func (dao *UserDao) FindUserByUserID(id uint) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).First(&user, "id=?", id).Error
	if err != nil {
		return nil, err
	}
	return
}

func (dao *UserDao) Create(user *model.User) (err error) {
	err = dao.DB.Model(&model.User{}).Create(&user).Error
	return
}

func (dao *UserDao) Update(user *model.User) (err error) {
	err = dao.DB.Model(&user).Save(user).Error
	return
}

func (dao *UserDao) UpdateAvatar(user *model.User, url string) (err error) {
	err = dao.DB.Model(&user).Update("avatar", url).Error
	return
}
func (dao *UserDao) UpdateTotpStatus(user *model.User, ok bool) (err error) {
	err = dao.DB.Model(&user).Update("totp_enable", ok).Error
	return err
}

func (dao *UserDao) UpdateOtpSecret(user *model.User, secret string, url string) (err error) {
	err = dao.DB.Model(&user).Update("totp_secret", secret).Error
	if err != nil {
		return err
	}
	err = dao.DB.Model(&user).Update("totp_url", url).Error
	return err
}

// FindUserFollow 查询该用户是否已经关注过，如果有就返回一个error，没有就返回有nil
func (dao *UserDao) FindUserFollow(uid uint, followerID uint, tableName string) (err error) {
	fan := model.Fans{}
	dao.DB.Raw("SELECT * FROM "+tableName+" WHERE (uid=? AND follower_id=?) AND deleted_at IS NULL LIMIT 1", uid, followerID).Scan(&fan)
	if fan.ID == 0 {
		return nil
	}
	return errors.New("record have exist")
}

func (dao *UserDao) Follow(fan *model.Fans) (err error) {
	_, err = dao.FindUserByUserID(fan.FollowerId)
	if err != nil {
		return err
	}
	index := fan.Uid / consts.EachUserRecordAFansTable
	err = dao.DB.Exec(consts.CreateNewFansTable(index)).Error
	if err != nil {
		return err
	}

	tableName := fmt.Sprintf("fans%d", index)
	err = dao.FindUserFollow(fan.Uid, fan.FollowerId, tableName)
	if err != nil {
		return err
	}
	err = dao.DB.Exec("INSERT INTO "+tableName+"(created_at,updated_at,deleted_at,uid,follower_id) VALUES (NOW(),NOW(),null,?,?)", fan.Uid, fan.FollowerId).Error
	return err
}

func (dao *UserDao) UnFollow(fan *model.Fans) (err error) {
	tableName := fmt.Sprintf("fans%d", fan.Uid/consts.EachUserRecordAFansTable)
	err = dao.FindUserFollow(fan.Uid, fan.FollowerId, tableName)
	if err == nil {
		return errors.New("record does not exist")
	}
	err = dao.DB.Exec("UPDATE "+tableName+"  SET deleted_at=? WHERE uid=? AND follower_id=? AND deleted_at IS NULL LIMIT 1", time.Now().Format("2006-01-02 15:04:05"), fan.Uid, fan.FollowerId).Error
	return err
}

func (dao *UserDao) FriendList(uid uint) (result []uint, err error) {
	index := uid / consts.EachUserRecordAFansTable
	tableName := fmt.Sprintf("fans%d", index)
	err = dao.Raw("SELECT f1.follower_id AS result FROM "+tableName+" f1 INNER JOIN "+tableName+" f2 ON f1.uid = f2.follower_id AND f1.follower_id = f2.uid WHERE f1.uid = ?", uid).Scan(&result).Error
	return result, err
}

func (dao *UserDao) FollowerList(uid uint) (result []uint, err error) {
	index := uid / consts.EachUserRecordAFansTable
	tableName := fmt.Sprintf("fans%d", index)
	err = dao.Raw("SELECT follower_id FROM "+tableName+" WHERE uid=? AND deleted_at IS NULL", uid).Scan(&result).Error
	return result, err
}

func (dao *UserDao) Delete(uid uint) (err error) {
	err = dao.DB.Delete(&model.User{}, "id=?", uid).Error
	return err
}
