package game

import (
  "log"
  "time"
  "encoding/json"
  "github.com/taemon1337/serf-cluster/pkg/constants"
)

type GameController struct {
  Name            string
}

type GameEngine struct {
  Controller    *GameController
  EventChan     chan GameEvent
  QueryChan     chan GameQuery
  Scoreboard    map[string]int
}

func NewGameEngine(name string) *GameEngine {
  return &GameEngine{
    Controller:   &GameController{Name: name},
    EventChan:    make(chan GameEvent, 0),
    QueryChan:    make(chan GameQuery, 0),
    Scoreboard:   map[string]int{},
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

  // send all teams to nodes
  log.Printf("setting game teams")
  if err = ge.SendEvent(NewGameEvent(constants.TEAM_ADD, []byte(""))); err != nil {
    return err
  }

  // query all nodes game mode
  resp, err := ge.SendQuery(NewGameQuery(constants.GAME_MODE, []byte(""), constants.NODE_TAGS))
  if err != nil {
    log.Printf("error querying game node: %s", err)
    return err
  }

  // check the game mode on each node was properly set
  for node, val := range resp {
    if string(val) != constants.GAME_MODE_DOMINATION {
      log.Printf("node %s game mode was set to %s, not %s", node, val, constants.GAME_MODE_DOMINATION)
    }
  }

  // start game
  if err := ge.SendEvent(NewGameEvent(constants.GAME_ACTION_BEGIN, []byte("Let the game begin!"))); err != nil {
    return err
  }

  for {
    if err := ge.SendEvent(NewGameEvent(constants.TEAM_HIT, []byte("blue:10"))); err != nil {
      return err
    }
    if err := ge.SendEvent(NewGameEvent(constants.TEAM_HIT, []byte("red:5"))); err != nil {
      return err
    }

    resp, err := ge.SendQuery(NewGameQuery(constants.TEAM_HIT, []byte(""), constants.NODE_TAGS))
    if err != nil {
      log.Printf("error querying team hit counts: %s", err)
      return err
    }

    totalhits := map[string]int{"node": 0}

    // check the game mode on each node was properly set
    for node, val := range resp {
      nodehits := map[string]int{}

      if err := json.Unmarshal(val, &nodehits); err != nil {
        log.Printf("cannot parse node hits: %s", err)
      } else {
        totalhits[node] = nodehits[constants.TAG_ROLE_NODE]
        for key,count := range nodehits {
          if _, ok := totalhits[key]; ok {
            totalhits[key] += count
          } else {
            totalhits[key] = count
          }
        }
      }
    }

    log.Printf("SCOREBOARD: %s", totalhits)
    ge.Scoreboard = totalhits

    time.Sleep(5 * time.Second)
  }

  if err := ge.SendEvent(NewGameEvent(constants.GAME_ACTION_END, []byte("The game is over!"))); err != nil {
    return err
  }

  return nil
}
