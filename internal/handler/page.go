package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) pageIndex(c *gin.Context) {
	jwt, _ := c.Cookie("jwt")
	authenticated := jwt != ""

	userData, err := h.getUserDataFromTokenClaims(c.Request.Context(), jwt)
	if err != nil {
		h.templates.ExecuteTemplate(c.Writer, "index", nil)
		return
	}

	lastArticles, err := h.services.Article.FindLastArticles(c.Request.Context(), 3)
	if err != nil {
		h.templates.ExecuteTemplate(c.Writer, "index", gin.H{"authenticated": authenticated})
		return
	}

	h.templates.ExecuteTemplate(c.Writer, "index", gin.H{"authenticated": authenticated, "userData": userData, "lastArticles": lastArticles})
}

func (h *Handler) pageSignUp(c *gin.Context) {
	h.templates.ExecuteTemplate(c.Writer, "signup", nil)
}

func (h *Handler) pageLogin(c *gin.Context) {
	h.templates.ExecuteTemplate(c.Writer, "login", nil)
}

func (h *Handler) pageCreateArticle(c *gin.Context) {
	jwt, _ := c.Cookie("jwt")
	authenticated := jwt != ""

	userData, err := h.getUserDataFromTokenClaims(c.Request.Context(), jwt)
	if err != nil {
		h.templates.ExecuteTemplate(c.Writer, "index", nil)
		return
	}

	h.templates.ExecuteTemplate(c.Writer, "create_article", gin.H{"authenticated": authenticated, "userData": userData})
}

func (h *Handler) pageProfile(c *gin.Context) {
	jwt, _ := c.Cookie("jwt")
	authenticated := jwt != ""

	userData, err := h.getUserDataFromTokenClaims(c.Request.Context(), jwt)
	if err != nil {
		h.templates.ExecuteTemplate(c.Writer, "index", nil)
		return
	}

	username := c.Param("username")

	user, err := h.services.User.FindByUsername(c.Request.Context(), username)
	if err != nil {
		h.templates.ExecuteTemplate(c.Writer, "profile", nil)
		return
	}

	userArticles, err := h.services.Article.FindAuthorArticles(c.Request.Context(), user.ID.String())
	if err != nil {
		c.Redirect(302, "/")
		return
	}

	h.templates.ExecuteTemplate(c.Writer, "profile", gin.H{"authenticated": authenticated, "userData": userData, "data": user, "dataArticles": userArticles})
}

func (h *Handler) pageArticle(c *gin.Context) {
	jwt, _ := c.Cookie("jwt")
	authenticated := jwt != ""

	userData, err := h.getUserDataFromTokenClaims(c.Request.Context(), jwt)
	if err != nil {
		h.templates.ExecuteTemplate(c.Writer, "index", nil)
		return
	}

	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		h.templates.ExecuteTemplate(c.Writer, "index", gin.H{"authenticated": authenticated})
		return
	}

	article, err := h.services.Article.FindByID(c.Request.Context(), int64(id))
	if err != nil {
		h.templates.ExecuteTemplate(c.Writer, "index", gin.H{"authenticated": authenticated})
		return
	}

	comments, err := h.services.Comment.FindArticleComments(c.Request.Context(), int64(id))
	if err != nil {
		h.templates.ExecuteTemplate(c.Writer, "index", gin.H{"authenticated": authenticated})
		return
	}

	h.templates.ExecuteTemplate(c.Writer, "article", gin.H{
		"authenticated": authenticated,
		"userData":      userData,
		"data":          article,
		"comments":      comments,
	})
}

func (h *Handler) pageSearch(c *gin.Context) {
	jwt, _ := c.Cookie("jwt")
	authenticated := jwt != ""

	userData, err := h.getUserDataFromTokenClaims(c.Request.Context(), jwt)
	if err != nil {
		h.templates.ExecuteTemplate(c.Writer, "index", nil)
		return
	}

	query := c.Query("q")

	searchResults, err := h.services.Article.Search(c.Request.Context(), query)
	if err != nil {
		h.templates.ExecuteTemplate(c.Writer, "index", gin.H{"authenticated": authenticated})
		return
	}

	h.templates.ExecuteTemplate(c.Writer, "search", gin.H{
		"authenticated": authenticated,
		"userData": userData,
		"data": searchResults,
		"query": query,
	})
}

func (h *Handler) pageSettings(c *gin.Context) {
	jwt, _ := c.Cookie("jwt")
	authenticated := jwt != ""

	userData, err := h.getUserDataFromTokenClaims(c.Request.Context(), jwt)
	if err != nil {
		h.templates.ExecuteTemplate(c.Writer, "index", nil)
		return
	}

	h.templates.ExecuteTemplate(c.Writer, "settings", gin.H{"authenticated": authenticated, "userData": userData})
}
