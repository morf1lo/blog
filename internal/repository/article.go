package repository

import (
	"github.com/morf1lo/blog/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArticleRepo struct {
	db *gorm.DB
}

func NewArticleRepo(db *gorm.DB) *ArticleRepo {
	return &ArticleRepo{db: db}
}

func (r *ArticleRepo) Create(article *model.Article) error {
	return r.db.Save(article).Error
}

func (r *ArticleRepo) FindByID(id int64) (*model.ArticleDto, error) {
	var article model.ArticleDto
	if err := r.db.Table("articles a").
		Select("a.id, a.user_id, a.title, a.text, a.date_added, u.avatar as user_avatar, u.username as username").
		Joins("JOIN users u ON u.id = a.user_id").
		Where("a.id = ?", id).
		First(&article).Error; err != nil {
		return nil, err
	}

	return &article, nil
}

func (r *ArticleRepo) UpdateByID(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.Article{}).Where("id = ?").Updates(updates).Error
}

func (r *ArticleRepo) Search(query string) ([]*model.Article, error) {
	var articles []*model.Article
	if err := r.db.Where("title ILIKE ?", "%"+query+"%").Limit(10).Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepo) FindAuthorArticles(authorID string) ([]*model.Article, error) {
	var articles []*model.Article
	if err := r.db.Table("articles a").Where("user_id = ?", authorID).Order("a.date_added desc").Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepo) FindLastArticles(limit int) ([]*model.ArticleDto, error) {
	var articles []*model.ArticleDto
	if err := r.db.Table("articles a").
		Select("a.id, a.user_id, a.title, a.text, a.date_added, u.avatar as user_avatar, u.username as username").
		Joins("JOIN users u ON u.id = a.user_id").
		Order("a.date_added desc").
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepo) Delete(id int64, userID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.User{}).Error
}
