# go-tinyserver
一个轻的go服务器，实现了绝大部分场景的需求。支持api suite方式

## api suite
调用api.Add()可以把一个api suite注册到服务中。一个 api suite 常常有相同的前缀，
共同的预处理逻辑，在逻辑上也更有紧密型。在每次服务调用过程时，都是一个全新的上下文，
所以suite本身不具备缓存的功能 。  

```
type SetUpper interface {
  SetUp(ctx context.Context, request *Request, apiRequest interface{}) bool
}

type TearDowner interface {
  TearDown(ctx context.Context, apiResponse interface{}, response *Response)
}

type URIMapper interface {
  MappingPreUri()string
}

type Suite interface {
  SetUpper
  TearDowner
  URIMapper
}

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
 假定  URIMapper 返回的值为 /doc
 生成的url为  host:port/doc/GetInfo 及  host:port/doc/getInfo;

 对于具体服务来说，最后的url还与具体的服务配置有关,
 比如 服务配置要求所有的都添加一个前缀/API, 则会变成 host:port/API/doc/GetInfo 及  host:port/API/doc/getInfo;

 每一个api的执行流程为：
    SetUpper(非post请求在此阶段可以把请求值放入Request的rawData中) ---- 请求数据的预处理及转换为request的类型
      --> api （当SetUpper为true时执行）
    --> TearDowner --- 响应数据的再次处理，需要把api.response的值序列化后放入到Response中，
                        无论SetUpper的返回值，此步都要执行
 以上任何一步出现panic或者主动调用 Request.Terminate()方法后，
 直接停止请求的处理，后续的流程都不再执行，并返回错误给请求方。


```

## api Register
如果注册一个普通api，可以调用api.Register()直接注册，uri的规则与底层http pattern一样。  
此处注册的api与具体的底层服务协议无关，无论是使用http 还是 fastcgi等，都可调用到注册的api。


## http server
Http 的服务器逻辑，配置文件的生成及读取使用了[go-config](https://github.com/xpwu/go-config)。  
RootUri: 表示此服务所有uri的前缀，与api中注册的uri一起组合成最后匹配的uri。  
调用http.Start()即可启动http服务。

## fastcgi server
fastcgi 服务器逻辑，配置文件的生成及读取使用了[go-config](https://github.com/xpwu/go-config)。  
调用fastcgi.Start()即可启动fastcgi服务。  
  

