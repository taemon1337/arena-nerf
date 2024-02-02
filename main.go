package main

import (
  "log"
  "fmt"
  "flag"

  "golang.org/x/sync/errgroup"

  "github.com/taemon1337/serf-cluster/pkg/node"
  "github.com/taemon1337/serf-cluster/pkg/game"
  "github.com/taemon1337/serf-cluster/pkg/config"
)

func main() {
  cfg := config.NewConfig(node.TAG_ROLE_NODE)
  var tags []string
  var role string
  var mode string
  var expect int
  var timeout int

  flag.StringVar(&role, "role", cfg.AgentConf.Tags["role"], fmt.Sprintf("set node role to %s or %s", node.TAG_ROLE_NODE, node.TAG_ROLE_CTRL))
  flag.StringVar(&cfg.AgentConf.NodeName, "name", cfg.AgentConf.NodeName, "name of this node in the cluster")
  flag.StringVar(&cfg.AgentConf.BindAddr, "bind", cfg.AgentConf.BindAddr, "address to bind listeners to")
  flag.StringVar(&cfg.AgentConf.AdvertiseAddr, "advertise", cfg.AgentConf.AdvertiseAddr, "address to advertise to cluster")
  flag.StringVar(&cfg.AgentConf.EncryptKey, "encrypt", cfg.AgentConf.EncryptKey, "encryption key")
  flag.Var((*config.AppendSliceValue)(&tags), "tag", "add tag to node with key=value")
  flag.Var((*config.AppendSliceValue)(&cfg.JoinAddrs), "join", "addresses to try to join automatically and repeatable until success")
  flag.StringVar(&mode, "mode", "", "set to the desired game mode to run a game from this node (must be a control node)")
  flag.IntVar(&expect, "expect", 1, "set to the expected number of game nodes (not including control node) to wait for before starting the game")
  flag.IntVar(&timeout, "wait", 5, "set the default number of seconds to timeout and/or wait for nodes")
  flag.Parse()

  parsedtags, err := config.UnmarshalTags(tags)
  if err != nil {
    log.Fatal(err)
  }

  cfg.AgentConf.Tags = parsedtags
  cfg.AgentConf.Tags["role"] = role

  n, err := node.NewNode(cfg)
  if err != nil {
    log.Fatal(err)
  }

  g := new(errgroup.Group)

  g.Go(func() error {
    return n.Start()
  })

  g.Go(func() error {
    return n.AutoJoin()
  })

  if mode != "" {
    n.AttachGameEngine(game.NewGameEngine(mode))

    g.Go(func() error {
      return n.ListenGameEngine()
    })

    g.Go(func() error {
      return n.RunGameEngine(expect, timeout)
    })
  }

  log.Fatal(g.Wait())
}
