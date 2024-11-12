package main

import (
	"deplagene/snippetbox/pkg/repository"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	Addr, StaticDir string
}

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	repository *repository.PGRepository
	cacheTemplate map[string]*template.Template
}

// Подключение к базе данных
const connectionString string = "postgres://postgres:postgres@localhost:5432/snippetbox_db"

func main() {
	// Конфигурация приложения
	cfg := new(Config)

	flag.StringVar(
		&cfg.Addr,
		"addr",
		":4000",
		"Сетевой адрес HTTP",
	)

	flag.StringVar(
		&cfg.StaticDir,
		"static-dir",
		"./ui/static",
		"Папка для статических файлов",
	)

	flag.Parse()

	// Логирование ( многоуровневое )
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, connErr := repository.New(connectionString)
	if connErr != nil {
		errorLog.Fatal(connErr)
	}

	cacheTemplate, cacheErr := newTemplateCache("./ui/html/")
	if cacheErr != nil {
		errorLog.Fatal(cacheErr)
	}
	
	app := application {
		errorLog: errorLog,
		infoLog: infoLog,
		repository: db,
		cacheTemplate: cacheTemplate,
 	}

	srv := &http.Server {
		Addr: cfg.Addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Запуск сервера на %s", cfg.Addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

// Open открывает файл path из файловой системы fs. Если path указывает
// на директорию, то Open проверяет, существует ли файл index.html в
// этой директории. Если файл не существует, то возвращается ошибка.
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}
	return f, nil
}