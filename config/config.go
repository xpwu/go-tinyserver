package config

import (
	"github.com/xpwu/go-cmd/exe"
	"path/filepath"
	"strings"
)

type server struct {
	ServerName         string
	Listen             string
	UseProtocol        string `conf:",http, https, fastcgi"`
	TLS                tls
}

type concurrent struct {
	MaxConnections int
}

func (s *server) RealListen() string {
	if !strings.Contains(s.Listen, ":") {
		return ":" + s.Listen
	}

	return s.Listen
}

type tls struct {
	PrivateKeyPEMFile string `conf:",support relative path, must PEM encode data"`
	CertPEMFile       string `conf:",support relative path, must PEM encode data"`
}

func (t *tls) RealKeyFile() string {
	if filepath.IsAbs(t.PrivateKeyPEMFile) {
		return t.PrivateKeyPEMFile
	}

	return filepath.Join(exe.Exe.AbsDir, t.PrivateKeyPEMFile)
}

func (t *tls) RealCertFile() string {
	if filepath.IsAbs(t.CertPEMFile) {
		return t.CertPEMFile
	}

	return filepath.Join(exe.Exe.AbsDir, t.CertPEMFile)
}

var Server = &server{
	Listen:             "80",
	UseProtocol:        "http",
	TLS: tls{
		PrivateKeyPEMFile: "",
		CertPEMFile:       "",
	},
}

var Con = &concurrent{MaxConnections: 10000}

//
//func init() {
//  //config.Unmarshal(Server)
//  //config.Unmarshal(Con)
//}
