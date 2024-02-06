package game

import (
  "fmt"
)

type GameEvent struct {
  EventName       string        `yaml:"event" json:"event"`
  Payload         []byte        `yaml:"payload" json:"payload"`
}

type GameQueryResponse struct {
  Answer          map[string][]byte
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

func (e GameEvent) String() string {
  return fmt.Sprintf("%s: %s", e.EventName, string(e.Payload))
}

func NewGameQuery(query string, payload []byte, tags map[string]string) GameQuery {
  return GameQuery{
    Query:     query,
    Payload:   payload,
    Tags:      tags,
    Response:  make(chan GameQueryResponse, 0),
  }
}

func NewGameQueryResponse(resp map[string][]byte, err error) GameQueryResponse {
  return GameQueryResponse{
    Answer:   resp,
    Error:    err,
  }
}


