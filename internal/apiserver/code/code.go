package code

// 业务错误码，1(项目编码)01(模块编码)00(错误编码)
const (
	Success         = 200
	BadRequest      = 400
	Unauthorized    = 401
	Forbidden       = 403
	NotFound        = 404
	TooManyRequests = 429
	Internal        = 500

	//财报模块错误码 START

	//文件不存在
	FileDoesNotExist = 10100
	//END
)

// 业务状态码和http状态码以及错误信息对应，可以根据需要自行添加
// 调用 WithCodeMsg(Success, "") 默认用定义的错误信息
func init() {
	register(Success, 200, "success")
	register(BadRequest, 400, "bad request")
	register(Unauthorized, 401, "unauthorized")
	register(Forbidden, 403, "forbidden")
	register(NotFound, 404, "not found")
	register(TooManyRequests, 429, "too many requests")
	register(Internal, 500, "internal server error")
	register(FileDoesNotExist, 200, "file does not exist")
}
