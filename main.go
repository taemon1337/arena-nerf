package main

import (
  "log"
  "flag"

  "golang.org/x/sync/errgroup"

  "github.com/taemon1337/serf-cluster/pkg/node"
  "github.com/taemon1337/serf-cluster/pkg/config"
)

func main() {
  cfg := config.NewConfig()

  flag.StringVar(&cfg.AgentConf.NodeName, "name", cfg.AgentConf.NodeName, "name of this node in the cluster")
  flag.StringVar(&cfg.AgentConf.BindAddr, "bind", cfg.AgentConf.BindAddr, "address to bind listeners to")
  flag.StringVar(&cfg.AgentConf.AdvertiseAddr, "advertise", cfg.AgentConf.AdvertiseAddr, "address to advertise to cluster")
  flag.StringVar(&cfg.AgentConf.EncryptKey, "encrypt", cfg.AgentConf.EncryptKey, "encryption key")
  flag.Var(&cfg.JoinAddrs, "join", "addresses to try to join automatically and repeatable until success")
  flag.Parse()

  n := node.NewNode(cfg)
  g := new(errgroup.Group)

  g.Go(func() error {
    return n.Start()
  })

  g.Go(func() error {
    return n.AutoJoin()
  })

  log.Fatal(g.Wait())
}
