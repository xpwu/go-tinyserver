package http

import (
  config2 "github.com/xpwu/go-config/config"
  "github.com/xpwu/go-xnet/xtcp"
  "regexp"
)

type config struct {
  Servers []*server
}

type server struct {
  Net      *xtcp.Net
  HostName []string
  RootUri  string
  nameReg  []*regexp.Regexp
}

var configValue = &config{
  Servers: []*server{
    {
      Net:      xtcp.DefaultNetConfig(),
      HostName: []string{"*"},
    },
  },
}

func init() {
  config2.Unmarshal(configValue)
}
