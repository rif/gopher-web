package gopher

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"net/http"
)

type Package struct {
	Name        string
	Url         string
	Description string
	Accepted    bool
}

type UpdateRequest struct {
	Name        string
	Url         string
	Description string
	Accepted    bool
}

type RemoveRequest struct {
	Url      string
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
	id := r.FormValue("packageId")
	q := datastore.NewQuery("Package")//.Filter("Accepted =", true)
	if id != "" {
	    q.Filter("Name =", id)
	}
	q.Order("Name")
	var packages []Package
	if _, err := q.GetAll(c, &packages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(packages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func add(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	p := Package{
		Name:        r.FormValue("name"),
		Url:         r.FormValue("url"),
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
