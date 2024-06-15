package repository

import (
	"github.com/morf1lo/blog/internal/model"

	"gorm.io/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(comment *model.Comment) error {
	return r.db.Save(comment).Error
}

func (r *CommentRepo) FindByID(id int64) (*model.Comment, error) {
	var comment model.Comment
	if err := r.db.Model(&model.Comment{}).Where("id = ?", id).First(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepo) FindArticleComments(articleID int64) ([]*model.CommentDto, error) {
	var comments []*model.CommentDto
	if err := r.db.Table("comments c").
		Select("c.id, c.user_id, c.article_id, c.body, c.date_added, u.avatar AS user_avatar, u.username AS username").
		Joins("JOIN users u ON u.id = c.user_id").
		Where("c.article_id = ?", articleID).
		Find(&comments).Error; err != nil {
			return nil, err
		}

	return comments, nil
}

func (r *CommentRepo) Delete(commentID int64) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", commentID).Delete(&model.Comment{}).Error
}
