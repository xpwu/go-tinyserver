package api

import (
  "context"
  "fmt"
  config2 "github.com/xpwu/go-config/config"
  "github.com/xpwu/go-log/log"
  "github.com/xpwu/go-tinyserver/reqID"
  "io/ioutil"
  "net/http"
  "strconv"
)

type API interface {
  Process(ctx context.Context, request *Request) *Response
}

func options(writer http.ResponseWriter, logger *log.Logger) {
 logger.Debug("OPTIONS method")

 writer.Header().Add("Access-Control-Max-Age", strconv.Itoa(24 * 3600))
 writer.Header().Add("Access-Control-Allow-Headers",
   " accept, content-type, _t, _i, _f, _l, _s,Accept-Language," +
   "Content-Language,Origin, No-Cache, X-Requested-With, If-Modified-Since," +
   " Pragma, Last-Modified, Cache-Control, Expires, Content-Type, " +
   "X-E4M-With,authorization,application/x-www-form-urlencoded,multipart/form-data,text/plain")
 writer.Header().Add("Access-Control-Allow-Methods", " OPTIONS, POST")

 logger.Info("OPTIONS end")
}


func Register(uri string, api API)  {
  if !config2.HasRead() {
    panic("config has not read, registering api must be after config.Read()")
  }
  // 由各Sever控制Hostname的访问
  //if config.Server.ServerName != "" {
  //  uri = config.Server.ServerName + uri
  //}

  // OPTIONS && path = * 的情况 底层已默认处理
  http.HandleFunc(uri, func(writer http.ResponseWriter, rawRequest *http.Request) {

    id := rawRequest.Header.Get(reqID.HeaderKey)
    if id == "" {
      id = reqID.GetID()
    }

    ctx := reqID.NewContext(rawRequest.Context(), id)
    ctx,logger := log.WithCtx(ctx)

    request := &Request{
      RawReq: rawRequest,
      Header: rawRequest.Header,
      URI:    rawRequest.URL.Path,
      ReqID:  id,
    }

    logger.PushPrefix(fmt.Sprintf("api:%s, reqid:%s", request.URI, request.ReqID))
    defer logger.PopPrefix()

    writer.Header().Add("Access-Control-Allow-Origin", "*")

    if rawRequest.Method == http.MethodOptions {
     options(writer, logger)
     return
    }

    if rawRequest.Method == http.MethodPost {
      request.RawData,_ = ioutil.ReadAll(rawRequest.Body)
    }

    logger.Debug("start.")

    response := NewResponse(request)
    func(){
      defer func() {
        if r := recover(); r != nil {
          response.HttpStatus = http.StatusInternalServerError
          response.HttpErrMsg = http.StatusText(response.HttpStatus)
          if rr,ok := r.(stopError); ok {
            // 打印真正调用terminate()的行号信息
            l := logger.AddSkipCallerDepth(4)
            l.PushPrefix("RequestTerminate!")
            l.Error(rr)
          } else {
            logger.Fatal(r)
          }
        }
      }()

      response = api.Process(ctx, request)
    }()

    for key,values := range response.Header {
      for _,value := range values {
        writer.Header().Add(key, value)
      }
    }

    // 成功时，可能没有被设置，默认为零值。
    if response.HttpStatus != 0 && response.HttpStatus != http.StatusOK {
      //http.Error(writer, response.HttpErrMsg, response.HttpStatus)
      writer.WriteHeader(response.HttpStatus)
      logger.Error("end. ", response.HttpStatus, " ", response.HttpErrMsg)
      return
    }

    writer.WriteHeader(http.StatusOK)

    if _,err := writer.Write(response.RawData);err != nil {
      logger.Error(err)
    }

    logger.Info("end. ")
  })
}

func Do404(api API) {
  Register("/", api)
}

type stopError error

// Deprecated:
func StopCurrentServer(err error) {
  panic(stopError(err))
}

func stopCurrentRequest(err error) {
  panic(stopError(err))
}
