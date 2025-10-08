package admin

import (
	"fmt"
	"medassist/internal/admin/dto"
	"medassist/internal/repository"
	"medassist/utils"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type AdminService interface {
	ApproveNurseRegister(approvedUserId string) (string, error)
	GetNurseDocumentsToAnalisys(nurseID string) ([]dto.DocumentInfoResponse, error)
	GetFileStream(fileID primitive.ObjectID) (*gridfs.DownloadStream, error)
	GetDashboardData() (dto.DashboardAdminDataResponse, error)
	RejectNurseRegister(rejectedNurseId string, rejectDescription dto.RejectDescription) (string, error)
	UserLists() (dto.UserListsResponse, error)
	UpdateUser(userId string, updates map[string]interface{}) (dto.UserTypeResponse, error)
	DeleteNurseOrUser(userId string) error
	UpdateVisit(visitId string, updates map[string]interface{}) (dto.VisitTypeResponse, error)
	DeleteVisit(visitId string) error
}

type adminService struct {
	userRepository  repository.UserRepository
	nurseRepository repository.NurseRepository
	visitRepository repository.VisitRepository
}

func NewAdminService(userRepository repository.UserRepository, nurseRepository repository.NurseRepository, visitRepository repository.VisitRepository) AdminService {
	return &adminService{userRepository: userRepository, nurseRepository: nurseRepository, visitRepository: visitRepository}
}

func (s *adminService) ApproveNurseRegister(approvedNurseId string) (string, error) {
	nurse, err := s.nurseRepository.FindNurseById(approvedNurseId)
	if err != nil {
		return "", err
	}

	if nurse.Hidden {
		return "", fmt.Errorf("Usuário hidden.")
	}

	if nurse.Role != "NURSE" {
		return "", fmt.Errorf("Usuário não é Nurse.")
	}

	nurseUpdates := bson.M{
		"verification_seal": true,
		"updated_at":        time.Now(),
	}

	//salve user com status true/false
	nurse, err = s.nurseRepository.UpdateNurseFields(approvedNurseId, nurseUpdates)
	if err != nil {
		return "", fmt.Errorf("Erro ao atualizar user.")
	}

	err = utils.SendEmailApprovedNurse(nurse.Email)
	if err != nil {
		return "", err
	}

	return "Enfermeiro(a) aprovado(a) com sucesso.", nil
}

func (s *adminService) GetNurseDocumentsToAnalisys(nurseID string) ([]dto.DocumentInfoResponse, error) {
	nurse, err := s.nurseRepository.FindNurseById(nurseID)
	if err != nil {
		return nil, err
	}

	if nurse.Role != "NURSE" {
		return nil, fmt.Errorf("o usuário com ID '%s' não é um enfermeiro", nurseID)
	}

	var documents []dto.DocumentInfoResponse

	baseURL := os.Getenv("DOWNLOAD_URL")

	if !nurse.LicenseDocumentID.IsZero() {
		documents = append(documents, dto.DocumentInfoResponse{
			Name:        "Documento de Licença (COREN)",
			Type:        "license_document",
			DownloadURL: fmt.Sprintf("%s/%s", baseURL, nurse.LicenseDocumentID.Hex()),
			ImageID:     nurse.LicenseDocumentID.Hex(),
		})
	}
	if !nurse.QualificationsID.IsZero() {
		documents = append(documents, dto.DocumentInfoResponse{
			Name:        "Certificado de Qualificações",
			Type:        "qualifications",
			DownloadURL: fmt.Sprintf("%s/%s", baseURL, nurse.QualificationsID.Hex()),
			ImageID:     nurse.QualificationsID.Hex(),
		})
	}
	if !nurse.GeneralRegisterID.IsZero() {
		documents = append(documents, dto.DocumentInfoResponse{
			Name:        "Documento de Identidade (RG)",
			Type:        "general_register",
			DownloadURL: fmt.Sprintf("%s/%s", baseURL, nurse.GeneralRegisterID.Hex()),
			ImageID:     nurse.GeneralRegisterID.Hex(),
		})
	}
	if !nurse.ResidenceComprovantId.IsZero() {
		documents = append(documents, dto.DocumentInfoResponse{
			Name:        "Comprovante de Residência",
			Type:        "residence_comprovant",
			DownloadURL: fmt.Sprintf("%s/%s", baseURL, nurse.ResidenceComprovantId.Hex()),
			ImageID:     nurse.ResidenceComprovantId.Hex(),
		})
	}
	if !nurse.ProfileImageID.IsZero() {
		documents = append(documents, dto.DocumentInfoResponse{
			Name:        "Foto de perfil",
			Type:        "profile_image",
			DownloadURL: fmt.Sprintf("%s/%s", baseURL, nurse.ProfileImageID.Hex()),
			ImageID:     nurse.ProfileImageID.Hex(),
		})
	}

	return documents, nil
}

func (s *adminService) GetFileStream(fileID primitive.ObjectID) (*gridfs.DownloadStream, error) {
	return s.userRepository.DownloadFileByID(fileID)
}

func (s *adminService) GetDashboardData() (dto.DashboardAdminDataResponse, error) {
	var response dto.DashboardAdminDataResponse

	allUsers, err := s.userRepository.FindAllUsers()
	if err != nil {
		return response, err
	}

	allNurses, err := s.nurseRepository.FindAllNurses()
	if err != nil {
		return response, err
	}

	allNursesNotVerified, err := s.nurseRepository.FindAllNursesNotVerified()
	if err != nil {
		return response, err
	}

	var nursesFields []dto.NursesFieldsForDashboardResponse
	for _, nurse := range allNursesNotVerified {
		nursesFields = append(nursesFields, dto.NursesFieldsForDashboardResponse{
			ID:   nurse.ID.Hex(),
			Name: nurse.Name,
		})
	}

	adminDashboardData := dto.DashboardAdminDataResponse{
		TotalNurses:         len(allNurses),
		TotalPatients:       len(allUsers),
		NumberVisits:        100, // ADD
		VisitsToday:         100, // ADD
		PendentApprovations: len(allNursesNotVerified),
		NursesFields:        nursesFields,
	}

	return adminDashboardData, nil
}

func (s *adminService) RejectNurseRegister(rejectedNurseId string, rejectDescription dto.RejectDescription) (string, error) {
	nurse, err := s.nurseRepository.FindNurseById(rejectedNurseId)
	if err != nil {
		return "", err
	}

	err = utils.SendEmailRegistrationRejected(nurse.Email, rejectDescription.Description)
	if err != nil {
		return "", err
	}

	//possivel funcao que salva esse acontecimento no historico

	return "Enfermeiro(a) rejeitado com sucesso.", nil
}

func (s *adminService) UserLists() (dto.UserListsResponse, error) {

	users, err := s.userRepository.FindAllUsers()
	if err != nil {
		return dto.UserListsResponse{}, err
	}

	nurses, err := s.nurseRepository.FindAllNurses()
	if err != nil {
		return dto.UserListsResponse{}, err
	}

	visits, err := s.visitRepository.FindAllVisits()
	if err != nil {
		return dto.UserListsResponse{}, err
	}

	var userLists dto.UserListsResponse
	for _, user := range users {
		if user.Role == "PATIENT" {
			userLists.Users = append(userLists.Users, dto.UserTypeResponse{
				ID:          user.ID.Hex(),
				Name:        user.Name,
				Email:       user.Email,
				Role:        user.Role,
				Phone:       user.Phone,
				Address:     user.Address,
				Cpf:         user.Cpf,
				Password:    "",
				Hidden:      user.Hidden,
				FirstAccess: user.FirstAccess,
			})
		}
	}

	for _, nurse := range nurses {
		userLists.Nurses = append(userLists.Nurses, dto.NurseTypeResponse{
			ID:               nurse.ID.Hex(),
			Name:             nurse.Name,
			Email:            nurse.Email,
			Phone:            nurse.Phone,
			Address:          nurse.Address,
			Cpf:              nurse.Cpf,
			Password:         "",
			Hidden:           nurse.Hidden,
			FirstAccess:      nurse.FirstAccess,
			Role:             nurse.Role,
			VerificationSeal: nurse.VerificationSeal,
			LicenseNumber:    nurse.LicenseNumber,
			Specialization:   nurse.Specialization,
			Shift:            nurse.Shift,
			Department:       nurse.Department,
			YearsExperience:  nurse.YearsExperience,
			Price:            nurse.Price,
			Bio:              nurse.Bio,
		})
	}

	for _, visit := range visits {
		userLists.Visits = append(userLists.Visits, dto.VisitTypeResponse{
			ID:           visit.ID.Hex(),
			Status:       visit.Status,
			PatientId:    visit.PatientId,
			PatientName:  visit.PatientName,
			PatientEmail: visit.PatientEmail,
			Description:  visit.Description,
			Reason:       visit.Reason,
			CancelReason: visit.CancelReason,
			NurseId:      visit.NurseId,
			NurseName:    visit.NurseName,
			VisitValue:   visit.VisitValue,
			VisitType:    visit.VisitType,
			VisitDate:    visit.VisitDate,
		})
	}

	return userLists, nil
}

func (s *adminService) UpdateUser(userId string, updates map[string]interface{}) (dto.UserTypeResponse, error) {

	if emailRaw, ok := updates["email"]; ok {
		email, ok := emailRaw.(string)
		if ok {
			normalizedEmail := strings.ToLower(email)

			_, err := utils.EmailRegex(email)
			if err != nil {
				return dto.UserTypeResponse{}, fmt.Errorf("Email no formato incorreto.")
			}

			existingUser, err := s.userRepository.FindUserByEmail(normalizedEmail)
			// se nao achar em user, busca em nurseRepositoru
			if err != nil {
				existingUser, err = s.nurseRepository.FindNurseByEmail(normalizedEmail)
			}

			if err == nil && existingUser.ID.Hex() != userId {
				return dto.UserTypeResponse{}, fmt.Errorf("Email já está em uso por outro usuário")
			}

			updates["email"] = normalizedEmail
		}
	}

	if existingUser, err := s.userRepository.FindUserById(userId); err == nil && existingUser.Role == "PATIENT" {
		updated, err := s.userRepository.UpdateUser(userId, updates)
		if err != nil {
			return dto.UserTypeResponse{}, fmt.Errorf("erro ao atualizar campos do usuario: %w", err)
		}
		return dto.UserTypeResponse{
			Name:        updated.Name,
			Email:       updated.Email,
			Role:        updated.Role,
			Phone:       updated.Phone,
			Address:     updated.Address,
			Cpf:         updated.Cpf,
			Password:    updated.Password,
			Hidden:      updated.Hidden,
			FirstAccess: updated.FirstAccess,
		}, nil
	}

	if _, err := s.nurseRepository.FindNurseById(userId); err == nil {
		updated, err := s.nurseRepository.UpdateNurse(userId, updates)
		if err != nil {
			return dto.UserTypeResponse{}, fmt.Errorf("erro ao atualizar campos do enfermeiro(a): %w", err)
		}
		return dto.UserTypeResponse{
			Name:        updated.Name,
			Email:       updated.Email,
			Role:        updated.Role,
			Phone:       updated.Phone,
			Address:     updated.Address,
			Cpf:         updated.Cpf,
			Password:    updated.Password,
			Hidden:      updated.Hidden,
			FirstAccess: updated.FirstAccess,
		}, nil
	}

	return dto.UserTypeResponse{}, fmt.Errorf("usuário não encontrado")
}

func (s *adminService) UpdateVisit(visitId string, updates map[string]interface{}) (dto.VisitTypeResponse, error) {
	updated, err := s.visitRepository.UpdateVisitFields(visitId, updates)
	if err != nil {
		return dto.VisitTypeResponse{}, fmt.Errorf("erro ao atualizar campos do enfermeiro(a): %w", err)
	}

	updatedVisit := dto.VisitTypeResponse{
		Status:       updated.Status,
		PatientId:    updated.PatientId,
		PatientName:  updated.PatientName,
		PatientEmail: updated.PatientEmail,
		Description:  updated.Description,
		Reason:       updated.Reason,
		CancelReason: updated.CancelReason,
		NurseId:      updated.NurseId,
		NurseName:    updated.NurseName,
	}

	return updatedVisit, nil

}

func (s *adminService) DeleteNurseOrUser(userId string) error {
	if existingUser, err := s.userRepository.FindUserById(userId); err == nil && existingUser.Role == "PATIENT" {
		err := s.userRepository.DeleteUser(userId)
		if err != nil {
			return fmt.Errorf("erro ao deletar usuario: %w", err)
		}
	}

	if _, err := s.nurseRepository.FindNurseById(userId); err == nil {
		err := s.nurseRepository.DeleteNurse(userId)
		if err != nil {
			return fmt.Errorf("erro ao deletar enfermeiro(a): %w", err)
		}
	}

	return nil
}

func (s *adminService) DeleteVisit(visitId string) error {
	err := s.visitRepository.DeleteVisit(visitId)
	if err != nil {
		return fmt.Errorf("erro ao deletar visita: %w", err)
	}

	return nil
}