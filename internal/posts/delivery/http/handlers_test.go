package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aliykh/log"
	"github.com/aliykh/reddit-feed/internal/posts/mock"
	"github.com/aliykh/reddit-feed/internal/posts/models"
	"github.com/aliykh/reddit-feed/pkg/customErrors"
	"github.com/aliykh/reddit-feed/pkg/helpers"
	"github.com/aliykh/reddit-feed/pkg/pagination"
	"github.com/aliykh/reddit-feed/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"reflect"
	"testing"
)

var router *gin.Engine

func init() {

	router = gin.Default()

	binding.Validator = new(helpers.DefaultValidator)

	engine := binding.Validator.Engine().(*validator.Validate)

	eng := en.New()
	uni := ut.New(eng, eng)
	customErrors.Trans, _ = uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(engine, customErrors.Trans)

}

func TestHandlers_Create(t *testing.T) {

	reqBody := &models.Post{
		Title:     "Title101",
		Link:      "https://github.com9",
		Subreddit: "/r/subreddit",
		Score:     new(int),
		Promoted:  new(bool),
		NSFW:      new(bool),
	}

	model := &models.Post{
		Title:     "Title101",
		Link:      "https://github.com9",
		Subreddit: "/r/subreddit",
		Score:     new(int),
		Promoted:  new(bool),
		NSFW:      new(bool),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostUC := mock.NewMockUseCase(ctrl)

	logger := log.NewFactory(log.Mock, "test")

	postHandlers := New(logger, mockPostUC)

	router.POST("/ok", postHandlers.Create)

	t.Run("passes", func(t *testing.T) {

		req, err := utils.MakeRequest(utils.POST, utils.JSON, "/ok", reqBody)
		require.NoError(t, err)

		mockPostUC.EXPECT().Create(context.Background(), gomock.Eq(model)).Return(model, nil)

		resp, err := utils.InvokeHandler(req, router)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		data := &models.Post{}
		err = json.Unmarshal(resp.Body, &data)
		require.NoError(t, err)
		reflect.DeepEqual(data, model)
	})

	t.Run("fails", func(t *testing.T) {

		req, err := utils.MakeRequest(utils.POST, utils.JSON, "/ok", reqBody)
		require.NoError(t, err)

		expectedErr := customErrors.New(http.StatusBadRequest, mongo.ErrNilValue)
		mockPostUC.EXPECT().Create(context.Background(), gomock.Eq(model)).Return(nil, expectedErr)

		httpResult, err := utils.InvokeHandler(req, router)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, httpResult.StatusCode)

		errData := &customErrors.ErrorResponse{}
		err = json.Unmarshal(httpResult.Body, &errData)
		require.NoError(t, err)
		require.Equal(t, expectedErr, errData)

	})

}

func TestHandlers_Create_Validation_Errs(t *testing.T) {

	type TestCase struct {
		ReqBody  *models.Post
		Expected *customErrors.ErrorResponse
	}

	cases := []TestCase{
		{
			ReqBody: &models.Post{
				//Title:     "Title101", // omit title
				Link:      "https://github.com9",
				Subreddit: "/r/subreddit",
				Score:     new(int),
				Promoted:  new(bool),
				NSFW:      new(bool),
			},
			Expected: &customErrors.ErrorResponse{
				ErrStatus: http.StatusBadRequest,
				Errors: []customErrors.ErrorValidation{
					{
						Field:   "title",
						Message: "title is a required field",
					},
				},
			},
		},
		{
			ReqBody: &models.Post{
				Title:     "Title101",
				Link:      "invalid url", // invalid url
				Subreddit: "/r/subreddit",
				Score:     new(int),
				Promoted:  new(bool),
				NSFW:      new(bool),
			},
			Expected: &customErrors.ErrorResponse{
				ErrStatus: http.StatusBadRequest,
				Errors: []customErrors.ErrorValidation{
					{
						Field:   "link",
						Message: "link must be a valid URL",
					},
				},
			},
		},
		{
			ReqBody: &models.Post{
				Title:     "Title101",
				Link:      "https://example.com",
				Subreddit: "/r/subreddit",
				Content:   "content", // link and content exists
				Score:     new(int),
				Promoted:  new(bool),
				NSFW:      new(bool),
			},
			Expected: &customErrors.ErrorResponse{
				ErrStatus: http.StatusBadRequest,
				ErrError:  "post cannot have both content and link fields",
			},
		},
	}

	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostUC := mock.NewMockUseCase(ctrl)

	logger := log.NewFactory(log.Mock, "test")

	postHandlers := New(logger, mockPostUC)

	router.POST("/validation/fail", postHandlers.Create)

	for _, c := range cases {

		req, err := utils.MakeRequest(utils.POST, utils.JSON, "/validation/fail", c.ReqBody)
		require.NoError(t, err)

		httpResult, err := utils.InvokeHandler(req, router)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, httpResult.StatusCode)

		resp := &customErrors.ErrorResponse{}
		err = json.Unmarshal(httpResult.Body, &resp)
		require.NoError(t, err)
		require.Equal(t, c.Expected, resp)

	}

}

func TestHandlers_Generate(t *testing.T) {

	// pagination
	query := &pagination.Query{
		Size: 2,
		Page: 0,
	}

	// get use case
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostUC := mock.NewMockUseCase(ctrl)


	//
	logger := log.NewFactory(log.Mock, "test")

	postHandlers := New(logger, mockPostUC)

	router.GET("/generate/ok", postHandlers.Generate)


	t.Run("ok", func(t *testing.T) {

		// model
		feed := &models.Feed{
			TotalCount: 10,
			TotalPages: 1,
			Page:       0,
			Size:       2,
			HasMore:    true,
			Posts: []*models.Post{
				{
					Id:        "6252a65cc511344c065986f3",
					Title:     "Title20",
					Author:    "t2_6wmjk11m",
					Link:      "https://www.example.com",
					Subreddit: "/r/subreddit",
					Score:     new(int),
					Promoted:  new(bool),
					NSFW:      new(bool),
				},
				{
					Id:        "1234",
					Title:     "Title1234",
					Author:    "t2_6wmjk11m",
					Link:      "https://www.example.com",
					Subreddit: "/r/subreddit12345",
					Score:     new(int),
					Promoted:  new(bool),
					NSFW:      new(bool),
				},
			},
		}

		mockPostUC.EXPECT().GenerateFeeds(context.Background(), gomock.Eq(query)).Return(feed, nil)

		req, err := utils.MakeRequest(utils.GET, utils.FORM, "/generate/ok", *query)
		require.NoError(t, err)

		resp, err := utils.InvokeHandler(req, router)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		resFeed := &models.Feed{}

		err = json.Unmarshal(resp.Body, &resFeed)
		require.NoError(t, err)
		reflect.DeepEqual(feed, resFeed)

	})


	t.Run("usecase fail", func(t *testing.T) {

		mockPostUC.EXPECT().GenerateFeeds(context.Background(), gomock.Eq(query)).Return(nil, errors.New("fails"))

		req, err := utils.MakeRequest(utils.GET, utils.FORM, "/generate/ok", *query)
		require.NoError(t, err)

		resp, err := utils.InvokeHandler(req, router)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		errData := &customErrors.ErrorResponse{}
		err = json.Unmarshal(resp.Body, &errData)
		require.NoError(t, err)
		require.Equal(t, customErrors.NewInternalServerError(), errData)

	})


	t.Run("query validation fails", func(t *testing.T) {

		req, err := utils.MakeRequest(utils.GET, utils.FORM, "/generate/ok", map[string]string{
			"size": "abs",
		})

		require.NoError(t, err)

		resp, err := utils.InvokeHandler(req, router)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)


		errData := &customErrors.ErrorResponse{}
		err = json.Unmarshal(resp.Body, &errData)
		require.NoError(t, err)
		require.Equal(t, customErrors.New(http.StatusBadRequest, errors.New(`strconv.ParseInt: parsing "abs": invalid syntax`)), errData)

	})



}
