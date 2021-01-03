package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type cacheHandler struct {
	*Server
}

// To implement interface: http.Handler
func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.EscapedPath(), "/")[2]
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m := r.Method
	if m == http.MethodPut {
		b, e := ioutil.ReadAll(r.Body)
		if e != nil {
			// The original code just ignores the error here, not sure why
			// Is empty body an error??
			log.Println(e)
			w.WriteHeader(http.StatusBadRequest)
		}
		if len(b) != 0 {
			e := h.Set(key, b)
			if e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}

	} else if m == http.MethodGet {
		b, e := h.Get(key)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(b) == 0 {
			// I don't like it. I think it should be able to store key with empty value
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(b)

	} else if m == http.MethodDelete {
		e := h.Del(key)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
}
