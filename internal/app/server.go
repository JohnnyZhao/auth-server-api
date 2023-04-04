package app

import "github.com/gin-gonic/gin"

type Server struct {
	api    *Api
	router *gin.Engine
}

func NewServer(api *Api) *Server {
	router := gin.Default()
	return &Server{api: api, router: router}
}

func (s *Server) Setup() {
	s.router.POST("/signup", s.api.CreateUser)
	s.router.Use(s.api.AuthMiddleware())
	{
		s.router.GET("/users/:user_id", s.api.GetByUserID)
		s.router.PATCH("/users/:user_id", s.api.PatchByID)
		s.router.POST("/close", s.api.DeleteByID)
	}
}

func (s *Server) Run(addr string) {
	s.router.Run(addr)
}
