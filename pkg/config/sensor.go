package config

type SensorConfig struct {
  Gpiochip      string
  HitPin        int
  LedPin        int
  Debounce      int
}

func NewSensorConfig(chip string, hitpin, ledpin, debouncetime int) *SensorConfig {
  return &SensorConfig{
    Gpiochip:     chip,
    HitPin:       hitpin,
    LedPin:       ledpin,
    Debounce:     debouncetime,
  }
}


