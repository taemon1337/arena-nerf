package node

import (
  "log"
  "time"
  "github.com/hashicorp/serf/serf"
)

func (n *Node) Active() bool {
  if n.agent != nil {
    srf := n.agent.Serf()
    if srf != nil {
      return srf.State() == serf.SerfAlive
    }
  }
  return false
}

func (n *Node) AutoJoin() error {
  for {
    if n.Active() {
      i, err := n.agent.Join(n.conf.JoinAddrs, n.conf.JoinReplay)
      if err != nil {
        log.Printf("error joining %s: %s", n.conf.JoinAddrs, err)
      }

      if i > 0 {
        log.Printf("successfully joined %d nodes", i)
        return nil
      }
    }
    time.Sleep(WAIT_TIME)
  }
}
