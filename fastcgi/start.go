package fastcgi

import (
  "context"
  "github.com/xpwu/go-tinyserver/api"
  "github.com/xpwu/go-xnet/xtcp"
  "log"
  "net/http"
  "net/http/fcgi"
)

func Start() {

  if !server.Net.Listen.On() {
    return
  }
  go runServer(server)
}

func runServer(s *serverConfig) {
  defer func() {
    if r := recover(); r != nil {
      log.Fatal(r)
    }
  }()

  serverMux := http.NewServeMux()

  for k,v := range api.AllHandlers() {
    serverMux.HandleFunc(k, v)
  }

  ln, err := xtcp.NetListen(&s.Net.Listen)
  if err != nil {
    panic(err)
  }
  defer func() {
    _ = ln.Close()
  }()

  if s.Net.TLS && s.Net.Listen.CanTLS() {
    ln, err = xtcp.NetListenTLS(ln, &s.Net.TlsFile)
    if err != nil {
      panic(err)
    }
    defer func() {
      _ = ln.Close()
    }()
  }

  ln, err = xtcp.NetListenConcurrent(context.Background(), ln, s.Net.MaxConnections)
  if err != nil {
    panic(err)
  }
  defer func() {
    _ = ln.Close()
  }()

  err = fcgi.Serve(ln, serverMux)
  if err != nil {
    panic(err)
  }
}
