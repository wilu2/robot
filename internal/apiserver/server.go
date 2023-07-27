package apiserver

import (
	"financial_statement/internal/apiserver/config"
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/job"
	"financial_statement/internal/apiserver/middleware"
	"financial_statement/internal/apiserver/routes"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/pkg/settings"
	genericOptions "financial_statement/internal/pkg/options"
	genericapiserver "financial_statement/internal/pkg/server"
	storage "financial_statement/internal/pkg/storage"
	"financial_statement/pkg/shutdown"

	"financial_statement/pkg/log"
)

type apiServer struct {
	cfg                 *config.Config
	gs                  *shutdown.GracefulShutdown
	mysqlOption         *genericOptions.MySQLOptions
	dm8Option           *genericOptions.DM8Options
	redisOptions        *genericOptions.RedisOptions
	localStorageOptions *genericOptions.LocalStorageOptions
	nfsStorageOptions   *genericOptions.NfsStorageOptions
	genericApiServer    *genericapiserver.GenericApiServer
}

// createAPIServer 创建 API Server 服务
func createAPIServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(shutdown.NewPosixSignalManager()) // 优雅关闭系统

	genericConfig, err := buildGenericConfig(cfg) // 创建 GenericAPIServer 配置
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New() // 初始化完成 HTTP 服务
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		cfg:                 cfg,
		gs:                  gs,
		genericApiServer:    genericServer,
		mysqlOption:         cfg.MySQLOptions,
		dm8Option:           cfg.DM8Options,
		redisOptions:        cfg.RedisOptions,
		nfsStorageOptions:   cfg.NfsStorageOptions,
		localStorageOptions: cfg.LocalStorageOptions,
	}

	return server, nil
}

// Run  业务代码路由地址
func (s *apiServer) Run() error {
	s.initDbStore()
	s.initRedisStore()
	s.initFileStorage()
	s.InitSetting()
	jobMgr := s.initJobMgr()
	middleware.InitJWTAuth()        // 初始化 JWT 认证
	middleware.InitRateMiddleware() // 初始化限流中间件

	svcCtx := svc.NewServiceContext(s.cfg, jobMgr)
	routes.Setup(s.genericApiServer.Engine, svcCtx)

	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		log.Info("Close api Server")
		s.genericApiServer.Close()
		return nil
	}))

	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}
	return s.genericApiServer.Run()
}

// buildGenericConfig 根据命令行参数创建 GenericAPIServer 配置
func buildGenericConfig(cfg *config.Config) (genericConfig *genericapiserver.Config, lastErr error) {
	genericConfig = genericapiserver.NewConfig()
	if lastErr = cfg.ServerOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}
	return
}

// initMysqlStore 数据库连接初始化
func (s *apiServer) initDbStore() {
	if s.cfg.ServerOptions.Db == "mysql" {
		log.Debugf("mysql option:%v", s.mysqlOption)
		_, err := dal.GetDbFactoryOr(s.mysqlOption)
		if err != nil {
			log.Errorf("init mysql store failed: %s", err.Error())
			return
		}
	} else {
		_, err := dal.GetDbFactoryOr(s.dm8Option)
		if err != nil {
			log.Errorf("init dameng8 store failed: %s", err.Error())
			return
		}
	}
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		log.Infof("Close db")
		dbIns, _ := dal.GetDbFactoryOr(nil)
		return dbIns.CloseDb()
	}))
}

// initFileStorage 文件存储初始化
func (s *apiServer) initFileStorage() {
	if s.cfg.ServerOptions.FileStorage == "local" {
		storage.NewFileStore(s.localStorageOptions)
	} else {
		storage.NewFileStore(s.nfsStorageOptions)
	}
}

// initRedisStore 初始化 Redis 数据库链接
func (s *apiServer) initRedisStore() {
	log.Debugf("redis option:%v", s.redisOptions)
	_, err := dal.GetRedisFactoryOr(s.redisOptions)
	if err != nil {
		log.Errorf("init redis store failed: %s", err.Error())
		return
	}
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		log.Infof("Close Redis db")
		dbIns, _ := dal.GetRedisFactoryOr(nil)
		return dbIns.CloseRedisDb()
	}))
}

func (s *apiServer) initJobMgr() *job.JobMgr {
	jobMgr := job.InitJobMgr()
	return jobMgr
}

func (s *apiServer) InitSetting() {
	settings.InitSetting()
}