package auth

import (
	"fmt"
	"testing"

	"medassist/internal/auth/dto"
	"medassist/internal/model"
	repmocks "medassist/internal/repository/mocks"
	"medassist/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestAuthService_LoginUser(t *testing.T) {
	t.Run("Erro_Email_Invalido_Banco", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		service := NewAuthService(mockUserRepo, mockNurseRepo)

		mockUserRepo.EXPECT().FindUserByEmail("invalido@test.com").Return(dto.AuthUser{}, fmt.Errorf("usuário não encontrado"))
		mockNurseRepo.EXPECT().FindNurseByEmail("invalido@test.com").Return(dto.AuthUser{}, fmt.Errorf("usuário não encontrado"))

		_, _, err := service.LoginUser(dto.LoginRequestDTO{
			Email:    "invalido@test.com",
			Password: "Asd@123()",
		})

		assert.Error(t, err)
		assert.EqualError(t, err, "email não cadastrado")
	})

	t.Run("Erro_Senha_Incorreta", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		service := NewAuthService(mockUserRepo, mockNurseRepo)

		hashedPassword, _ := utils.HashPassword("SenhaCorreta@123")

		fakeUser := dto.AuthUser{
			ID:       primitive.NewObjectID(),
			Email:    "test@test.com",
			Password: string(hashedPassword),
			Role:     "PATIENT",
		}

		mockUserRepo.EXPECT().FindUserByEmail("test@test.com").Return(fakeUser, nil)

		_, _, err := service.LoginUser(dto.LoginRequestDTO{
			Email:    "test@test.com",
			Password: "SenhaErrada@123",
		})

		assert.Error(t, err)
		assert.EqualError(t, err, "Credenciais inválidas. Tente novamente.")
	})

	t.Run("Sucesso_Login_Paciente", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		service := NewAuthService(mockUserRepo, mockNurseRepo)

		hashedPassword, _ := utils.HashPassword("SenhaCorreta@123")

		fakeUser := dto.AuthUser{
			ID:       primitive.NewObjectID(),
			Email:    "test@test.com",
			Password: string(hashedPassword),
			Role:     "PATIENT",
			Hidden:   false,
		}

		mockUserRepo.EXPECT().FindUserByEmail("test@test.com").Return(fakeUser, nil)

		token, userObj, err := service.LoginUser(dto.LoginRequestDTO{
			Email:    "test@test.com",
			Password: "SenhaCorreta@123",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, "PATIENT", userObj.Role)
	})

	t.Run("Sucesso_Login_Enfermeiro", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		service := NewAuthService(mockUserRepo, mockNurseRepo)

		hashedPassword, _ := utils.HashPassword("SenhaCorreta@123")

		fakeNurse := dto.AuthUser{
			ID:               primitive.NewObjectID(),
			Email:            "nurse@test.com",
			Password:         string(hashedPassword),
			Role:             "NURSE",
			VerificationSeal: true,
			Hidden:           false,
		}

		mockUserRepo.EXPECT().FindUserByEmail("nurse@test.com").Return(dto.AuthUser{}, fmt.Errorf("usuário não encontrado"))
		mockNurseRepo.EXPECT().FindNurseByEmail("nurse@test.com").Return(fakeNurse, nil)
		
		mockNurseRepo.EXPECT().UpdateNurseFields(fakeNurse.ID.Hex(), gomock.Any()).Return(model.Nurse{}, nil)

		token, userObj, err := service.LoginUser(dto.LoginRequestDTO{
			Email:    "nurse@test.com",
			Password: "SenhaCorreta@123",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, "NURSE", userObj.Role)
	})
}

func TestAuthService_FirstLoginAdmin(t *testing.T) {
	t.Run("Erro_Admin_Ja_Existente", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		service := NewAuthService(mockUserRepo, mockNurseRepo)

		t.Setenv("ADMIN_PASSWORD", "Adm@123")
		t.Setenv("ADMIN_NAME", "Admin Test")
		t.Setenv("ADMIN_EMAIL", "admin@admin.com")

		mockUserRepo.EXPECT().UserExistsByEmail("admin@admin.com").Return(true, nil)

		err := service.FirstLoginAdmin()

		assert.Error(t, err)
		assert.EqualError(t, err, "o usuário já existe")
	})

	t.Run("Sucesso_Criacao_Admin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		service := NewAuthService(mockUserRepo, mockNurseRepo)

		t.Setenv("ADMIN_PASSWORD", "Adm@123")
		t.Setenv("ADMIN_NAME", "Admin Test")
		t.Setenv("ADMIN_EMAIL", "admin@admin.com")

		mockUserRepo.EXPECT().UserExistsByEmail("admin@admin.com").Return(false, nil)
		mockUserRepo.EXPECT().CreateUser(gomock.AssignableToTypeOf(&model.User{})).Return(nil)

		err := service.FirstLoginAdmin()

		assert.NoError(t, err)
	})
}
