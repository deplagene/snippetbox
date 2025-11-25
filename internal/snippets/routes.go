package snippets

import (
	"log"
	"net/http"
	"strconv"

	"github.com/deplagene/snippetbox/internal/database"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	s SnippetService
}

func NewHandler(s SnippetService) *Handler {
	return &Handler{s: s}
}

func (r *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/", r.home)
	rg.GET("/snippet", r.snippetView)
	rg.GET("/snippet/create", r.snippetCreateForm)
	rg.POST("/snippet/create", r.snippetCreate)
	rg.POST("/snippet/delete", r.snippetDelete)
}

func (r *Handler) home(c *gin.Context) {
	session := sessions.Default(c)
	userID, ok := session.Get("userID").(uuid.UUID)

	var snippets []database.Snippet
	var err error

	if ok {
		snippets, err = r.s.GetLatestForUser(userID)
	} else {
		snippets, err = r.s.GetLatest()
	}

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

	session := sessions.Default(c)
	userID, ok := session.Get("userID").(uuid.UUID)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/user/login")
		return
	}

	snippetID, err := r.s.Create(form.Title, form.Content, userID, expires)
	if err != nil {
		log.Printf("Error creating snippet: %v", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Redirect(http.StatusSeeOther, "/snippet?id="+snippetID.String())
}

func (r *Handler) snippetDelete(c *gin.Context) {
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

	session := sessions.Default(c)
	userID, ok := session.Get("userID").(uuid.UUID)
	if !ok || snippet.UserID.Bytes != userID {
		c.String(http.StatusForbidden, "You are not authorized to delete this snippet")
		return
	}

	if err := r.s.Delete(id); err != nil {
		log.Printf("Error deleting snippet: %v", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}