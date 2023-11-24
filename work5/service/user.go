package service

import (
	"context"
	"errors"
	"five/config"
	"five/pkg/ctl"
	"five/pkg/e"
	"five/pkg/log"
	"five/pkg/myutils"
	"five/repository/db/dao"
	"five/repository/db/model"
	"five/repository/redis"
	"five/types"
	"gorm.io/gorm"
	"sync"
)

type UserService struct {
}

var UserSrv *UserService
var userOnce sync.Once

// GetUserSrv 获取UserService单例
func GetUserSrv() *UserService {
	userOnce.Do(func() {
		UserSrv = &UserService{}
	})
	return UserSrv
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *types.UserRegisterReq) (resp interface{}, err error) {
	code := e.SUCCESS
	userDao := dao.GetUserDao(ctx)

	// 1. 检查用户名是否存在
	user, err := userDao.FindUserByUserName(req.UserName)
	if err != nil && err != gorm.ErrRecordNotFound {
		code = e.CreateUserFailed
		log.LogrusObj.Errorln(err)
		return ctl.RespError(code, err), err
	} else if err == nil {
		code = e.UserExist
		log.LogrusObj.Errorln(err)
		return ctl.RespError(code, err), err
	}

	// 2. 创建用户
	user = &model.User{}
	err = user.SetPassword(req.Password)
	if err != nil {
		code = e.SetPasswordFailed
		log.LogrusObj.Errorln(err)
		return ctl.RespError(code, err), err
	}

	// 3. 设置默认头像
	user.Avatar = config.Config.Local.DefaultAvatarPath
	user.UserName = req.UserName
	user.ID = uint(redis.RedisClient.Incr("user_id").Val())
	// 4. 创建用户
	err = userDao.CreateUser(user)
	if err != nil {
		code = e.CreateUserFailed
		log.LogrusObj.Errorln(err)
		return ctl.RespError(code, err), err
	}

	// 5. 返回用户信息
	data := &types.UserInfoResp{
		ID:       user.ID,
		UserName: user.UserName,
		Avatar:   user.Avatar,
	}
	return ctl.RespSuccessWithData(code, data), nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *types.UserLoginReq) (resp interface{}, err error) {
	code := e.SUCCESS
	userDao := dao.GetUserDao(ctx)

	// 1. 检查用户名是否存在
	user, err := userDao.FindUserByUserName(req.UserName)
	if err != nil {
		code = e.UserNotExist
		log.LogrusObj.Errorln(err)
		return ctl.RespError(code, err), err
	}

	// 2. 检查密码是否正确
	if !user.CheckPassword(req.Password) {
		code = e.PasswordError
		log.LogrusObj.Errorln(err)
		return ctl.RespError(code, err), err
	}

	// 3. 生成token
	aToken, rToken, err := myutils.GenerateToken(user.UserName, user.ID)
	if err != nil {
		code = e.GenerateTokenFailed
		log.LogrusObj.Errorln(err)
		return ctl.RespError(code, err), err
	}
	token := &types.TokenResp{
		AccessToken:  aToken,
		RefreshToken: rToken,
	}

	// 4. 返回用户信息
	data := types.UserInfoResp{
		ID:       user.ID,
		UserName: user.UserName,
		Avatar:   user.Avatar,
		Token:    token,
	}

	return ctl.RespSuccessWithData(code, data), nil
}

// CreateGroup 创建群组
func (s *UserService) CreateGroup(ctx context.Context, req *types.CreateGroupReq) (resp interface{}, err error) {
	code := e.SUCCESS
	groupDao := dao.GetGroupDao(ctx)

	// 1. 检查群组名是否存在
	group, err := groupDao.FindGroupByName(req.GroupName)
	if err != nil && err != gorm.ErrRecordNotFound {
		code = e.CreateGroupFailed
		log.LogrusObj.Errorln(err)
		return ctl.RespError(code, err), err
	} else if err == nil {
		code = e.GroupExist
		err = errors.New("group exist")
		log.LogrusObj.Errorln(err)
		return ctl.RespError(code, err), err
	}

	// 2. 创建群组
	group = &model.Group{}
	group.GroupName = req.GroupName
	err = groupDao.CreateGroup(group)
	if err != nil {
		code = e.CreateGroupFailed
		log.LogrusObj.Errorln(err)
		return ctl.RespError(code, err), err
	}

	// 3. 返回群组信息
	data := &types.GroupInfoResp{
		ID:        group.ID,
		GroupName: group.GroupName,
	}
	return ctl.RespSuccessWithData(code, data), nil
}
