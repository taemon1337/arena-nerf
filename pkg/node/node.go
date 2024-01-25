package node

import (
  "log"
  "fmt"
  "time"
  "errors"

  "github.com/hashicorp/serf/serf"
  "github.com/hashicorp/serf/cmd/serf/command/agent"

  "github.com/taemon1337/serf-cluster/pkg/config"
)

var (
  TAG_ROLE_NODE = "node"
  TAG_ROLE_CTRL = "ctrl"
  ERR_INVALID_CONFIG = errors.New("invalid node config")
  EVENT_ALIVE = "alive"
  COALESCE = false
  WAIT_TIME = 10 * time.Second
)

type Node struct {
  conf        *config.Config          `yaml:"config" json:"config"`
  agent       *agent.Agent            `yaml:"-" json:"-"`
}

func NewNode(cfg *config.Config) *Node {
  return &Node{
    conf:   cfg,
    agent:  nil,
  }
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
  if n.Name() == "" {
    return ERR_INVALID_CONFIG
  }

  log.Printf("starting node %s", n.Name())

  a, err := agent.Create(n.conf.AgentConf, n.conf.SerfConf, nil)
  if err != nil {
    return err
  }

  n.agent = a
  log.Printf("starting %s %s", n.Role(), n.Name())

  err = n.agent.Start()
  if err != nil {
    return errors.New(fmt.Sprintf("cannot start node - %s", err))
  }

  for {
    select {
      case <-time.After(WAIT_TIME):
        if n.Role() == TAG_ROLE_CTRL {
          for _, member := range n.Serf().Members() {
            log.Printf("%s: status=%s, ip=%s, tags=%s", member.Name, member.Status, member.Addr, member.Tags)
          }
        } else {
          err := n.agent.UserEvent(EVENT_ALIVE, []byte(fmt.Sprintf("%s: I am alive!", n.Name())), COALESCE)
          if err != nil {
            log.Printf("cannot send alive event: %s", err)
          }
        }
    }
  }

  return nil
}
