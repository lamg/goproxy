package main

import (
	"flag"
	gp "github.com/lamg/goproxy"
	"log"
	"net"
	h "net/http"
)

func main() {
	var addr string
	flag.StringVar(&addr, "a", ":8081", "Address to serve the proxy")
	flag.Parse()
	rgs := []string{"10.2.9.0/24"}
	iprgs := make([]*net.IPNet, len(rgs))
	var e error
	for i := 0; e == nil && i != len(rgs); i++ {
		_, iprgs[i], e = net.ParseCIDR(rgs[i])
	}
	if e == nil {
		proxy := gp.NewProxyHttpServer()
		proxy.OnRequest().DoFunc(
			func(r *h.Request, ctx *gp.ProxyCtx) (n *h.Request,
				p *h.Response) {
				n = r
				host, _, e := net.SplitHostPort(r.RemoteAddr)
				if e == nil {
					ni := net.ParseIP(host)
					ok := false
					for i := 0; !ok && i != len(iprgs); i++ {
						ok = iprgs[i].Contains(ni)
					}
					if ok {
						p = nil
					} else {
						p = gp.NewResponse(r, gp.ContentTypeText,
							h.StatusForbidden, "Not allowed IP range")
					}
				} else {
					p = gp.NewResponse(r, gp.ContentTypeText,
						h.StatusBadRequest, "Malformed request")
				}
				return
			})
		log.Fatalln(h.ListenAndServe(addr, proxy))
	}
}
