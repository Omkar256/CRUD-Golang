package routes

import (
	blog "crud-golang/internal/blogs"
	"crud-golang/pkg/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
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
	w.Write([]byte("This is Golang CRUD app"))
}

func createBlog(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpan("createBlog")
	defer span.Finish()

	var data blog.Blog
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	error_channel := make(chan error)
	go database.InsertintoBlogs(data, error_channel)
	err = <-error_channel
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
}

func getBlogs(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpan("getBlogs")
	defer span.Finish()

	dataChannel := make(chan blog.Blogs)
	go database.GetAllBlogs(dataChannel)
	data := <-dataChannel
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getBlogbyID(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpan("getBlogsbyID")
	defer span.Finish()

	id, err := parseId(r)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	dataChannel := make(chan blog.Blog)
	errChannel := make(chan error)
	go database.GetBlogbyID(id, dataChannel, errChannel)
	select {
	case data := <-dataChannel:
		json.NewEncoder(w).Encode(data)
		return
	case err := <-errChannel:
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
}

func updateBlogbyID(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpan("updateBlogbyID")
	defer span.Finish()

	data := blog.Blog{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	errorChannel := make(chan error)
	go database.UpdateBlogbyID(data, errorChannel)
	err = <-errorChannel
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(200)
}

func deleteBlogbyID(w http.ResponseWriter, r *http.Request) {
	span := tracer.StartSpan("deleteBlogbyID")
	defer span.Finish()

	id, err := parseId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	errorChannel := make(chan error)
	go database.DeleteBlogbyID(id, errorChannel)
	err = <-errorChannel

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
