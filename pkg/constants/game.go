package constants

import (
  "time"
)

var (
  TAG_ROLE_NODE = "node"
  TAG_ROLE_CTRL = "ctrl"
  EVENT_ALIVE = "alive"
  COALESCE = false
  WAIT_TIME = 10 * time.Second
  GAME_START = "game:start"
  GAME_STOP = "game:stop"
  GAME_MODE = "game:mode"
  GAME_MODE_DOMINATION = "domination"
  GAME_MODE_NONE = "none"
  // game actions
  GAME_ACTION_BEGIN = "game:begin"
  GAME_ACTION_END = "game:end"

  // game states
  GAME_STATE_INIT = "game:init"
  GAME_STATE_ACTIVE = "game:active"
  GAME_STATE_OVER = "game:over"

  HITS_TOTAL = "hits:total"
  NODE_WINNER = "node:winner"
  NODE_READY = "node:ready"
  NODE_IS_READY = "true"
  NODE_IS_NOT_READY = "false"
  NODE_TAGS = map[string]string{"role": "node"}
)
