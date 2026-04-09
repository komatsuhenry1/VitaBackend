package admin

import (
	"fmt"
	"testing"

	authDto "medassist/internal/auth/dto"
	"medassist/internal/model"
	repmocks "medassist/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestAdminService_DeleteVisit(t *testing.T) {
	t.Run("Erro_Falha_No_Repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		mockVisitRepo := repmocks.NewMockVisitRepository(ctrl)

		service := NewAdminService(mockUserRepo, mockNurseRepo, mockVisitRepo)

		mockVisitRepo.EXPECT().DeleteVisit("123").Return(fmt.Errorf("visita não encontrada"))

		err := service.DeleteVisit("123")

		assert.Error(t, err)
		assert.EqualError(t, err, "erro ao deletar visita: visita não encontrada")
	})

	t.Run("Sucesso_Visita_Deletada", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		mockVisitRepo := repmocks.NewMockVisitRepository(ctrl)

		service := NewAdminService(mockUserRepo, mockNurseRepo, mockVisitRepo)

		mockVisitRepo.EXPECT().DeleteVisit("123").Return(nil)

		err := service.DeleteVisit("123")

		assert.NoError(t, err)
	})
}

func TestAdminService_UserLists(t *testing.T) {
	t.Run("Erro_UserRepository_Falha", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		mockVisitRepo := repmocks.NewMockVisitRepository(ctrl)

		service := NewAdminService(mockUserRepo, mockNurseRepo, mockVisitRepo)

		mockUserRepo.EXPECT().FindAllUsers().Return(nil, fmt.Errorf("banco caiu"))

		_, err := service.UserLists()

		assert.Error(t, err)
	})

	t.Run("Sucesso_Busca_Filtrada", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		mockVisitRepo := repmocks.NewMockVisitRepository(ctrl)

		service := NewAdminService(mockUserRepo, mockNurseRepo, mockVisitRepo)

		fakeUser1 := model.User{
			Role: "PATIENT",
		}
		fakeUser2 := model.User{
			Role: "ADMIN",
		}

		mockUserRepo.EXPECT().FindAllUsers().Return([]model.User{fakeUser1, fakeUser2}, nil)
		mockNurseRepo.EXPECT().FindAllNurses().Return([]model.Nurse{{Name: "Nurse"}}, nil)
		mockVisitRepo.EXPECT().FindAllVisits().Return([]model.Visit{}, nil)

		lists, err := service.UserLists()

		assert.NoError(t, err)
		assert.Len(t, lists.Users, 1)
		assert.Len(t, lists.Nurses, 1)
		assert.Len(t, lists.Visits, 0)
	})
}

func TestAdminService_DeleteNurseOrUser(t *testing.T) {
	t.Run("Sucesso_Deleta_Paciente", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		mockVisitRepo := repmocks.NewMockVisitRepository(ctrl)

		service := NewAdminService(mockUserRepo, mockNurseRepo, mockVisitRepo)

		fakeUser := model.User{Role: "PATIENT"}
		
		mockUserRepo.EXPECT().FindUserById("usr123").Return(fakeUser, nil)
		mockNurseRepo.EXPECT().FindNurseById("usr123").Return(model.Nurse{}, fmt.Errorf("não achou nurse"))
		
		mockUserRepo.EXPECT().DeleteUser("usr123").Return(nil)

		err := service.DeleteNurseOrUser("usr123")
		assert.NoError(t, err)
	})
}

func TestAdminService_UpdateUser(t *testing.T) {
	t.Run("Erro_Email_Em_Uso", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		mockVisitRepo := repmocks.NewMockVisitRepository(ctrl)

		service := NewAdminService(mockUserRepo, mockNurseRepo, mockVisitRepo)

		updates := map[string]interface{}{"email": "existente@test.com"}
		
		fakeUser := authDto.AuthUser{ID: primitive.NewObjectID()}

		mockUserRepo.EXPECT().FindUserByEmail("existente@test.com").Return(fakeUser, nil)

		_, err := service.UpdateUser("id-1234", updates)
		assert.Error(t, err)
		assert.EqualError(t, err, "Email já está em uso por outro usuário")
	})

	t.Run("Erro_Usuario_Nao_Encontrado_Ambos_Repos", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		mockVisitRepo := repmocks.NewMockVisitRepository(ctrl)

		service := NewAdminService(mockUserRepo, mockNurseRepo, mockVisitRepo)

		mockUserRepo.EXPECT().FindUserById("id-1234").Return(model.User{}, fmt.Errorf("n existe"))
		mockNurseRepo.EXPECT().FindNurseById("id-1234").Return(model.Nurse{}, fmt.Errorf("n existe"))

		_, err := service.UpdateUser("id-1234", map[string]interface{}{"name": "Novo"})
		
		assert.Error(t, err)
		assert.EqualError(t, err, "usuário não encontrado")
	})

	t.Run("Sucesso_Atualiza_Paciente_Encontrado", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUserRepo := repmocks.NewMockUserRepository(ctrl)
		mockNurseRepo := repmocks.NewMockNurseRepository(ctrl)
		mockVisitRepo := repmocks.NewMockVisitRepository(ctrl)

		service := NewAdminService(mockUserRepo, mockNurseRepo, mockVisitRepo)

		fakeUser := model.User{Role: "PATIENT"}

		mockUserRepo.EXPECT().FindUserById("id-1234").Return(fakeUser, nil)
		
		updatedUser := model.User{Name: "Atualizado", Role: "PATIENT"}
		mockUserRepo.EXPECT().UpdateUser("id-1234", map[string]interface{}{"name": "Atualizado"}).Return(updatedUser, nil)

		resp, err := service.UpdateUser("id-1234", map[string]interface{}{"name": "Atualizado"})
		
		assert.NoError(t, err)
		assert.Equal(t, "Atualizado", resp.Name)
	})
}
