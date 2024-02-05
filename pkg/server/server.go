package server

import (
  "time"
  "net/http"
  "crypto/tls"
  "github.com/gin-gonic/gin"
)

type Server struct {
  Router      *gin.Engine
  TLS         *tls.Config
}

func (s *Server) ListenAndServe(addr string) error {
  srv := &http.Server{
    Addr:           addr,
    Handler:        s.Router,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
    TLSConfig:      s.TLS,
  }

  if s.TLS == nil {
    return srv.ListenAndServe()
  } else {
    return srv.ListenAndServeTLS("", "")
  }
}

func NewServer() *Server {
  g := gin.Default()

  g.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "OK",
    })
  })

  return &Server{
    Router:   g, 
    TLS:      nil,
  }
}

