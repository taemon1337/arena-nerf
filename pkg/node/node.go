package node

import (
  "log"
  "time"
  "strconv"
  "strings"
  "encoding/json"
  "golang.org/x/sync/errgroup"

  "github.com/hashicorp/serf/serf"

  "github.com/taemon1337/serf-cluster/pkg/config"
  "github.com/taemon1337/serf-cluster/pkg/constants"
  "github.com/taemon1337/serf-cluster/pkg/connector"
  "github.com/taemon1337/serf-cluster/pkg/sensor"
)

type Node struct {
  conn          *connector.Connector
  conf          *config.Config
  sensor        *sensor.Sensor
  mode          string
  state         string
  hits          map[string]int
  teams         map[string]string
}

func NewNode(cfg *config.Config) *Node {
  return &Node{
    conn:     connector.NewConnector(cfg),
    conf:     cfg,
    mode:     constants.GAME_MODE_NONE,
    state:    constants.GAME_STATE_INIT,
    sensor:   nil,
    hits:     map[string]int{cfg.AgentConf.NodeName: 0},
    teams:    map[string]string{},
  }
}

func (n *Node) Start() error {
  g := new(errgroup.Group)

  err := n.conn.Connect()
  if err != nil {
    return err
  }

  n.conn.RegisterEventHandler(n)

  if n.conf.Sensor {
    log.Printf("starting %s sensor", n.conf.AgentConf.NodeName)
    n.sensor = sensor.NewSensor(n.conf.AgentConf.NodeName)

    g.Go(func () error {
      return n.sensor.Start()
    })
  }

  g.Go(func () error {
    return n.conn.Join()
  })

  g.Go(func () error {
    return n.Listen()
  })

  return g.Wait()
}

func (n *Node) Ready() bool {
  return n.conn.IsConnected()
}

func (n *Node) Listen() error {
  for {
    time.Sleep(10 * time.Second)
  }
}

func (n *Node) HandleEvent(evt serf.Event) {
  if evt.EventType() == serf.EventUser {
    e := evt.(serf.UserEvent)
    switch e.Name {
      case constants.GAME_MODE:
        n.mode = string(e.Payload)
      case constants.GAME_ACTION_BEGIN:
        n.state = constants.GAME_STATE_ACTIVE
      case constants.GAME_ACTION_END:
        n.state = constants.GAME_STATE_OVER
      case n.NodeEventName(constants.TEAM_HIT):
        if n.state != constants.GAME_STATE_ACTIVE {
          log.Printf("game is not active - no hits allowed")
          return
        }

        parts := strings.Split(string(e.Payload), constants.SPLIT)
        if len(parts) < 2 {
          log.Printf("cannot parse team hit from %s - should be <team>:<count>", string(e.Payload))
        } else {
          hits, err := strconv.Atoi(parts[1])
          if err != nil {
            log.Printf("cannot parse team hit from %s - %s", string(e.Payload), err)
          } else {
            n.hits[parts[0]] += hits
            n.hits[n.conf.AgentConf.NodeName] += hits
          }
        }
      case constants.TEAM_ADD:
        teams := strings.Split(string(e.Payload), constants.SPLIT)
        for _, team := range teams {
          if _, ok := n.teams[team]; !ok {
            log.Printf("adding team %s", team)
            n.teams[team] = "" // only add team if it doesn't exist
          }
        }
      case constants.GAME_ACTION_RESET:
        if n.state == constants.GAME_STATE_ACTIVE {
          log.Printf("cannot reset an active game, must stop it first")
          return
        }

        n.state = constants.GAME_STATE_INIT
        n.hits = map[string]int{n.conf.AgentConf.NodeName: 0}
      default:
        log.Printf("warn: unrecognized event - %s", e.Name)
    }
  }
  if evt.EventType() == serf.EventQuery {
    var err error = nil
    q := evt.(*serf.Query)
    switch q.Name {
      case constants.NODE_READY:
        err = q.Respond([]byte(constants.NODE_IS_READY))
      case constants.GAME_MODE:
        err = q.Respond([]byte(n.mode))
      case constants.TEAM_HIT:
        data, err := json.Marshal(n.hits)
        if err != nil {
          log.Printf("cannot marshal node hits: %s", err)
        } else {
          err = q.Respond(data)
        }
      default:
        log.Printf("warn: unrecognized query - %s", q.Name)
    }

    if err != nil {
      log.Printf("error responding to query %s: %s", q.Name, err)
    }
  }
}

func (n *Node) NodeEventName(action string) string {
  return strings.Join([]string{n.conf.AgentConf.NodeName, action}, constants.SPLIT)
}
