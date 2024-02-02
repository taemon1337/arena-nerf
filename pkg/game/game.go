package game

import (
  "log"
  "time"
  "github.com/hashicorp/serf/serf"
)

var (
  GAME_START = "game:start"
  GAME_STOP = "game:stop"
  GAME_MODE = "game:mode"
  HITS_TOTAL = "hits:total"
  NODE_WINNER = "node:winner"
  NODE_READY = "node:ready"
  NODE_TAGS = map[string]string{"role": "node"}
)

type GameEvent struct {
  EventName       string        `yaml:"event" json:"event"`
  Payload         []byte        `yaml:"payload" json:"payload"`
}

type GameQueryResponse struct {
  Answer          *serf.QueryResponse
  Error           error
}

type GameQuery struct {
  Query           string                  `yaml:"query" json:"query"`
  Payload         []byte                  `yaml:"payload" json:"payload"`
  Tags            map[string]string       `yaml:"tags" json:"tags"`
  Response        chan GameQueryResponse  `yaml:"-" json:"-"`
}

func NewGameEvent(event string, payload []byte) GameEvent {
  return GameEvent{
    EventName: event,
    Payload:   payload,
  }
}

func NewGameQuery(query string, payload []byte, tags map[string]string) GameQuery {
  return GameQuery{
    Query:     query,
    Payload:   payload,
    Tags:      tags,
    Response:  make(chan GameQueryResponse, 0),
  }
}

func NewGameQueryResponse(resp *serf.QueryResponse, err error) GameQueryResponse {
  return GameQueryResponse{
    Answer:   resp,
    Error:    err,
  }
}

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

func (ge *GameEngine) Start() error {
  return ge.SendEvent(NewGameEvent(GAME_START, []byte("Let the game begin!")))
}

func (ge *GameEngine) Stop() error {
  return ge.SendEvent(NewGameEvent(GAME_STOP, []byte("Stop the game!")))
}

func (ge *GameEngine) SendEvent(e GameEvent) error {
  ge.EventChan <- e
  return nil
}

func (ge *GameEngine) SendQuery(q GameQuery) ([]interface{}, error) {
  ge.QueryChan <- q
  return nil, nil
}

func (ge *GameEngine) WaitForNodes(expect int, timeout int) error {
  for {
    // wait for ready
    resp, err := ge.SendQuery(NewGameQuery(NODE_READY, []byte(""), NODE_TAGS))
    if err != nil {
      log.Printf("error query readiness of nodes: %s", err)
      return err
    }

    if len(resp) >= expect {
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
  err := ge.WaitForNodes(expect, timeout)
  if err != nil {
    return err
  }

  // set game mode
  log.Printf("setting game mode")
  if err = ge.SendEvent(NewGameEvent(GAME_MODE, []byte("domination"))); err != nil {
    return err
  }

  // start game
  if err := ge.Start(); err != nil {
    return err
  }

//  scoreboard := map[string]string{} // node: winner

  for {
    time.Sleep(10 * time.Second)
  }
  return nil
}
