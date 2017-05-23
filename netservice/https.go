package netservice

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
)

func NewHttpsServer() *HttpsServer {
	return &HttpsServer{}
}

//certificate config
type CertificateConfig struct {
	Domain   string
	CertFile string
	KeyFile  string
}

//https server
type HttpsServer struct {
	http.Server
	//If GetCertificate is set, the leaf certificate is returned by calling this function.
	//	func(clientInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
	//		fmt.Println(clientInfo.ServerName)
	//		x509Cert, err := tls.LoadX509KeyPair("config/ca/server.crt", "config/ca/server.key")
	//		result, _ := x509.ParseCertificate(x509Cert.Certificate[0])
	//		return nil, err
	//	}
	GetCertificate func(*tls.ClientHelloInfo) (*tls.Certificate, error)
}

//If the GetCertificate field is not set, defaultGetCertificate will be used as the default value
func (self *HttpsServer) defaultGetCertificate(clientInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if x509Cert, ok := self.TLSConfig.NameToCertificate[clientInfo.ServerName]; ok {
		return x509Cert, nil
	}
	clientInfo.Conn.Close()
	return nil, errors.New("Did't find the specified digital certificate")
}

//Add the certificate to the tls.Config.Certificates list, and add the domain name mapping
func (self *HttpsServer) AddDomainCertificateItem(domain, certFile, keyFile string) error {
	if domain == "" || certFile == "" || keyFile == "" {
		errMsg := fmt.Sprintf(
			"The parameters of the addDomainCertificate function are incorrect!,Domain:%s,certFile:%s,keyFile:%s",
			domain,
			certFile,
			keyFile)
		return errors.New(errMsg)
	}
	if x509Cert, err := tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		return err
	} else {
		if self.TLSConfig == nil {
			self.TLSConfig = &tls.Config{}
		}
		if self.TLSConfig.NameToCertificate == nil {
			self.TLSConfig.NameToCertificate = make(map[string]*tls.Certificate)
		}
		self.TLSConfig.Certificates = append(self.TLSConfig.Certificates, x509Cert)
		self.TLSConfig.NameToCertificate[domain] = &x509Cert
	}
	return nil
}

//Add the certificate to the tls.Config.Certificates list, and add the domain name mapping
func (self *HttpsServer) AddDomainCertificateConfig(config []*CertificateConfig) error {
	for _, v := range config {
		if err := self.AddDomainCertificateItem(v.Domain, v.CertFile, v.KeyFile); err != nil {
			return err
		}
	}
	return nil
}

//Check the legitimacy of https access
func (self *HttpsServer) checkValidHttpsReq(host string) bool {
	if _, ok := self.TLSConfig.NameToCertificate[host]; ok {
		return true
	}
	return false
}

//Start the https server
//If the two parameters certFile and keyFile are empty,
//you must call the addDomainCertificate function or the addDomainCertificate function
//to add a list of digital certificates to multiple certificates list
func (self *HttpsServer) RunHttpsService(addr, certFile, keyFile string, handler http.Handler) error {
	self.Addr = addr
	self.Handler = handler
	if self.TLSConfig == nil {
		self.TLSConfig = &tls.Config{}
	}
	if self.GetCertificate != nil {
		self.TLSConfig.GetCertificate = self.GetCertificate
	} else {
		self.TLSConfig.GetCertificate = self.defaultGetCertificate
	}
	if self.GetCertificate == nil && self.TLSConfig == nil {
		return errors.New("RunHttpsService:No Https configuration,Please call AddDomainCertificateConfig  AddDomainCertificateItem or  AddDomainCertificateItem function......")
	}
	return self.ListenAndServeTLS(certFile, keyFile)
}
