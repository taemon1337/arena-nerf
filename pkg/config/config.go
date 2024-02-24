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
  SensorConf      *SensorConfig   `yaml:"sensor_conf" json:"sensor_conf"`
  JoinAddrs       []string        `yaml:"join_addrs" json:"join_addrs"`
  Teams           []string        `yaml:"teams" json:"teams"`
  JoinReplay      bool            `yaml:"join_replay" json:"join_replay"`
  ExpectNodes     int             `yaml:"expect_nodes" json:"expect_nodes"`
  Timeout         int             `yaml:"timeout" json:"timeout"`
  Webserver       bool            `yaml:"webserver" json:"webserver"`
  WebAddr         string          `yaml:"webaddr" json:"webaddr"`
  Gametime        string          `yaml:"gametime" json:"gametime"`
  AllowApiActions bool            `yaml:"allow_api_actions" json:"allow_api_actions"`
  Logdir          string          `yaml:"logdir" json:"logdir"`
}

func NewConfig(role string) *Config {
  ac := agent.DefaultConfig()
  sc := serf.DefaultConfig()
  joinaddrs := Getenv("SERF_JOIN_ADDRS", "127.0.0.1")
  joinreplay := Getenv("SERF_JOIN_REPLAY", "") // default is false, set to 'true|True|TRUE' otherwise
  ac.NodeName = Getenv("SERF_NAME", GetHostname())
  ac.BindAddr = Getenv("SERF_BIND_ADDR", "0.0.0.0")
  ac.AdvertiseAddr = Getenv("SERF_ADVERTISE_ADDR", "")
  ac.EncryptKey = Getenv("SERF_ENCRYPT_KEY", "")
  ac.Tags["role"] = Getenv("SERF_ROLE", role) // role is stored as tag

  return &Config{
    AgentConf:        ac,
    SerfConf:         sc,
    SensorConf:       DefaultSensorConfig(),
    JoinAddrs:        strings.Split(joinaddrs, ","),
    JoinReplay:       (joinreplay == "true" || joinreplay == "True" || joinreplay == "TRUE"),
    Teams:            []string{},
    ExpectNodes:      3,
    Timeout:          10,
    Webserver:        false,
    WebAddr:          ":8080",
    Gametime:         "5m",
    AllowApiActions:  false,
    Logdir:           "",
  }
}

func (c *Config) Validate() error {
  ac := c.AgentConf
  sc := c.SerfConf

  var bindIP string
  var bindPort int
  var advertIP string
  var advertPort int
  var err error

  if ac.BindAddr != "" {
    bindIP, bindPort, err = ac.AddrParts(ac.BindAddr)
    if err != nil {
      return err
    }
  }

  if ac.AdvertiseAddr != "" {
    advertIP, advertPort, err = ac.AddrParts(ac.AdvertiseAddr)
    if err != nil {
      return err
    }
  }

  encryptKey, err := ac.EncryptBytes()
  if err != nil {
    return err
  }

  // https://github.com/hashicorp/serf/blob/master/cmd/serf/command/agent/command.go#L320
  sc.Tags = ac.Tags
  sc.NodeName = ac.NodeName
  sc.MemberlistConfig.BindAddr = bindIP
  sc.MemberlistConfig.BindPort = bindPort
  sc.MemberlistConfig.AdvertiseAddr = advertIP
  sc.MemberlistConfig.AdvertisePort = advertPort
  sc.MemberlistConfig.SecretKey = encryptKey
  sc.ProtocolVersion = uint8(ac.Protocol)
  sc.SnapshotPath = ac.SnapshotPath
  sc.MemberlistConfig.EnableCompression = ac.EnableCompression
  sc.QuerySizeLimit = ac.QuerySizeLimit
  sc.UserEventSizeLimit = ac.UserEventSizeLimit
  sc.EnableNameConflictResolution = !ac.DisableNameResolution
  sc.RejoinAfterLeave = ac.RejoinAfterLeave

  return nil
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
