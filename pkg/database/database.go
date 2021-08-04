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

func GetAllBlogs() blog.Blogs {
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
	return ret
}

func GetBlogbyID(id int) (b blog.Blog, err error) {
	err = db.QueryRow("SELECT id, title, content FROM blogs where id = ?", id).Scan(&b.ID, &b.Title, &b.Content)
	if err != nil {
		return b, err
	}
	return b, nil
}

func InsertintoBlogs(data blog.Blog) error {
	query := fmt.Sprintf("INSERT INTO blogs(title, content) VALUES ( '%s', '%s' )", data.Title, data.Content)
	insert, err := db.Query(query)
	if err != nil {
		return err
	}
	defer insert.Close()
	return nil
}

func UpdateBlogbyID(data blog.Blog) error {
	_, err := db.Exec("Update blogs set title = ?, content = ? where id = ?", data.Title, data.Content, data.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteBlogbyID(id int) error {
	_, err := db.Exec("delete from blogs where id = ?", id)
	return err;
}
