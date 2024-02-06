package sensor

import (
  "log"
  "time"
  "github.com/taemon1337/serf-cluster/pkg/game"
  "github.com/stianeikeland/go-rpio/v4"
)

var (
  PIN10 = rpio.Pin(10)
)

type Sensor struct {
  id        string
  emit      chan game.GameEvent
  listen    chan game.GameEvent
}

func NewSensor(id string) *Sensor {
  return &Sensor{
    id:       id,
    emit:     make(chan game.GameEvent),
    listen:   make(chan game.GameEvent),
  }
}

func (s *Sensor) Start() error {
  if err := rpio.Open(); err != nil {
    log.Printf("error opening rpio: %s", err)
    return err
  }

  defer rpio.Close()

  PIN10.Output()

  for {
    PIN10.Toggle()
    time.Sleep(10 * time.Second)
  }

  return nil
}
