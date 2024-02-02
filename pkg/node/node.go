package node

import (
  "log"
  "time"
  "golang.org/x/sync/errgroup"

  "github.com/hashicorp/serf/serf"

  "github.com/taemon1337/serf-cluster/pkg/config"
  "github.com/taemon1337/serf-cluster/pkg/connector"
)

type Node struct {
  conn          *connector.Connector
}

func NewNode(cfg *config.Config) *Node {
  return &Node{
    conn:   connector.NewConnector(cfg),
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

func (n *Node) Listen() error {
  for {
    time.Sleep(10 * time.Second)
  }
}

func (n *Node) HandleEvent(e serf.Event) {
  if e.EventType() == serf.EventUser {
    log.Printf("EVENT: %s", e)
  }
  if e.EventType() == serf.EventQuery {
    log.Printf("QUERY: %s", e)
  }
}
