package register

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mentorlink/internal/domain/model"
	"mentorlink/internal/domain/response"
	"mentorlink/internal/handlers/mocks"
	"mentorlink/internal/lib/logger/slogdiscard"
	"mentorlink/internal/storage/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegisterHandler(t *testing.T) {
	cases := []struct {
		name           string
		email          string
		password       string
		repeatPassword string
		role           string
		expectedStatus int
		respError      string
		mockError      error
	}{
		{
			name:           "Success",
			email:          "valid@mail.com",
			password:       "securePassword123",
			repeatPassword: "securePassword123",
			role:           "user",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Passwords mismatch",
			email:          "test@mail.com",
			password:       "password123",
			repeatPassword: "password456",
			role:           "user",
			expectedStatus: http.StatusBadRequest,
			respError:      "server error",
		},
		{
			name:           "User already exists",
			email:          "exists@mail.com",
			password:       "password123",
			repeatPassword: "password123",
			role:           "user",
			expectedStatus: http.StatusConflict,
			respError:      "user already exists",
		},
		{
			name:           "Invalid role",
			email:          "test@mail.com",
			password:       "password123",
			repeatPassword: "password123",
			role:           "invalid_role",
			expectedStatus: http.StatusBadRequest,
			respError:      "server error",
		},
		{
			name:           "Internal error",
			email:          "test@mail.com",
			password:       "password123",
			repeatPassword: "password123",
			role:           "user",
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			respError:      "failed to create user",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			userCreaterMock := mocks.NewUserCreater(t)
			newMentorMock := mocks.NewNewMentor(t)

			// Настройка моков только для кейсов с обращением к БД
			if tc.name != "Passwords mismatch" && tc.name != "Invalid role" {
				if tc.name == "User already exists" {
					userCreaterMock.On("GetByEmail", tc.email).
						Return(&model.User{Email: tc.email}, nil)
				} else {
					userCreaterMock.On("GetByEmail", tc.email).
						Return(nil, db.ErrUserNotFound)
				}
			}

			if tc.expectedStatus == http.StatusCreated || tc.mockError != nil {
				userCreaterMock.On("CreateUser", mock.Anything).
					Return(tc.mockError).Once()
			}

			handler := Register(context.Background(), slogdiscard.NewDiscardLogger(), userCreaterMock, newMentorMock)

			body := fmt.Sprintf(
				`{"email": "%s", "password": "%s", "repeat_password": "%s", "role": "%s"}`,
				tc.email, tc.password, tc.repeatPassword, tc.role,
			)

			req, err := http.NewRequest(
				http.MethodPost,
				"/auth/register",
				bytes.NewBufferString(body),
			)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// Проверка статуса
			require.Equal(t, tc.expectedStatus, rr.Code)

			// Проверка ошибки
			if tc.respError != "" {
				var resp response.Response
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
				require.Equal(t, tc.respError, resp.Error)
			}

			// Проверка вызовов моков
			if tc.name != "Passwords mismatch" && tc.name != "Invalid role" {
				userCreaterMock.AssertExpectations(t)
			}
		})
	}
}
