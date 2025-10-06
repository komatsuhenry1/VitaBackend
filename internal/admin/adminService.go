package admin

import (
	"fmt"
	"medassist/internal/admin/dto"
	"medassist/internal/repository"
	"medassist/utils"
	"os"
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
}

type adminService struct {
	userRepository  repository.UserRepository
	nurseRepository repository.NurseRepository
}

func NewAdminService(userRepository repository.UserRepository, nurseRepository repository.NurseRepository) AdminService {
	return &adminService{userRepository: userRepository, nurseRepository: nurseRepository}
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
		"updated_at":         time.Now(),
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

	// 2. Monta a URL base para os downloads. Em um ambiente real, isso viria de uma variável de ambiente.
	baseURL := os.Getenv("DOWNLOAD_URL")

	// 3. Verifica cada campo de documento e, se existir, adiciona à lista de resposta.
	if !nurse.LicenseDocumentID.IsZero() {
		documents = append(documents, dto.DocumentInfoResponse{
			Name:        "Documento de Licença (COREN)",
			Type:        "license_document",
			DownloadURL: fmt.Sprintf("%s/%s", baseURL, nurse.LicenseDocumentID.Hex()),
			ImageID: nurse.LicenseDocumentID.Hex(),
		})
	}
	if !nurse.QualificationsID.IsZero() {
		documents = append(documents, dto.DocumentInfoResponse{
			Name:        "Certificado de Qualificações",
			Type:        "qualifications",
			DownloadURL: fmt.Sprintf("%s/%s", baseURL, nurse.QualificationsID.Hex()),
			ImageID: nurse.QualificationsID.Hex(),
		})
	}
	if !nurse.GeneralRegisterID.IsZero() {
		documents = append(documents, dto.DocumentInfoResponse{
			Name:        "Documento de Identidade (RG)",
			Type:        "general_register",
			DownloadURL: fmt.Sprintf("%s/%s", baseURL, nurse.GeneralRegisterID.Hex()),
			ImageID: nurse.GeneralRegisterID.Hex(),
		})
	}
	if !nurse.ResidenceComprovantId.IsZero() {
		documents = append(documents, dto.DocumentInfoResponse{
			Name:        "Comprovante de Residência",
			Type:        "residence_comprovant",
			DownloadURL: fmt.Sprintf("%s/%s", baseURL, nurse.ResidenceComprovantId.Hex()),
			ImageID: nurse.ResidenceComprovantId.Hex(),
		})
	}

	return documents, nil
}

func (s *adminService) GetFileStream(fileID primitive.ObjectID) (*gridfs.DownloadStream, error) {
	return s.userRepository.DownloadFileByID(fileID)
}

func (s *adminService) GetDashboardData() (dto.DashboardAdminDataResponse, error) {
	var response dto.DashboardAdminDataResponse

	// nursesIDsPendents, err := s.nurseRepository.GetIdsNursesPendents()
	// if err != nil {
	// 	return response, err
	// }

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
