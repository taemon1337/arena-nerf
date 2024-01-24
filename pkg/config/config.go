package config

import (
  "os"
  "log"

  "github.com/hashicorp/serf/serf"
  "github.com/hashicorp/serf/cmd/serf/command/agent"
)

type Config struct {
  NodeName        string          `yaml:"node_name" json:"node_name"`
  AgentConf       *agent.Config   `yaml:"agent_conf" json:"agent_conf"`
  SerfConf        *serf.Config    `yaml:"serf_conf" json:"serf_conf"`
}

func NewConfig() *Config {
  return &Config{
    NodeName: GetHostname(),
    AgentConf:  agent.DefaultConfig(),
    SerfConf:   serf.DefaultConfig(),
  }
}

func GetHostname() string {
  hostname, err := os.Hostname()
  if err != nil {
    log.Printf("cannot get hostname - %s", err)
    return ""
  }

  return hostname
}
