package database

import (
	blog "crud-golang/internal/blogs"
	"database/sql"
	"fmt"
)

var db *sql.DB

func SetDB(d *sql.DB) {
	db = d
}

func GetAllBlogs(dataChannel chan blog.Blogs) {
	results, err := db.Query("SELECT id, title, content FROM blogs")
	if err != nil {
		panic(err.Error())
	}
	ret := make([]blog.Blog, 0)
	for results.Next() {
		var b blog.Blog
		err = results.Scan(&b.ID, &b.Title, &b.Content)
		if err != nil {
			panic(err.Error())
		}
		ret = append(ret, b)
	}
	dataChannel <- ret
}

func GetBlogbyID(id int, dataChannel chan blog.Blog, errChannel chan error) {
	var b blog.Blog
	err := db.QueryRow("SELECT id, title, content FROM blogs where id = ?", id).Scan(&b.ID, &b.Title, &b.Content)
	if err != nil {
		errChannel <- err
		return
	}
	dataChannel <- b
}

func InsertintoBlogs(data blog.Blog, error_channel chan error) {
	query := fmt.Sprintf("INSERT INTO blogs(title, content) VALUES ( '%s', '%s' )", data.Title, data.Content)
	insert, err := db.Query(query)
	error_channel <- err
	defer insert.Close()
}

func UpdateBlogbyID(data blog.Blog, errorChannel chan error) {
	_, err := db.Exec("Update blogs set title = ?, content = ? where id = ?", data.Title, data.Content, data.ID)
	errorChannel <- err
}

func DeleteBlogbyID(id int, errorChannel chan error) {
	_, err := db.Exec("DELETE from blogs where id = ?", id)
	errorChannel <- err
}
