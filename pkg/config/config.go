package config

import (
  "os"
  "log"

  "github.com/hashicorp/serf/serf"
  "github.com/hashicorp/serf/cmd/serf/command/agent"
)

type Config struct {
  AgentConf       *agent.Config   `yaml:"agent_conf" json:"agent_conf"`
  SerfConf        *serf.Config    `yaml:"serf_conf" json:"serf_conf"`
}

func NewConfig() *Config {
  ac := agent.DefaultConfig()
  sc := serf.DefaultConfig()

  ac.NodeName = GetHostname()

  return &Config{
    AgentConf:  ac,
    SerfConf:   sc,
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
