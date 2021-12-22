package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"
)

//controller层用于实现 处理路由，参数校验，请求转发

// SignUpHandler 注册业务
func SignUpHandler(c *gin.Context) {
	//1.获取参数，参数校验
	var p models.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误,返回日志
		zap.L().Error("SignUp with invalid param", zap.Error(err))

		//获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok { //非validator.ValidationErrors类型的错误直接返回错误消息(不翻译)
			ResponseError(c, CodeInvalidParam)
			return
		}
		//是Validator.ValidationErrors类型的错误(即校验器错误)，则进行翻译
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}

	//2.业务处理: 更新数据库
	if err := logic.SignUp(&p); err != nil {
		//fmt.Println(err)
		zap.L().Error("logic.SignUp() failed!", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist) //用户已存在，返回 用户已存在
		} else {
			ResponseError(c, CodeServerBusy) //数据库处理出错，返回ServerBusy
		}
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler 登录业务
func LoginHandler(c *gin.Context) {
	//1.参数校验
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误,返回日志
		zap.L().Error("Login with invalid param", zap.Error(err))

		errs, ok := err.(validator.ValidationErrors)
		if !ok { //非validator.ValidationErrors类型的错误直接返回错误消息(不翻译)
			ResponseError(c, CodeInvalidParam)
			return
		}
		//是Validator.ValidationErrors类型的错误(即校验器错误)，则进行翻译
		ResponseErrorWithMsg(c, CodeInvalidParam, errs.Translate(trans))
		return
	}

	//2. 业务处理：查询数据库，核验密码，返回token
	token, err := logic.Login(&p)
	if err != nil {
		zap.L().Error("login.Login() failed! ", zap.String("username: ", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
		} else {
			ResponseError(c, CodeInvalidPassword)
		}
		return
	}

	//3 返回响应，这里把token返回给了前端
	ResponseSuccess(c, token)

}
