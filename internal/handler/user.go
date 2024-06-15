package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) usersSetAvatar(c *gin.Context) {
	user := h.getUser(c)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	if err := h.services.User.UpdateAvatar(c, user.ID, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "error": nil})
}
