package node

import (
  "log"
  "errors"

  "github.com/hashicorp/serf/cmd/serf/command/agent"
  "github.com/taemon1337/serf-cluster/pkg/config"
)

var (
  ERR_INVALID_CONFIG = errors.New("invalid node config")
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

func (n *Node) Start() error {
  if n.conf.NodeName == "" {
    return ERR_INVALID_CONFIG
  }

  log.Printf("Config: %s", n.conf)

  a, err := agent.Create(n.conf.AgentConf, n.conf.SerfConf, nil)
  if err != nil {
    return err
  }

  n.agent = a
  log.Printf("starting serf agent.")

  return n.agent.Start()
}
