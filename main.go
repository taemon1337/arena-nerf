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

  flag.StringVar(&cfg.NodeName, "name", cfg.NodeName, "The name of this node in the cluster")

  n := node.NewNode(cfg)
  g := new(errgroup.Group)

  g.Go(func() error {
    return n.Start()
  })
  
  log.Fatal(g.Wait())
}
