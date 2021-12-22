package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"web_app/models"
)

const secret = "link98"

// CheckUserExist
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err = dbSqlx.Get(&count, sqlStr, username); err != nil {
		return
	}
	if count > 0 {
		return ErrorUserExist
	}
	return

}

func encryptPassword(oldPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oldPassword)))
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)

	sqlStr := `insert into user(user_id, username, password) values (?,?,?)`
	_, err = dbSqlx.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func CheckUserPWD(user *models.User) (err error) {

	oldPassword := encryptPassword(user.Password)
	sqlStr := `select user_id, username, password from user where username=?`
	err = dbSqlx.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err //查询出错时err不为nil， 如果user不存在时则只会返回空值？不会返回nil
	}

	if oldPassword != user.Password {
		return ErrorInvalidPassword
	}
	return
}

func GetUserById(userId int64) (user *models.User, err error) {
	user = new(models.User)
	sqlstr := `select user_id, username from user where user_id = ?`
	err = dbSqlx.Get(user, sqlstr, userId)
	return
}
