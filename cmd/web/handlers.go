package main

import (
	"deplagene/snippetbox/pkg/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// home обрабатывает запросы на главную страницу.
//
// Он проверяет, является ли запрос на главную страницу,
// и если это не так, то возвращает код 404.
// Затем он парсит файлы HTML и генерирует текст HTML
// с помощью функции Execute.
func(app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return // обязательно.
	}

	s, err := app.repository.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home-page.html", &templateData{
		Snippets: s,
	})
}

// showSnippet обрабатывает запрос на страницу отдельной заметки.
//
// Он извлекает id из GET-параметра "id" и
// проверяет, является ли id положительным числом.
// Если id не является положительным числом, то
// возвращает код 404.
// Затем он выводит текст, содержащий id
// заметки.
func(app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return // обязательно.
	}

	s, err := app.repository.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show-page.html", &templateData{
		Snippet: s,
	})
}

// createSnippet обрабатывает запрос на создание заметки.
// Он позволяет только POST-запросы.
func(app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "Тестовая заметка"
	content := "Шла Ева по шоссе и сосала сушку."
	expires := "321"
 
	id, err := app.repository.Create(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
