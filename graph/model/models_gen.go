// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AddBookInput struct {
	Title  string   `json:"title"`
	Author string   `json:"author"`
	Genre  []string `json:"genre"`
	Price  float64  `json:"price"`
}

type Books struct {
	ID     string   `json:"_id"`
	Title  string   `json:"title"`
	Author string   `json:"author"`
	Genre  []string `json:"genre"`
	Price  float64  `json:"price"`
}

type Mutation struct {
}

type Query struct {
}

type RemoveBookResponse struct {
	DeletedBookID string `json:"deletedBookId"`
}

type UpdateBookInput struct {
	Title  *string   `json:"title,omitempty"`
	Author *string   `json:"author,omitempty"`
	Genre  []*string `json:"genre,omitempty"`
	Price  *float64  `json:"price,omitempty"`
}
