package fastcgi

import (
  "github.com/xpwu/go-config/configs"
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
  configs.Unmarshal(configValue)
}

