package handler

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/bb4ssttardio/RESTapi_todo-app/pkg/service"
	mock_service "github.com/bb4ssttardio/RESTapi_todo-app/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockAuthorization, token string) {
				r.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "Invalid Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth header"}`,
		},
		// {
		// 	name:        "Invalid Header Value",
		// 	headerName:  "Authorization",
		// 	headerValue: "Bearr token",
		// 	token:       "token",
		// 	mockBehavior: func(r *mock_service.MockAuthorization, token string) {
		// 	},
		// 	expectedStatusCode:   401,
		// 	expectedResponseBody: `{"message":"invalid auth header"}`,
		// },
		{
			name:                 "Empty Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"token is empty"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockAuthorization, token string) {
				r.EXPECT().ParseToken(token).Return(0, errors.New("invalid token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid token"}`,
		},
	}

	for _, testCase := range testTable {
		// Init Dependencies
		c := gomock.NewController(t)
		defer c.Finish()

		auth := mock_service.NewMockAuthorization(c)
		testCase.mockBehavior(auth, testCase.token)

		services := &service.Service{Authorization: auth}
		handler := Handler{services}

		// test Server
		r := gin.New()
		r.GET("/protected", handler.userIdentity, func(c *gin.Context) {
			id, _ := c.Get(userCtx)
			c.String(200, fmt.Sprintf("%d", id.(int)))
		})

		// test req

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set(testCase.headerName, testCase.headerValue)

		// make req

		r.ServeHTTP(w, req)

		assert.Equal(t, w.Code, testCase.expectedStatusCode)
		assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
	}
}
