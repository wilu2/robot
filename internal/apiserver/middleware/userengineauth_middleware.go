package middleware

import (
	"github.com/gin-gonic/gin"
)

func UserEngineAuthMiddleware(c *gin.Context) {
	// var (
	// 	u = c.Keys["user"].(*model.User)
	// )
	// userObj := &model.User{}
	// dbIns, _ := dal.GetMySQLFactoryOr(nil)
	// dbIns.GetDb().Where(&model.User{ID: u.ID, Status: consts.StateEnable}).First(
	// 	userObj)
	// if userObj.ID == 0 { // 存在用户删除，jwt 还存在的情况
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }
	c.Next()
}
