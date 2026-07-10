package app

import (
	"blog_platform/config"
	"blog_platform/internal/routes"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Server struct {
	E   *echo.Echo
	DB  *gorm.DB
	Cfg *config.Config
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	e := echo.New()
	routes.RegisterRoutes(e, db, cfg)
	return &Server{E: e, DB: db, Cfg: cfg}
}

func (s *Server) Start(addr string) error {
	return s.Start(addr)
}
