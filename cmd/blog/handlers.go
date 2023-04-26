package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type indexPageData struct {
	Title    string
	Subtitle string
	RPosts   []*RecentPosts
	FPosts   []*FeaturedPosts
}

type FeaturedPosts struct {
	Title       string `db:"title"`
	Subtitle    string `db:"subtitle"`
	PublishDate string `db:"publish_date"`
	Author      string `db:"author"`
	AuthorUrl   string `db:"author_url"`
	ImageUrl    string `db:"image_url"`
	Featured    string `db:"featured"`
	PostID      string `db:"post_id"`
	PostURL     string // URL ордера, на который мы будем переходить для конкретного поста
}

type RecentPosts struct {
	Title       string `db:"title"`
	Subtitle    string `db:"subtitle"`
	PublishDate string `db:"publish_date"`
	Author      string `db:"author"`
	AuthorUrl   string `db:"author_url"`
	ImageUrl    string `db:"image_url"`
	Featured    string `db:"featured"`
	PostID      string `db:"post_id"`
	PostURL     string // URL ордера, на который мы будем переходить для конкретного поста
}

type PostData struct {
	Title    string `db:"title"`
	Subtitle string `db:"subtitle"`
	ImageUrl string `db:"image_url"`
	Content  string `db:"content"`
}

func index(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// posts, err := orders(db)
		fposts, err := fposts(db)
		rposts, err := rposts(db)
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
			Title:    "Blog for traveling",
			Subtitle: "My best blog for adventures and burgers",
			RPosts:   rposts,
			FPosts:   fposts,
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

func order(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orderIDStr := mux.Vars(r)["PostID"] // Получаем orderID в виде строки из параметров урла

		orderID, err := strconv.Atoi(orderIDStr) // Конвертируем строку orderID в число
		if err != nil {
			http.Error(w, "Invalid order id", 403)
			log.Println(err)
			return
		}

		order, err := orderByID(db, orderID)
		if err != nil {
			if err == sql.ErrNoRows {
				// sql.ErrNoRows возвращается, когда в запросе к базе не было ничего найдено
				// В таком случае мы возвращем 404 (not found) и пишем в тело, что ордер не найден
				http.Error(w, "Order not found", 404)
				log.Println(err)
				return
			}

			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		ts, err := template.ParseFiles("pages/post.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		err = ts.Execute(w, order)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")
	}
}

func fposts(db *sqlx.DB) ([]*FeaturedPosts, error) {
	const query = `
		SELECT
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url,
			featured,
			post_id
		FROM 
			post
		WHERE featured = 1
	`
	// Такое объединение строк делается только для таблицы order, т.к. это зарезерированное слово в SQL, наряду с SELECT, поэтому его нужно заключить в ``

	var orders []*FeaturedPosts // Заранее объявляем массив с результирующей информацией

	err := db.Select(&orders, query) // Делаем запрос в базу данных
	if err != nil {                  // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	for _, order := range orders {
		order.PostURL = "/post/" + order.PostID // Формируем исходя из ID ордера в базе
	}

	return orders, nil
}

func rposts(db *sqlx.DB) ([]*RecentPosts, error) {
	const query = `
		SELECT
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url,
			featured,
			post_id
		FROM 
			post
		WHERE featured = 0
	`
	// Такое объединение строк делается только для таблицы order, т.к. это зарезерированное слово в SQL, наряду с SELECT, поэтому его нужно заключить в ``

	var orders []*RecentPosts // Заранее объявляем массив с результирующей информацией

	err := db.Select(&orders, query) // Делаем запрос в базу данных
	if err != nil {                  // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	for _, order := range orders {
		order.PostURL = "/post/" + order.PostID // Формируем исходя из ID ордера в базе
	}

	return orders, nil
}

func orderByID(db *sqlx.DB, orderID int) (PostData, error) {
	const query = `
		SELECT
			title,
			subtitle,
			image_url,
			content
		FROM
			post
		WHERE
			post_id = ?
	`
	// В SQL-запросе добавились параметры, как в шаблоне. ? означает параметр, который мы передаем в запрос ниже

	var order PostData

	// Обязательно нужно передать в параметрах orderID
	err := db.Get(&order, query, orderID)
	if err != nil {
		return PostData{}, err
	}

	return order, nil
}
