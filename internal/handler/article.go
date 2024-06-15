package handler

import (
	"net/http"
	"github.com/morf1lo/blog/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) articlesCreate(c *gin.Context) {
	var article model.Article
	if err := c.ShouldBind(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := h.getUser(c)
	article.UserID = user.ID

	if err := h.services.Article.Create(c.Request.Context(), &article); err != nil {
		c.Redirect(302, "/")
		return
	}

	c.Redirect(302, "/")
}
