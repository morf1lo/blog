package handler

import (
	"net/http"
	"github.com/morf1lo/blog/internal/model"
	"github.com/morf1lo/blog/pkg/auth"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) authSignUp(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := h.services.Authorization.SignUp(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := auth.SendToken(c, jwt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(302, "/")
}

func (h *Handler) authSignIn(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := h.services.Authorization.SignIn(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := auth.SendToken(c, jwt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(302, "/")
}

func (h *Handler) authActivate(c *gin.Context) {
	activationLink := c.Param("link")

	if err := h.services.Authorization.Activate(activationLink); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(302, "/")
}

type requestSaveResetPasswordToken struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *Handler) authSaveResetPasswordToken(c *gin.Context) {
	var request requestSaveResetPasswordToken
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	token := model.ResetPasswordToken{
		UserEmail: request.Email,
		Expiry: time.Now().Add(time.Minute * 2),
	}

	if err := h.services.Authorization.SaveResetPasswordToken(&token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type requestResetPasswordToken struct {
	Token string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"min=8,max=24"`
}

func (h *Handler) authResetPasswordToken(c *gin.Context) {
	var request requestResetPasswordToken
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.services.Authorization.ResetPassword(request.Token, request.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) authSignOut(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "localhost", true, true)
	c.Redirect(302, "/")
}
