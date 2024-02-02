package config

import (
  "os"
  "log"
  "strings"

  "github.com/hashicorp/serf/serf"
  "github.com/hashicorp/serf/cmd/serf/command/agent"
)

type Config struct {
  AgentConf       *agent.Config   `yaml:"agent_conf" json:"agent_conf"`
  SerfConf        *serf.Config    `yaml:"serf_conf" json:"serf_conf"`
  JoinAddrs       []string        `yaml:"join_addrs" json:"join_addrs"`
  JoinReplay      bool            `yaml:"join_replay" json:"join_replay"`
  ExpectNodes     int             `yaml:"expect_nodes" json:"expect_nodes"`
  Timeout         int             `yaml:"timeout" json:"timeout"`
}

func NewConfig(role string) *Config {
  ac := agent.DefaultConfig()
  sc := serf.DefaultConfig()
  joinaddrs := Getenv("SERF_JOIN_ADDRS", "127.0.0.1")
  joinreplay := Getenv("SERF_JOIN_REPLAY", "") // default is false, set to 'true|True|TRUE' otherwise
  ac.NodeName = Getenv("SERF_NAME", GetHostname())
  ac.BindAddr = Getenv("SERF_BIND_ADDR", "")
  ac.AdvertiseAddr = Getenv("SERF_ADVERTISE_ADDR", "")
  ac.EncryptKey = Getenv("SERF_ENCRYPT_KEY", "")
  ac.Tags["role"] = Getenv("SERF_ROLE", role) // role is stored as tag

  sc.Tags = ac.Tags
  sc.NodeName = ac.NodeName

  return &Config{
    AgentConf:    ac,
    SerfConf:     sc,
    JoinAddrs:    strings.Split(joinaddrs, ","),
    JoinReplay:   (joinreplay == "true" || joinreplay == "True" || joinreplay == "TRUE"),
    ExpectNodes:  3,
    Timeout:      10,
  }
}

func Getenv(key, val string) string {
  a, exists := os.LookupEnv(key)
  if a != "" && exists {
    return a
  }
  return val // default
}

func GetHostname() string {
  hostname, err := os.Hostname()
  if err != nil {
    log.Printf("cannot get hostname - %s", err)
    return ""
  }

  return hostname
}
