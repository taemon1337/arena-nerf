package controller

import (
  "os"
  "fmt"
  "strings"
  "net/http"
  "path/filepath"
  "github.com/gin-gonic/gin"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type PayloadForm struct {
  Payload       string      `form:"payload"`
}

func (ctrl *Controller) Router() {
  api := ctrl.server.Router.Group("api")
  v1 := api.Group("v1")
  {
    v1.GET("/games/:uuid", ctrl.ApiGameStats())
    v1.POST("/do/:action", ctrl.ApiAction())
  }
}

func (ctrl *Controller) ApiGameStats() func (*gin.Context) {
  return func (c *gin.Context) {
    uuid := c.Param("uuid")
    switch uuid {
      case "all":
        // send listing of all archived game logs
        entries, err := os.ReadDir(ctrl.conf.Logdir)
        if err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{"error": err})
          return
        }
        uuids := []string{"current"}
        for _, e := range entries {
          if strings.HasSuffix(e.Name(), ".json") {
            uuids = append(uuids, strings.TrimSuffix(e.Name(), ".json"))
          }
        }
        c.JSON(http.StatusOK, gin.H{
          "games": uuids,
        })
        return
      case "current":
        // send current game stats
        c.JSON(http.StatusOK, gin.H{
          "stats": ctrl.game.GameStats,
        })
        return
      default:
        // send archived log file
        http.ServeFile(c.Writer, c.Request, filepath.Join(ctrl.conf.Logdir, uuid + ".json"))
        return
    }
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
      c.JSON(http.StatusInternalServerError, gin.H{"error": err})
    } else {
      c.JSON(http.StatusOK, gin.H{
        "message": fmt.Sprintf("%s command sent", c.Param("action")),
      })
    }
  }
}
