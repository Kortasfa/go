package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // Импортируем для возможности подключения к MySQL
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

const (
	port         = ":3000"
	dbDriverName = "mysql"
)

func main() {
	db, err := openDB() // Открываем соединение к базе данных в самом начале
	if err != nil {
		log.Fatal(err)
	}

	dbx := sqlx.NewDb(db, dbDriverName) // Расширяем стандартный клиент к базе

	// Обязательно подключить github.com/gorilla/router в импортах
	router := mux.NewRouter()
	router.HandleFunc("/home", index(dbx)) // Передаём клиент к базе данных в ф-ию обработчик запроса

	router.HandleFunc("/admin", admin()) // Страничка админа, бд не нужно

	router.HandleFunc("/login", login()) // Страничка админа, бд не нужно

	router.HandleFunc("/api/post", createPost(dbx)).Methods(http.MethodPost) //Обработка json запроса

	// Указываем orderID поста в URL для перехода на конкретный пост
	router.HandleFunc("/post/{PostID}", post(dbx))

	// Правим отдачу статического контента, ввиду переезда на новый роутер
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Start server")
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}

func openDB() (*sql.DB, error) {
	// Здесь прописываем соединение к базе данных
	return sql.Open(dbDriverName, "root:root123321@tcp(localhost:3306)/blog?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true")
}
