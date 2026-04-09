package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"medassist/internal/admin/dto"
	"medassist/internal/admin/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAdminHandler_AdminDashboard(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Erro_400_Service_Retorna_Erro", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdminService := mocks.NewMockAdminService(ctrl)
		handler := NewAdminHandler(mockAdminService)

		router := gin.Default()
		router.GET("/admin/dashboard", handler.AdminDashboard)

		mockAdminService.EXPECT().GetDashboardData().Return(dto.DashboardAdminDataResponse{}, fmt.Errorf("erro no banco"))

		req, _ := http.NewRequest(http.MethodGet, "/admin/dashboard", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Sucesso_200_Retorna_Dashboard", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdminService := mocks.NewMockAdminService(ctrl)
		handler := NewAdminHandler(mockAdminService)

		router := gin.Default()
		router.GET("/admin/dashboard", handler.AdminDashboard)

		mockData := dto.DashboardAdminDataResponse{
			TotalNurses: 10,
		}

		mockAdminService.EXPECT().GetDashboardData().Return(mockData, nil)

		req, _ := http.NewRequest(http.MethodGet, "/admin/dashboard", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, true, response["success"])
	})
}

func TestAdminHandler_ApproveNurseRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Erro_400_Servico_Aprovacao_Falha", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdminService := mocks.NewMockAdminService(ctrl)
		handler := NewAdminHandler(mockAdminService)

		router := gin.Default()
		router.PATCH("/admin/approve/:id", handler.ApproveNurseRegister)

		mockAdminService.EXPECT().ApproveNurseRegister("123").Return("", fmt.Errorf("Erro ao atualizar"))

		req, _ := http.NewRequest(http.MethodPatch, "/admin/approve/123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Sucesso_200_Aprova_Enfermeiro", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdminService := mocks.NewMockAdminService(ctrl)
		handler := NewAdminHandler(mockAdminService)

		router := gin.Default()
		router.PATCH("/admin/approve/:id", handler.ApproveNurseRegister)

		mockAdminService.EXPECT().ApproveNurseRegister("123").Return("Enfermeiro(a) aprovado(a) com sucesso.", nil)

		req, _ := http.NewRequest(http.MethodPatch, "/admin/approve/123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAdminHandler_UsersManagement(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Sucesso_200_Retorna_Listas", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdminService := mocks.NewMockAdminService(ctrl)
		handler := NewAdminHandler(mockAdminService)

		router := gin.Default()
		router.GET("/admin/users", handler.UsersManagement)

		mockAdminService.EXPECT().UserLists().Return(dto.UserListsResponse{}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/admin/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAdminHandler_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Erro_400_Protegido", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdminService := mocks.NewMockAdminService(ctrl)
		handler := NewAdminHandler(mockAdminService)

		router := gin.Default()
		router.PATCH("/admin/user/:id", handler.UpdateUser)

		body := []byte(`{"id": "hack_tentativa"}`)
		req, _ := http.NewRequest(http.MethodPatch, "/admin/user/123", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Sucesso_200_Atualiza_User", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdminService := mocks.NewMockAdminService(ctrl)
		handler := NewAdminHandler(mockAdminService)

		router := gin.Default()
		router.PATCH("/admin/user/:id", handler.UpdateUser)

		mockAdminService.EXPECT().UpdateUser("123", map[string]interface{}{"name": "Novo"}).Return(dto.UserTypeResponse{}, nil)

		body := []byte(`{"name": "Novo"}`)
		req, _ := http.NewRequest(http.MethodPatch, "/admin/user/123", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
