package game

import (
  "os"
  "path/filepath"
  "encoding/json"
)

func (ge *GameEngine) Logfile() string {
  return filepath.Join(ge.Logdir, ge.GameStats.Uuid + ".json")
}

func (ge *GameEngine) Logstats() error {
  data, err := json.Marshal(ge)
  if err != nil {
    return err
  }

  return os.WriteFile(ge.Logfile(), data, os.ModePerm)
}


