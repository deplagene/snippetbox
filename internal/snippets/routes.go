package snippets

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s: s}
}

func (r *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/", r.home)
	rg.GET("/snippet", r.snippetView)
	rg.GET("/snippet/create", r.snippetCreateForm)
	rg.POST("/snippet/create", r.snippetCreate)
}

func (r *Handler) home(c *gin.Context) {
	snippets, err := r.s.GetLatest()
	if err != nil {
		log.Printf("Error getting latest snippets: %v", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.HTML(http.StatusOK, "home-page.html", gin.H{
		"Snippets": snippets,
	})
}

func (r *Handler) snippetView(c *gin.Context) {
	idStr := c.Query("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid snippet ID")
		return
	}

	snippet, err := r.s.GetById(id)
	if err != nil {
		log.Printf("Error getting snippet by id %s: %v", idStr, err)
		c.String(http.StatusNotFound, "Snippet not found")
		return
	}

	c.HTML(http.StatusOK, "show-page.html", gin.H{
		"Snippet": snippet,
	})
}

func (r *Handler) snippetCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create-page.html", gin.H{
		"Form": &SnippetCreateForm{},
	})
}

func (r *Handler) snippetCreate(c *gin.Context) {
	var form SnippetCreateForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "Invalid form data")
		return
	}

	if !form.IsValid() {
		c.HTML(http.StatusOK, "create-page.html", gin.H{
			"Form": form,
		})
		return
	}

	expires, _ := strconv.Atoi(form.Expires)

	snippetID, err := r.s.Create(form.Title, form.Content, expires)
	if err != nil {
		log.Printf("Error creating snippet: %v", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Redirect(http.StatusSeeOther, "/snippet?id="+snippetID.String())
}
