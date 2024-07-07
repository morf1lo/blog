package handler

import (
	"context"
	"html/template"

	"github.com/morf1lo/blog/internal/model"
	"github.com/morf1lo/blog/internal/service"
	"github.com/morf1lo/blog/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	services *service.Service
	templates *template.Template
}

func New(services *service.Service) *Handler {
	h := &Handler{services: services}
	h.loadTemplates()
	return h
}

func (h *Handler) loadTemplates() {
	h.templates = template.Must(template.ParseGlob("templates/*.html"))
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Static("/static", "static")
	router.Static("/public", "public")

	router.GET("/", h.pageIndex)
	router.GET("/create", h.pageCreateArticle)
	router.GET("/register", h.pageSignUp)
	router.GET("/login", h.pageLogin)
	router.GET("/profile/:username", h.pageProfile)
	router.GET("/article/:id", h.pageArticle)
	router.GET("/search", h.pageSearch)
	router.GET("/settings", h.authMiddleware, h.pageSettings)

	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.authSignUp)
		auth.POST("/signin", h.authSignIn)
		auth.GET("/activate/:link", h.authActivate)
		auth.POST("/reset", h.authSaveResetPasswordToken)
		auth.PATCH("/reset", h.authResetPasswordToken)
		auth.GET("/signout", h.authSignOut)
	}

	users := router.Group("/users")
	{
		users.POST("/avatar", h.authMiddleware, h.usersSetAvatar)
	}

	articles := router.Group("/articles")
	{
		articles.POST("/create", h.authMiddleware, h.articlesCreate)
	}

	comments := router.Group("/comments")
	{
		comments.POST("/:article_id", h.authMiddleware, h.commentsCreate)
		comments.DELETE("/:id", h.authMiddleware, h.commentsDelete)
	}

	return router
}

func (h *Handler) getUserDataFromTokenClaims(ctx context.Context, token string) (*model.User, error) {
	claims, err := auth.GetTokenClaims(token)
	if err != nil {
		return nil, err
	}

	idString := claims["sub"].(string)
	id, err := uuid.Parse(idString)
	if err != nil {
		return nil, err
	}

	user, err := h.services.User.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) getUser(c *gin.Context) *model.User {
	userReq, _ := c.Get("user")

	user, ok := userReq.(model.User)
	if !ok {
		return nil
	}

	return &user
}
