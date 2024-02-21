package sensor

import (
  "os"
  "log"
  "sync"
  "time"
  "syscall"
  "os/signal"
  "github.com/taemon1337/arena-nerf/pkg/game"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/gpiod"
)

var (
  OFF int = 0
  ON int = 1
  BLINK_DELAY = 100 * time.Millisecond
)

type Sensor struct {
  id        string
  conf      *config.SensorConfig
  emit      chan game.GameEvent
  listen    chan game.GameEvent
  hitchan   chan gpiod.LineEvent
  hitlock   sync.Mutex
  hittime   time.Time
}

func NewSensor(id string, cfg *config.SensorConfig) *Sensor {
  return &Sensor{
    id:       id,
    conf:     cfg,
    emit:     make(chan game.GameEvent),
    listen:   make(chan game.GameEvent),
    hitchan:  make(chan gpiod.LineEvent),
    hitlock:  sync.Mutex{},
    hittime:  time.Now(),
  }
}

func (s *Sensor) ProcessEvent(evt gpiod.LineEvent) {
  s.hitlock.Lock()
  defer s.hitlock.Unlock()

  debounce_duration := time.Duration(s.conf.Debounce) * time.Millisecond

  if time.Since(s.hittime) < debounce_duration {
    return // ignore since within debounce window
  }

  s.hittime = time.Now() // hittime is last debounced hit time
}

func (s *Sensor) Start() error {
  echan := make(chan gpiod.LineEvent, 6)

  // event channel buffer
  eh := func(evt gpiod.LineEvent) {
    select {
    case echan <- evt:
    default:
      log.Printf("event chan overflow - discarding event")
    }
  }

  led, err := gpiod.RequestLine(s.conf.Gpiochip, s.conf.LedPin, gpiod.AsOutput(OFF))
  if err != nil {
    return err
  }

  hit, err := gpiod.RequestLine(s.conf.Gpiochip, s.conf.HitPin, gpiod.WithPullUp, gpiod.WithRisingEdge, gpiod.WithEventHandler(eh))
  if err != nil {
    return err
  }

  // start by blinking led
  s.Blink(5, led)
  time.Sleep(2 * time.Second)

  defer func() {
    led.Reconfigure(gpiod.AsInput)
    hit.Reconfigure(gpiod.AsInput)
    led.Close()
    hit.Close()
  }()

  quit := make(chan os.Signal, 1)
  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
  defer signal.Stop(quit)

  done := false
  for !done {
    select {
    case evt := <-echan:
      go s.ProcessEvent(evt)
    case <-quit:
      log.Printf("stopping...")
      done = true
    }
  }
  return nil
}

func (s *Sensor) Blink(times int, led *gpiod.Line) {
  for i := 0; i < times; i++ {
    s.BlinkOnce(led)
    time.Sleep(BLINK_DELAY)
  }
}

func (s *Sensor) BlinkOnce(led *gpiod.Line) {
  led.SetValue(ON)
  time.Sleep(BLINK_DELAY)
  led.SetValue(OFF)
}
