package api

import (
"context"
"net/http"

"github.com/gin-gonic/gin"
"github.com/khryptorgraphics/ollama-nova/internal/inference"
)

type Server struct {
engine *inference.Engine
router *gin.Engine
}

func NewServer(engine *inference.Engine) *Server {
return &Server{
engine: engine,
router: gin.Default(),
}
}

func (s *Server) SetupRoutes() {
s.router.POST("/api/generate", s.handleGenerate)
s.router.GET("/api/models", s.handleListModels)
s.router.GET("/health", s.handleHealth)
}

func (s *Server) handleGenerate(c *gin.Context) {
var req inference.Request
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

resp, err := s.engine.Process(context.Background(), &req)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, resp)
}

func (s *Server) handleListModels(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"models": []string{"llama2", "mistral"}})
}

func (s *Server) handleHealth(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func (s *Server) Start(addr string) error {
s.SetupRoutes()
return s.router.Run(addr)
}
