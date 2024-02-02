package controller

import (
  "log"
  "time"
  "golang.org/x/sync/errgroup"

  "github.com/hashicorp/serf/serf"

  "github.com/taemon1337/serf-cluster/pkg/game"
  "github.com/taemon1337/serf-cluster/pkg/config"
  "github.com/taemon1337/serf-cluster/pkg/constants"
  "github.com/taemon1337/serf-cluster/pkg/connector"
)

type Controller struct {
  conn        *connector.Connector
  conf        *config.Config
  game        *game.GameEngine
}

func NewController(cfg *config.Config, ge *game.GameEngine) *Controller {
  return &Controller{
    conn:     connector.NewConnector(cfg),
    conf:     cfg,
    game:     ge,
  }
}

func (c *Controller) Start() error {
  g := new(errgroup.Group)

  err := c.conn.Connect()
  if err != nil {
    return err
  }

  c.conn.RegisterEventHandler(c)

  g.Go(func () error {
    time.Sleep(5 * time.Second)
    return c.conn.Join()
  })

  g.Go(func () error {
    return c.ListenToGameEvents()
  })

  g.Go(func () error {
    return c.game.Run(c.conf.ExpectNodes, c.conf.Timeout)
  })

  return g.Wait()
}

func (c *Controller) ListenToGameEvents() error {
  for {
    select {
      case e := <-c.game.EventChan:
        err := c.conn.UserEvent(e.EventName, e.Payload, constants.COALESCE)
        if err != nil {
          log.Printf("could not send game event %s: %s", e.EventName, err)
        }
      case q := <-c.game.QueryChan:
        data := map[string][]byte{}
        resp, err := c.conn.Query(q.Query, q.Payload, &serf.QueryParam{FilterTags: q.Tags})
        if err != nil {
          q.Response <- game.NewGameQueryResponse(data, err)
        }
        for r := range resp.ResponseCh() {
          data[r.From] = r.Payload
        }
        q.Response <- game.NewGameQueryResponse(data, err)
    }
  }
}

func (c *Controller) HandleEvent(e serf.Event) {
  if e.EventType() == serf.EventUser {
    log.Printf("EVENT: %s", e)
  }
  if e.EventType() == serf.EventQuery {
    log.Printf("QUERY: %s", e)
  }
}

