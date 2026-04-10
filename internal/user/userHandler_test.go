package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"medassist/internal/user/dto"
	"medassist/internal/user/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_GetAllNurses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Erro_500_Service_Retorna_Erro", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserService := mocks.NewMockUserService(ctrl)
		handler := NewUserHandler(mockUserService)
		router := gin.Default()
		
		router.Use(func(c *gin.Context) {
			c.Set("claims", jwt.MapClaims{"sub": "id-paciente-123"})
			c.Next()
		})
		router.GET("/user/all_nurses", handler.GetAllNurses)

		mockUserService.EXPECT().GetAllNurses("id-paciente-123").Return([]dto.AllNursesListDto{}, fmt.Errorf("Erro simulado"))

		req, _ := http.NewRequest(http.MethodGet, "/user/all_nurses", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, false, response["success"])
	})

	t.Run("Sucesso_200_Retorna_Enfermeiros", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserService := mocks.NewMockUserService(ctrl)
		handler := NewUserHandler(mockUserService)
		router := gin.Default()
		
		router.Use(func(c *gin.Context) {
			c.Set("claims", jwt.MapClaims{"sub": "id-paciente-123"})
			c.Next()
		})
		router.GET("/user/all_nurses", handler.GetAllNurses)

		mockResponse := []dto.AllNursesListDto{{ID: "123", Name: "Nurse 1"}}
		mockUserService.EXPECT().GetAllNurses("id-paciente-123").Return(mockResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/user/all_nurses", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, true, response["success"])
	})
}

func TestUserHandler_GetNurseProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Erro_400_Falha_Buscar_Perfil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserService := mocks.NewMockUserService(ctrl)
		handler := NewUserHandler(mockUserService)
		router := gin.Default()
		router.GET("/user/nurse/:id", handler.GetNurseProfile)

		mockUserService.EXPECT().GetNurseProfile("666").Return(dto.NurseProfileResponseDTO{}, fmt.Errorf("Erro db"))

		req, _ := http.NewRequest(http.MethodGet, "/user/nurse/666", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Sucesso_200_Retorna_Perfil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserService := mocks.NewMockUserService(ctrl)
		handler := NewUserHandler(mockUserService)
		router := gin.Default()
		router.GET("/user/nurse/:id", handler.GetNurseProfile)

		mockResponse := dto.NurseProfileResponseDTO{Name: "Nurse Ok"}
		mockUserService.EXPECT().GetNurseProfile("777").Return(mockResponse, nil)

		req, _ := http.NewRequest(http.MethodGet, "/user/nurse/777", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUserHandler_AddReview(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Erro_400_JSON_Invalido", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserService := mocks.NewMockUserService(ctrl)
		handler := NewUserHandler(mockUserService)
		router := gin.Default()
		router.POST("/user/review/:id", handler.AddReview)

		req, _ := http.NewRequest(http.MethodPost, "/user/review/visit-123", bytes.NewBuffer([]byte(`invalido`)))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Sucesso_200_Adiciona_Review", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserService := mocks.NewMockUserService(ctrl)
		handler := NewUserHandler(mockUserService)
		router := gin.Default()
		
		router.Use(func(c *gin.Context) {
			c.Set("claims", jwt.MapClaims{"sub": "id-paciente"})
			c.Next()
		})
		router.POST("/user/review/:id", handler.AddReview)

		reviewDTO := dto.ReviewDTO{Rating: 5, Comment: "Ótimo!"}
		mockUserService.EXPECT().AddReview("id-paciente", "visit-123", reviewDTO).Return(nil)

		body, _ := json.Marshal(reviewDTO)
		req, _ := http.NewRequest(http.MethodPost, "/user/review/visit-123", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
