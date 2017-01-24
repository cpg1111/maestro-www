package www

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	http.Server
}

func New(cert, key, host string, port int) (*Server, error) {
	mux := http.NewServeMux()
	ws := NewWS(4096, DefaultWSActions)
	mux.HandleFunc("/ws", ws.HandleUpgrades)
	srv := &Server{
		http.Server{
			Addr:    fmt.Sprintf("%s:%d", host, port),
			Handler: mux,
		},
	}
	if len(cert) > 0 && len(key) > 0 {
		keypair, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		tlsConf := &tls.Config{
			Certificates: []tls.Certificate{keypair},
		}
		srv.TLSConfig = tlsConf
	}
	return srv
}

func (s *Server) Run() error {
	log.Println("listening on %s", s.Addr)
	return s.ListenAndServe()
}
