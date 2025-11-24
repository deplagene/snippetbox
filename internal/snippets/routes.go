package snippets

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Router handles the snippet routes
type Router struct {
	s *Service
}

// NewRouter creates a new snippet router
func NewRouter(s *Service) *Router {
	return &Router{s: s}
}

// RegisterRoutes registers the snippet routes with the Gin router
func (r *Router) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/", r.home)
	rg.GET("/snippet/view", r.snippetView) // query param: ?id=...
	rg.GET("/snippet/create", r.snippetCreateForm)
	rg.POST("/snippet/create", r.snippetCreate)
}

func (r *Router) home(c *gin.Context) {
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

func (r *Router) snippetView(c *gin.Context) {
	idStr := c.Query("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid snippet ID")
		return
	}

	snippet, err := r.s.GetById(id)
	if err != nil {
		log.Printf("Error getting snippet by id %s: %v", id, err)
		// A more specific error check for "not found" would be better,
		// but for now, this is sufficient.
		c.String(http.StatusNotFound, "Snippet not found")
		return
	}

	c.HTML(http.StatusOK, "show-page.html", gin.H{
		"Snippet": snippet,
	})
}

func (r *Router) snippetCreateForm(c *gin.Context) {
	// For simplicity, we don't have a dedicated create form page yet.
	// We can add one later if needed. For now, we'll just show the home page.
	// Or better, let's create a simple one.
	c.HTML(http.StatusOK, "create-page.html", nil)
}

func (r *Router) snippetCreate(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	expiresStr := c.PostForm("expires")

	if title == "" || content == "" || expiresStr == "" {
		c.String(http.StatusBadRequest, "Title, content, and expires are required")
		return
	}

	expires, err := strconv.Atoi(expiresStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid expires value, must be a number of days")
		return
	}

	snippet, err := r.s.Create(title, content, expires)
	if err != nil {
		log.Printf("Error creating snippet: %v", err)
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Redirect(http.StatusSeeOther, "/snippet/view?id="+snippet.SnippetID.String())
}
