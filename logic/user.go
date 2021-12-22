package logic

import (
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 生成uid
	userID := snowflake.GenID()
	user := models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	//密码加密，保存入数据库
	return mysql.InsertUser(&user)
}

// Login
func Login(p *models.ParamLogin) (token string, err error) {
	user := models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 这里传递的是user的指针，所以user的值会在函数调用后被改变(指结构体中的密码会改变)
	// 这里说传递的是指针，只是CheckUSerPWD函数声明时接收的参数为指针，而调用时，这是传参应该传 user 还是 &user ？？？
	if err = mysql.CheckUserPWD(&user); err != nil {
		return "", err //包含密码错误和用户不存在两种错误
	}

	// 当登录成功时，生成JWT
	// 生成JWT成功则返回对应的token和nil
	// 生成JWT失败则返回空字符串和对应的错误
	return jwt.GenToken(user.UserID, user.Username)

}
