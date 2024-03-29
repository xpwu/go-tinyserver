package http

import (
  "context"
  "fmt"
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

/**
  必须先通过Add方法或者Register方法注册API，再启动服务
*/

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

  serverMux := http.NewServeMux()

  for k, v := range api.AllHandlers() {
    serverMux.HandleFunc(path.Join("/", s.RootUri, k), v)
  }
  // 404。 HandleFunc 不能对同一个pattern多次注册，所以先判断再注册
  // 只是s.RootUri 后的url不匹配
  rootUri := path.Join("/", s.RootUri, "")
  if _, exist := api.AllHandlers()[rootUri]; !exist {
    serverMux.HandleFunc(rootUri, _404)
  }
  // s.RootUri 就不能匹配的情况，同时要注意 s.RootUri 为空或者 / 的情况，此种情况就是前面的_404情况，不能再注册
  if _, exist := api.AllHandlers()["/"]; rootUri != "/" && !exist {
    serverMux.HandleFunc("/", rootUri404(s.RootUri))
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

      //  if !found {
      //    goto notFound
      //  }
      //
      //  // 不把RootUri看着服务级别的权限控制，视为location的匹配
      //  //if s.RootUri != "" {
      //  //  p := r.URL.Path
      //  //  if !path.IsAbs(p) {
      //  //    p = "/" + p
      //  //  }
      //  //  if !strings.HasPrefix(p, s.RootUri) {
      //  //    found = false
      //  //    goto notFound
      //  //  }
      //  //}
      //
      //notFound:
      if !found {
        _, logger := log.WithCtx(r.Context())
        logger.Error(fmt.Sprintf("404 not found. HostName is not mattched. HostNames in config are %s, but url is %s",
          strings.Join(s.HostName, ","), r.URL))
        http.NotFound(w, r)
        return
      }

      serverMux.ServeHTTP(w, r)

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

func rootUri404(rootUri string) func(writer http.ResponseWriter, rawRequest *http.Request) {
  return func(writer http.ResponseWriter, rawRequest *http.Request) {
    _, logger := log.WithCtx(rawRequest.Context())
    logger.Error(fmt.Sprintf("404 not found. RootUri is not mattched. RootUri in config is %s, but url is %s",
      rootUri, rawRequest.URL))

    http.NotFound(writer, rawRequest)
  }
}

func _404(writer http.ResponseWriter, rawRequest *http.Request) {
  _, logger := log.WithCtx(rawRequest.Context())
  logger.Error(fmt.Sprintf("404 not found. URI is not mattched. url is %s", rawRequest.URL))

  http.NotFound(writer, rawRequest)
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
