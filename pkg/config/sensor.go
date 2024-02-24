package config

type SensorConfig struct {
  Device        string
  Gpiochip      string
  HitPin        string
  LedPin        string
  Debounce      int
}

func NewSensorConfig(device, chip, hitpin, ledpin string, debouncetime int) *SensorConfig {
  return &SensorConfig{
    Device:       device,
    Gpiochip:     chip,
    HitPin:       hitpin,
    LedPin:       ledpin,
    Debounce:     debouncetime,
  }
}

func DefaultSensorConfig() *SensorConfig {
  return &SensorConfig{
    Device:       "",
    Gpiochip:     "gpiochip0",
    HitPin:       "",
    LedPin:       "",
    Debounce:     100,
  }
}
