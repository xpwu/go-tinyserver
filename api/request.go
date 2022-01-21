package api

import (
	"fmt"
	"net/http"
	"net/url"
)

type Request struct {
	RawReq *http.Request
	// post 的boby 数据
	RawData []byte
	Header  http.Header
	// url ? 后的query数据
	Query url.Values
	URI   string
	ReqID string
}

//var normalReqHeader = map[string]bool{"Accept": true, "Accept-Charset": true, "Accept-Encoding": true,
//	"Accept-Language": true, "Cache-Control": true, "Connection": true, "Cookie": true,
//	"Content-Length": true, "Content-Type": true, "Host": true, "Referer": true, "User-Agent": true,
//	reqID.HeaderKey: true, "Sec-Fetch-Mode": true, "Sec-Fetch-Site": true, "Sec-Fetch-User": true,
//	"Upgrade-Insecure-Requests": true}

func (r *Request) String() string {
	//header := r.Header.Clone()
	//for key := range normalReqHeader {
	//  header.Del(key)
	//}
	return fmt.Sprintf("api:%s, reqid:%s", r.URI, r.ReqID)
}

func (r *Request) Terminate(err error) {
	stopCurrentRequest(err)
}
