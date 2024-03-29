package game

import (
  "log"
  "fmt"
  "time"
  "strings"
  "slices"
  "math/rand"
  "encoding/json"
  "github.com/google/uuid"
  "github.com/taemon1337/arena-nerf/pkg/constants"
  "github.com/taemon1337/arena-nerf/pkg/config"
)

type GameController struct {
  Name            string
}

type GameStats struct {
  Uuid          string          `yaml:"id" json:"id"`
  Name          string          `yaml:"name" json:"name"`
  Mode          string          `yaml:"mode" json:"mode"`
  StartAt       time.Time       `yaml:"start_at" json:"start_at"`
  EndAt         time.Time       `yaml:"end_at" json:"end_at"`
  Length        string          `yaml:"length" json:"length"`
  Completed     bool            `yaml:"completed" json:"completed"`
  Status        string          `yaml:"status" json:"status"`
  Events        []string        `yaml:"events" json:"events"`
  Teams         []string        `yaml:"teams" json:"teams"`
  Nodes         []string        `yaml:"nodes" json:"nodes"`
  Nodeboard     map[string]int  `yaml:"nodeboard" json:"nodeboard"`
  Scoreboard    map[string]int  `yaml:"scoreboard" json:"scoreboard"`
  Winner        string          `yaml:"winner" json:"winner"`
  Highscore     int             `yaml:"highscore" json:"highscore"`
}

type GameEngine struct {
  Controller    *GameController `yaml:"-" json:"-"`
  EventChan     chan GameEvent  `yaml:"-" json:"-"`
  QueryChan     chan GameQuery  `yaml:"-" json:"-"`
  Logdir        string          `yaml:"-" json:"-"`
  GameStats     *GameStats      `yaml:"stats" json:"stats"`
}

func NewGameEngine(name, mode string, cfg *config.Config) *GameEngine {
  scoreboard := map[string]int{}
  nodeboard := map[string]int{}
  for _, team := range cfg.Teams {
    scoreboard[team] = 0
  }

  return &GameEngine{
    Controller:   &GameController{Name: name},
    EventChan:    make(chan GameEvent, 0),
    QueryChan:    make(chan GameQuery, 0),
    Logdir:       cfg.Logdir,
    GameStats:    &GameStats{
      Uuid:         uuid.New().String(),
      Name:         name,
      Mode:         mode,
      StartAt:      time.Time{},
      EndAt:        time.Time{},
      Length:       cfg.Gametime,
      Completed:    false,
      Status:       constants.GAME_STATE_INIT,
      Teams:        cfg.Teams,
      Nodes:        []string{},
      Nodeboard:    nodeboard,
      Scoreboard:   scoreboard,
      Winner:       "",
      Highscore:    0,
    },
  }
}

func (ge *GameEngine) Active() bool {
  return ge.GameStats.Completed
}

func (ge *GameEngine) SendAction(action, payload string) error {
  return ge.SendEvent(NewGameEvent(action, []byte(payload)))
}

func (ge *GameEngine) SendEvent(e GameEvent) error {
  ge.EventChan <- e
  ge.GameStats.Events = append(ge.GameStats.Events, e.String())
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
      log.Printf("nodes ready: %d", readycount)
      break // got expected amount node responses indicating readiness
    } else {
      log.Printf("waiting for %d ready nodes [%d/%d]...", expect, readycount, expect)
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
  if err = ge.SendEvent(NewGameEvent(constants.GAME_MODE, []byte(ge.GameStats.Mode))); err != nil {
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
    if !slices.Contains(ge.GameStats.Nodes, node) {
      ge.GameStats.Nodes = append(ge.GameStats.Nodes, node)
    }
    if string(val) != constants.GAME_MODE_DOMINATION {
      log.Printf("node %s game mode was set to %s, not %s", node, val, constants.GAME_MODE_DOMINATION)
    }
  }

  // start game
  if err := ge.SendEvent(NewGameEvent(constants.GAME_ACTION_BEGIN, []byte("Let the game begin!"))); err != nil {
    return err
  }

  ge.GameStats.StartAt = time.Now()
  ge.GameStats.Status = constants.GAME_STATE_ACTIVE

  gameduration := ge.ComputeGameDuration(ge.GameStats.StartAt, ge.GameStats.Length)

  for {
    select {
      // generate random hits (for testing only, should be configurable)
      case <-time.After(2 * time.Second):
        if err := ge.RandomTeamHits(); err != nil {
          log.Printf("error sending random team hit: %s", err)
        }

        scoreboard, nodeboard, err := ge.GetScoreboard()
        if err != nil {
          log.Printf("error querying scoreboard: %s", err)
          continue
        }

        ge.GameStats.Scoreboard = scoreboard
        ge.GameStats.Nodeboard = nodeboard

        if time.Now().After(ge.GameStats.StartAt.Add(gameduration)) {
          log.Printf("game time expired.")
          return ge.EndGame()
        }
    }
  }

  return nil
}

func (ge *GameEngine) EndGame() error {
  ge.GameStats.EndAt = time.Now()
  ge.GameStats.Status = constants.GAME_STATE_OVER
  ge.GameStats.Completed = true

  // inform nodes game ended
  if err := ge.SendEvent(NewGameEvent(constants.GAME_ACTION_END, []byte("The game is over!"))); err != nil {
    return err
  }

  // compute final score
  scoreboard, nodeboard, err := ge.GetScoreboard()
  if err != nil {
    log.Printf("error querying final scoreboard: %s", err)
    return err
  }

  ge.GameStats.Scoreboard = scoreboard
  ge.GameStats.Nodeboard = nodeboard
  log.Printf("Final Score: %s", scoreboard)

  for team, count := range scoreboard {
    if count > ge.GameStats.Highscore {
      ge.GameStats.Winner = team
      ge.GameStats.Highscore = count
    }
  }

  log.Printf("The winning team is %s with a score of %d", ge.GameStats.Winner, ge.GameStats.Highscore)

  if err := ge.SendEvent(NewGameEvent(constants.TEAM_WINNER, []byte(ge.GameStats.Winner))); err != nil {
    log.Printf("error sending team winner: %s", err)
    return err
  }

  if ge.Logdir != "" {
    err := ge.Logstats()
    if err != nil {
      log.Printf("could not write game log: %s", err)
    } else {
      log.Printf("saved game log to %s", ge.Logfile())
    }
  }

  return nil
}

func (ge *GameEngine) GetScoreboard() (map[string]int, map[string]int, error) {
  scoreboard := map[string]int{}
  nodeboard := map[string]int{}

  resp, err := ge.SendQuery(NewGameQuery(constants.TEAM_HIT, []byte(""), constants.NODE_TAGS))
  if err != nil {
    log.Printf("error querying team hit counts: %s", err)
    return scoreboard, nodeboard, err
  }

  // accumulate each node response
  for node, val := range resp {
    nodehits := map[string]int{}

    if err := json.Unmarshal(val, &nodehits); err != nil {
      log.Printf("cannot parse node hits: %s", err)
    } else {
      for key,count := range nodehits {
        isnode := slices.Contains(ge.GameStats.Nodes, key)
        isteam := slices.Contains(ge.GameStats.Teams, key)

        if isteam {
          if _, ok := scoreboard[key]; ok {
            scoreboard[key] += count
          } else {
            scoreboard[key] = count
          }
        }

        if isnode {
          if _, ok := nodeboard[key]; ok {
            nodeboard[key] += count
          } else {
            nodeboard[key] = count
          }
        }

        if !isteam && !isnode {
          msg := fmt.Sprintf("unrecognized team|node %s found in response from node %s", key, node)
          log.Printf(msg)
          ge.GameStats.Events = append(ge.GameStats.Events, msg)
        }
      }
    }
  }

  return scoreboard, nodeboard, nil
}


func (ge *GameEngine) ComputeGameDuration(start time.Time, ts string) time.Duration {
  dur, err := time.ParseDuration(ts)
  if err != nil {
    log.Printf("could not parse game length '%s', will use 5m - %s", ts, err)
    return time.Duration(5 * time.Minute)
  }

  return dur
}

func (ge *GameEngine) RandomTeamHit(hits int) error {
  node := ge.RandomNode()
  team := ge.RandomTeam()
  evt := strings.Join([]string{node, constants.TEAM_HIT}, constants.SPLIT)
  pay := fmt.Sprintf("%s%s%d", team, constants.SPLIT, hits)
  if err := ge.SendEvent(NewGameEvent(evt, []byte(pay))); err != nil {
    log.Printf("error sending test %s game event: %s", team, err)
  }

  return nil
}

func (ge *GameEngine) RandomTeamHits() error {
  for i := 1; i <= rand.Intn(10); i++ {
    if err := ge.RandomTeamHit(rand.Intn(5)); err != nil {
      return err
    }
  }
  return nil
}

func (ge *GameEngine) RandomTeam() string {
  return ge.GameStats.Teams[rand.Intn(len(ge.GameStats.Teams))]
}

func (ge *GameEngine) RandomNode() string {
  return ge.GameStats.Nodes[rand.Intn(len(ge.GameStats.Nodes))]
}
