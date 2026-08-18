package api

import (
	"io/fs"
	"net/http"
)

type Server struct {
	web     fs.FS
	example fs.FS
	mux     *http.ServeMux
}

func New(webFS, exampleFS fs.FS) *Server {
	s := &Server{web: webFS, example: exampleFS}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/design", s.handleDesign)
	mux.HandleFunc("/api/response", s.handleResponse)
	mux.HandleFunc("/api/zplane", s.handleZPlane)
	mux.HandleFunc("/api/apply", s.handleApply)
	mux.HandleFunc("/api/meta", s.handleMeta)
	mux.HandleFunc("/example/", s.handleExample)
	mux.HandleFunc("/", s.handleStatic)
	s.mux = mux
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
