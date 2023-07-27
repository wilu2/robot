package setting

import (
	"context"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/setting"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/pkg/settings"
	"financial_statement/internal/apiserver/code"
	"financial_statement/pkg/errors"
	"financial_statement/pkg/log"
	"github.com/gin-gonic/gin"
	
)

type SettingUpdateLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewSettingUpdateLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) SettingUpdateLogic {
	return SettingUpdateLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// SettingUpdate 修改配置信息
func (l *SettingUpdateLogic) SettingUpdate(req *setting.UpdateSettingReq) (err error) {
	//验证管理员权限
	u := l.ginCtx.Keys["user"].(*model.User)
	if u.Account != "admin"{
		err = errors.WithCodeMsg(code.BadRequest, "非管理员无法更改配置")
		return
	}

	//更新配置信息
	newSetting := make(map[string]int32,5)
	newSetting["failed_logins"] = req.FailedLogin
	newSetting["locked_time"] = req.LockedTime
	newSetting["valid_time"] = req.ValidTime
	newSetting["session_time"] = req.SessionTime
	newSetting["not_login_time"] = req.NotLoginTime
	error := settings.UpdateSetting(newSetting)
	if error != nil{
		log.Errorf("update settings error: %s", error)
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	return
}
