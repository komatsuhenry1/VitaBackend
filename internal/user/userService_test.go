package user

import (
	"fmt"
	"testing"

	"medassist/internal/model"
	repmocks "medassist/internal/repository/mocks"
	"medassist/internal/user/dto"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_GetAllNurses(t *testing.T) {
	t.Run("Erro_Falha_Encontrar_Paciente", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repmocks.NewMockUserRepository(ctrl)
		nurseRepo := repmocks.NewMockNurseRepository(ctrl)
		visitRepo := repmocks.NewMockVisitRepository(ctrl)
		reviewRepo := repmocks.NewMockReviewRepository(ctrl)

		service := NewUserService(userRepo, nurseRepo, visitRepo, reviewRepo, nil)

		userRepo.EXPECT().FindUserById("id-invalido").Return(model.User{}, fmt.Errorf("Erro db"))

		_, err := service.GetAllNurses("id-invalido")
		
		assert.Error(t, err)
		assert.EqualError(t, err, "Erro ao buscar id de paciente.")
	})

	t.Run("Sucesso_Retorna_Enfermeiros_Da_Cidade", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repmocks.NewMockUserRepository(ctrl)
		nurseRepo := repmocks.NewMockNurseRepository(ctrl)
		visitRepo := repmocks.NewMockVisitRepository(ctrl)
		reviewRepo := repmocks.NewMockReviewRepository(ctrl)

		service := NewUserService(userRepo, nurseRepo, visitRepo, reviewRepo, nil)

		fakeUser := model.User{City: "São Paulo"}
		userRepo.EXPECT().FindUserById("paciente-sp").Return(fakeUser, nil)

		fakeNurses := []dto.AllNursesListDto{{Name: "Nurse 1"}}
		nurseRepo.EXPECT().GetAllNurses("São Paulo").Return(fakeNurses, nil)

		resp, err := service.GetAllNurses("paciente-sp")
		
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "Nurse 1", resp[0].Name)
	})
}

func TestUserService_GetOnlineNurses(t *testing.T) {
	t.Run("Erro_Falha_Busca_Paciente_Retorna_Lista_Vazia", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repmocks.NewMockUserRepository(ctrl)
		nurseRepo := repmocks.NewMockNurseRepository(ctrl)
		visitRepo := repmocks.NewMockVisitRepository(ctrl)
		reviewRepo := repmocks.NewMockReviewRepository(ctrl)

		service := NewUserService(userRepo, nurseRepo, visitRepo, reviewRepo, nil)

		userRepo.EXPECT().FindUserById("id-invalido").Return(model.User{}, fmt.Errorf("Erro db"))

		resp, err := service.GetOnlineNurses("id-invalido")
		
		assert.NoError(t, err)
		assert.Len(t, resp, 0)
	})

	t.Run("Sucesso_Retorna_Enfermeiros_Online_Na_Area", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repmocks.NewMockUserRepository(ctrl)
		nurseRepo := repmocks.NewMockNurseRepository(ctrl)
		visitRepo := repmocks.NewMockVisitRepository(ctrl)
		reviewRepo := repmocks.NewMockReviewRepository(ctrl)

		service := NewUserService(userRepo, nurseRepo, visitRepo, reviewRepo, nil)

		fakeUser := model.User{City: "Santos", Latitude: 12.3, Longitude: 45.6}
		userRepo.EXPECT().FindUserById("paciente-santos").Return(fakeUser, nil)

		fakeNurses := []dto.AllNursesListDto{{Name: "Nurse Online"}}
		nurseRepo.EXPECT().GetAllOnlineNurses("Santos", 12.3, 45.6).Return(fakeNurses, nil)

		resp, err := service.GetOnlineNurses("paciente-santos")
		
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "Nurse Online", resp[0].Name)
	})
}

func TestUserService_GetNurseProfile(t *testing.T) {
	t.Run("Erro_Enfermeiro_Com_Cadastro_Incompleto", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repmocks.NewMockUserRepository(ctrl)
		nurseRepo := repmocks.NewMockNurseRepository(ctrl)
		visitRepo := repmocks.NewMockVisitRepository(ctrl)
		reviewRepo := repmocks.NewMockReviewRepository(ctrl)

		service := NewUserService(userRepo, nurseRepo, visitRepo, reviewRepo, nil)

		fakeNurse := model.Nurse{
			Name: "Incompleto",
			MaxPatientsPerDay: 0,
		}
		
		nurseRepo.EXPECT().FindNurseById("nurse-1").Return(fakeNurse, nil)

		_, err := service.GetNurseProfile("nurse-1")
		
		assert.Error(t, err)
		assert.EqualError(t, err, "O enfermeiro ainda não preencheu os dados necessários para ser visto por pacientes.")
	})

	t.Run("Sucesso_Montagem_De_Perfil_Completo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepo := repmocks.NewMockUserRepository(ctrl)
		nurseRepo := repmocks.NewMockNurseRepository(ctrl)
		visitRepo := repmocks.NewMockVisitRepository(ctrl)
		reviewRepo := repmocks.NewMockReviewRepository(ctrl)

		service := NewUserService(userRepo, nurseRepo, visitRepo, reviewRepo, nil)

		fakeNurse := model.Nurse{
			Name: "Nurse Completo",
			MaxPatientsPerDay: 5,
			DaysAvailable: []string{"Monday"},
			Services: []string{"Curativo"},
			AvailableNeighborhoods: []string{"Centro"},
		}
		
		nurseRepo.EXPECT().FindNurseById("nurse-ok").Return(fakeNurse, nil)
		reviewRepo.EXPECT().FindAverageRatingByNurseId("nurse-ok").Return(4.5, nil)
		
		fakeReviews := []model.Review{{PatientName: "John", Rating: 5, Comment: "Foi mt bom"}}
		reviewRepo.EXPECT().FindAllNurseReviews("nurse-ok").Return(fakeReviews, nil)

		resp, err := service.GetNurseProfile("nurse-ok")
		
		assert.NoError(t, err)
		assert.Equal(t, "Nurse Completo", resp.Name)
		assert.Equal(t, 4.5, resp.Rating)
		assert.Len(t, resp.Reviews, 1)
		assert.Equal(t, "John", resp.Reviews[0].PatientName)
	})
}
