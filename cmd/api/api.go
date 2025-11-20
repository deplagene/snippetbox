package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type API struct {
	Port string
	db   *pgxpool.Pool
}

func NewAPI(port string, db *pgxpool.Pool) *API {
	return &API{
		Port: port,
		db:   db,
	}
}

func (a *API) Run() error {
	router := gin.Default()
	return router.Run(a.Port)
}
