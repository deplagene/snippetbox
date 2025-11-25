package users

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	s UserService
}

func NewHandler(s UserService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/signup", h.SignupForm)
	rg.POST("/signup", h.Signup)
	rg.GET("/login", h.LoginForm)
	rg.POST("/login", h.Login)
}

func (h *Handler) SignupForm(c *gin.Context) {
	c.HTML(http.StatusOK, "signup-page.html", gin.H{
		"Form": &RegisterForm{},
	})
}

func (h *Handler) Signup(c *gin.Context) {
	var form RegisterForm

	log.Printf("%s", form)
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "Invalid form data")
		return
	}

	if !form.IsValid() {
		c.HTML(http.StatusOK, "signup-page.html", gin.H{
			"Form": form,
		})
		return
	}

	if _, err := h.s.Register(form.Name, form.Email, form.Password); err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Redirect(http.StatusSeeOther, "/user/login")
}

func (h *Handler) LoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login-page.html", gin.H{
		"Form": &LoginForm{},
	})
}

func (h *Handler) Login(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "Invalid form data")
		return
	}

	if !form.IsValid() {
		c.HTML(http.StatusOK, "login-page.html", gin.H{
			"Form": form,
		})
		return
	}

	user, err := h.s.Login(form.Email, form.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	session := sessions.Default(c)
	session.Set("userID", user.UserID)
	if err := session.Save(); err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}