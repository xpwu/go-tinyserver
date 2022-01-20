package server

import (
  "github.com/xpwu/go-tinyserver/config"
  "net"
  "net/http"
  "net/http/fcgi"
  "sync"
  "time"
)

type conChanT chan byte

var con struct {
  once    sync.Once;
  conChan conChanT
}

func conChan() conChanT {
  con.once.Do(func() {
    con.conChan = make(conChanT, config.Con.MaxConnections)
  })

  return con.conChan
}

type conConnection struct {
  net.Conn
}

func (cc *conConnection)Close() error {
  err := cc.Conn.Close()
  <-conChan()

  return err
}

type conListener struct {
  net.Listener
  close chan struct{}
}

func newConL(listener net.Listener) net.Listener {
  return &conListener{
    Listener: listener,
    close:    make(chan struct{}),
  }
}

func (cl *conListener) Accept() (net.Conn, error) {
  select {
  case conChan() <- 0:
  case <-cl.close:
    //return nil, errors.New("accept closed")
    return nil, &net.OpError{Op: "accept", Net: "tcp", Source: nil, Addr: cl.Addr(), Err: nil}
  }

  conn,err := cl.Listener.Accept()
  if err != nil {
    return nil, err
  }

  return &conConnection{conn}, nil
}

func (cl *conListener) Close() error {
  close(cl.close)
  return cl.Listener.Close()
}

func conConnState(c net.Conn, s http.ConnState) {
  switch s{
  case http.StateNew:
    conChan()<-0
  case http.StateClosed:
    <-conChan()
  }
}

var protocols = map[string]func(){
  "http": func() {
    server := &http.Server{
      Addr: config.Server.RealListen(),
      Handler: nil,
      ConnState:conConnState,
      IdleTimeout:30*time.Second,
    }
    err := server.ListenAndServe()
    if err != nil {
      panic(err)
    }
  },

  "https": func() {
    server := &http.Server{
      Addr: config.Server.RealListen(),
      Handler: nil,
      ConnState:conConnState,
      IdleTimeout:30*time.Second,
    }
    err := server.ListenAndServeTLS(config.Server.TLS.RealCertFile(),
      config.Server.TLS.RealKeyFile())
    if err != nil {
      panic(err)
    }
  },

  "fastcgi": func() {
    ln, err := net.Listen("tcp", config.Server.RealListen())
    if err != nil {
      panic(err)
    }
    ln = newConL(ln)

    defer func() {
      _ = ln.Close()
    }()

    err = fcgi.Serve(ln, nil)
    if err != nil {
      panic(err)
    }
  },
}

func start() {
  run, ok := protocols[config.Server.UseProtocol]
  if !ok {
    panic("not support protocol: " + config.Server.UseProtocol)
  }
  run()
}
