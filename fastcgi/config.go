package fastcgi

import (
  "github.com/xpwu/go-config/configs"
  "github.com/xpwu/go-xnet/xtcp"
)


type serverConfig struct {
  *xtcp.Net
}

var server = &serverConfig{
  xtcp.DefaultNetConfig(),
}

func init() {
  configs.Unmarshal(server)
}

