package reqID

import (
  "context"
  "github.com/google/uuid"
)

type reqIDKey struct {}

const (
  HeaderKey = "X-Req-Id"
)


func WithCtx(parent context.Context) (ctx context.Context, id string) {
  switch value := parent.Value(reqIDKey{}).(type) {
  case string:
    return parent, value
  default:
    id = GetID()
    ctx = context.WithValue(parent, reqIDKey{},id)
    return
  }
}

func NewContext(parent context.Context, id string) context.Context {
  if value,ok := parent.Value(reqIDKey{}).(string); ok && value == id {
    return parent
  }
  return context.WithValue(parent, reqIDKey{}, id)
}

func GetID() string {
   // todo panic
  return uuid.New().String()
}


