package node

import (
  "log"
  "fmt"
  "time"
  "errors"

  "github.com/hashicorp/serf/serf"
  "github.com/hashicorp/serf/cmd/serf/command/agent"

  "github.com/taemon1337/serf-cluster/pkg/game"
  "github.com/taemon1337/serf-cluster/pkg/config"
)

var (
  TAG_ROLE_NODE = "node"
  TAG_ROLE_CTRL = "ctrl"
  ERR_INVALID_CONFIG = errors.New("invalid node config")
  ERR_NOT_CONNECTED = errors.New("agent not connected")
  EVENT_ALIVE = "alive"
  COALESCE = false
  WAIT_TIME = 10 * time.Second
)

type Node struct {
  conf        *config.Config          `yaml:"config" json:"config"`
  agent       *agent.Agent            `yaml:"-" json:"-"`
  engine      *game.GameEngine        `yaml:"-" json:"-"`
}

func NewNode(cfg *config.Config) (*Node, error) {
  if cfg.AgentConf.NodeName == "" {
    return nil, ERR_INVALID_CONFIG
  }

  a, err := agent.Create(cfg.AgentConf, cfg.SerfConf, nil)
  if err != nil {
    return nil, err
  }

  return &Node{
    conf:   cfg,
    agent:  a,
    engine: nil,
  }, nil
}

func (n *Node) Name() string {
  return n.conf.AgentConf.NodeName
}

func (n *Node) Role() string {
  role, ok := n.conf.AgentConf.Tags["role"]
  if ok {
    return role
  }
  return ""
}

func (n *Node) Serf() *serf.Serf {
  return n.agent.Serf()
}

func (n *Node) Start() error {
  log.Printf("starting %s %s", n.Role(), n.Name())

  err := n.agent.Start()
  if err != nil {
    return errors.New(fmt.Sprintf("cannot start node - %s", err))
  }

  n.agent.RegisterEventHandler(n)

  for {
    select {
      case <-time.After(WAIT_TIME):
        if n.Role() == TAG_ROLE_CTRL {
          for _, member := range n.Serf().Members() {
            log.Printf("%s: status=%s, ip=%s, tags=%s", member.Name, member.Status, member.Addr, member.Tags)
          }
        }
    }
  }

  return nil
}

func (n *Node) MemberCount() int {
  return n.Serf().NumNodes()
}

func (n *Node) SendEvent(e string, payload []byte) error {
  if n.agent != nil {
    return n.agent.UserEvent(e, payload, COALESCE)
  }
  return ERR_NOT_CONNECTED
}

func (n *Node) SendQuery(q string, payload []byte, tags map[string]string) (*serf.QueryResponse, error) {
  if n.agent != nil {
    return n.agent.Query(q, payload, &serf.QueryParam{RequestAck: true, FilterTags: tags})
  }
  return nil, ERR_NOT_CONNECTED
}

func (n *Node) HandleEvent(e serf.Event) {
  if e.EventType() == serf.EventUser {
    log.Printf("EVENT: %s", e)
  }
  if e.EventType() == serf.EventQuery {
    log.Printf("QUERY: %s", e)
  }
}

