package main

import (
	"html/template"
	"log"
	"net/http"
	"github.com/jmoiron/sqlx"
)

type indexPageData struct {
	Title         string
	Subtitle      string
	FeaturedPosts []FeaturedPosts
	RecentPosts []RecentPosts
}

type FeaturedPosts struct {
	Title       string `db:"title"`
	Subtitle    string `db:"subtitle"`
	PublishDate string `db:"publish_date"`
	Author      string `db:"author"`
	AuthorUrl   string `db:"author_url"`
	ImageUrl    string `db:"image_url"`
	Featured    string `db:"featured"`
}

type RecentPosts struct {
	Title       string `db:"title"`
	Subtitle    string `db:"subtitle"`
	PublishDate string `db:"publish_date"`
	Author      string `db:"author"`
	AuthorUrl   string `db:"author_url"`
	ImageUrl    string `db:"image_url"`
	Featured    string `db:"featured"`
}


func index(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := featuredPosts(db)
		rposts, err := mostRecent(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
			log.Println(err)
			return // Не забываем завершить выполнение ф-ии
		}

		ts, err := template.ParseFiles("pages/index.html") // Главная страница блога
		if err != nil {
			http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
			log.Println(err)
			return // Не забываем завершить выполнение ф-ии
		}

		data := indexPageData{
			Title:         "Blog for traveling",
			Subtitle:      "My best blog for adventures and burgers",
			FeaturedPosts: posts,
			RecentPosts: rposts,
		}

		err = ts.Execute(w, data) // Заставляем шаблонизатор вывести шаблон в тело ответа
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")
	}
}

func featuredPosts(db *sqlx.DB) ([]FeaturedPosts, error) {
	const query = `
		SELECT
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url,
			featured
		FROM
			post
		WHERE featured = 1
	` // Составляем SQL-запрос для получения записей для секции featured-posts

	var posts []FeaturedPosts // Заранее объявляем массив с результирующей информацией

	err := db.Select(&posts, query) // Делаем запрос в базу данных
	if err != nil {                 // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	return posts, nil
}

func mostRecent(db *sqlx.DB) ([]RecentPosts, error) {
	const query = `
		SELECT
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url,
			featured
		FROM
			post
		WHERE featured = 0
	` // Составляем SQL-запрос для получения записей для секции featured-posts

	var rposts []RecentPosts // Заранее объявляем массив с результирующей информацией

	err := db.Select(&rposts, query) // Делаем запрос в базу данных
	if err != nil {                 // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	return rposts, nil
}
