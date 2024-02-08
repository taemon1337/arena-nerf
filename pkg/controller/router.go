package controller

import (
  "os"
  "log"
  "fmt"
  "strings"
  "net/http"
  "path/filepath"
  "github.com/gin-gonic/gin"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type PayloadForm struct {
  Payload       string      `yaml:"payload" json:"payload"`
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
    action := c.Param("action")

    if !ctrl.conf.AllowApiActions {
      c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s", constants.ERR_API_ACTIONS_NOT_ALLOWED)})
      return
    }

    var payForm PayloadForm
    c.ShouldBindJSON(&payForm)
    err := ctrl.ActionFromUi(action, payForm.Payload)

    if err != nil {
      log.Printf("cannot perform api action %s: %s", action, err)
      c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s", err)})
      return
    }

    c.JSON(http.StatusOK, gin.H{
      "message": fmt.Sprintf("%s action sent", action),
    })
  }
}

func (ctrl *Controller) ActionFromUi(action, payload string) error {
  switch action {
    case "ui:game:mode":
      if !ctrl.game.GameStats.Completed {
        return constants.ERR_ONGOING_GAME
      } 
      return ctrl.game.SendAction(constants.GAME_MODE, payload)
// TODO: how to properly async start the game and attach to waitgroup
//    case "ui:game:begin":
//      go ctrl.game.Run(ctrl.conf.ExpectNodes, ctrl.conf.Timeout)
//      return ctrl.game.SendAction(constants.GAME_ACTION_BEGIN, "web: Start the game!")
    case "ui:game:end":
      return ctrl.game.EndGame()
    case "ui:team:add":
      return ctrl.game.SendAction(constants.TEAM_ADD, payload)
    case "ui:team:del":
      return ctrl.game.SendAction(constants.TEAM_DEL, payload)
    default:
      return constants.ERR_UI_ACTION_NOT_ALLOWED
  }
}
