package api

import (
  "context"
  "errors"
  "fmt"
  "github.com/xpwu/go-log/log"
  "net/http"
  "path"
  "reflect"
  "strings"
)

type base struct {
  newSuite     SuiteCreator
  method       reflect.Method
  requestType  reflect.Type
  responseType reflect.Type
}

func (base *base) Process(ctx context.Context, request *Request) (response *Response) {

  response = NewResponse(request)
  response.HttpStatus = http.StatusInternalServerError

  ctx, _ = log.WithCtx(ctx)

  apiReq := reflect.New(base.requestType)
  var apiRes interface{} = &struct{}{}

  suite := base.newSuite()

  if suite.SetUp(ctx, request, apiReq.Interface()) {

    in := []reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(ctx), apiReq}
    out := base.method.Func.Call(in)
    if len(out) != 1 {
      request.Terminate(errors.New("len(out) != 1, 不满足suite api 的输出要求，请检查代码"))
    }
    apiRes = out[0].Interface()

    response.HttpStatus = http.StatusOK
  }

  suite.TearDown(ctx, apiRes, response)

  return
}

// 参见suite的注释说明

const preAPI = "API"

func Add(newSuite SuiteCreator) {
  add(newSuite, Register)
}

func errStr(method, reason string) string {
  return fmt.Sprintf(
    "%s has %s prefix, but %s. NOT add this method to server",
    method, preAPI, reason)
}

func add(newSuite SuiteCreator, register func(uri string, api API)) {
  logger := log.NewLogger()

  suite := newSuite()
  typ := reflect.TypeOf(suite)

  logger.PushPrefix(fmt.Sprintf("add suite(name:%s) to server. ", typ.Elem().Name()))
  defer logger.PopPrefix()

  for i := 0; i < typ.NumMethod(); i++ {
    // Method must be exported.
    method := typ.Method(i)
    if method.PkgPath != "" {
      continue
    }

    mName := method.Name
    if !strings.HasPrefix(mName, preAPI) {
      continue
    }
    mName = mName[len(preAPI):]

    mType := method.Type

    // Method needs three ins: receiver, ctx, request.
    if mType.NumIn() != 3 {
      log.Warning(errStr(method.Name, "it does have two parameters"))
      continue
    }

    // first must implement Context
    argType := mType.In(1)
    if !argType.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
      log.Warning(errStr(method.Name, "the first parameter is not a context.Context"))
      continue
    }
    // second must be ptr
    reqType := mType.In(2)
    if reqType.Kind() != reflect.Ptr {
      log.Warning(errStr(method.Name, "the second parameter is not a pointer type"))
      continue
    }
    reqType = reqType.Elem()

    // return must be prt
    if mType.NumOut() != 1 {
      log.Warning(errStr(method.Name, "it does have onw return value"))
      continue
    }
    returnType := mType.Out(0)
    if returnType.Kind() != reflect.Ptr {
      log.Warning(errStr(method.Name, "the return parameter is not a pointer type"))
      continue
    }
    returnType = returnType.Elem()

    b := &base{newSuite: newSuite, method: method,
      requestType: reqType, responseType: returnType}

    // 大小写都加入
    register(path.Join(suite.MappingPreUri(), mName), b)
    if 'A' <= mName[0] && mName[0] <= 'Z' {
      register(path.Join(suite.MappingPreUri(), strings.ToLower(mName[0:1])+mName[1:]), b)
    }
  }
}

