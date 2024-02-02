package game

import (
  "log"
  "time"
  "github.com/taemon1337/serf-cluster/pkg/constants"
)

type GameController struct {
  Name            string
}

type GameEngine struct {
  Controller    *GameController
  EventChan     chan GameEvent
  QueryChan     chan GameQuery
}

func NewGameEngine(name string) *GameEngine {
  return &GameEngine{
    Controller:   &GameController{Name: name},
    EventChan:    make(chan GameEvent, 0),
    QueryChan:    make(chan GameQuery, 0),
  }
}

func (ge *GameEngine) SendEvent(e GameEvent) error {
  ge.EventChan <- e
  return nil
}

func (ge *GameEngine) SendQuery(q GameQuery) (map[string][]byte, error) {
  ge.QueryChan <- q
  resp := <-q.Response // block for response
  return resp.Answer, resp.Error
}

func (ge *GameEngine) WaitForNodes(expect int, timeout int) error {
  for {
    // wait for ready
    resp, err := ge.SendQuery(NewGameQuery(constants.NODE_READY, []byte(""), constants.NODE_TAGS))
    if err != nil {
      log.Printf("error query readiness of nodes: %s", err)
      return err
    }

    readycount := 0

    for _, val := range resp {
      if string(val) == constants.NODE_IS_READY {
        readycount += 1
      }
    }

    if readycount >= expect {
      log.Printf("nodes ready: %d", len(resp))
      break // got expected amount node responses indicating readiness
    } else {
      log.Printf("waiting for %d ready nodes...", expect)
      time.Sleep(time.Duration(timeout))
    }
  }
  return nil
}

// node/station cababilities
// node:query:mode - do you have this game mode?
// node:get:mode - get the current game mode of the node
// node:point:count - how many points on this node
// team:point:count - a team points count
// node:winner - the team with the most points

func (ge *GameEngine) Run(expect, timeout int) error {

  // wait for nodes to be ready
  log.Printf("waiting for %d nodes", expect)
  err := ge.WaitForNodes(expect, timeout)
  if err != nil {
    return err
  }

  // set game mode
  log.Printf("setting game mode")
  if err = ge.SendEvent(NewGameEvent(constants.GAME_MODE, []byte(constants.GAME_MODE_DOMINATION))); err != nil {
    return err
  }

  resp, err := ge.SendQuery(NewGameQuery(constants.GAME_MODE, []byte(""), constants.NODE_TAGS))
  if err != nil {
    log.Printf("error querying game node: %s", err)
    return err
  }

  for node, val := range resp {
    if string(val) != constants.GAME_MODE_DOMINATION {
      log.Printf("node %s game mode was set to %s, not %s", node, val, constants.GAME_MODE_DOMINATION)
    }
  }

  // start game
  if err := ge.SendEvent(NewGameEvent(constants.GAME_ACTION_BEGIN, []byte("Let the game begin!"))); err != nil {
    return err
  }

//  scoreboard := map[string]string{} // node: winner

  for {
    time.Sleep(10 * time.Second)
  }

  if err := ge.SendEvent(NewGameEvent(constants.GAME_ACTION_END, []byte("The game is over!"))); err != nil {
    return nil
  }

  return nil
}
