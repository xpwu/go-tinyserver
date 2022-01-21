package api

import (
	"context"
	"fmt"
	"github.com/xpwu/go-config/configs"
	"github.com/xpwu/go-log/log"
	"github.com/xpwu/go-reqid/reqid"
	"io/ioutil"
	"net/http"
	"strconv"
)

type API interface {
	Process(ctx context.Context, request *Request) *Response
}

const allowOrigin = "Access-Control-Allow-Origin"

func defaultOptions(ctx context.Context, writer http.ResponseWriter) {
  ctx, logger := log.WithCtx(ctx)
	logger.Debug("OPTIONS method")

	writer.Header().Add(allowOrigin, "*")
	writer.Header().Add("Access-Control-Max-Age", strconv.Itoa(24*3600))
	writer.Header().Add("Access-Control-Allow-Headers",
		" accept, content-type, _t, _i, _f, _l, _s,Accept-Language,"+
			"Content-Language,Origin, No-Cache, X-Requested-With, If-Modified-Since,"+
			" Pragma, Last-Modified, Cache-Control, Expires, Content-Type, "+
			"X-E4M-With,authorization,application/x-www-form-urlencoded,multipart/form-data,text/plain")
	writer.Header().Add("Access-Control-Allow-Methods", " OPTIONS, POST, GET")

	logger.Info("OPTIONS end")
}

func Register(uri string, api API) {
  RegisterApiAndOpt(uri, api, defaultOptions)
}

// api GET POST 调用的方法
// opt OPTIONS 调用的方法
func RegisterApiAndOpt(uri string, api API, opt func(ctx context.Context, write http.ResponseWriter)) {
	if !configs.HasRead() {
		panic("config has not read, registering api must be after config.Read()")
	}
	// 由各Sever控制Hostname的访问
	//if config.Server.ServerName != "" {
	//  uri = config.Server.ServerName + uri
	//}

	// OPTIONS && path = * 的情况 底层已默认处理

	http.HandleFunc(uri, func(writer http.ResponseWriter, rawRequest *http.Request) {

		ctx := rawRequest.Context()

		id := rawRequest.Header.Get(reqid.HeaderKey)
		if id == "" {
			id = reqid.RandomID()
		}

		ctx, id = reqid.WithCtx(ctx)
		ctx, logger := log.WithCtx(ctx)

		request := &Request{
			RawReq: rawRequest,
			Header: rawRequest.Header,
			URI:    rawRequest.URL.Path,
			ReqID:  id,
		}

		logger.PushPrefix(fmt.Sprintf("api:%s, reqid:%s", request.URI, request.ReqID))
		defer logger.PopPrefix()

		if rawRequest.Method == http.MethodOptions {
      opt(ctx, writer)
			return
		}

    request.Query = rawRequest.URL.Query()

		if rawRequest.Method == http.MethodPost {
			request.RawData, _ = ioutil.ReadAll(rawRequest.Body)
		}

		logger.Debug("start.")

		response := NewResponse(request)
		func() {
			defer func() {
				if r := recover(); r != nil {
					response.HttpStatus = http.StatusInternalServerError
					response.HttpErrMsg = http.StatusText(response.HttpStatus)
					if rr, ok := r.(stopError); ok {
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

		hasAllowOrigin := false
		for key, values := range response.Header {
		  if key == allowOrigin {
		    hasAllowOrigin = true
      }
			for _, value := range values {
				writer.Header().Add(key, value)
			}
		}
		if !hasAllowOrigin {
      writer.Header().Add(allowOrigin, "*")
    }

		// 成功时，可能没有被设置，默认为零值。
		if response.HttpStatus != 0 && response.HttpStatus != http.StatusOK {
			//http.Error(writer, response.HttpErrMsg, response.HttpStatus)
			writer.WriteHeader(response.HttpStatus)
			logger.Error("end. ", response.HttpStatus, " ", response.HttpErrMsg)
			return
		}

		writer.WriteHeader(http.StatusOK)

		if _, err := writer.Write(response.RawData); err != nil {
			logger.Error(err)
		}

		logger.Info("end. ")
	})
}

//func Do404(api API) {
//	Register("/", api)
//}

type stopError error

// Deprecated:
//func StopCurrentServer(err error) {
//	panic(stopError(err))
//}

func stopCurrentRequest(err error) {
	panic(stopError(err))
}
