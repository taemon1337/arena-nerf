package controller

import (
  "log"
  "time"
  "golang.org/x/sync/errgroup"

  "github.com/hashicorp/serf/serf"

  "github.com/taemon1337/serf-cluster/pkg/config"
  "github.com/taemon1337/serf-cluster/pkg/connector"
)

type Controller struct {
  conn        *connector.Connector
}

func NewController(cfg *config.Config) *Controller {
  return &Controller{
    conn:     connector.NewConnector(cfg),
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
    return c.Listen()
  })

  return g.Wait()
}

func (c *Controller) Listen() error {
  for {
    select {
      case <-time.After(config.WAIT_TIME):
        for _, member := range c.conn.Serf().Members() {
          log.Printf("%s: status=%s, ip=%s, tags=%s", member.Name, member.Status, member.Addr, member.Tags)
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

