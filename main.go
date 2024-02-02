package main

import (
  "log"
  "fmt"
  "flag"

  "github.com/taemon1337/serf-cluster/pkg/node"
  "github.com/taemon1337/serf-cluster/pkg/game"
  "github.com/taemon1337/serf-cluster/pkg/config"
  "github.com/taemon1337/serf-cluster/pkg/constants"
  "github.com/taemon1337/serf-cluster/pkg/controller"
)

func main() {
  cfg := config.NewConfig(constants.TAG_ROLE_NODE)
  var tags []string
  var role string
  var mode string

  flag.StringVar(&role, "role", cfg.AgentConf.Tags["role"], fmt.Sprintf("set node role to %s or %s", constants.TAG_ROLE_NODE, constants.TAG_ROLE_CTRL))
  flag.StringVar(&cfg.AgentConf.NodeName, "name", cfg.AgentConf.NodeName, "name of this node in the cluster")
  flag.StringVar(&cfg.AgentConf.BindAddr, "bind", cfg.AgentConf.BindAddr, "address to bind listeners to")
  flag.StringVar(&cfg.AgentConf.AdvertiseAddr, "advertise", cfg.AgentConf.AdvertiseAddr, "address to advertise to cluster")
  flag.StringVar(&cfg.AgentConf.EncryptKey, "encrypt", cfg.AgentConf.EncryptKey, "encryption key")
  flag.Var((*config.AppendSliceValue)(&tags), "tag", "add tag to node with key=value")
  flag.Var((*config.AppendSliceValue)(&cfg.JoinAddrs), "join", "addresses to try to join automatically and repeatable until success")
  flag.StringVar(&mode, "mode", "", "set to the desired game mode to run a game from this node (must be a control node)")
  flag.IntVar(&cfg.ExpectNodes, "expect", cfg.ExpectNodes, "set to the expected number of game nodes (not including control node) to wait for before starting the game")
  flag.IntVar(&cfg.Timeout, "wait", cfg.Timeout, "set the default number of seconds to timeout and/or wait for nodes")
  flag.Parse()

  parsedtags, err := config.UnmarshalTags(tags)
  if err != nil {
    log.Fatal(err)
  }

  cfg.AgentConf.Tags = parsedtags
  if role != "" {
    cfg.AgentConf.Tags["role"] = role
  } else {
    role = cfg.AgentConf.Tags["role"]
  }

  log.Printf("Node: %s", cfg.AgentConf.NodeName)
  log.Printf("Role: %s", role)
  log.Printf("Join: %s", cfg.JoinAddrs)

  if role == constants.TAG_ROLE_CTRL {
    if mode == "" {
      log.Fatal("Must set game mode for a control node.")
    }
    log.Fatal(controller.NewController(cfg, game.NewGameEngine(mode)).Start())
  } else if role == constants.TAG_ROLE_NODE {
    log.Fatal(node.NewNode(cfg).Start())
  } else {
    log.Fatal("invalid role")
  }
}
