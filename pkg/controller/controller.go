package controller

import (
  "log"
  "time"
  "errors"
  "strings"
  "golang.org/x/sync/errgroup"

  "github.com/hashicorp/serf/serf"

  "github.com/taemon1337/serf-cluster/pkg/game"
  "github.com/taemon1337/serf-cluster/pkg/server"
  "github.com/taemon1337/serf-cluster/pkg/config"
  "github.com/taemon1337/serf-cluster/pkg/constants"
  "github.com/taemon1337/serf-cluster/pkg/connector"
)

type Controller struct {
  conn        *connector.Connector
  conf        *config.Config
  server      *server.Server
  game        *game.GameEngine
}

func NewController(cfg *config.Config) *Controller {
  return &Controller{
    conn:     connector.NewConnector(cfg),
    conf:     cfg,
    server:   nil,
    game:     nil,
  }
}

func (c *Controller) Start() error {
  g := new(errgroup.Group)

  if c.conf.Webserver {
    c.server = server.NewServer()
    c.Router()

    g.Go(func () error {
      return c.server.ListenAndServe(c.conf.WebAddr)
    })
  }

  err := c.conn.Connect()
  if err != nil {
    return err
  }

  c.conn.RegisterEventHandler(c)

  g.Go(func () error {
    time.Sleep(5 * time.Second)
    return c.conn.Join()
  })

  return g.Wait()
}

func (c *Controller) RunGame(ge *game.GameEngine) error {
  if ge != nil && ge.Active() {
    return errors.New("already running a game, must clear current game first before running another")
  }
  c.game = ge

  g := new(errgroup.Group)

  g.Go(func () error {
    return c.ListenToGameEvents(ge)
  })

  g.Go(func () error {
    return ge.Run(c.conf.ExpectNodes, c.conf.Timeout)
  })

  return nil
}

func (c *Controller) ListenToGameEvents(ge *game.GameEngine) error {
  for {
    select {
      case e := <-ge.EventChan:
        var err error = nil
        switch e.EventName {
          case constants.GAME_SHUTDOWN:
            return constants.ERR_SHUTDOWN
          case constants.TEAM_ADD:
            // inject teams from config into event
            err = c.conn.UserEvent(e.EventName, []byte(strings.Join(c.conf.Teams, constants.SPLIT)), constants.COALESCE)
          default:
            err = c.conn.UserEvent(e.EventName, e.Payload, constants.COALESCE)
        }
        if err != nil {
          log.Printf("could not send game event %s: %s", e.EventName, err)
        }
      case q := <-ge.QueryChan:
        switch q.Query {
          case constants.TEAM_QUERY:
            // if team query, respond with all registered teams
            data := map[string][]byte{}
            for _, team := range c.conf.Teams {
              data[team] = []byte{}
            }
            q.Response <- game.NewGameQueryResponse(data, nil)
          case constants.RANDOM_NODE:
            data := map[string][]byte{}
            nodes := c.conn.Serf().Members()
            data[nodes[0].Name] = []byte{}
            q.Response <- game.NewGameQueryResponse(data, nil)
          default:
            // by default send all queries from game engine to all nodes
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
}

func (c *Controller) HandleEvent(e serf.Event) {
  if e.EventType() == serf.EventUser {
    log.Printf("EVENT: %s", e)
  }
  if e.EventType() == serf.EventQuery {
    log.Printf("QUERY: %s", e)
  }
}

