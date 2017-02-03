package www

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/cpg1111/maestrod/datastore"

	"github.com/cpg1111/maestro-www/auth"
)

type Server struct {
	http.Server
}

// TODO load a config to choose a datastore
func New(cert, key, host string, port int, dStore datastore.Datastore) (*Server, error) {
	mux := http.NewServeMux()
	ws := NewWS(4096, DefaultWSActions)
	authHandler := auth.New(dStore)
	mux.HandleFunc("/ws", ws.HandleUpgrades)
	mux.Handle("/auth", authHandler)
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
