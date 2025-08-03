package main

import (
	"crypto/x509"
	"crypto/tls"
	"fmt"
)


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

func printLogger(data *ReqLogger){
	fmt.Println(data.timestamp + " | " + data.source_ip +" | " + data.method + " | " + data.url + " | " + data.user_agent)
}