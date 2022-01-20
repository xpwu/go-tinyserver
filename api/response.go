package api

import "net/http"

type Response struct {
  request    *Request
  RawData    []byte
  HttpStatus int
  HttpErrMsg string
  Header     http.Header
}

func NewResponse(r *Request) *Response {
  return &Response{
    HttpStatus: http.StatusOK,
    Header:     http.Header{},
    RawData:    make([]byte, 0),
    request:    r,
  }
}

//func NewWithResWriter(r http.ResponseWriter) *Response {
//  return &Response{HttpStatus:http.StatusOK, Header:r.Header(), RawData:make([]byte, 0)}
//}

func (r *Response) Request() *Request {
  return r.request
}
