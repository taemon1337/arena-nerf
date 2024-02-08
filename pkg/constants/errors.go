package constants

import (
  "errors"
)

var (
  ERR_EXISTING_CONNECTION = errors.New("already connected")
  ERR_NO_AGENT_CONFIG = errors.New("no node agent config")
  ERR_INVALID_CONFIG = errors.New("invalid node config")
  ERR_NOT_CONNECTED = errors.New("agent not connected")
  ERR_SHUTDOWN = errors.New("game shutdown called")
  ERR_API_ACTIONS_NOT_ALLOWED = errors.New("no api actions are allowed")
  ERR_UI_ACTION_NOT_ALLOWED = errors.New("invalid action - only certain ui actions are allowed")
  ERR_ONGOING_GAME = errors.New("invalid action - a game is ongoing, stop the game first")
)
