package models

import (
	"errors"
	"github.com/aliykh/reddit-feed/pkg/customErrors"
	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Feed struct {
	TotalCount int64     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Posts      []*Post `json:"posts"`
}

type Post struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	Title     string `json:"title" bson:"title" binding:"required"`
	Author    string `json:"author" bson:"author"`
	Link      string `json:"link,omitempty" bson:"link,omitempty"`
	Subreddit string `json:"subreddit" bson:"subreddit" binding:"required"`
	Content   string `json:"content,omitempty" bson:"content,omitempty"`
	Score     *int   `json:"score" bson:"score" binding:"required"`
	Promoted  *bool  `json:"promoted" bson:"promoted" binding:"required"`
	NSFW      *bool  `json:"nsfw" bson:"nsfw" binding:"required"`
}

func (p Post) CheckValidity() error {

	if p.Link != "" && p.Content != "" {
		return customErrors.New(http.StatusBadRequest, errors.New("post cannot have both content and link fields"))
	} else if p.Link == "" && p.Content == "" {
		return customErrors.New(http.StatusBadRequest, errors.New("post should have one of the following fields: link or content but not both"))
	}

	if p.Link != "" {
		errs := binding.Validator.Engine().(*validator.Validate).Var(p.Link, "url")
		if errs != nil {
			return &customErrors.ErrorResponse{
				ErrStatus: http.StatusBadRequest,
				Errors: []customErrors.ErrorValidation{
					{
						Field:   "link",
						Message: "link must be a valid URL",
					},
				},
			}
		}
	}

	return nil
}

var chars = []byte("abcdefghijklmnopqrstuvwxyz0123456789")

func (p *Post) GenerateAuthorName() {
	inBytes := make([]byte, 0, 11)
	inBytes = append(inBytes, []byte("t2_")...)
	inBytes = append(inBytes, uniuri.NewLenChars(8, chars)...)
	p.Author = string(inBytes)
}
