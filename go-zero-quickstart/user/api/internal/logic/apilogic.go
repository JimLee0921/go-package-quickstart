// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"go-zero-quickstart/user/api/internal/svc"
	"go-zero-quickstart/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiLogic {
	return &ApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiLogic) Api(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return &types.Response{Message: "JimLee"}, nil
}
