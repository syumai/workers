package model

type Article struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt uint64 `json:"createdAt"`
}

type CreateArticleRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type CreateArticleResponse struct {
	Article Article `json:"article"`
}

type ListArticlesResponse struct {
	Articles []Article `json:"articles"`
}
