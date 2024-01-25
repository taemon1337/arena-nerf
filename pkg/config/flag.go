package config

import (
  "strings"
)

type JoinAddrList []string

func (i *JoinAddrList) String() string {
  return strings.Join(*i, ",")
}

func (i *JoinAddrList) Set(val string) error {
  *i = append(*i, val)
  return nil
}
