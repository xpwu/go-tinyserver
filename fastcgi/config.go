package fastcgi

import (
  config2 "github.com/xpwu/go-config/config"
  "github.com/xpwu/go-xnet/xtcp"
)

type config struct {
  Servers []*server
}

type server struct {
  Net *xtcp.Net
}

var configValue = &config{
  Servers:[]*server{
    {Net: xtcp.DefaultNetConfig()},
  },
}

func init() {
  config2.Unmarshal(configValue)
}

