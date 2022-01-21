package http

import (
  "context"
  "github.com/xpwu/go-log/log"
  "github.com/xpwu/go-log/log/level"
  "github.com/xpwu/go-tinyserver/api"
  "github.com/xpwu/go-xnet/xhttp"
  "net"
  "net/http"
  "path"
  "strings"
  "time"
)

func Start() {
  if !server.Net.Listen.On() {
    return
  }

  go runServer(server)
}

func escapeReg(str string) string {
  str = strings.Replace(str, ".", "\\.", -1)
  str = strings.Replace(str, "?", "\\?", -1)
  str = "^" + strings.Replace(str, "*", ".*", -1) + "$"

  return str
}

func runServer(s *serverConfig) {
  defer func() {
    if r := recover(); r != nil {
      log.Fatal(r)
    }
  }()

  log.Info("server(http) listen " + s.Net.Listen.LogString())

  ctx, logger := log.WithCtx(context.Background())

  for k,v := range api.AllHandlers() {
    http.HandleFunc(path.Join("/", s.RootUri, k), v)
  }

  srv := &http.Server{
    Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      host := stripHostPort(r.Host)

      found := false

      for _, h := range s.HostName {
        if strings.HasPrefix(host, h) {
          found = true
          break
        }
      }
      if len(s.HostName) == 0 {
        found = true
      }

      if !found {
        goto notFound
      }

      // 不把RootUri看着服务级别的权限控制，视为location的匹配
      //if s.RootUri != "" {
      //  p := r.URL.Path
      //  if !path.IsAbs(p) {
      //    p = "/" + p
      //  }
      //  if !strings.HasPrefix(p, s.RootUri) {
      //    found = false
      //    goto notFound
      //  }
      //}

    notFound:
      if !found {
        http.NotFound(w, r)
        return
      }

      http.DefaultServeMux.ServeHTTP(w, r)

    }),

    Addr:        s.Net.Listen.String(),
    ErrorLog:    log.NewSysLog(logger, level.ERROR),
    IdleTimeout: 30 * time.Second,
  }

  err := xhttp.SeverAndBlock(ctx, srv, s.Net)

  if err != nil {
    panic(err)
  }
}

// copy from net/http/server.go:2235
// stripHostPort returns h without any trailing ":<port>".
func stripHostPort(h string) string {
  // If no port on host, return unchanged
  if strings.IndexByte(h, ':') == -1 {
    return h
  }
  host, _, err := net.SplitHostPort(h)
  if err != nil {
    return h // on error, return unchanged
  }
  return host
}

