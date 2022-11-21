package server

//
//import (
//	"bytes"
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"golang.org/x/crypto/bcrypt"
//	"math/rand"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//	"time"
//	"user-service/models"
//	"user-service/server/jwt"
//)
//
//func TestSignup(t *testing.T) {
//	newReq := &models.User{
//		Name:        "Tolu",
//		PhoneNumber: "+2348163608141",
//		Email:       "toluwase@gmail.com",
//		Password:    "12345678",
//	}
//	noEmail := &models.User{
//		PhoneNumber: "+2348163608141",
//		Password:    "12345678",
//	}
//	noPhone := &models.User{
//		Name:  "Tolu",
//		Email: "toluwase@gmail.com",
//	}
//
//	user := &models.User{
//		Name:        newReq.Name,
//		PhoneNumber: newReq.PhoneNumber,
//		Email:       newReq.Email,
//	}
//	user.ID = 6
//	user.CreatedAt = time.Now().Unix()
//	user.UpdatedAt = time.Now().Unix()
//
//	cases := []struct {
//		Name            string
//		Request         *models.User
//		ExpectedCode    int
//		ExpectedMessage string
//		ExpectedError   string
//		mockDB          func(ctrl *mocks.MockAuthRepository, service *mocks.MockAuthService)
//		checkResponse   func(recorder *httptest.ResponseRecorder)
//	}{
//		{
//			Name:            "Test Signup with correct details",
//			Request:         newReq,
//			ExpectedCode:    http.StatusCreated,
//			ExpectedMessage: "Signup successful, check your email for verification",
//			ExpectedError:   "",
//			mockDB: func(ctrl *mocks.MockAuthRepository, service *mocks.MockAuthService) {
//				service.EXPECT().SignupUser(newReq)
//			},
//		},
//		{
//			Name:            "Test Signup with no email",
//			Request:         noEmail,
//			ExpectedCode:    http.StatusBadRequest,
//			ExpectedMessage: "",
//			ExpectedError:   "Email is invalid: toluwase.tt.com",
//			mockDB: func(ctrl *mocks.MockAuthRepository, service *mocks.MockAuthService) {
//				service.EXPECT().SignupUser(noEmail).
//					Return(&models.User{}, nil).AnyTimes()
//			},
//		},
//		{
//			Name:            "Test Signup with invalid fields",
//			Request:         noPhone,
//			ExpectedCode:    http.StatusBadRequest,
//			ExpectedMessage: "",
//			ExpectedError:   "Email is invalid: toluwase.tt.com",
//			mockDB: func(ctrl *mocks.MockAuthRepository, service *mocks.MockAuthService) {
//				service.EXPECT().SignupUser(noPhone).
//					Return(&models.User{}, nil).AnyTimes()
//			},
//		},
//		{
//			Name:            "Test Signup with duplicate email address",
//			Request:         &models.User{Name: "Tolu", PhoneNumber: "08141", Email: "toluwase@gmail.com"},
//			ExpectedCode:    http.StatusBadRequest,
//			ExpectedMessage: "",
//			ExpectedError:   "user already exists",
//			mockDB: func(ctrl *mocks.MockAuthRepository, service *mocks.MockAuthService) {
//				service.EXPECT().SignupUser(newReq).
//					Return(&models.User{}, nil).AnyTimes()
//			},
//		},
//	}
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	mockAuthRepo := mocks.NewMockAuthRepository(ctrl)
//	authService := mocks.NewMockAuthService(ctrl)
//	testServer.handler.AuthService = authService
//	testServer.handler.AuthRepository = mockAuthRepo
//	for _, c := range cases {
//		t.Run(c.Name, func(t *testing.T) {
//			// FIXME: refactor this test
//			c.mockDB(mockAuthRepo, authService)
//			data, err := json.Marshal(c.Request)
//			require.NoError(t, err)
//			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewReader(data))
//			require.NoError(t, err)
//			recorder := httptest.NewRecorder()
//			testServer.router.ServeHTTP(recorder, req)
//			assert.Equal(t, recorder.Code, c.ExpectedCode)
//			assert.Contains(t, recorder.Body.String(), c.ExpectedMessage)
//
//		})
//	}
//}
//
//func TestLoginHandler(t *testing.T) {
//	// generate a random user
//	user, password := randomUser(t)
//
//	// test cases
//	testCases := []struct {
//		name          string
//		reqBody       interface{}
//		loginRequest  *models.LoginRequest
//		loginResponse *models.LoginResponse
//		buildStubs    func(service *mocks.MockAuthService, request *models.LoginRequest, response *models.LoginResponse)
//		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
//	}{
//		{
//			name: "success case",
//			reqBody: gin.H{
//				"email":    user.Email,
//				"password": password,
//			},
//			loginRequest: &models.LoginRequest{
//				Email:    user.Email,
//				Password: password,
//			},
//			loginResponse: &models.LoginResponse{
//				UserResponse: models.UserResponse{
//					ID:          user.ID,
//					Name:        user.Name,
//					PhoneNumber: user.PhoneNumber,
//					Email:       user.Email,
//				},
//				AccessToken: "",
//			},
//			buildStubs: func(service *mocks.MockAuthService, request *models.LoginRequest, response *models.LoginResponse) {
//				service.EXPECT().LoginUser(request).Times(1).Return(response, nil)
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusOK, recorder.Code)
//			},
//		},
//		{
//			name: "invalid password case",
//			reqBody: gin.H{
//				"email":    user.Email,
//				"password": "invalid password",
//			},
//			loginRequest: &models.LoginRequest{
//				Email:    user.Email,
//				Password: "invalid password",
//			},
//			loginResponse: &models.LoginResponse{
//				UserResponse: models.UserResponse{
//					ID:          user.ID,
//					Name:        user.Name,
//					PhoneNumber: user.PhoneNumber,
//					Email:       user.Email,
//				},
//				AccessToken: "",
//			},
//			buildStubs: func(service *mocks.MockAuthService, request *models.LoginRequest, response *models.LoginResponse) {
//				service.EXPECT().LoginUser(request).Times(1).Return(nil, errors.ErrInvalidPassword)
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusUnauthorized, recorder.Code)
//				require.Contains(t, recorder.Body.String(), "invalid password")
//			},
//		},
//		{
//			name: "bad request case",
//			reqBody: gin.H{
//				"email": user.Email,
//			},
//			loginRequest:  nil,
//			loginResponse: nil,
//			buildStubs: func(service *mocks.MockAuthService, request *models.LoginRequest, response *models.LoginResponse) {
//				service.EXPECT().LoginUser(request).Times(0)
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusBadRequest, recorder.Code)
//				require.Contains(t, recorder.Body.String(), "Bad Request")
//			},
//		},
//		{
//			name: "invalid email case",
//			reqBody: gin.H{
//				"email":    "user@email.com",
//				"password": password,
//			},
//			loginRequest: &models.LoginRequest{
//				Email:    "user@email.com",
//				Password: password,
//			},
//			loginResponse: &models.LoginResponse{
//				UserResponse: models.UserResponse{
//					ID:          user.ID,
//					Name:        user.Name,
//					PhoneNumber: user.PhoneNumber,
//					Email:       user.Email,
//				},
//				AccessToken: "",
//			},
//			buildStubs: func(service *mocks.MockAuthService, request *models.LoginRequest, response *models.LoginResponse) {
//				service.EXPECT().LoginUser(request).Times(1).Return(nil, errors.New("invalid email", http.StatusUnprocessableEntity))
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
//				require.Contains(t, recorder.Body.String(), "invalid email")
//			},
//		},
//		{
//			name: "internal server error case",
//			reqBody: gin.H{
//				"email":    user.Email,
//				"password": password,
//			},
//			loginRequest: &models.LoginRequest{
//				Email:    user.Email,
//				Password: password,
//			},
//			loginResponse: &models.LoginResponse{
//				UserResponse: models.UserResponse{
//					ID:          user.ID,
//					Name:        user.Name,
//					PhoneNumber: user.PhoneNumber,
//					Email:       user.Email,
//				},
//				AccessToken: "",
//			},
//			buildStubs: func(service *mocks.MockAuthService, request *models.LoginRequest, response *models.LoginResponse) {
//				service.EXPECT().LoginUser(request).Times(1).Return(nil, errors.ErrInternalServerError)
//			},
//			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusInternalServerError, recorder.Code)
//				require.Contains(t, recorder.Body.String(), "internal server error")
//			},
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//			mockService := mocks.NewMockAuthService(ctrl)
//			testServer.handler.AuthService = mockService
//			tc.buildStubs(mockService, tc.loginRequest, tc.loginResponse)
//
//			jsonFile, err := json.Marshal(tc.reqBody)
//			require.NoError(t, err)
//			recorder := httptest.NewRecorder()
//			req, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(string(jsonFile)))
//			require.NoError(t, err)
//			testServer.router.ServeHTTP(recorder, req)
//			tc.checkResponse(t, recorder)
//		})
//	}
//}
//
//func randomUser(t *testing.T) (user models.User, password string) {
//	password = RandomString(6)
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	require.NoError(t, err)
//
//	user = models.User{
//		Model: models.Model{
//			ID:        6,
//			CreatedAt: RandomInt(1, 100),
//			UpdatedAt: RandomInt(1, 100),
//			DeletedAt: RandomInt(1, 100),
//		},
//		Name:           RandomOwnerName(),
//		HashedPassword: string(hashedPassword),
//		PhoneNumber:    RandomOwnerName(),
//		Email:          RandomEmail(),
//	}
//	return
//}
//
//func Test_Logout(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	auth := mocks.NewMockAuthService(ctrl)
//	repo := mocks.NewMockAuthRepository(ctrl)
//
//	conf, err := config.Load()
//	if err != nil {
//		t.Error(err)
//	}
//	user := &models.User{
//		Name:          "Tolu",
//		PhoneNumber:   "+2348163608141",
//		Email:         "toluwase@gmail.com",
//		Password:      "12345678",
//		IsEmailActive: true,
//	}
//	conf.JWTSecret = "testSecret"
//	token, err := jwt.GenerateToken(user.Email, conf.JWTSecret)
//
//	s := &Server{
//		Config:         conf,
//		AuthRepository: repo,
//		AuthService:    auth,
//	}
//
//	repo.EXPECT().AddToBlackList(&models.BlackList{Email: user.Email, Token: token}).Return(nil)
//	repo.EXPECT().TokenInBlacklist(token).Return(false)
//	repo.EXPECT().FindUserByEmail(user.Email).Return(user, nil)
//
//	r := s.setupRouter()
//	resp := httptest.NewRecorder()
//	req, _ := http.NewRequest(http.MethodGet, "/api/v1/logout", strings.NewReader(string(user.Email)))
//	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
//
//	r.ServeHTTP(resp, req)
//	fmt.Println(resp.Body.String())
//	assert.Equal(t, 200, resp.Code)
//}
//
//func init() {
//	rand.Seed(time.Now().UnixNano())
//}
//
//const alphabet = "abcdefghhijklmnopqrstuvwxyz"
//
//// RandomInt generates a random integer between min and max
//func RandomInt(min, max int64) int64 {
//	return min + rand.Int63n(max-min+1)
//}
//
//// RandomString geerates a random string of length n
//func RandomString(n int) string {
//	var sb strings.Builder
//	k := len(alphabet)
//
//	for i := 0; i < n; i++ {
//		c := alphabet[rand.Intn(k)]
//		sb.WriteByte(c)
//	}
//
//	return sb.String()
//}
//
//// RandomOwnerName generates random account owner names for testing
//func RandomOwnerName() string {
//	return RandomString(6)
//}
//
//// RandomEmail generates a random email address
//func RandomEmail() string {
//	return fmt.Sprintf("%s@email.com", RandomString(6))
//}
