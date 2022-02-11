package http

import (
  "github.com/xpwu/go-config/configs"
  "github.com/xpwu/go-xnet/xtcp"
)

type serverConfig struct {
  Net      *xtcp.Net
  HostName []string `conf:",leftmost match, []: allow all host name"`
  RootUri  string `conf:",match_uri = RootUri + api.RegisterUri"`
}

var server = &serverConfig{
  Net:      xtcp.DefaultNetConfig(),
  HostName: []string{},
  RootUri: "/",
}

func init() {
  configs.Unmarshal(server)
}
