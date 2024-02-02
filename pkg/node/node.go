package node

import (
  "log"
  "time"
  "golang.org/x/sync/errgroup"

  "github.com/hashicorp/serf/serf"

  "github.com/taemon1337/serf-cluster/pkg/config"
  "github.com/taemon1337/serf-cluster/pkg/constants"
  "github.com/taemon1337/serf-cluster/pkg/connector"
)

type Node struct {
  conn          *connector.Connector
  mode          string
  state         string
}

func NewNode(cfg *config.Config) *Node {
  return &Node{
    conn:   connector.NewConnector(cfg),
    mode:   constants.GAME_MODE_NONE,
    state:  constants.GAME_STATE_INIT,
  }
}

func (n *Node) Start() error {
  g := new(errgroup.Group)

  err := n.conn.Connect()
  if err != nil {
    return err
  }

  n.conn.RegisterEventHandler(n)

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
      default:
        log.Printf("warn: unrecognized query - %s", q.Name)
    }

    if err != nil {
      log.Printf("error responding to query %s: %s", q.Name, err)
    }
  }
}
