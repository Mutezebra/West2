package types

type UserRegisterReq struct {
	UserName string `json:"user_name" form:"user_name" binding:"required,max=20,min=1"`
	Password string `json:"password" form:"password" binding:"required,max=20,min=6"`
}

type UserLoginReq struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type CreateGroupReq struct {
	GroupName string `json:"group_name" form:"group_name" binding:"required,max=14,min=1"`
}

type UserInfoResp struct {
	ID       uint       `json:"id,omitempty"`
	UserName string     `json:"user_name,omitempty"`
	Avatar   string     `json:"avatar,omitempty"`
	Token    *TokenResp `json:"token,omitempty"`
}
