package api

import (
	"deplagene/snippetbox"
	"deplagene/snippetbox/internal/database"
	"deplagene/snippetbox/internal/snippets"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type API struct {
	Addr string
	db   *pgxpool.Pool
}

func NewAPI(addr string, db *pgxpool.Pool) *API {
	return &API{
		Addr: addr,
		db:   db,
	}
}

func (a *API) Run() error {
	querier := database.New(a.db)
	snippetSvc := snippets.NewService(querier)
	snippetRouter := snippets.NewHandler(snippetSvc)

	router := gin.Default()

	templates, err := NewTemplateRenderer()
	if err != nil {
		return err
	}
	router.HTMLRender = templates

	staticFS, err := fs.Sub(snippetbox.UIFS, "ui/static")
	if err != nil {
		return err
	}
	router.StaticFS("/static", http.FS(staticFS))

	snippetRouter.RegisterRoutes(router.Group("/"))

	return router.Run(a.Addr)
}
