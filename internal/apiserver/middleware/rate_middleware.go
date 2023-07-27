package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"financial_statement/internal/apiserver/dal"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var ratePathMap = map[string]string{
	"/api/v1/email/check_exist":  "20-5m", // 5 分钟限制 20 次
	"/api/v1/email/confirmation": "20-5m", // 邮件链接激活
	"/api/v1/email/send":         "20-5m", // 发送邮件
	"/api/v1/user/login":         "20-5m", // 登陆接口
	"/api/v1/user/register":      "20-5m", // 注册接口
	"/api/v1/user/reset_pwd":     "20-5m", // 重置密码接口
}

type Dispatcher struct {
	redisClient redis.UniversalClient    // redis 连接池
	prefix      string                   // redis key 前缀
	limitMap    map[string]int           // uri 限制的频次
	intervalMap map[string]time.Duration // uri 限制的时间间隔
}

var dispatcher Dispatcher

// RateMiddleware 限流中间件
func RateMiddleware(c *gin.Context) {
	urlPath := c.FullPath()
	if _, ok := ratePathMap[urlPath]; !ok {
		c.Next()
		return
	}
	redisKey := dispatcher.prefix + ":" + urlPath + ":" + c.ClientIP()
	if curr, _ := dispatcher.redisClient.Get(context.Background(), redisKey).Int(); curr == 0 {
		dispatcher.redisClient.Set(context.Background(), redisKey, "1", dispatcher.intervalMap[urlPath])
	} else if curr <= dispatcher.limitMap[urlPath] {
		dispatcher.redisClient.IncrBy(context.Background(), redisKey, 1)
	} else {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"code":    "429",
			"message": "too many requests",
		})
		return
	}
	c.Next()
}

// ParseCommand 解析命令
func ParseCommand(command string) (limit int, interval time.Duration, err error) {
	values := strings.Split(command, "-")
	limit, _ = strconv.Atoi(values[0])
	interval, err = time.ParseDuration(values[1])
	return
}

// InitRateMiddleware 初始化限流中间件
func InitRateMiddleware() {
	dbIns, _ := dal.GetRedisFactoryOr(nil)
	client := dbIns.GetRedisDb()
	dispatcher = Dispatcher{
		redisClient: client,
		prefix:      "rate",
		limitMap:    make(map[string]int),
		intervalMap: make(map[string]time.Duration),
	}
	for uri, command := range ratePathMap {
		limit, interval, err := ParseCommand(command)
		if err != nil {
			continue
		}
		dispatcher.limitMap[uri] = limit
		dispatcher.intervalMap[uri] = interval
	}
}
