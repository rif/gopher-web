package gopher

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/user"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
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
	http.HandleFunc("/admin/", admin)
	http.HandleFunc("/admin/accept/", accept)
	http.HandleFunc("/admin/reject/", reject)
	http.HandleFunc("/admin/acceptremoval/", acceptremoval)
	http.HandleFunc("/admin/rejectremoval/", rejectremoval)
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
	t, err := template.ParseFiles("app/base.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	url, _ := user.LogoutURL(c, "/")

	acceptQuery := datastore.NewQuery("Package").Filter("Accepted =", false).Order("Added")
	var packages []*Package
	keys, err := acceptQuery.GetAll(c, &packages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	removeQuery := datastore.NewQuery("RemoveRequest")
	var removerequests []*RemoveRequest
	removeKeys, err := removeQuery.GetAll(c, &removerequests)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		"user":           u,
		"url":            url,
		"pacakges":       packages,
		"keys":           keys,
		"removeRequests": removerequests,
		"removeKeys":     removeKeys,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func accept(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := r.FormValue("id")
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pkg := &Package{}
	key := datastore.NewKey(c, "Package", "", intID, nil)
	err = datastore.Get(c, key, pkg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// check if the package is already present
	acceptQuery := datastore.NewQuery("Package").
		Filter("Accepted =", true).
		Filter("Repo =", pkg.Repo)

	var packages []*Package
	keys, err := acceptQuery.GetAll(c, &packages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(packages) > 1 {
		// just print an error to let admin know
		c.Errorf("More tha one package for repo: %v", pkg.Repo)
	}

	if len(packages) > 0 {
		// update the package and delete
		oldKey := keys[0]
		oldPkg := packages[0]
		oldPkg.Name = pkg.Name
		oldPkg.Description = pkg.Description
		oldPkg.IsLibrary = pkg.IsLibrary
		oldPkg.Category = pkg.Category
		oldPkg.Updated = time.Now()
		if _, err = datastore.Put(c, oldKey, oldPkg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = datastore.Delete(c, key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// accept the new package
		pkg.Accepted = true
		if _, err = datastore.Put(c, key, pkg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	memcache.Delete(c, pkg.Repo)
	memcache.Delete(c, ALL_QUERY)
}

func reject(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := r.FormValue("id")
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	key := datastore.NewKey(c, "Package", "", intID, nil)
	if err := datastore.Delete(c, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func acceptremoval(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	repo := r.FormValue("repo")

	c.Errorf("REPO: %v", repo)

	removeQuery := datastore.NewQuery("Package").
		Filter("Accepted =", true).
		Filter("Repo = ", repo).
		KeysOnly()

	keys, err := removeQuery.GetAll(c, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Errorf("FOUND: %v", keys)
	for _, key := range keys {
		c.Errorf("DELETING: %v", key)
		if err := datastore.Delete(c, key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// delete removal request
	id := r.FormValue("id")
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	key := datastore.NewKey(c, "RemoveRequest", "", intID, nil)
	if err := datastore.Delete(c, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	memcache.Delete(c, repo)
	memcache.Delete(c, ALL_QUERY)
}

func rejectremoval(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := r.FormValue("id")
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	key := datastore.NewKey(c, "RemoveRequest", "", intID, nil)
	if err := datastore.Delete(c, key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
