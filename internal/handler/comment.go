package handler

import (
	"net/http"
	"strconv"

	"github.com/morf1lo/blog/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) commentsCreate(c *gin.Context) {
	user := h.getUser(c)

	articleIDString := c.Param("article_id")
	articleID, err := strconv.Atoi(articleIDString)
	if err != nil  {
		c.Redirect(302, "/")
		return
	}

	var comment model.Comment
	if err := c.ShouldBind(&comment); err != nil {
		c.Redirect(302, "/")
		return
	}

	comment.UserID = user.ID
	comment.ArticleID = int64(articleID)

	if err := h.services.Comment.Create(c.Request.Context(), &comment); err != nil {
		c.Redirect(302, "/")
		return
	}

	c.Redirect(302, "/article/" + articleIDString)
}

func (h *Handler) commentsDelete(c *gin.Context) {
	user := h.getUser(c)

	commentIDString := c.Param("id")
	commentID, err := strconv.Atoi(commentIDString)
	if err != nil  {
		c.Redirect(302, "/")
		return
	}

	if err := h.services.Comment.Delete(c.Request.Context(), int64(commentID), user.ID); err != nil {
		c.Redirect(302, "/")
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
