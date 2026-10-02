package mocktest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// fakeAPI is a minimal in-memory runZero API used to exercise resources end to end
// through Terraform without real credentials.
type fakeAPI struct {
	t           *testing.T
	server      *httptest.Server
	mu          sync.Mutex
	seq         int
	clock       int64
	collections []string
	objects     map[string]map[string]any // keyed by "<collection>/<id>"
	orgKey      map[string]any
	requests    []string
}

func newFakeAPI(t *testing.T) *fakeAPI {
	t.Helper()
	f := &fakeAPI{
		t:           t,
		clock:       1700000000,
		collections: []string{"/org/sites", "/org/tasks", "/account/groups"},
		objects:     map[string]map[string]any{},
	}
	f.server = httptest.NewServer(http.HandlerFunc(f.handle))
	t.Cleanup(f.server.Close)
	return f
}

// providerConfig returns a provider block that points at the fake API.
func (f *fakeAPI) providerConfig() string {
	return fmt.Sprintf(`
provider "runzero" {
  api_url       = %q
  org_token     = "OT_fake"
  account_token = "CT_fake"
}
`, f.server.URL)
}

func (f *fakeAPI) nextID(prefix string) string {
	f.seq++
	return fmt.Sprintf("%s-%04d", prefix, f.seq)
}

func (f *fakeAPI) tick() int64 {
	f.clock++
	return f.clock
}

func (f *fakeAPI) get(key string) (map[string]any, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	obj, ok := f.objects[key]
	return obj, ok
}

func (f *fakeAPI) remove(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.objects, key)
}

func (f *fakeAPI) count(collection string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for k := range f.objects {
		if strings.HasPrefix(k, collection+"/") {
			n++
		}
	}
	return n
}

func (f *fakeAPI) received(request string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, r := range f.requests {
		if r == request {
			return true
		}
	}
	return false
}

func (f *fakeAPI) setOrgKey(key map[string]any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.orgKey = key
}

func (f *fakeAPI) handle(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.requests = append(f.requests, r.Method+" "+r.URL.Path)

	auth := r.Header.Get("Authorization")
	if auth != "Bearer OT_fake" && auth != "Bearer CT_fake" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var body map[string]any
	if r.Body != nil && r.ContentLength != 0 {
		_ = json.NewDecoder(r.Body).Decode(&body)
	}

	p := r.URL.Path

	switch {
	case p == "/org/key":
		f.handleOrgKey(w, r)
		return
	case p == "/org/key/rotate" && r.Method == http.MethodPatch:
		if f.orgKey == nil {
			writeStatus(w, http.StatusNotFound)
			return
		}
		f.orgKey["token"] = f.nextID("OT_rotated")
		f.orgKey["counter"] = int64(0)
		writeJSON(w, f.orgKey)
		return
	case strings.HasPrefix(p, "/org/sites/") && strings.HasSuffix(p, "/scan") && r.Method == http.MethodPut:
		siteID := strings.TrimSuffix(strings.TrimPrefix(p, "/org/sites/"), "/scan")
		if _, ok := f.objects["/org/sites/"+siteID]; !ok {
			writeStatus(w, http.StatusNotFound)
			return
		}
		now := f.tick()
		task := map[string]any{
			"id":         f.nextID("task"),
			"name":       body["scan-name"],
			"site_id":    siteID,
			"type":       "scan",
			"status":     "new",
			"created_at": now,
			"updated_at": now,
		}
		f.objects["/org/tasks/"+task["id"].(string)] = task
		writeJSON(w, task)
		return
	case strings.HasPrefix(p, "/org/tasks/") && strings.HasSuffix(p, "/stop") && r.Method == http.MethodPost:
		id := strings.TrimSuffix(strings.TrimPrefix(p, "/org/tasks/"), "/stop")
		if task, ok := f.objects["/org/tasks/"+id]; ok {
			task["status"] = "stopped"
		}
		writeStatus(w, http.StatusOK)
		return
	}

	for _, coll := range f.collections {
		if p == coll {
			f.handleCollection(w, r, coll, body)
			return
		}
		if strings.HasPrefix(p, coll+"/") {
			f.handleItem(w, r, p, body)
			return
		}
	}

	f.t.Errorf("fake API: unhandled request %s %s", r.Method, p)
	writeStatus(w, http.StatusNotImplemented)
}

func (f *fakeAPI) handleCollection(w http.ResponseWriter, r *http.Request, coll string, body map[string]any) {
	switch r.Method {
	case http.MethodGet:
		items := []map[string]any{}
		for k, v := range f.objects {
			if strings.HasPrefix(k, coll+"/") {
				items = append(items, v)
			}
		}
		writeJSON(w, items)
	case http.MethodPut, http.MethodPost:
		// Some collections update via PUT with the id in the body.
		if id, ok := body["id"].(string); ok && id != "" {
			f.handleItem(w, &http.Request{Method: http.MethodPatch}, coll+"/"+id, body)
			return
		}
		now := f.tick()
		obj := map[string]any{}
		for k, v := range body {
			obj[k] = v
		}
		obj["id"] = f.nextID(strings.TrimPrefix(coll[strings.LastIndex(coll, "/"):], "/"))
		obj["created_at"] = now
		obj["updated_at"] = now
		f.objects[coll+"/"+obj["id"].(string)] = obj
		writeJSON(w, obj)
	default:
		writeStatus(w, http.StatusMethodNotAllowed)
	}
}

func (f *fakeAPI) handleItem(w http.ResponseWriter, r *http.Request, key string, body map[string]any) {
	obj, ok := f.objects[key]
	if !ok {
		writeStatus(w, http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, obj)
	case http.MethodPatch:
		for k, v := range body {
			if k != "id" {
				obj[k] = v
			}
		}
		obj["updated_at"] = f.tick()
		writeJSON(w, obj)
	case http.MethodDelete:
		delete(f.objects, key)
		writeStatus(w, http.StatusNoContent)
	default:
		writeStatus(w, http.StatusMethodNotAllowed)
	}
}

func (f *fakeAPI) handleOrgKey(w http.ResponseWriter, r *http.Request) {
	if f.orgKey == nil {
		writeStatus(w, http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, f.orgKey)
	case http.MethodDelete:
		f.orgKey = nil
		writeStatus(w, http.StatusNoContent)
	default:
		writeStatus(w, http.StatusMethodNotAllowed)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func writeStatus(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
