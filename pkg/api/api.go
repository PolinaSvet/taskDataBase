package api

import (
	"GoNews/pkg/storage"
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

// Программный интерфейс сервера GoNews
type API struct {
	db     storage.Interface
	router *mux.Router
}

// Конструктор объекта API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {

	api.router.HandleFunc("/", api.templateHandler).Methods(http.MethodGet, http.MethodOptions)

	api.router.HandleFunc("/posts", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/posts", api.addPostHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/posts", api.updatePostHandler).Methods(http.MethodPut, http.MethodOptions)
	api.router.HandleFunc("/posts", api.deletePostHandler).Methods(http.MethodDelete, http.MethodOptions)

	api.router.HandleFunc("/authors", api.authorsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/authors", api.addAuthorHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/authors", api.updateAuthorHandler).Methods(http.MethodPut, http.MethodOptions)
	api.router.HandleFunc("/authors", api.deleteAuthorHandler).Methods(http.MethodDelete, http.MethodOptions)

	// Регистрация обработчика для статических файлов (шаблонов)
	api.router.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("ui"))))
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.router
}

// Базовый маршрут.
func (api *API) templateHandler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("ui/html/base.html", "ui/html/routes.html"))

	// Отправляем HTML страницу с данными
	if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Получение всех публикаций.
func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {

	posts, err := api.db.Posts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

// 1) Post
// Добавление публикации.
func (api *API) addPostHandler(w http.ResponseWriter, r *http.Request) {

	var p storage.Post
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = api.db.AddPost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Обновление публикации.
func (api *API) updatePostHandler(w http.ResponseWriter, r *http.Request) {

	var p storage.Post
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = api.db.UpdatePost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Удаление публикации.
func (api *API) deletePostHandler(w http.ResponseWriter, r *http.Request) {

	var p storage.Post
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = api.db.DeletePost(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// 2) Authors
// Получение всех авторов.
func (api *API) authorsHandler(w http.ResponseWriter, r *http.Request) {

	authors, err := api.db.Authors()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(authors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

// Добавление автора.
func (api *API) addAuthorHandler(w http.ResponseWriter, r *http.Request) {

	var p storage.Author
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = api.db.AddAuthor(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Обновление автора.
func (api *API) updateAuthorHandler(w http.ResponseWriter, r *http.Request) {

	var p storage.Author
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = api.db.UpdateAuthor(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Удаление автора.
func (api *API) deleteAuthorHandler(w http.ResponseWriter, r *http.Request) {

	var p storage.Author
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = api.db.DeleteAuthor(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
