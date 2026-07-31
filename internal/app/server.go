package app

import (
	"blog_platform/config"
	"blog_platform/internal/middleware"
	"blog_platform/internal/routes"
	"blog_platform/internal/workers"

	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Server struct {
	E   *echo.Echo
	DB  *gorm.DB
	Cfg *config.Config
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	e := echo.New()
	e.HTTPErrorHandler = middleware.GlobalErrorHandler
	routes.RegisterRoutes(e, db, cfg)
	return &Server{E: e, DB: db, Cfg: cfg}
}

func (s *Server) Start(addr string) error {

	viewWorker := workers.NewViewWorker(s.DB)

	c := cron.New()
	c.AddFunc("@every 10m", viewWorker.ProcessViews)

	c.Start()
	return s.E.Start(addr)
}
