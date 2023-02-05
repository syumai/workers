package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mailru/easyjson"
	_ "github.com/syumai/workers/cloudflare/d1" // register driver
	"github.com/syumai/workers/examples/d1-blog-server/app/model"
	"github.com/syumai/workers/internal/jsutil"
)

type articleHandler struct{}

var _ http.Handler = (*articleHandler)(nil)

func NewArticleHandler() http.Handler {
	return &articleHandler{}
}

func (h *articleHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	db, err := sql.Open("d1", "BlogDB")
	if err != nil {
		h.handleErr(w, http.StatusInternalServerError, fmt.Sprintf("failed to initialize DB: %v", err))
	}
	switch req.Method {
	case http.MethodGet:
		h.listArticles(w, db)
		return
	case http.MethodPost:
		h.createArticle(w, req, db)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("not found"))
}

func (h *articleHandler) handleErr(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(msg))
}

func (h *articleHandler) createArticle(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var createArticleReq model.CreateArticleRequest
	if err := easyjson.UnmarshalFromReader(req.Body, &createArticleReq); err != nil {
		h.handleErr(w, http.StatusBadRequest,
			"request format is invalid")
		return
	}

	now := time.Now().Unix()
	article := model.Article{
		ID:        jsutil.NewUUID(),
		Title:     createArticleReq.Title,
		Body:      createArticleReq.Body,
		CreatedAt: uint64(now),
	}

	_, err := db.Exec(`
INSERT INTO articles (id, title, body, created_at)
VALUES (?, ?, ?, ?)
   `, article.ID, article.Title, article.Body, article.CreatedAt)
	if err != nil {
		log.Println(err)
		h.handleErr(w, http.StatusInternalServerError,
			"failed to save article")
		return
	}

	res := model.CreateArticleResponse{
		Article: article,
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := easyjson.MarshalToWriter(&res, w); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode response: %w\n", err)
	}
}

func (h *articleHandler) listArticles(w http.ResponseWriter, db *sql.DB) {
	rows, err := db.Query(`
SELECT id, title, body, created_at FROM articles
ORDER BY created_at DESC;
   `)
	if err != nil {
		h.handleErr(w, http.StatusInternalServerError,
			"failed to load article")
		return
	}

	articles := []model.Article{}
	for rows.Next() {
		var (
			id, title, body string
			createdAt       float64 // number value is always retrieved as float64.
		)
		err = rows.Scan(&id, &title, &body, &createdAt)
		if err != nil {
			break
		}
		articles = append(articles, model.Article{
			ID:        id,
			Title:     title,
			Body:      body,
			CreatedAt: uint64(createdAt),
		})
	}
	res := model.ListArticlesResponse{
		Articles: articles,
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := easyjson.MarshalToWriter(&res, w); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode response: %w\n", err)
	}
}
