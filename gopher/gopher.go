package gopher

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/user"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ALL_QUERY = "all"
)

type Package struct {
	Name        string
	Repo        string
	Description string
	Accepted    bool
}

type RemoveRequest struct {
	Repo   string
	Reason string
}

func init() {
	http.HandleFunc("/api/pkg", pkg)
	http.HandleFunc("/admin/", admin)
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
			var packages []Package
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
		p := &Package{}
		if err := json.Unmarshal(body, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Package", nil), p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		memcache.Delete(c, p.Repo)
		memcache.Delete(c, ALL_QUERY)
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

func admin(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Hello, %v!", u)
}
