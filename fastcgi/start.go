package fastcgi

import (
  "context"
  "github.com/xpwu/go-xnet/xtcp"
  "log"
  "net/http/fcgi"
)

func Start() {
  for _, s := range configValue.Servers {
    if !s.Net.Listen.On() {
      continue
    }
    go runServer(s)
  }
}

func runServer(s *server) {
  defer func() {
    if r := recover(); r != nil {
      log.Fatal(r)
    }
  }()

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

  err = fcgi.Serve(ln, nil)
  if err != nil {
    panic(err)
  }
}
