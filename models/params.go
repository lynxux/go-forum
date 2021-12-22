package models

//定义请求的参数结构体

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamSignUp 注册时的参数
type ParamSignUp struct {
	// binding:"required" 表示强制非空,属于validator
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录时的参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVoteData struct {
	//user_id: get from the current user
	PostID string `json:"post_id" binding:"required"`
	//这里的Direction不添加required字段，因为会把0识别为null，而不是取消投票
	Direction int8 `json:"direction,string" binding:"oneof=1 0 -1"` //1:agree  -1:disagree  0:cancel
}

type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"` // 可以为空
	Page        int64  `json:"page" form:"page" example:"1"`
	Size        int64  `json:"size" form:"size" example:"10"`
	Order       string `json:"order" form:"order" example:"score"`
}
