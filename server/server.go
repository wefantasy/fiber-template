package server

import (
	v1 "app/api/http/v1"
	"app/api/http/v1/auth"
	"app/conf"
	"app/db"
	"app/i18n"
	"app/log"
	"app/middleware"
	"app/scheduler"
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "app/docs"
	"app/repo"
	"app/serv"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func NewServer() (*Server, error) {
	// 初始化
	conf.Initialize()
	log.Initialize()
	db.Initialize()
	i18n.Initialize()
	scheduler.Initialize()

	// 初始化数据库
	userRepo := repo.NewUserRepo()
	repos := []repo.BaseRepo{userRepo}

	// 初始化服务
	userService := serv.NewUserService(userRepo)
	services := []serv.BaseServ{userService}

	// 初始化API
	commonController := v1.NewCommonController(userService)
	controllers := []v1.BaseContro{commonController}
	authUserController := auth.NewUserController(userService)
	authControllers := []v1.BaseContro{authUserController}

	// 初始化Fiber
	app := fiber.New(fiber.Config{
		ServerHeader: conf.AppName,
		AppName:      conf.AppName,
	})
	app.Use(middleware.Recover())
	app.Use(middleware.RequestId())
	app.Use(middleware.Logger())
	app.Use(cors.New())
	app.Use(middleware.Limiter())
	app.Use(middleware.CircuitBreaker())
	app.Use(middleware.ErrorParse())
	app.Use(middleware.Swagger())
	app.Use(healthcheck.New())
	app.Hooks().OnRoute(middleware.HookRoute)
	app.Get("/metrics", monitor.New())

	return &Server{
		engine:          app,
		repos:           repos,
		services:        services,
		controllers:     controllers,
		authControllers: authControllers,
	}, nil
}

type Server struct {
	engine          *fiber.App
	repos           []repo.BaseRepo
	services        []serv.BaseServ
	controllers     []v1.BaseContro
	authControllers []v1.BaseContro
}

func (s *Server) Close() {
	if db.DB != nil {
		if err := db.DB.Close(); err != nil {
			log.Error(err)
		}
	}
	if conf.Redis.Enable {
		if err := db.RDB.Close(); err != nil {
			log.Error(err)
		}
	}
}

func (s *Server) initRouter() {
	app := s.engine
	api := app.Group("/api/v1")
	controllers := make([]string, 0, len(s.controllers))
	for _, router := range s.controllers {
		router.RegisterRoute(api)
		controllers = append(controllers, router.Name())
	}
	log.Infof("api enabled controllers: %v", controllers)

	// 注册需要鉴权的路由
	authApi := app.Group("/api/v1")
	// authApi.Use(middleware.JwtAuth())
	authControllers := make([]string, 0, len(s.authControllers))
	for _, router := range s.authControllers {
		router.RegisterRoute(authApi)
		authControllers = append(authControllers, router.Name())
	}
	log.Infof("auth api enabled controllers: %v", authControllers)
}

func (s *Server) Run() error {
	defer s.Close()

	s.initRouter()
	addr := conf.Server.Address + ":" + conf.Server.Port
	log.Infof("Start server on: http://%s", addr)
	go func() {
		if err := s.engine.Listen(addr); err != nil && err.Error() != "http: Server closed" {
			log.Panicf("Failed to start server, %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	dTime, err := strconv.Atoi(conf.Server.Port)
	if err != nil {
		log.Panic(err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(dTime)*time.Second)
	defer cancel()
	ch := <-sig
	log.Infof("Receive signal: %s", ch)
	return s.engine.ShutdownWithContext(ctx)
}
