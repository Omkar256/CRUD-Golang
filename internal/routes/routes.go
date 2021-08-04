package routes

import (
	blog "crud-golang/internal/blogs"
	"crud-golang/pkg/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func INIT_Routes() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/create", createBlog).Methods("POST")
	router.HandleFunc("/get", getBlogs).Methods("GET")
	router.HandleFunc("/get/{id}", getBlogbyID).Methods("GET")
	router.HandleFunc("/update", updateBlogbyID).Methods("PUT")
	router.HandleFunc("/delete/{id}", deleteBlogbyID).Methods("DELETE")

	http.ListenAndServe(":3333", router)

}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is basic CRUD app"))
}

func createBlog(w http.ResponseWriter, r *http.Request) {
	var data blog.Blog
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.InsertintoBlogs(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
}

func getBlogs(w http.ResponseWriter, r *http.Request) {
	data := database.GetAllBlogs()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getBlogbyID(w http.ResponseWriter, r *http.Request) {
	id, err := parseId(r)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	data, err := database.GetBlogbyID(id)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Row not found",
		})
		return
	}
	json.NewEncoder(w).Encode(data)
}

func updateBlogbyID(w http.ResponseWriter, r *http.Request) {
	data := blog.Blog{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateBlogbyID(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
}

func deleteBlogbyID(w http.ResponseWriter, r *http.Request) {
	id, err := parseId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteBlogbyID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
}

func parseId(r *http.Request) (id int, err error) {
	vars := mux.Vars(r)
	id, err = strconv.Atoi(vars["id"])
	if err != nil {
		return -1, err
	}
	return id, nil
}
