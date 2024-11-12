package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// serverError записывает информацию об ошибке в лог и возвращает
// ошибку 500 Internal Server Error в HTTP-ответе
func (app *application) serverError(w http.ResponseWriter, err error)  {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError отправляет ошибку клиенту, используя подходящий HTTP-код
func (app *application) clientError(w http.ResponseWriter, code int)  {
	http.Error(w, http.StatusText(code), code)
}

// notFound отправляет ошибку 404 Not Found в HTTP-ответе
func (app *application) notFound(w http.ResponseWriter)  {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.cacheTemplate[name]
	if !ok {
		app.serverError(w, fmt.Errorf("Шаблона %s не существует", name))
		return
	}

	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}