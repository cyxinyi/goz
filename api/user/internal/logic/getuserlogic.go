// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"awesomeProject3/go-zero/api/user/internal/svc"
	"awesomeProject3/go-zero/api/user/internal/types"
	"awesomeProject3/go-zero/rpc/user/user"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserRequest) (resp *types.GetUserResponse, err error) {
	// todo: add your logic here and delete this line

	userInfo, errUser := l.svcCtx.UserRpc.GetUser(l.ctx, &user.GetUserRequest{
		Id: req.Id,
	})

	if errUser != nil {
		return nil, errUser
	}

	return &types.GetUserResponse{
		Id:   userInfo.Id,
		Name: userInfo.Name,
	}, err
}
