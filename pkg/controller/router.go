package controller

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func (ctrl *Controller) Router() {
  api := ctrl.server.Router.Group("api")
  v1 := api.Group("v1")
  {
    v1.GET("/scoreboard", ctrl.Scoreboard())
  }
}

func (ctrl *Controller) Scoreboard() func (*gin.Context) {
  return func (c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "scoreboard": ctrl.game.Scoreboard,
    })
  }
}
