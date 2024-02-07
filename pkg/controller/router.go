package controller

import (
  "fmt"
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/taemon1337/serf-cluster/pkg/constants"
)

type PayloadForm struct {
  Payload       string      `form:"payload"`
}

func (ctrl *Controller) Router() {
  api := ctrl.server.Router.Group("api")
  v1 := api.Group("v1")
  {
    current := v1.Group("/current")
    {
      current.GET("/stats", ctrl.ApiGameStats())
      current.POST("/action/:action", ctrl.ApiAction())
    }
  }
}

func (ctrl *Controller) ApiGameStats() func (*gin.Context) {
  return func (c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "stats": ctrl.game.GameStats,
    })
  }
}

func (ctrl *Controller) ApiAction() func (*gin.Context) {
  return func (c *gin.Context) {
    var payForm PayloadForm
    err := constants.ERR_API_ACTIONS_NOT_ALLOWED
    if ctrl.conf.AllowApiActions {
      c.ShouldBind(&payForm)
      err = ctrl.game.Send(c.Param("action"), payForm.Payload)
    }

    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error": err,
      })
    } else {
      c.JSON(http.StatusOK, gin.H{
        "message": fmt.Sprintf("%s command sent", c.Param("action")),
      })
    }
  }
}
