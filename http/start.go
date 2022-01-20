package http

import (
  "context"
  "github.com/xpwu/go-log/log"
  "github.com/xpwu/go-log/log/level"
  "github.com/xpwu/go-xnet/xhttp"
  "net"
  "net/http"
  "regexp"
  "strings"
  "time"
)

func Start() {
  for _, s := range configValue.Servers {
    if !s.Net.Listen.On() {
      continue
    }
    go runServer(s)
  }
}

func escapeReg(str string) string {
  str = strings.Replace(str, ".", "\\.", -1)
  str = strings.Replace(str, "?", "\\?", -1)
  str = "^" + strings.Replace(str, "*", ".*", -1) + "$"

  return str
}

func runServer(s *server) {
  defer func() {
    if r := recover(); r != nil {
      log.Fatal(r)
    }
  }()

  log.Info("server(http) listen " + s.Net.Listen.LogString())

  ctx, logger := log.WithCtx(context.Background())

  s.nameReg = make([]*regexp.Regexp, len(s.HostName))
  for i, name := range s.HostName {
    s.nameReg[i] = regexp.MustCompile(escapeReg(name))
  }

  srv := &http.Server{
    Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      host := stripHostPort(r.Host)

      found := false

      for _, reg := range s.nameReg {
        if host == reg.FindString(host) {
          found = true
          break
        }
      }
      if !found {
        goto notFound
      }

      if s.RootUri != "" {
        if !strings.HasPrefix(r.URL.Path, s.RootUri) {
          found = false
          goto notFound
        }
        r.URL.Path = strings.TrimPrefix(r.URL.Path, s.RootUri)
      }

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
