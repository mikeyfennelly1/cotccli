package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStreamHierarchy_PrintsTree(t *testing.T) {
	const fixture = `[
	  {
	    "id": 1,
	    "name": "test",
	    "children": [
	      {
	        "id": 3,
	        "name": "test.child",
	        "children": [],
	        "sources": []
	      },
	      {
	        "id": 4,
	        "name": "test.child2",
	        "children": [],
	        "sources": []
	      }
	    ],
	    "sources": []
	  }
	]`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/reporting/streams/hierarchy" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fixture))
	}))
	defer srv.Close()

	client := NewReportingClient(srv.URL)
	if err := client.GetStreamHierarchy(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
