package api

import (
	"deplagene/snippetbox"
	"deplagene/snippetbox/internal/database"
	"deplagene/snippetbox/internal/snippets"
	"html/template"
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
	// Initialize dependencies
	querier := database.New(a.db)
	snippetSvc := snippets.NewService(querier)
	snippetRouter := snippets.NewRouter(snippetSvc)

	// Set up router
	router := gin.Default()

	// Load templates
	templates, err := template.ParseFS(snippetbox.UIFS, "ui/html/*.html")
	if err != nil {
		return err
	}
	router.SetHTMLTemplate(templates)

	// Serve static files
	staticFS, err := fs.Sub(snippetbox.UIFS, "ui/static")
	if err != nil {
		return err
	}
	router.StaticFS("/static", http.FS(staticFS))

	// Register routes
	snippetRouter.RegisterRoutes(router.Group("/"))

	return router.Run(a.Addr)
}
