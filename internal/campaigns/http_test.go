package campaigns

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPRequiresActorForEveryOrderRoute(t *testing.T) {
	h := Handler(Service{}, map[string]string{"private-token": "advertiser-a"})
	for _, path := range []string{"/v1/orders", "/v1/orders/missing", "/v1/orders/missing/refunds"} {
		r := httptest.NewRequest(http.MethodPost, path, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 401 {
			t.Fatalf("%s: expected unauthorized, got %d", path, w.Code)
		}
	}
}
func TestRejectTrailingJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{} {}`))
	var v map[string]any
	if decode(w, r, &v) == nil {
		t.Fatal("trailing JSON accepted")
	}
}
