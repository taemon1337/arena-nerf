package constants

import (
  "time"
)

var (
  SPLIT = ":"
  TAG_ROLE_NODE = "node"
  TAG_ROLE_CTRL = "ctrl"
  EVENT_ALIVE = "alive"
  COALESCE = false
  WAIT_TIME = 10 * time.Second

  // game mode
  GAME_MODE = "game:mode"
  GAME_MODE_DOMINATION = "domination"
  GAME_MODE_NONE = "none"

  // game actions
  GAME_ACTION_BEGIN = "game:begin"
  GAME_ACTION_END = "game:end"
  GAME_ACTION_RESET = "game:reset"
  GAME_ACTION_HIT = "game:hit"

  // game states
  GAME_STATE_INIT = "game:init"
  GAME_STATE_ACTIVE = "game:active"
  GAME_STATE_OVER = "game:over"

  // team
  TEAM_ADD = "team:add"
  TEAM_DEL = "team:del"

  NODE_HIT = "node:hit"
  TEAM_HIT = "team:hit"

  NODE_WINNER = "node:winner"
  NODE_READY = "node:ready"
  NODE_IS_READY = "true"
  NODE_IS_NOT_READY = "false"
  NODE_TAGS = map[string]string{"role": "node"}
)
