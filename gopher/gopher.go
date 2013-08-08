package gopher

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ALL_QUERY = "all"
)

type Package struct {
	Name        string
	Repo        string
	Description string
	IsLibrary   bool
	Category    string
	Accepted    bool
	Added       time.Time
	Updated     time.Time
}

type RemoveRequest struct {
	Repo   string
	Reason string
}

func init() {
	http.HandleFunc("/api/pkg", pkg)
}

func pkg(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c := appengine.NewContext(r)
		id := r.FormValue("repo")
		var jsonResult []byte
		if item, err := memcache.Get(c, id); err == memcache.ErrCacheMiss {
			q := datastore.NewQuery("Package").Filter("Accepted =", true)
			if id != ALL_QUERY {
				q = q.Filter("Repo =", id)
			}
			q = q.Order("Name")
			var packages []*Package
			if _, err := q.GetAll(c, &packages); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var result interface{}
			if id != ALL_QUERY && len(packages) > 0 {
				result = packages[0]
			} else {
				result = packages
			}
			if jsonResult, err = json.Marshal(result); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			memcache.Add(c, &memcache.Item{Key: id, Value: jsonResult})
		} else {
			jsonResult = item.Value
		}
		fmt.Fprint(w, string(jsonResult))
	case "POST":
		c := appengine.NewContext(r)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		p := &Package{Added: time.Now(), Updated: time.Now()}
		if err := json.Unmarshal(body, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Package", nil), p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, `{"status": "ok"}`)
	case "DELETE":
		c := appengine.NewContext(r)
		rr := &RemoveRequest{
			Repo:   r.FormValue("repo"),
			Reason: r.FormValue("reason"),
		}
		if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "RemoveRequest", nil), rr); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, `{"status": "ok"}`)
	}

}
