package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"medassist/internal/auth/dto"
	"medassist/internal/auth/mocks"
	"medassist/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthHandler_LoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Deve_retornar_erro_400_quando_o_body_estiver_mal_formatado_ou_faltando_campos", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.POST("/auth/login", handler.LoginUser)

		body := []byte(`{"email": "test@email.com"}`)
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, false, response["success"])
		assert.Equal(t, "Requisição inválida", response["message"])
	})

	t.Run("Deve_retornar_erro_400_quando_as_credenciais_estiverem_invalidas", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.POST("/auth/login", handler.LoginUser)

		loginDTO := dto.LoginRequestDTO{
			Email:    "test@email.com",
			Password: "senha_errada_123",
		}
		mockAuthService.EXPECT().LoginUser(loginDTO).Return("", dto.AuthUser{}, fmt.Errorf("Credenciais inválidas. Tente novamente."))

		body, _ := json.Marshal(loginDTO)
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, false, response["success"])
		assert.Equal(t, "Credenciais inválidas. Tente novamente.", response["message"])
	})

	t.Run("Deve_retornar_200_com_token_e_dados_do_usuario_em_caso_de_sucesso", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.POST("/auth/login", handler.LoginUser)

		loginDTO := dto.LoginRequestDTO{
			Email:    "test@email.com",
			Password: "senha_correta_123",
		}
		fakeAuthUserResponse := dto.AuthUser{
			Email: "test@email.com",
			Role:  "PATIENT",
		}
		mockAuthService.EXPECT().LoginUser(loginDTO).Return("fake_jwt_token_123", fakeAuthUserResponse, nil)

		body, _ := json.Marshal(loginDTO)
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, true, response["success"])
		assert.Equal(t, "Usuário logado com sucesso.", response["message"])
	})
}

func TestAuthHandler_UserRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Erro_400_Formulario_multipart_nao_enviado", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.POST("/auth/user", handler.UserRegister)

		req, _ := http.NewRequest(http.MethodPost, "/auth/user", bytes.NewBuffer([]byte(`{"name":"test"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Sucesso_200_Usuario_criado_com_sucesso", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.POST("/auth/user", handler.UserRegister)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("name", "Henry")
		_ = writer.WriteField("email", "henry@test.com")
		_ = writer.Close()

		mockAuthService.EXPECT().
			UserRegister(gomock.Any(), gomock.Any()).
			Return(model.User{Email: "henry@test.com", Name: "Henry"}, nil)

		req, _ := http.NewRequest(http.MethodPost, "/auth/user", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, true, response["success"])
	})
}

func TestAuthHandler_SendCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Erro_400_Email_vazio", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.PATCH("/auth/code", handler.SendCode)

		dtoObj := dto.EmailAuthRequestDTO{Email: ""}
		mockAuthService.EXPECT().SendCodeToEmail(dtoObj).Return(dto.CodeResponseDTO{}, fmt.Errorf("email invalido"))

		body := []byte(`{}`)
		req, _ := http.NewRequest(http.MethodPatch, "/auth/code", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Sucesso_200_Codigo_enviado", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.PATCH("/auth/code", handler.SendCode)

		dtoObj := dto.EmailAuthRequestDTO{Email: "test@email.com"}
		mockAuthService.EXPECT().SendCodeToEmail(dtoObj).Return(dto.CodeResponseDTO{Code: 123456}, nil)

		body, _ := json.Marshal(dtoObj)
		req, _ := http.NewRequest(http.MethodPatch, "/auth/code", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthHandler_ValidateCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Erro_400_Codigo_invalido", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.POST("/auth/validate", handler.ValidateCode)

		dtoObj := dto.InputCodeDto{Email: "test@test.com", Code: 00000}
		mockAuthService.EXPECT().ValidateUserCode(dtoObj).Return("", dto.AuthUser{}, fmt.Errorf("codigo incorreto"))

		body, _ := json.Marshal(dtoObj)
		req, _ := http.NewRequest(http.MethodPost, "/auth/validate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Sucesso_200_Codigo_validado", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.POST("/auth/validate", handler.ValidateCode)

		dtoObj := dto.InputCodeDto{Email: "test@test.com", Code: 123456}
		fakeAuthUserResponse := dto.AuthUser{Email: "test@test.com", Role: "PATIENT"}

		mockAuthService.EXPECT().ValidateUserCode(dtoObj).Return("token-123", fakeAuthUserResponse, nil)

		body, _ := json.Marshal(dtoObj)
		req, _ := http.NewRequest(http.MethodPost, "/auth/validate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthHandler_FirstLoginAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Erro_500_Erro_ao_criar_admin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.POST("/auth/adm", handler.FirstLoginAdmin)

		mockAuthService.EXPECT().FirstLoginAdmin().Return(fmt.Errorf("admin ja existe"))

		req, _ := http.NewRequest(http.MethodPost, "/auth/adm", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Sucesso_200_Admin_inicial_criado", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.POST("/auth/adm", handler.FirstLoginAdmin)

		mockAuthService.EXPECT().FirstLoginAdmin().Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/auth/adm", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthHandler_SendEmailForgotPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Sucesso_200_Email_de_recuperacao_enviado", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.POST("/auth/email", handler.SendEmailForgotPassword)

		dtoObj := dto.ForgotPasswordRequestDTO{Email: "test@test.com"}
		mockAuthService.EXPECT().SendEmailForgotPassword(dtoObj).Return(nil)

		body, _ := json.Marshal(dtoObj)
		req, _ := http.NewRequest(http.MethodPost, "/auth/email", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthHandler_ResetPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Sucesso_200_Senha_resetada_com_sucesso", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.POST("/auth/reset-password", handler.ResetPassword)

		dtoObj := dto.ResetPasswordDTO{Token: "token123", NewPassword: "newpassword123"}
		mockAuthService.EXPECT().ResetPassword(dtoObj).Return(nil)

		body, _ := json.Marshal(dtoObj)
		req, _ := http.NewRequest(http.MethodPost, "/auth/reset-password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthHandler_ChangePasswordLogged(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Sucesso_200_Senha_de_usuario_logado_alterada", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		
		router := gin.Default()
		router.Use(func(c *gin.Context) {
            c.Set("claims", jwt.MapClaims{"sub": "id-mongo"})
            c.Next()
        })
		router.PATCH("/auth/logged/password", handler.ChangePasswordLogged)

		dtoObj := dto.ChangePasswordBothRequestDTO{
			Password: "old_password",
			NewPassword: "new_password_123",
			TwoFactor: false,
		}
		
		mockAuthService.EXPECT().ChangePasswordLogged(dtoObj, "id-mongo").Return(true, nil)

		body, _ := json.Marshal(dtoObj)
		req, _ := http.NewRequest(http.MethodPatch, "/auth/logged/password", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, true, response["success"])
		assert.Equal(t, "Senha atualizada com sucesso.", response["message"])
	})
}

func TestAuthHandler_ChangePasswordUnlogged(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Sucesso_200_Senha_deslogado_alterada", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockAuthService := mocks.NewMockAuthService(ctrl)
		handler := NewAuthHandler(mockAuthService)
		router := gin.Default()
		router.PATCH("/auth/password/:id", handler.ChangePasswordUnlogged)

		dtoObj := dto.UpdatedPasswordByNewPassword{NewPassword: "newpassword123"}
		mockAuthService.EXPECT().ChangePasswordUnlogged(dtoObj, "123").Return(nil)

		body, _ := json.Marshal(dtoObj)
		req, _ := http.NewRequest(http.MethodPatch, "/auth/password/123", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
