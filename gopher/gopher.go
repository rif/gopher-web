package gopher

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"encoding/json"
	"fmt"
	"net/http"
)

type Package struct {
	Name        string
	Repo        string
	Description string
	Accepted    bool
}

type UpdateRequest struct {
	Name        string
	Repo        string
	Description string
	Accepted    bool
}

type RemoveRequest struct {
	Repo     string
	Reason   string
	Accepted bool
}

func init() {
	http.HandleFunc("/api/query", query)
	http.HandleFunc("/api/add", add)
	http.HandleFunc("/api/remove", remove)
}

func query(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := r.FormValue("repo")
	var jsonResult []byte
	if item, err := memcache.Get(c, id); err == memcache.ErrCacheMiss {
		q := datastore.NewQuery("Package").Filter("Accepted =", true)
		if id != "all" {
			q.Filter("Repo =", id)
		}
		q.Order("Name")
		var packages []Package
		if _, err := q.GetAll(c, &packages); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var result interface{}
		if id != "all" && len(packages) > 0 {
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

}

func add(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	p := Package{
		Name:        r.FormValue("name"),
		Repo:        r.FormValue("repo"),
		Description: r.FormValue("description"),
		Accepted:    false,
	}
	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Package", nil), &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "OK")
}

func remove(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	g := Package{
		Name: r.FormValue("content"),
	}
	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Package", nil), &g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
