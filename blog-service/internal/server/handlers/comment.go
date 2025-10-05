package handlers

import (
	"blog-service/internal/db/mongo"
	"blog-service/internal/db/postgres"
	"blog-service/internal/server/models"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

type CommentHandler struct {
	MongoDB    *mongo.Client
	PostgresDB *postgres.Client
}

var (
	CommentDeleteRe = regexp.MustCompile(`/comment/\d+`)
)

func (h *CommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodDelete && CommentDeleteRe.MatchString(r.URL.Path):
		h.CommentDelete(w, r)
		return
	case r.Method == http.MethodPost:
		h.CommentCreate(w, r)
		return
	}
}

func (h *CommentHandler) CommentDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *CommentHandler) CommentCreate(w http.ResponseWriter, r *http.Request) {
	var comment models.CommentCreateDTO
	var invalidArticleErr *models.InvalidArticleError
	var paramErr *models.ParamError
	var unauthorizedErr *models.UnauthorizedError

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, paramErr.Error(), http.StatusBadRequest)
		return
	}

	err = models.CreateComment(r.Context(), h.PostgresDB, h.MongoDB, &comment)
	if err != nil {
		switch {
		case errors.As(err, &invalidArticleErr):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.As(err, &paramErr):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.As(err, &unauthorizedErr):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
