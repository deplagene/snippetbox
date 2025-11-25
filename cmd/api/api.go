package api

import (
	"io/fs"
	"net/http"

	"github.com/deplagene/snippetbox"
	"github.com/deplagene/snippetbox/internal/database"
	"github.com/deplagene/snippetbox/internal/snippets"
	"github.com/deplagene/snippetbox/internal/users"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

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
	userSvc := users.NewService(querier)
	userRouter := users.NewHandler(userSvc)

	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

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
	userRouter.RegisterRoutes(router.Group("/user"))

	return router.Run(a.Addr)
}
