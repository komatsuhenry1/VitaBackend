package nurse

import (
	"fmt"
	"medassist/internal/model"
	"medassist/internal/nurse/dto"
	"medassist/internal/repository"
	"medassist/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type NurseService interface {
	UpdateAvailablityNursingService(userId string) (model.Nurse, error)
	GetAllVisits(nurseId string) (dto.NurseVisitsListsDto, error)
	ConfirmOrCancelVisit(nurseId, visitId, reason string) (string, error)
	GetPatientProfile(patientId string) (dto.PatientProfileResponseDTO, error)
	NurseDashboardData(nurseId string) (dto.NurseDashboardDataResponseDTO, error)

}

type nurseService struct {
	userRepository  repository.UserRepository
	nurseRepository repository.NurseRepository
	visitRepository repository.VisitRepository
}

func NewNurseService(userRepository repository.UserRepository, nurseRepository repository.NurseRepository, visitRepository repository.VisitRepository) NurseService {
	return &nurseService{userRepository: userRepository, nurseRepository: nurseRepository, visitRepository: visitRepository}
}

func (s *nurseService) UpdateAvailablityNursingService(nurseId string) (model.Nurse, error) {

	//busca o user
	nurse, err := s.nurseRepository.FindNurseById(nurseId)
	if err != nil {
		return model.Nurse{}, fmt.Errorf("erro ao buscar user by id")
	}

	if nurse.Online {
		nurse.Online = false
	} else {
		nurse.Online = true
	}

	nurseUpdates := bson.M{
		"online":    nurse.Online,
		"updatedAt": time.Now(),
	}

	//salve user com status true/false
	nurse, err = s.nurseRepository.UpdateNurseFields(nurseId, nurseUpdates)
	if err != nil {
		return model.Nurse{}, fmt.Errorf("erro ao atualizar user")
	}

	return nurse, nil
}

func (s *nurseService) GetAllVisits(nurseId string) (dto.NurseVisitsListsDto, error) {
	visits, err := s.visitRepository.FindAllVisitsForNurse(nurseId)
	if err != nil {
		return dto.NurseVisitsListsDto{}, err
	}

	pendingVisits := make([]dto.VisitDto, 0)
	confirmedVisits := make([]dto.VisitDto, 0)
	completedVisits := make([]dto.VisitDto, 0)

	for _, visit := range visits {
		fmt.Println(visit)
		visitDto := dto.VisitDto{
			ID:          visit.ID.Hex(),
			Description: visit.Description,
			Reason:      visit.Reason,
			VisitType:   visit.VisitType,
			VisitValue:  visit.VisitValue,
			CreatedAt:   visit.CreatedAt.Format("02/01/2006 15:04"),
			Date:        visit.VisitDate.Format("02/01/2006 15:04"),
			Status:      visit.Status,
			PatientName: visit.PatientName,
			PatientId:   visit.PatientId,
			NurseName:   visit.NurseName,
		}

		switch visit.Status {
		case "PENDING":
			pendingVisits = append(pendingVisits, visitDto)
		case "CONFIRMED":
			confirmedVisits = append(confirmedVisits, visitDto)
		case "COMPLETED":
			completedVisits = append(completedVisits, visitDto)
		}
	}

	allVisitsDto := dto.NurseVisitsListsDto{
		Pending:   pendingVisits,
		Confirmed: confirmedVisits,
		Completed: completedVisits,
	}

	return allVisitsDto, nil
}

func (s *nurseService) ConfirmOrCancelVisit(nurseId, visitId, reason string) (string, error) {
	visit, err := s.visitRepository.FindVisitById(visitId)
	if err != nil {
		return "", err
	}

	var response string
	if visit.Status == "CONFIRMED" {
		visit.CancelReason = reason
		visit.Status = "PENDING"
		response = "Visita cancelada com sucesso."

		utils.SendEmailVisitCanceledWithReason("komatsuhenry@gmail.com", visit.NurseName, visit.VisitDate.Format("02/01/2006 15:04"), reason)
	} else if visit.Status == "PENDING" {
		visit.CancelReason = ""
		visit.Status = "CONFIRMED"
		response = "Visita confirmada com sucesso."

		utils.SendEmailVisitApproved("komatsuhenry@gmail.com", visit.NurseName, visit.VisitDate.Format("02/01/2006 15:04"), visit.VisitValue)
	}

	visitUpdates := bson.M{
		"status":        visit.Status,
		"cancel_reason": visit.CancelReason,
		"updated_at":    time.Now(),
	}

	_, err = s.visitRepository.UpdateVisitFields(visitId, visitUpdates)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (s *nurseService) GetPatientProfile(patientId string) (dto.PatientProfileResponseDTO, error) {

	patient, err := s.userRepository.FindUserById(patientId)
	if err != nil {
		return dto.PatientProfileResponseDTO{}, err
	}

	patientProfile := dto.PatientProfileResponseDTO{
		ID:          patient.ID,
		Name:        patient.Name,
		Email:       patient.Email,
		Phone:       patient.Phone,
		Address:     patient.Address,
		Cpf:         patient.Cpf,
		Password:    patient.Password,
		Hidden:      patient.Hidden,
		Role:        patient.Role,
		FirstAccess: patient.FirstAccess,
		CreatedAt:   patient.CreatedAt,
		TempCode:    patient.TempCode,
		UpdatedAt:   patient.UpdatedAt,
	}

	return patientProfile, nil
}

func (s *nurseService) NurseDashboardData(nurseId string) (dto.NurseDashboardDataResponseDTO, error) {
	var dashboardData dto.NurseDashboardDataResponseDTO

	nurse, err := s.nurseRepository.FindNurseById(nurseId)
	if err != nil {
		return dashboardData, err
	}

	visits, err := s.visitRepository.FindAllVisitsForNurse(nurseId)
	if err != nil {
		return dashboardData, err
	}

	// stats
	stats := dto.StatsDTO{
		PatientsAttended: 10,
		AppointmentsToday: 3,
		AverageRating: 4.5,
		MonthlyEarnings: 1000,
	}

	//itera pela lista de visits e mapeia para o dto VisitDto
	visitsDto := make([]dto.VisitDto, 0)
	for _, visit := range visits {
		visitsDto = append(visitsDto, dto.VisitDto{
			ID: visit.ID.Hex(),
			Description: visit.Description,
			Reason: visit.Reason,
			VisitType: visit.VisitType,
			VisitValue: visit.VisitValue,
			CreatedAt: visit.CreatedAt.Format("02/01/2006 15:04"),
			Date: visit.VisitDate.Format("02/01/2006 15:04"),
			Status: visit.Status,
			PatientName: visit.PatientName,
			PatientId: visit.PatientId,
			NurseName: visit.NurseName,
		})
	}

	//history


	//nurseProfile
	nurseProfile := dto.NurseProfileResponseDTO{
		Name: nurse.Name,
		Email: nurse.Email,
		Phone: nurse.Phone,
		Coren: nurse.LicenseNumber,
		ExperienceYears: nurse.YearsExperience,
		Department: nurse.Department,
		Bio: nurse.Bio,
	}

	nurseAvailability := dto.AvailabilityDTO{
		IsAvailable: nurse.Online,
		StartTime: nurse.StartTime,
		EndTime: nurse.EndTime,
		Specialization: nurse.Specialization,
	}

	dashboardData = dto.NurseDashboardDataResponseDTO{
		Online: nurse.Online,
		Stats: stats,
		Visits: visitsDto,
		// Patients: patientsDto,
		// History: historyDto,
		Profile: nurseProfile,
		Availability: nurseAvailability,
	}

	return dashboardData, nil


}