package api

import (
  "context"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "reflect"
  "testing"
)

type suitB struct {

}
func (s *suitB) APIGetInfo(ctx context.Context, request *getUserInfoRequest) *getUserInfoResponse {
  return &getUserInfoResponse{}
}

func (s *suitB) APIGetInfo0(ctx context.Context, request getUserInfoRequest) *getUserInfoResponse {
  return &getUserInfoResponse{}
}

func (s *suitB) APIGetInfo1(ctx context.Context, request *getUserInfoRequest) *getUserInfoResponse {
  return &getUserInfoResponse{}
}

func (s *suitB) GetInfo(ctx context.Context, request *getUserInfoRequest) *getUserInfoResponse {
  return &getUserInfoResponse{}
}

type suitT struct {
  suitB
}

func (s *suitT) SetUp(ctx context.Context, request *Request, apiRequest interface{}) bool {
  panic("implement me")
}

func (s *suitT) TearDown(ctx context.Context, apiResponse interface{}, response *Response) {
  panic("implement me")
}

func (s *suitT) MappingPreUri() string {
  return "/test"
}

type getUserInfoRequest struct{ Uid string }
type getUserInfoResponse struct {
  Name string
  Age  int
}

func (s *suitT) APIGetUserInfo(ctx context.Context, request *getUserInfoRequest) *getUserInfoResponse {
  return &getUserInfoResponse{}
}

func (s *suitT) APIGetUserInfo1(request *getUserInfoRequest) *getUserInfoResponse {
  return &getUserInfoResponse{}
}

func (s *suitT) APIGetUserInfo2(ctx context.Context) *getUserInfoResponse {
  return &getUserInfoResponse{}
}

func (s *suitT) APIGetUserInfo3(request getUserInfoRequest) *getUserInfoResponse {
  return &getUserInfoResponse{}
}

func (s *suitT) APIGetUserInfo4(ctx context.Context, request *getUserInfoRequest) getUserInfoResponse {
  return getUserInfoResponse{}
}

func (s *suitT) GetUserInfo5(ctx context.Context, request *getUserInfoRequest) *getUserInfoResponse {

  return &getUserInfoResponse{}
}

func (s *suitT) APIGetInfo1(ctx context.Context, request *getUserInfoRequest) *getUserInfoResponse {
  return &getUserInfoResponse{}
}

func (s *suitT) APIgetInfo2(ctx context.Context, request *getUserInfoRequest) *getUserInfoResponse {
  return &getUserInfoResponse{}
}

func getMethod(t *testing.T, newSuite SuiteCreator, name string) reflect.Method {
  m, ok := reflect.TypeOf(newSuite()).MethodByName(name)
  require.Truef(t, ok, "type(%v) can not find method by name(%s)",
    reflect.TypeOf(newSuite()).Elem().Name(), name)
  return m
}

func TestAdd(t *testing.T) {
  a := assert.New(t)

  newSuite := func() Suite { return &suitT{}}

  expectes := map[string]*struct {
    b         *base
    hasTested bool
  }{
    `/test/GetUserInfo`:
    {
      &base{
        newSuite,
        getMethod(t, newSuite, "APIGetUserInfo"),
        reflect.TypeOf(getUserInfoRequest{}),
        reflect.TypeOf(getUserInfoResponse{}),
      },
      false,
    },

    `/test/getUserInfo`:
    {
      &base{
        newSuite,
        getMethod(t, newSuite, "APIGetUserInfo"),
        reflect.TypeOf(getUserInfoRequest{}),
        reflect.TypeOf(getUserInfoResponse{}),
      },
      false,
    },

    `/test/GetInfo1`:
    {
      &base{
        newSuite,
        getMethod(t, newSuite, "APIGetInfo1"),
        reflect.TypeOf(getUserInfoRequest{}),
        reflect.TypeOf(getUserInfoResponse{}),
      },
      false,
    },

    `/test/getInfo1`:
    {
      &base{
        newSuite,
        getMethod(t, newSuite, "APIGetInfo1"),
        reflect.TypeOf(getUserInfoRequest{}),
        reflect.TypeOf(getUserInfoResponse{}),
      },
      false,
    },

    `/test/getInfo2`:
    {
      &base{
        newSuite,
        getMethod(t, newSuite, "APIgetInfo2"),
        reflect.TypeOf(getUserInfoRequest{}),
        reflect.TypeOf(getUserInfoResponse{}),
      },
      false,
    },

    `/test/getInfo`:
    {
      &base{
        newSuite,
        getMethod(t, newSuite, "APIGetInfo"),
        reflect.TypeOf(getUserInfoRequest{}),
        reflect.TypeOf(getUserInfoResponse{}),
      },
      false,
    },

    `/test/GetInfo`:
    {
      &base{
        newSuite,
        getMethod(t, newSuite, "APIGetInfo"),
        reflect.TypeOf(getUserInfoRequest{}),
        reflect.TypeOf(getUserInfoResponse{}),
      },
      false,
    },
  }

  add(newSuite, func(uri string, api API) {
    e, ok := expectes[uri]
    if !a.Truef(ok, "uri(%s) is not expect", uri) {
      return
    }

    if !a.Falsef(e.hasTested, "uri(%s) is duplicated", uri) {
      return
    }

    e.hasTested = true
    b := api.(*base)

    a.Equalf(e.b.newSuite(), b.newSuite(), "new suite error of uri(%s)", uri)
    a.Equalf(e.b.method.Name, b.method.Name, "method error of uri(%s)", uri)
    a.Equalf(e.b.requestType, b.requestType, "requestType error of uri(%s)", uri)
    a.Equalf(e.b.responseType, b.responseType, "responseType error of uri(%s)", uri)
  })

  for k, e := range expectes {
    a.Truef(e.hasTested, "not register uri(%s)", k)
  }

}
