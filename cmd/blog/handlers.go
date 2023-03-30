package main

import (
	"html/template"
	"log"
	"net/http"
)

type indexPage struct {
	FeaturedPosts []featuredPostData
	MostRecent    []mostRecentData
}

type postPage struct {
}

type featuredPostData struct {
	Title                   string
	Subtitle                string
	BackgroundImageModifier string
	AuthorAvatar            string
	AuthorName              string
	PublishDate             string
	PostButton              string
	ButtonClass             string
}

type mostRecentData struct {
	PhotoImage   string
	Title        string
	Subtitle     string
	AuthorName   string
	AuthorAvatar string
	PublishDate  string
}

func index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/index.html") // Главная страница блога
	if err != nil {
		http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
		log.Println(err.Error())                    // Используем стандартный логгер для вывода ошбики в консоль
		return                                      // Не забываем завершить выполнение ф-ии
	}

	data := indexPage{
		FeaturedPosts: featuredPosts(),
		MostRecent:    mostRecent(),
	}

	err = ts.Execute(w, data) // Заставляем шаблонизатор вывести шаблон в тело ответа
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/post.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	data := postPage{}

	err = ts.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func featuredPosts() []featuredPostData {
	return []featuredPostData{
		{
			Title:                   "The Road Ahead",
			Subtitle:                "The road ahead might be paved - it might not be.",
			BackgroundImageModifier: "featured__post_background-road",
			AuthorAvatar:            "/static/images/mat_vogels.jpg",
			AuthorName:              "Mat Vogels",
			PublishDate:             "September 25, 2015",
		},
		{
			Title:                   "From Top Down",
			Subtitle:                "Once a year, go someplace you’ve never been before.",
			BackgroundImageModifier: "featured__post_background-fromtop",
			AuthorAvatar:            "/static/images/william_wong.jpg",
			AuthorName:              "William Wong",
			PublishDate:             "September 25, 2015",
			PostButton:              "ADVENTURE",
			ButtonClass:             "post__button",
		},
	}
}

func mostRecent() []mostRecentData {
	return []mostRecentData{
		{
			PhotoImage:   "/static/images/still_standing_tall.jpg",
			Title:        "Still Standing Tall",
			Subtitle:     "Life begins at the end of your comfort zone.",
			AuthorName:   "William Wong",
			AuthorAvatar: "/static/images/william_wong.jpg",
			PublishDate:  "9/25/2015",
		},
		{
			PhotoImage:   "/static/images/sunny_side_up.jpg",
			Title:        "Sunny Side Up",
			Subtitle:     "No place is ever as bad as they tell you it’s going to be.",
			AuthorName:   "Mat Vogels",
			AuthorAvatar: "/static/images/mat_vogels.jpg",
			PublishDate:  "9/25/2015",
		},
		{
			PhotoImage:   "/static/images/water_falls.jpg",
			Title:        "Water Falls",
			Subtitle:     "We travel not to escape life, but for life not to escape us.",
			AuthorName:   "William Wong",
			AuthorAvatar: "/static/images/william_wong.jpg",
			PublishDate:  "9/25/2015",
		},
		{
			PhotoImage:   "/static/images/through_the_mist.jpg",
			Title:        "Through the Mist",
			Subtitle:     "Travel makes you see what a tiny place you occupy in the world.",
			AuthorName:   "William Wong",
			AuthorAvatar: "/static/images/william_wong.jpg",
			PublishDate:  "9/25/2015",
		},
		{
			PhotoImage:   "/static/images/still_standing_tall.jpg",
			Title:        "Still Standing Tall",
			Subtitle:     "Life begins at the end of your comfort zone.",
			AuthorName:   "William Wong",
			AuthorAvatar: "../../static/images/william_wong.jpg",
			PublishDate:  "9/25/2015",
		},
		{
			PhotoImage:   "/static/images/awaken_early.jpg",
			Title:        "Awaken Early",
			Subtitle:     "Not all those who wander are lost.",
			AuthorName:   "Mat Vogels",
			AuthorAvatar: "/static/images/mat_vogels.jpg",
			PublishDate:  "9/25/2015",
		},
	}
}
