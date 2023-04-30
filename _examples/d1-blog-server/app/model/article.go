//go:generate easyjson .
package model

//easyjson:json
type Article struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt uint64 `json:"createdAt"`
}

//easyjson:json
type CreateArticleRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

//easyjson:json
type CreateArticleResponse struct {
	Article Article `json:"article"`
}

//easyjson:json
type ListArticlesResponse struct {
	Articles []Article `json:"articles"`
}
