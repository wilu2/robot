package settings

import (
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/pkg/log"
)

var setting map[string]int32

func InitSetting(){
	setting = make(map[string]int32, 5)
	dbIns, _:= dal.GetDbFactoryOr(nil)
	setting_kv := make([]model.Setting,5,5)
	result := dbIns.GetDb().Model(&model.Setting{}).Find(&setting_kv)
	if result.Error != nil{
		log.Errorf("create captcha error: %s", result.Error)
		return
	}
	for i := range setting_kv{
		setting[setting_kv[i].Key] = setting_kv[i].Setting
	}
}

func GetSetting() map[string]int32{
	return setting
}

func UpdateSetting(newSetting map[string]int32) error{
	dbIns, _:= dal.GetDbFactoryOr(nil)
	for k,v :=range newSetting{
		if setting[k] != v{
			result := dbIns.GetDb().Model(&model.Setting{}).Where("`key` = ?", k).Update("setting", v)
			if result.Error != nil{
				log.Errorf("create captcha error: %s", result.Error)
				return result.Error
			}
			setting[k] = v
		}
	}
	return nil
}