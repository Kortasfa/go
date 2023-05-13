package main

import (
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"encoding/base64"
	"encoding/json" // Импортируем библиотеку для работы с JSON

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

type createPostRequest struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	PublishDate string `json:"publish_date"`
	Author      string `json:"author"`
	AuthorUrl   string `json:"author_url"`
	PostImg     string `json:"post_img"`
	Content     string `json:"content"`
}

type PostData struct {
	Title    string `db:"title"`
	Subtitle string `db:"subtitle"`
	ImageUrl string `db:"image_url"`
	Content  string `db:"content"`
}

func index(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// posts, err := posts(db)
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

func post(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := mux.Vars(r)["PostID"] // Получаем postID в виде строки из параметров урла

		postID, err := strconv.Atoi(postIDStr) // Конвертируем строку postID в число
		if err != nil {
			http.Error(w, "Invalid post id", 403)
			log.Println(err)
			return
		}

		post, err := postByID(db, postID)
		if err != nil {
			if err == sql.ErrNoRows {
				// sql.ErrNoRows возвращается, когда в запросе к базе не было ничего найдено
				// В таком случае мы возвращем 404 (not found) и пишем в тело, что ордер не найден
				http.Error(w, "post not found", 404)
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

		err = ts.Execute(w, post)
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
	// Такое объединение строк делается только для таблицы post, т.к. это зарезерированное слово в SQL, наряду с SELECT, поэтому его нужно заключить в ``

	var posts []*FeaturedPosts // Заранее объявляем массив с результирующей информацией

	err := db.Select(&posts, query) // Делаем запрос в базу данных
	if err != nil {                 // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	for _, post := range posts {
		post.PostURL = "/post/" + post.PostID // Формируем исходя из ID ордера в базе
	}

	return posts, nil
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
	// Такое объединение строк делается только для таблицы post, т.к. это зарезерированное слово в SQL, наряду с SELECT, поэтому его нужно заключить в ``

	var posts []*RecentPosts // Заранее объявляем массив с результирующей информацией

	err := db.Select(&posts, query) // Делаем запрос в базу данных
	if err != nil {                 // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	for _, post := range posts {
		post.PostURL = "/post/" + post.PostID // Формируем исходя из ID ордера в базе
	}

	return posts, nil
}

func postByID(db *sqlx.DB, postID int) (PostData, error) {
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

	var post PostData

	// Обязательно нужно передать в параметрах postID
	err := db.Get(&post, query, postID)
	if err != nil {
		return PostData{}, err
	}

	return post, nil
}

func admin() func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("pages/admin.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}
		err = ts.Execute(w, 0)
    }
}

func createPost(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        reqData, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Failed to read request body", http.StatusBadRequest)
            return
        }

        var req createPostRequest
        err = json.Unmarshal(reqData, &req)
        if err != nil {
            http.Error(w, "Failed to parse request body", http.StatusBadRequest)
            return
        }

        err = savePost(db, req)
        if err != nil {
            http.Error(w, "Failed to save post", http.StatusInternalServerError)
            return
        }

        // Handle successful post creation
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Post created successfully"))
    }
}

func savePost(db *sqlx.DB, req createPostRequest) error {
    const query = `
        INSERT INTO
            post
        (
            title,
            subtitle,
            publish_date,
            author,
            author_url,
            image_url,
            content
        )
        VALUES
        (
            ?, ?, ?, ?, ?, ?, ?
        )
    `
	
    _, err := db.Exec(query, req.Title, req.Subtitle, req.PublishDate, req.Author, req.AuthorUrl, req.Title, req.Content)

	img, err := base64.StdEncoding.DecodeString(req.PostImg)
	file, err := os.Create("static/img/" + req.Title) 
	// создаем файл с именем переданным от фронта в папке static/img
		
	_, err = file.Write(img) // Записываем контент картинки в файл

    return err
}
