package campaigns

import (
	"adweb3.local/platform/web"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

func Handler(s Service, tokens map[string]string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/orders", func(w http.ResponseWriter, r *http.Request) {
		var in Order
		if e := decode(w, r, &in); e != nil {
			return
		}
		o, e := s.CreateOrder(r.Context(), r.Header.Get("X-Actor"), r.Header.Get("Idempotency-Key"), in)
		respond(w, o, e)
	})
	mux.HandleFunc("GET /v1/orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		o, e := s.Get(r.Context(), r.Header.Get("X-Actor"), r.PathValue("id"))
		respond(w, o, e)
	})
	mux.HandleFunc("POST /v1/orders/{id}/refunds", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Amount int64 `json:"amount,string"`
		}
		if e := decode(w, r, &in); e != nil {
			return
		}
		id, e := s.RequestRefund(r.Context(), r.Header.Get("X-Actor"), r.PathValue("id"), r.Header.Get("Idempotency-Key"), in.Amount)
		respond(w, map[string]string{"refundId": id, "status": "pending"}, e)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/campaigns.js" {
			name := "templates/orders.html"
			kind := "text/html; charset=utf-8"
			if r.URL.Path == "/campaigns.js" {
				name = "static/campaigns.js"
				kind = "text/javascript"
			}
			data, e := web.Files.ReadFile(name)
			if e != nil {
				http.Error(w, "unavailable", 500)
				return
			}
			w.Header().Set("Content-Type", kind)
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'unsafe-inline'; frame-ancestors 'none'; base-uri 'none'")
			w.Write(data)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		actor := tokens[strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")]
		if actor == "" {
			http.Error(w, "unauthorized", 401)
			return
		}
		r.Header.Set("X-Actor", actor)
		mux.ServeHTTP(w, r)
	})
}
func decode(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if e := d.Decode(v); e != nil {
		http.Error(w, "invalid JSON", 400)
		return e
	}
	var extra any
	if e := d.Decode(&extra); e != io.EOF {
		http.Error(w, "trailing JSON", 400)
		return ErrInvalid
	}
	return nil
}
func respond(w http.ResponseWriter, v any, e error) {
	w.Header().Set("Content-Type", "application/json")
	if e != nil {
		status := 500
		switch {
		case errors.Is(e, ErrInvalid):
			status = 400
		case errors.Is(e, ErrForbidden):
			status = 404
		case errors.Is(e, ErrConflict):
			status = 409
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": http.StatusText(status)})
		return
	}
	_ = json.NewEncoder(w).Encode(v)
}
func Server(addr string, h http.Handler) *http.Server {
	return &http.Server{Addr: addr, Handler: h, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
}
