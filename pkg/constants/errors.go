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
)


