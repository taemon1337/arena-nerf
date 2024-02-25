package sensor

import (
  "log"
  "sync"
  "time"
  "strconv"
  "strings"
  "github.com/taemon1337/arena-nerf/pkg/game"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
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
  led       *gpiod.Line
  hit       *gpiod.Line
  Gamechan  chan game.GameEvent
  hitchan   chan gpiod.LineEvent
  hitlock   sync.Mutex
  hittime   time.Time
}

func NewSensor(id string, cfg *config.SensorConfig) *Sensor {
  return &Sensor{
    id:       id,
    conf:     cfg,
    led:      nil,
    hit:      nil,
    Gamechan: make(chan game.GameEvent, 2),
    hitchan:  make(chan gpiod.LineEvent, 6),
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
  log.Printf("HIT DEBOUNCE REACHED")
  s.Gamechan <- game.NewGameEvent(constants.TARGET_HIT, []byte("blue:1")) // blue should be replaced with current LED target color
}

func (s *Sensor) Start() error {
  ledpin, err := s.ParseGpioPin(s.conf.LedPin)
  if err != nil {
    return err
  }

  hitpin, err := s.ParseGpioPin(s.conf.HitPin)
  if err != nil {
    return err
  }

  log.Printf("Sensor Hit pin: %d, LED pin: %d", hitpin, ledpin)

  // event channel buffer
  eh := func(evt gpiod.LineEvent) {
    select {
    case s.hitchan <- evt:
    default:
      log.Printf("event chan overflow - discarding event")
    }
  }

  led, err := gpiod.RequestLine(s.conf.Gpiochip, ledpin, gpiod.AsOutput(OFF))
  if err != nil {
    log.Printf("cannot request gpiod %d led line: %s", ledpin, err)
    return err
  }

  hit, err := gpiod.RequestLine(s.conf.Gpiochip, hitpin, gpiod.WithPullUp, gpiod.WithRisingEdge, gpiod.WithEventHandler(eh))
  if err != nil {
    log.Printf("cannot request gpiod %d hit line: %s", hitpin, err)
    return err
  }

  s.led = led
  s.hit = hit

  // start by blinking led
  log.Printf("Blinking LED 5 times...")
  time.Sleep(3 * time.Second)
  s.Blink(5)
  time.Sleep(1 * time.Second)

  defer func() {
    led.Reconfigure(gpiod.AsInput)
    hit.Reconfigure(gpiod.AsInput)
    led.Close()
    hit.Close()
  }()

  done := false
  for !done {
    select {
    case evt := <-s.hitchan:
      log.Printf("HIT: %s", evt)
      s.ProcessEvent(evt)
    }
  }
  return constants.ERR_SENSOR_STOPPED
}

func (s *Sensor) Listen() error {
  for {
    select {
      case e := <-s.Gamechan:
        log.Printf("SENSOR GAME EVENT RECEIVED: %s", e)
        switch e.EventName {
          case constants.TARGET_HIT:

          case constants.TEAM_HIT:
            parts := strings.Split(string(e.Payload), constants.SPLIT)
            if len(parts) < 2 {
              log.Printf("cannot parse sensor team hit from %s - should be <team>:<count>", string(e.Payload))
            } else {
              hits, err := strconv.Atoi(parts[1])
              if err != nil {
                log.Printf("cannot parse sensor team hit from %s - %s", string(e.Payload), err)
                continue
              } else {
                s.Blink(hits) // blink the number of hits times
              }
            }
          default:
            log.Printf("Sensor Received Event (no action found): %s", e)
        }
    }
  }
  return constants.ERR_SENSOR_STOPPED
}

// game event for this node read on game event channel
func (s *Sensor) NodeTeamHit(action string, payload []byte) {
  s.Gamechan <- game.NewGameEvent(action, payload)
}

func (s *Sensor) Blink(times int) {
  for i := 0; i < times; i++ {
    s.BlinkOnce()
    time.Sleep(BLINK_DELAY)
  }
}

func (s *Sensor) BlinkOnce() {
  if s.led != nil {
    s.led.SetValue(ON)
    time.Sleep(BLINK_DELAY)
    s.led.SetValue(OFF)
  } else {
    log.Printf("warn: led not setup - is nil")
  }
}
