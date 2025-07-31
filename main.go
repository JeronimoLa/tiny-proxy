package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net/http"
	"fmt"
	"net/http/httputil"

	"github.com/elazarl/goproxy"
)

func main() {
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()

	cert, err := parseCA(_caCert, _caKey)
	if err != nil {
		log.Fatal(err)
	}
	// Use default transport as a base
	baseTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	// defines how the proxy should handle HTTPS CONNECT requests (when a client wants to establish a TLS tunnel through the proxy)
	customCaMitm := &goproxy.ConnectAction{
		Action:    goproxy.ConnectMitm, // Determines the basic action the proxy take for the CONNECT request
		TLSConfig: goproxy.TLSConfigFromCA(cert), // lets you dynamically generate or return a custom TLS configuration to be used  when the proxy acts asa TLS server
	}
	var customAlwaysMitm goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return customCaMitm, host
	}

	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest().HandleConnect(customAlwaysMitm)
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ctx.RoundTripper = goproxy.RoundTripperFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Response, error) {
			clonedTransport := baseTransport.Clone()
			hostname := req.URL.Hostname()
			clonedTransport.TLSClientConfig = &tls.Config{
				ServerName: hostname, // This fixes the TLS error
			}
			return clonedTransport.RoundTrip(req)
		})
		// fmt.Println(req.Method, req.Header, req.URL.String())
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))

		return req, nil
	})

	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx)*http.Response{
		// fmt.Println(ctx.Req.Host,"->",r.Header.Get("Content-Type"))
		// b, err := httputil.DumpResponse(r, true)
		// if err != nil {
			// log.Fatalln(err)
		// }

		// fmt.Println(string(b))
		// b, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// fmt.Println(string(b))
		return r
	})
	proxy.Verbose = *verbose

	log.Fatal(http.ListenAndServe(*addr, proxy))
}


func parseCA(caCert, caKey []byte) (*tls.Certificate, error) {
	parsedCert, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return nil, err
	}
	if parsedCert.Leaf, err = x509.ParseCertificate(parsedCert.Certificate[0]); err != nil {
		return nil, err
	}
	return &parsedCert, nil
}