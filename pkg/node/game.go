package node

import (
  "time"

  "github.com/hashicorp/serf/serf"
  "github.com/taemon1337/serf-cluster/pkg/game"
)

func (n *Node) AttachGameEngine(ge *game.GameEngine) {
  n.engine = ge
}

func (n *Node) RunGameEngine(expect, timeout int) error {
  return n.engine.Run(expect, timeout)
}

func (n *Node) ListenGameEngine() error {
  for {
    time.Sleep(5 * time.Second)
    if n.agent != nil && n.engine != nil {
      select {
        case q := <-n.engine.QueryChan:
          resp, err := n.agent.Query(q.Query, q.Payload, &serf.QueryParam{RequestAck: true, FilterTags: q.Tags})
          q.Response <- game.NewGameQueryResponse(resp, err)
      }
    } else {
      time.Sleep(5 * time.Second)
    }
  }
  return nil
}
