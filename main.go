package main

import (
  "log"
  "fmt"
  "time"
  "flag"

  "golang.org/x/sync/errgroup"
  "github.com/taemon1337/serf-cluster/pkg/game"
  "github.com/taemon1337/serf-cluster/pkg/node"
  "github.com/taemon1337/serf-cluster/pkg/config"
  "github.com/taemon1337/serf-cluster/pkg/constants"
  "github.com/taemon1337/serf-cluster/pkg/controller"
)

func main() {
  cfg := config.NewConfig(constants.TAG_ROLE_NODE)
  var tags []string
  var role string
  var mode string
  var startgame bool

  flag.StringVar(&role, "role", cfg.AgentConf.Tags["role"], fmt.Sprintf("set node role to %s or %s", constants.TAG_ROLE_NODE, constants.TAG_ROLE_CTRL))
  flag.StringVar(&cfg.AgentConf.NodeName, "name", cfg.AgentConf.NodeName, "name of this node in the cluster")
  flag.StringVar(&cfg.AgentConf.BindAddr, "bind", cfg.AgentConf.BindAddr, "address to bind listeners to")
  flag.StringVar(&cfg.AgentConf.AdvertiseAddr, "advertise", cfg.AgentConf.AdvertiseAddr, "address to advertise to cluster")
  flag.StringVar(&cfg.AgentConf.EncryptKey, "encrypt", cfg.AgentConf.EncryptKey, "encryption key")
  flag.Var((*config.AppendSliceValue)(&tags), "tag", "add tag to node with key=value")
  flag.Var((*config.AppendSliceValue)(&cfg.JoinAddrs), "join", "addresses to try to join automatically and repeatable until success")
  flag.StringVar(&mode, "mode", "", "set to the desired game mode to run a game from this node (must be a control node)")
  flag.StringVar(&cfg.Gametime, "gametime", cfg.Gametime, "set to the length of the game (assuming no winner)")
  flag.StringVar(&cfg.Logdir, "logdir", cfg.Logdir, "directory to write game log files to")
  flag.BoolVar(&cfg.Webserver, "server", cfg.Webserver, "set to true to start the web server (only on control node)")
  flag.BoolVar(&cfg.Sensor, "sensor", cfg.Sensor, "set to true when this node is a Raspberry Pi with sensors to start sensor functions")
  flag.BoolVar(&startgame, "start", false, "set to true to immediate start game (assuming control node and game mode set)")
  flag.BoolVar(&cfg.AllowApiActions, "allow-api-actions", cfg.AllowApiActions, "set to true to allow API actions (only for control webserver)")
  flag.StringVar(&cfg.WebAddr, "addr", cfg.WebAddr, "the web server address to listen on")
  flag.Var((*config.AppendSliceValue)(&cfg.Teams), "team", "register/add a team name")
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
  log.Printf("Server: %s", cfg.Webserver)
  log.Printf("Logdir: %s", cfg.Logdir)
  log.Printf("Sensor: %s", cfg.Sensor)

  g := new(errgroup.Group)

  switch role {
    case constants.TAG_ROLE_CTRL:
      if mode == "" {
        log.Fatal("Must set game mode for a control node.")
      }

      ctrl := controller.NewController(cfg)

      g.Go(func () error {
        return ctrl.Start()
      })

      if startgame {
        time.Sleep(7 * time.Second)
        ctrl.RunGame(game.NewGameEngine(mode, mode, cfg))
      }
    case constants.TAG_ROLE_NODE:
      g.Go(func () error {
        return node.NewNode(cfg).Start()
      })
    default:
      log.Fatal("invalid role")
  }
  
  log.Fatal(g.Wait())
}
