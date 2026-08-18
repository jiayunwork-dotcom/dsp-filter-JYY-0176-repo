package api

import (
	"io/fs"
	"net/http"
	"sort"
	"strings"
)

type MetaResponse struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Methods  []string `json:"methods"`
	Examples []string `json:"examples"`
}

func (s *Server) handleMeta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, http.StatusMethodNotAllowed, "method_not_allowed", "use GET")
		return
	}
	entries, err := fs.ReadDir(s.example, ".")
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	examples := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			examples = append(examples, e.Name())
		}
	}
	sort.Strings(examples)
	writeJSON(w, http.StatusOK, MetaResponse{
		Name:     "dsp-filter",
		Version:  "1.0.0",
		Methods:  []string{"POST /api/design", "POST /api/response", "POST /api/zplane"},
		Examples: examples,
	})
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")
	if name == "" {
		name = "index.html"
	}
	data, err := fs.ReadFile(s.web, name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	writeFile(w, name, data)
}

func (s *Server) handleExample(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/example/")
	if name == "" || strings.Contains(name, "..") {
		http.NotFound(w, r)
		return
	}
	data, err := fs.ReadFile(s.example, name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	writeFile(w, name, data)
}

func writeFile(w http.ResponseWriter, name string, data []byte) {
	ctype := "text/plain; charset=utf-8"
	switch {
	case strings.HasSuffix(name, ".html"):
		ctype = "text/html; charset=utf-8"
	case strings.HasSuffix(name, ".js"):
		ctype = "application/javascript; charset=utf-8"
	case strings.HasSuffix(name, ".css"):
		ctype = "text/css; charset=utf-8"
	case strings.HasSuffix(name, ".json"):
		ctype = "application/json; charset=utf-8"
	}
	w.Header().Set("Content-Type", ctype)
	w.Write(data)
}
