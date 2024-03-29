package connector

import (
  "log"
  "time"

  "github.com/hashicorp/serf/serf"
  "github.com/hashicorp/serf/cmd/serf/command/agent"

  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type Connector struct {
  agent       *agent.Agent
  conf        *config.Config
}

func NewConnector(cfg *config.Config) *Connector {
  return &Connector{
    agent:      nil,
    conf:       cfg,
  }
}

func (c *Connector) IsConnected() bool {
  return c.agent != nil // TODO: improve connection detection
}

func (c *Connector) Connect() error {
  if c.conf.AgentConf == nil {
    return constants.ERR_NO_AGENT_CONFIG
  }

  c.conf.SerfConf.Tags = c.conf.AgentConf.Tags
  c.conf.SerfConf.NodeName = c.conf.AgentConf.NodeName

  if c.IsConnected() {
    return constants.ERR_EXISTING_CONNECTION
  }

  a, err := agent.Create(c.conf.AgentConf, c.conf.SerfConf, nil)
  if err != nil {
    return err
  }

  c.agent = a

  // start serf
  err = c.agent.Start()
  if err != nil {
    return nil
  }

  return nil
}

func (c *Connector) Join() error {
  for {
    i, err := c.agent.Join(c.conf.JoinAddrs, c.conf.JoinReplay)
    if err != nil {
      log.Printf("error joining %s: %s", c.conf.JoinAddrs, err)
    }

    if i > 0 {
      log.Printf("successfully joined %d nodes", i)
      return nil
    }
    time.Sleep(constants.WAIT_TIME)
  }
}

func (c *Connector) Query(name string, payload []byte, params *serf.QueryParam) (*serf.QueryResponse, error) {
  return c.agent.Query(name, payload, params)
}

func (c *Connector) UserEvent(name string, payload []byte, coalesce bool) error {
  return c.agent.UserEvent(name, payload, coalesce)
}

func (c *Connector) RegisterEventHandler(eh agent.EventHandler) {
  c.agent.RegisterEventHandler(eh)
}

func (c *Connector) Serf() *serf.Serf {
  return c.agent.Serf()
}



