package setting

import (
	"context"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/setting"
	"financial_statement/internal/pkg/settings"
	"github.com/gin-gonic/gin"
)

type SettingInfoLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewSettingInfoLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) SettingInfoLogic {
	return SettingInfoLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// SettingInfo 查看配置
func (l *SettingInfoLogic) SettingInfo() (resp setting.GetSettingResp, err error) {
	setting_map := settings.GetSetting()
	resp.FailedLogin = setting_map["failed_logins"]
	resp.LockedTime = setting_map["locked_time"]
	resp.NotLoginTime = setting_map["not_login_time"]
	resp.SessionTime = setting_map["session_time"]
	resp.ValidTime = setting_map["valid_time"]
	return
}
