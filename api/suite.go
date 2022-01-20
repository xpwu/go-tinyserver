package api

import (
  "context"
  "encoding/json"
  "encoding/xml"
  "github.com/xpwu/go-log/log"
)

type SetUpper interface {
  // request 输入参数；apiRequest 输出参数，也即时具体api的输入参数
  SetUp(ctx context.Context, request *Request, apiRequest interface{}) bool
}

type TearDowner interface {
  // apiResponse 输入参数，也即是具体api的返回参数；response 输出参数
  TearDown(ctx context.Context, apiResponse interface{}, response *Response)
}

type URIMapper interface {
  MappingPreUri()string
}


/**

Suite 一簇api，一个Suite中可以定义多个外部接口

 Suite接口的实现类型中满足以下要求的方法即可成为提供服务的网络接口：
 1、API为方法名字的前缀，两个输入参数，一个返回值的方法，
 2、第一个参数 满足 context.Context接口，
 3、第二个参数是一个ptr, 为具体的请求值，一般是一个struct 的指针类型
 4、有一个返回值，返回值也必须是ptr，返回值即为此接口的响应值，一般是一个struct 的指针类型
 5、具体的请求值类型与返回值类型由用户自定义

 具体的对外接口的uri由两部分拼接组成
 1、URIMapper 返回的前缀
 2、方法名去掉前缀API后剩下的部分，如果首字母是大写，会分别生成大小写的两个uri，如果是小写，则只会生成一个
 以上的uri加上服务本身设置的host及Scheme 即是外部访问的url。在生成uri时，上面两部分之间如果需要分隔符，会自动添加。
   在实际服务时，如果在服务的实际配置中，有另外的服务级别的url的处理流程，则这里生成的uri代表处理后的uri


 比如： APIGetInfo(ctx Context, request *Request1) *Response1
 假定  URIMapper 返回的值为 /doc/
 生成的url为  host:port/doc/GetInfo 及  host:port/doc/getInfo;

 对于具体服务来说，最后的url还与具体的服务配置有关,
 比如 服务配置要求所有的都添加一个前缀/API, 则会变成 host:port/API/doc/GetInfo 及  host:port/API/doc/getInfo;

 每一个api的执行流程为：
    SetUpper(非post请求在此阶段可以把请求值放入Request的rawData中) ---- 请求数据的预处理及转换为request的类型
      --> api （当SetUpper为true时执行）
    --> TearDowner --- 响应数据的再次处理，需要把api.response的值序列化后放入到Response中，
                        无论SetUpper的返回值，此步都要执行
 以上任何一步出现panic或者主动调用StopCurrentServer<Deprecated> 或者主动调用 Request.Terminate()方法后，
 直接停止请求的处理，后续的流程都不再执行，并返回错误给请求方。



 suite的reset问题：

 <2020.3.10 之前的方案
 在suite的实现类型中可能会使用存储变量供SetUpper 或者 TearDowner 或者接口自己存储值，这类变量的生命周期延续该suite的所有
 接口中，如果变量的值在每次执行接口或者其他执行过程中需要reset，则需要使用方自行完成，tinyGos 不负责变量的reset，目前go的支持
 也不能很准确的判断出哪些变量需要reset。上层可以在SetUpper开始执行时或者TearDowner执行结束时或者执行过程中reset。如果应该
 reset的变量而没有reset可能会影响到下一个接口执行时的逻辑。底层也不能直接在每次执行接口时直接清零，因为suite本身在初始化时可以
 有构造参数，如果清零，构造值将会出错。>

 现在方案：
 每次执行api时，会调用 NewSuite 新建一个Suite 再调用具体的api


*/


type Suite interface {
  SetUpper
  TearDowner
  URIMapper
}

type SuiteCreator func() Suite

type PostJsonSetUpper struct {
  Request *Request
}

func (p *PostJsonSetUpper) SetUp(ctx context.Context, request *Request, apiRequest interface{}) bool {
  _,logger := log.WithCtx(ctx)
  err := json.Unmarshal(request.RawData, apiRequest)
  if err != nil {
    logger.Error(err)
    request.Terminate(err)
  }
  p.Request = request
  return true
}

type PostXmlSetUpper struct {}

func (p *PostXmlSetUpper) SetUp(ctx context.Context, request *Request, apiRequest interface{}) bool {
  _,logger := log.WithCtx(ctx)
  err := xml.Unmarshal(request.RawData, apiRequest)
  if err != nil {
    logger.Error(err)
    request.Terminate(err)
  }
  return true
}

type PostJsonTearDowner struct {}

func (p *PostJsonTearDowner) TearDown(ctx context.Context, apiResponse interface{}, response *Response)  {
  _,logger := log.WithCtx(ctx)
  d,err := json.Marshal(apiResponse)

  if err != nil {
    logger.Error(err)
    response.Request().Terminate(err)
  }
  response.RawData = d
}

type PostXmlTearDowner struct {}

func (p *PostXmlTearDowner) TearDown(ctx context.Context, apiResponse interface{}, response *Response)  {
  _,logger := log.WithCtx(ctx)
  d,err := xml.Marshal(apiResponse)

  if err != nil {
    logger.Error(err)
    response.Request().Terminate(err)
  }
  response.RawData = d
}

type RootURIMapper struct {}

func (r *RootURIMapper) MappingPreUri() string {
  return "/"
}

