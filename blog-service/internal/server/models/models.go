package models

import (
	"blog-service/internal/db/mongo"
	mongomodels "blog-service/internal/db/mongo/models"
	"blog-service/internal/db/postgres"
	pgmodels "blog-service/internal/db/postgres/models"

	"context"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArticleUpdateDTO struct {
	Title       string `json:"title,omitempty"`
	Content     string `json:"content,omitempty"`
	Category    string `json:"category,omitempty"`
	PublisherID int    `json:"publisherId,omitempty"`
}

/*
	 TODO: After auth is implemented, get the publisher name and ID
		   from the token instead of the request
*/
type ArticleCreateDTO struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type ArticleGetDTO struct {
	CreatedAt     time.Time        `bson:"created_at" json:"createdAt"`
	Title         string           `bson:"title" json:"title"`
	Content       string           `bson:"content" json:"content"`
	PublisherName string           `bson:"publisher_name" json:"publisherName"`
	Category      string           `bson:"category" json:"category"`
	ID            string           `json:"id"`
	Comments      []CommentsGetDTO `json:"comments"`
	PublisherID   int              `bson:"publisher_id" json:"publisherId"`
}

type ArticleCreateResponse struct {
	ID string `json:"id"`
}

type CommentsGetDTO struct {
	pgmodels.Comment
	Likes int `json:"likes"`
}

/*
	 TODO: After auth is implemented, get the userID from the token
		   instead of the request
*/
type CommentCreateDTO struct {
	Content   string `json:"content"`
	ArticleID string `json:"articleId"`
}

func GetCommentsByArticleID(ctx context.Context, db *postgres.Client, id string, limit, page int) ([]CommentsGetDTO, error) {
	// Get comments count
	commCount, err := db.GetCommentsCount(ctx, id)
	if err != nil {
		return []CommentsGetDTO{}, err
	}

	ceil := math.Ceil(float64(commCount) / float64(limit)) // calculate max number of pages

	// if page requested is bigger than the max page, set the page to the max page
	if ceil < float64(page) {
		page = int(ceil)
	}

	comments, err := db.GetComments(ctx, id, limit, page)
	if err != nil {
		return []CommentsGetDTO{}, err
	}

	commentsRes := make([]CommentsGetDTO, 0, limit)
	for _, com := range comments {
		likes, err := db.GetCommentLikeCount(ctx, com.ID)
		if err != nil {
			return []CommentsGetDTO{}, err
		}

		commentsRes = append(commentsRes, CommentsGetDTO{
			com,
			likes,
		})
	}

	return commentsRes, nil
}

func GetArticleByID(ctx context.Context, db *mongo.Client, id string) (*ArticleGetDTO, error) {
	articleOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	article, err := db.FindArticleByID(ctx, &articleOID)
	if err != nil {
		return nil, err
	}

	return &ArticleGetDTO{
		CreatedAt:     article.CreatedAt,
		Title:         article.Title,
		Content:       article.Content,
		PublisherName: article.PublisherName,
		Category:      article.Category,
		ID:            id,
		PublisherID:   article.PublisherID,
	}, nil
}

func CreateArticle(ctx context.Context, db *mongo.Client, article *ArticleCreateDTO) (string, error) {
	userClaims := GetClaimsFromContext(ctx)
	if userClaims == nil {
		return "", &UnauthorizedError{}
	}

	if article.Title == "" || article.Content == "" || article.Category == "" {
		return "", &ParamError{}
	}

	articleToInsert := mongomodels.ArticleDB{
		CreatedAt:     time.Now(),
		Title:         article.Title,
		Content:       article.Content,
		Category:      article.Category,
		PublisherName: userClaims.Username,
		PublisherID:   userClaims.ID,
	}

	res, err := db.InsertArticle(ctx, &articleToInsert)
	if err != nil {
		return "", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert returned result to objectid")
	}

	return oid.Hex(), nil
}

func CreateComment(ctx context.Context, pgdb *postgres.Client, mdb *mongo.Client, comment *CommentCreateDTO) error {
	userClaims := GetClaimsFromContext(ctx)
	if userClaims == nil {
		return &UnauthorizedError{}
	}

	if comment.ArticleID == "" || comment.Content == "" {
		return &ParamError{}
	}

	articleOID, err := primitive.ObjectIDFromHex(comment.ArticleID)
	if err != nil {
		return err
	}

	_, err = mdb.FindArticleByID(ctx, &articleOID)
	if err != nil {
		return &InvalidArticleError{}
	}

	commentToInsert := pgmodels.Comment{
		CreatedAt: time.Now(),
		Content:   comment.Content,
		ArticleID: comment.ArticleID,
		UserID:    userClaims.ID,
	}

	_, err = pgdb.CreateComment(ctx, commentToInsert)

	return err
}
