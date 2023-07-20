package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kripsy/gophermart/internal/auth/internal/config"
	"github.com/kripsy/gophermart/internal/auth/internal/mocks"
	models "github.com/kripsy/gophermart/internal/auth/internal/models"
	"github.com/kripsy/gophermart/internal/auth/internal/usecase"
	"github.com/stretchr/testify/assert"
)

type TestParams struct {
	ctx context.Context
	cfg *config.Config
}

func getParamsForTest() *TestParams {
	// l, _ := logger.InitLogger("Debug")
	// ctx := logger.ContextWithLogger(context.Background(), l)
	ctx := context.Background()
	cfg := config.InitConfig()

	tp := &TestParams{
		ctx: ctx,
		cfg: cfg,
	}
	return tp
}

func TestRegisterUserHandler(t *testing.T) {
	paramTest := getParamsForTest()
	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name       string
		body       string
		methodType string
		want       want
	}{
		// TODO: Add test cases.
		{
			name: "success save",

			body: `{
				"login": "root",
				"password": "qwerty"
			}`,
			methodType: "POST",
			want: want{
				contentType: "application/json",
				statusCode:  200,
			},
		},
		{
			name: "uncorrect request format",
			body: `{
				"username": "root",
				"password": "qwerty",
			}`,
			methodType: "POST",
			want: want{
				contentType: "application/json",
				statusCode:  400,
			},
		},
		{
			name: "login conflict",
			body: `{
				"login": "root2",
				"password": "qwerty"
			}`,
			methodType: "POST",
			want: want{
				contentType: "application/json",
				statusCode:  409,
			},
		},
		{
			name: "internal server error",
			body: `{
				"login": "internalerroruser",
				"password": "qwerty"
			}`,
			methodType: "POST",
			want: want{
				contentType: "application/json",
				statusCode:  500,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepository(ctrl)
			if tt.want.statusCode == 200 {
				repo.EXPECT().IsUserExists(gomock.Any(), "root").Return(false, nil)
				repo.EXPECT().GetNextUserID(gomock.Any()).Return(5, nil)
				repo.EXPECT().RegisterUser(gomock.Any(), "root", gomock.Any(), gomock.Any()).Return(nil)
			} else {
				repo.EXPECT().IsUserExists(gomock.Any(), "root").Return(false, nil).AnyTimes()
				repo.EXPECT().IsUserExists(gomock.Any(), "root2").Return(true, models.NewUserExistsError("root2")).AnyTimes()
				repo.EXPECT().GetNextUserID(gomock.Any()).Return(5, nil).AnyTimes()
				repo.EXPECT().RegisterUser(gomock.Any(), "root", gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				repo.EXPECT().IsUserExists(gomock.Any(), "internalerroruser").Return(false, errors.New("")).AnyTimes()
			}

			body := strings.NewReader(tt.body)

			uc, err := usecase.InitUseCases(paramTest.ctx, repo, paramTest.cfg)
			assert.NoError(t, err)
			request := httptest.NewRequest(tt.methodType, "/", body)
			w := httptest.NewRecorder()
			ht, _ := InitHandler(paramTest.ctx, uc)
			h := ht.RegisterUserHandler
			h(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			if result.StatusCode == 200 {
				assert.NotEmpty(t, result.Header.Get("Authorization"), "Header shouldn't be empty")
			}
		})
	}
}

func TestHandler_LoginUserHandler(t *testing.T) {
	type fields struct {
		ctx context.Context
		uc  *usecase.UseCase
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				ctx: tt.fields.ctx,
				uc:  tt.fields.uc,
			}
			h.LoginUserHandler(tt.args.w, tt.args.r)
		})
	}
}
