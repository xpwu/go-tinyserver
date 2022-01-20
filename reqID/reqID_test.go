package reqID

import (
  "github.com/stretchr/testify/assert"
  "net/textproto"
  "testing"
)

func TestReqId(t *testing.T) {
  assert.Equal(t, HeaderKey, textproto.CanonicalMIMEHeaderKey(HeaderKey))
}

func TestInterface(t *testing.T) {
  a := reqIDKey{}
  c := reqIDKey{}
  if a != c {
    t.Error("not equal")
  }
  assert.Equal(t, reqIDKey{}, reqIDKey{})
}
