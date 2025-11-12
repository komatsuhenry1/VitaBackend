package nurse

import (
	"fmt"
	"math"
	"medassist/internal/model"
	"medassist/internal/nurse/dto"
	"medassist/internal/repository"
	userDTO "medassist/internal/user/dto"
	"medassist/utils"
	"time"

	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NurseService interface {
	UpdateAvailablityNursingService(userId string) (model.Nurse, error)
	GetAllVisits(nurseId string) (dto.NurseVisitsListsDto, error)
	ConfirmOrCancelVisit(nurseId, visitId, reason string) (string, error)
	GetPatientProfile(patientId string) (dto.PatientProfileResponseDTO, error)
	NurseDashboardData(nurseId string) (dto.NurseDashboardDataResponseDTO, error)
	UpdateNurseFields(id string, updates map[string]interface{}) (dto.NurseUpdateResponseDTO, error)
	DeleteNurse(nurseId string, deleteAccountPasswordDto dto.DeleteAccountPasswordDto) error
	GetAvailabilityInfo(nurseId string) (dto.AvailabilityResponseDTO, error)
	GetNurseProfile(nurseId string) (userDTO.NurseProfileResponseDTO, error)
	GetNurseVisitInfo(nurseId, visitId string) (dto.NurseVisitInfo, error)
	VisitServiceConfirmation(nurseId, visitId, confirmationCode string) error
	TurnOfflineOnLogout(nurseId string) error
	RejectVisit(nurseId, visitId string) error
	AddReview(nurseId, visitId string, reviewDto dto.ReviewDTO) error
	GetMyNurseProfile(nurseId string) (userDTO.NurseProfileResponseDTO, error)
	CreateStripeOnboardingLink(nurseId string) (dto.StripeOnboardingResponseDTO, error)
	AddPrescription(nurseId string, visitId string, prescriptions []string) (dto.NurseVisitInfo, error)
}

type nurseService struct {
	userRepository   repository.UserRepository
	nurseRepository  repository.NurseRepository
	visitRepository  repository.VisitRepository
	reviewRepository repository.ReviewRepository
	stripeRepository repository.StripeRepository
}

func NewNurseService(userRepository repository.UserRepository, nurseRepository repository.NurseRepository, visitRepository repository.VisitRepository, reviewRepository repository.ReviewRepository, stripeRepository repository.StripeRepository) NurseService {
	return &nurseService{userRepository: userRepository, nurseRepository: nurseRepository, visitRepository: visitRepository, reviewRepository: reviewRepository, stripeRepository: stripeRepository}
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
		"online":     nurse.Online,
		"updated_at": time.Now(),
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
		// Retorna um DTO vazio em caso de erro
		return dto.NurseVisitsListsDto{}, err
	}

	pendingVisits := make([]dto.VisitDto, 0)
	confirmedVisits := make([]dto.VisitDto, 0)
	completedVisits := make([]dto.VisitDto, 0)
	rejectedVisits := make([]dto.VisitDto, 0)

	visitsToday := make([]dto.VisitDto, 0)

	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		location = time.UTC
	}
	now := time.Now().In(location)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	tomorrowStart := todayStart.Add(24 * time.Hour)

	for _, visit := range visits {
		patient, err := s.userRepository.FindUserById(visit.PatientId)
		if err != nil {
			return dto.NurseVisitsListsDto{}, err
		}

		visitReview, err := s.reviewRepository.FindReviewByVisitId(visit.ID.Hex())
		if err != nil {
			// Se o erro NÃO FOR "documento não encontrado", é um erro real.
			if err != mongo.ErrNoDocuments {
				// É um erro inesperado (ex: DB offline), então paramos.
				return dto.NurseVisitsListsDto{}, err
			}
			// Se o erro FOR "mongo.ErrNoDocuments", está tudo bem.
			// 'visitReview' continuará sendo um 'model.Review{}' vazio (zero-value).
			// 'visitReview.Rating' será 0 (ou 0.0), que é o que queremos.
		}

		visitDto := dto.VisitDto{
			ID:             visit.ID.Hex(),
			Description:    visit.Description,
			Reason:         visit.Reason,
			VisitType:      visit.VisitType,
			VisitValue:     visit.VisitValue,
			CreatedAt:      visit.CreatedAt.Format("02/01/2006 15:04"),
			Date:           visit.VisitDate.Format("02/01/2006 15:04"),
			Status:         visit.Status,
			Rating:         visitReview.Rating,
			PatientName:    visit.PatientName,
			PatientImageID: patient.ProfileImageID.Hex(),
			PatientId:      visit.PatientId,
			NurseName:      visit.NurseName,
		}

		switch visit.Status {
		case "PENDING":
			pendingVisits = append(pendingVisits, visitDto)
		case "CONFIRMED":
			confirmedVisits = append(confirmedVisits, visitDto)
		case "COMPLETED":
			completedVisits = append(completedVisits, visitDto)
		case "REJECTED":
			fmt.Print("oi", visitDto)
			rejectedVisits = append(rejectedVisits, visitDto)
		}

		visitDate := visit.VisitDate.In(location)
		isValidStatus := visit.Status == "CONFIRMED"
		isToday := (visitDate.Equal(todayStart) || visitDate.After(todayStart)) && visitDate.Before(tomorrowStart)

		if isValidStatus && isToday {
			visitsToday = append(visitsToday, visitDto)
		}
	}

	allVisitsDto := dto.NurseVisitsListsDto{
		Pending:     pendingVisits,
		Confirmed:   confirmedVisits,
		Completed:   completedVisits,
		Rejected:    rejectedVisits,
		VisitsToday: visitsToday,
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
		visit.Status = "REJECTED"
		response = "Visita que estava confirmada foi cancelada com sucesso."
		utils.SendEmailVisitCanceledWithReason("komatsuhenry@gmail.com", visit.NurseName, visit.VisitDate.Format("02/01/2006 15:04"), reason)
	} else if visit.Status == "PENDING" {
		visit.CancelReason = ""
		visit.Status = "CONFIRMED"
		response = "Visita que estava pendente foi confirmada com sucesso."

		utils.SendEmailVisitApproved("komatsuhenry@gmail.com", visit.NurseName, visit.VisitDate.Format("02/01/2006 15:04"), visit.VisitValue)
	} else if visit.Status == "REJECTED" {
		visit.CancelReason = ""
		visit.Status = "CONFIRMED"
		response = "Visita que estava rejeitada foi confirmada com sucesso."

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

	completedVisitsObjs, err := s.visitRepository.FindAllCompletedVisitsForPatient(patientId)
	if err != nil {
		return dto.PatientProfileResponseDTO{}, fmt.Errorf("Erro ao buscar visitas do paciente.")
	}

	// FAZER ISSO DEPOIS

	// rating, err := h.reviewRepository.FindAverageRatingByNurseId(nurseId)
	// if err != nil {
	// 	return userDTO.NurseProfileResponseDTO{}, err
	// }

	reviews, err := s.reviewRepository.FindAllPatientReviews(patientId)
	if err != nil {
		return dto.PatientProfileResponseDTO{}, fmt.Errorf("Erro ao buscar avaliações do paciente.")
	}

	var dtoReviews []dto.Reviews

	for _, review := range reviews {
		dtoReviews = append(dtoReviews, dto.Reviews{
			NurseName: review.NurseName,
			Rating:    review.Rating,
			Comment:   review.Comment,
		})
	}

	patientProfile := dto.PatientProfileResponseDTO{
		ID:             patient.ID,
		Name:           patient.Name,
		Email:          patient.Email,
		Phone:          patient.Phone,
		Address:        patient.Address,
		Cpf:            patient.Cpf,
		Rating:         5,
		VisitCount:     len(completedVisitsObjs),
		Reviews:        dtoReviews,
		Password:       patient.Password,
		Hidden:         patient.Hidden,
		Role:           patient.Role,
		ProfileImageID: patient.ProfileImageID.Hex(),
		CreatedAt:      patient.CreatedAt,
		TempCode:       patient.TempCode,
		UpdatedAt:      patient.UpdatedAt,
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
		PatientsAttended:  10,
		AppointmentsToday: 3,
		AverageRating:     4.5,
		MonthlyEarnings:   1000,
	}

	//itera pela lista de visits e mapeia para o dto VisitDto
	visitsDto := make([]dto.VisitDto, 0)
	for _, visit := range visits {
		visitsDto = append(visitsDto, dto.VisitDto{
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
		})
	}

	//nurseProfile
	nurseProfile := dto.NurseProfileResponseDTO{
		Name:            nurse.Name,
		Email:           nurse.Email,
		Phone:           nurse.Phone,
		Coren:           nurse.Coren,
		ExperienceYears: nurse.YearsExperience,
		Department:      nurse.Department,
		Bio:             nurse.Bio,
	}

	nurseAvailability := dto.AvailabilityDTO{
		IsAvailable:    nurse.Online,
		StartTime:      nurse.StartTime,
		EndTime:        nurse.EndTime,
		Specialization: nurse.Specialization,
	}

	dashboardData = dto.NurseDashboardDataResponseDTO{
		Online: nurse.Online,
		Stats:  stats,
		Visits: visitsDto,
		// Patients: patientsDto,
		// History: historyDto,
		Profile:      nurseProfile,
		Availability: nurseAvailability,
	}

	return dashboardData, nil

}

func (h *nurseService) GetNurseProfile(nurseId string) (userDTO.NurseProfileResponseDTO, error) {
	nurse, err := h.nurseRepository.FindNurseById(nurseId)
	if err != nil {
		return userDTO.NurseProfileResponseDTO{}, err
	}

	schedule, err := h.visitRepository.FindAllVisitsForNurse(nurseId)
	if err != nil {
		return userDTO.NurseProfileResponseDTO{}, err
	}

	var totalEarnings float64 = 0.0
	var totalPatients int = 0

	for _, visit := range schedule {
		if visit.Status == "COMPLETED" {
			totalPatients += 1
			totalEarnings += visit.VisitValue
		}
	}
	// =======================================================

	nurseRatingAvg, err := h.reviewRepository.FindAverageRatingByNurseId(nurseId)
	if err != nil {
		return userDTO.NurseProfileResponseDTO{}, fmt.Errorf("Erro ao buscar média de avaliação de enfermeiro(a).")
	}

	reviews, err := h.reviewRepository.FindAllNurseReviews(nurseId)
	if err != nil {
		return userDTO.NurseProfileResponseDTO{}, fmt.Errorf("Erro ao buscar avaliações do enfemeiro(a).")
	}

	var dtoReviews []userDTO.Reviews

	// 3. Iterar sobre as reviews do banco (model.Review)
	//    e converter para o DTO (userDTO.Reviews)
	for _, review := range reviews {
		dtoReviews = append(dtoReviews, userDTO.Reviews{
			PatientName: review.PatientName,
			Rating:      review.Rating,
			Comment:     review.Comment,
		})
	}

	nurseProfile := userDTO.NurseProfileResponseDTO{
		ID:             nurse.ID.Hex(),
		Name:           nurse.Name,
		Experience:     nurse.YearsExperience,
		Rating:         nurseRatingAvg,
		Shift:          nurse.Shift,
		Image:          nurse.ProfileImageID.Hex(),
		Location:       nurse.Address,
		Phone:          nurse.Phone,
		Online:         nurse.Online,
		Coren:          nurse.Coren,
		ProfileImageID: nurse.ProfileImageID.Hex(),
		Schedule:       schedule,
		TotalPatients:  totalPatients,
		Earnings:       totalEarnings,

		Department: nurse.Department,
		Bio: nurse.Bio,
		Qualifications: nurse.Qualifications,
		Specialization: nurse.Specialization,
		MaxPatientsPerDay: nurse.MaxPatientsPerDay,
		Price: nurse.Price,
		Services: nurse.Services,
		DaysAvailable: nurse.DaysAvailable,
		StartTime: nurse.StartTime,
		EndTime: nurse.EndTime,
		Neighborhoods: nurse.AvailableNeighborhoods,
		StripeAccountId: nurse.StripeAccountId,

		Reviews:        dtoReviews,
	}

	return nurseProfile, nil
}

func (s *nurseService) UpdateNurseFields(id string, updates map[string]interface{}) (dto.NurseUpdateResponseDTO, error) {
	if emailRaw, ok := updates["email"]; ok {
		email, ok := emailRaw.(string)
		if ok {
			normalizedEmail := strings.ToLower(email)

			_, err := utils.EmailRegex(email)
			if err != nil {
				return dto.NurseUpdateResponseDTO{}, fmt.Errorf("Email no formato incorreto.")
			}

			existingUser, err := s.nurseRepository.FindNurseByEmail(normalizedEmail)
			if err == nil && existingUser.ID.Hex() != id {
				return dto.NurseUpdateResponseDTO{}, fmt.Errorf("Email já está em uso por outro usuário")
			}

			updates["email"] = normalizedEmail
		}
	}

	nurse, err := s.nurseRepository.UpdateNurseFields(id, updates)
	if err != nil {
		return dto.NurseUpdateResponseDTO{}, fmt.Errorf("erro ao atualizar campos do usuario: %w", err)
	}

	return dto.NurseUpdateResponseDTO{
		ID:        nurse.ID.Hex(),
		Name:      nurse.Name,
		Email:     nurse.Email,
		Hidden:    nurse.Hidden,
		Role:      nurse.Role,
		CreatedAt: nurse.CreatedAt,
		UpdatedAt: nurse.UpdatedAt,
	}, nil
}

func (s *nurseService) DeleteNurse(nurseId string, deleteAccountPasswordDto dto.DeleteAccountPasswordDto) error {

	nurse, err := s.nurseRepository.FindNurseById(nurseId)
	if err != nil {
		return fmt.Errorf("Erro ao buscar id de paciente.")
	}

	if !utils.ComparePassword(nurse.Password, deleteAccountPasswordDto.Password) {
		return fmt.Errorf("Credenciais inválidas. Tente novamente.")
	}

	err = s.nurseRepository.DeleteNurse(nurseId)
	if err != nil {
		return fmt.Errorf("erro ao deletar enfermeiro: %w", err)
	}

	return nil
}

func (s *nurseService) GetAvailabilityInfo(nurseId string) (dto.AvailabilityResponseDTO, error) {
	nurse, err := s.nurseRepository.FindNurseById(nurseId)
	if err != nil {
		return dto.AvailabilityResponseDTO{}, fmt.Errorf("erro ao buscar enfermeiro: %w", err)
	}

	availabilityResponseDto := dto.AvailabilityResponseDTO{
		Department:             nurse.Department,
		Bio:                    nurse.Bio,
		Online:                 nurse.Online,
		StartTime:              nurse.StartTime,
		EndTime:                nurse.EndTime,
		Specialization:         nurse.Specialization,
		Price:                  nurse.Price,
		MaxPatientsPerDay:      nurse.MaxPatientsPerDay,
		DaysAvailable:          nurse.DaysAvailable,
		Services:               nurse.Services,
		AvailableNeighborhoods: nurse.AvailableNeighborhoods,
		Qualifications:         nurse.Qualifications,
	}

	return availabilityResponseDto, nil
}

func (s nurseService) GetNurseVisitInfo(nurseId, visitId string) (dto.NurseVisitInfo, error) {
	visit, err := s.visitRepository.FindVisitById(visitId)
	if err != nil {
		return dto.NurseVisitInfo{}, fmt.Errorf("Erro ao buscar id da visita.")
	}

	if visit.NurseId != nurseId {
		return dto.NurseVisitInfo{}, fmt.Errorf("Essa visita é pertencente à outro enfermeiro.")
	}

	// today := time.Now()
	// visitDate := visit.VisitDate

	// if today.Year() != visitDate.Year() || today.Month() != visitDate.Month() || today.Day() != visitDate.Day() {
	// 	return dto.NurseVisitInfo{}, fmt.Errorf("Esta visita não está agendada para hoje.")
	// }

	patient, err := s.userRepository.FindUserById(visit.PatientId)
	if err != nil {
		return dto.NurseVisitInfo{}, fmt.Errorf("Erro ao buscar id de enfermeiro(a).")
	}

	patientDto := dto.PatientInfoDto{
		ID:             patient.ID.Hex(),
		Name:           patient.Name,
		Email:          patient.Email,
		Phone:          patient.Phone,
		CEP:            patient.CEP,
		Street:         patient.Street,
		Number:         patient.Number,
		Complement:     patient.Complement,
		Neighborhood:   patient.Neighborhood,
		City:           patient.City,
		UF:             patient.UF,
		Latitude:       patient.Latitude,
		Longitude:      patient.Longitude,
		Cpf:            patient.Cpf,
		ProfileImageID: patient.ProfileImageID.Hex(),
	}

	visitDto := dto.VisitInfoDto{
		ID:            visit.ID.Hex(),
		Status:        visit.Status,
		PatientId:     visit.PatientId,
		PatientName:   visit.PatientName,
		Description:   visit.Description,
		Reason:        visit.Reason,
		CancelReason:  visit.CancelReason,
		VisitValue:    visit.VisitValue,
		VisitType:     visit.VisitType,
		Prescriptions: visit.Prescriptions,
		VisitDate:     visit.VisitDate.Format("02/01/2006 15:04"),
		CreatedAt:     visit.CreatedAt.Format("02/01/2006 15:04"),
		UpdatedAt:     visit.UpdatedAt.Format("02/01/2006 15:04"),
	}

	visitInfo := dto.NurseVisitInfo{
		Visit:   visitDto,
		Patient: patientDto,
	}

	return visitInfo, nil
}

func (s *nurseService) VisitServiceConfirmation(nurseId, visitId, confirmationCode string) error {

	visit, err := s.visitRepository.FindVisitById(visitId)
	if err != nil {
		return fmt.Errorf("Erro ao buscar id da visita.")
	}

	if visit.Status != "CONFIRMED" {
		return fmt.Errorf("O atendimento aindão não foi confirmado pelo enfermeiro(a).")
	}

	if visit.NurseId != nurseId {
		return fmt.Errorf("Essa visita é pertencente à outro enfermeiro.")
	}

	//validacao de so conseguir confirmar servico no dia da visita
	// today := time.Now()
	// visitDate := visit.VisitDate

	// if today.Year() != visitDate.Year() || today.Month() != visitDate.Month() || today.Day() != visitDate.Day() {
	// 	return fmt.Errorf("Esta visita não está agendada para hoje.")
	// }

	if visit.ConfirmationCode != confirmationCode {
		return fmt.Errorf("Código de confirmação inválido.")
	}

	nurse, err := s.nurseRepository.FindNurseById(nurseId)
	if err != nil {
		return fmt.Errorf("Erro ao localizar dados do enfermeiro: %w", err)
	}

	if nurse.StripeAccountId == "" {
		return fmt.Errorf("Este enfermeiro(a) não possui uma conta de pagamentos configurada.")
	}

	if visit.PaymentIntentID == "" {
		return fmt.Errorf("Pagamento original não encontrado para esta visita.")
	}

	// 3. Calcular valor do repasse (Ex: comissão da plataforma de 10%)
	commissionRate := 0.10 // 10%
	amountToTransfer := visit.VisitValue * (1.0 - commissionRate)
	amountInCents := int64(math.Round(amountToTransfer * 100)) // valor em cents

	transfer, err := s.stripeRepository.CreateTransfer(
		amountInCents,
		nurse.StripeAccountId, // Destino (conta do enfermeiro)
		visit.PaymentIntentID, // Origem (pagamento do paciente)
	)
	if err != nil {
		// Se o repasse falhar, NÃO confirme a visita.
		return fmt.Errorf("Erro ao processar repasse para o enfermeiro: %w", err)
	}

	visitUpdates := bson.M{
		"status":      "COMPLETED",
		"updated_at":  time.Now(),
		"transfer_id": transfer.ID,
	}

	_, err = s.visitRepository.UpdateVisitFields(visitId, visitUpdates)
	if err != nil {
		return fmt.Errorf("Erro ao atualizar status da visita para completar serviço.")
	}

	//logica de liberar dinheiro retido para enfermerio

	return nil
}

func (s *nurseService) TurnOfflineOnLogout(nurseId string) error {

	nurseUpdates := bson.M{
		"online":     false,
		"updated_at": time.Now(),
	}

	//salve user com status true/false
	_, _ = s.nurseRepository.UpdateNurseFields(nurseId, nurseUpdates)

	return nil
}

func (s *nurseService) RejectVisit(nurseId, visitId string) error {
	visit, err := s.visitRepository.FindVisitById(visitId)
	if err != nil {
		return fmt.Errorf("Erro ao buscar id da visita.")
	}

	if visit.Status != "PENDING" {
		return fmt.Errorf("A visita não está pendente para ser rejeitada.")
	}

	if visit.NurseId != nurseId {
		fmt.Errorf("Essa visita é pertencente à outro enfermeiro.")
	}

	visitUpdates := bson.M{
		"status":     "REJECTED",
		"updated_at": time.Now(),
	}

	_, err = s.visitRepository.UpdateVisitFields(visitId, visitUpdates)
	if err != nil {
		fmt.Errorf("Erro ao atualizar status da visita para rejeitada.")
	}

	return nil
}

func (s *nurseService) AddReview(nurseId, visitId string, reviewDto dto.ReviewDTO) error {

	visit, err := s.visitRepository.FindVisitById(visitId)
	if err != nil {
		return fmt.Errorf("Erro ao buscar id da visita.")
	}

	//valida se a visita ja foi avaliada.

	if visit.Status != "COMPLETED" {
		return fmt.Errorf("A visita ainda não foi completada. Portanto não é possível deixar uma avaliação.")
	}

	if visit.NurseId != nurseId {
		return fmt.Errorf("Essa visita é pertencente à outro enfermeiro.")
	}

	patientObjectID, err := primitive.ObjectIDFromHex(visit.PatientId)
	if err != nil {
		return fmt.Errorf("Erro ao converter patientId em objectID.")
	}

	nurseObjectID, err := primitive.ObjectIDFromHex(visit.NurseId)
	if err != nil {
		return fmt.Errorf("Erro ao converter nurseId em objectID.")
	}

	review := model.Review{
		ID:          primitive.NewObjectID(),
		VisitId:     visit.ID,
		NurseId:     nurseObjectID,
		NurseName:   visit.NurseName,
		PatientId:   patientObjectID,
		PatientName: visit.PatientName,
		Rating:      reviewDto.Rating,
		Comment:     reviewDto.Comment,
		ReviewType:  "NURSE",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.reviewRepository.CreateReview(review)
	if err != nil {
		return fmt.Errorf("Erro ao criar review: %w", err)
	}

	return nil

}

func (h *nurseService) GetMyNurseProfile(nurseId string) (userDTO.NurseProfileResponseDTO, error) {
	nurse, err := h.nurseRepository.FindNurseById(nurseId)
	if err != nil {
		return userDTO.NurseProfileResponseDTO{}, err
	}

	rating, err := h.reviewRepository.FindAverageRatingByNurseId(nurseId)
	if err != nil {
		return userDTO.NurseProfileResponseDTO{}, err
	}

	nurseProfile := userDTO.NurseProfileResponseDTO{
		ID:              nurse.ID.Hex(),
		Name:            nurse.Name,
		Specialization:  nurse.Specialization,
		Experience:      nurse.YearsExperience,
		Rating:          rating,
		Price:           nurse.Price,
		Shift:           nurse.Shift,
		Department:      nurse.Department,
		TwoFactor:       nurse.TwoFactor,
		Image:           nurse.ProfileImageID.Hex(),
		Location:        nurse.Address,
		Neighborhood:    nurse.Neighborhood,
		Phone:           nurse.Phone,
		Online:          nurse.Online,
		Coren:           nurse.Coren,
		Bio:             nurse.Bio,
		Qualifications:  nurse.Qualifications,
		Services:        nurse.Services,
		DaysAvailable:   nurse.DaysAvailable,
		StartTime:       nurse.StartTime,
		EndTime:         nurse.EndTime,
		StripeAccountId: nurse.StripeAccountId,
		ProfileImageID:  nurse.ProfileImageID.Hex(),
	}

	return nurseProfile, nil
}

func (s *nurseService) CreateStripeOnboardingLink(nurseId string) (dto.StripeOnboardingResponseDTO, error) {
	var response dto.StripeOnboardingResponseDTO

	nurse, err := s.nurseRepository.FindNurseById(nurseId)
	if err != nil {
		return response, fmt.Errorf("Enfermeiro não encontrado: %w", err)
	}

	stripeAccountId := nurse.StripeAccountId

	if stripeAccountId == "" {
		newAccountId, err := s.stripeRepository.CreateExpressAccount(nurse.Email)
		if err != nil {
			return response, fmt.Errorf("Erro ao criar conta Stripe: %w", err)
		}

		err = s.nurseRepository.UpdateStripeAccountId(nurseId, newAccountId)
		if err != nil {
			return response, fmt.Errorf("Erro ao salvar ID da conta Stripe no banco: %w", err)
		}

		stripeAccountId = newAccountId
	}

	// 3. Crie o link de onboarding (onetime link)
	// O 'stripeAccountId' aqui será o antigo (se já existia) ou o novo (se acabamos de criar)
	linkUrl, err := s.stripeRepository.CreateAccountLink(stripeAccountId)
	if err != nil {
		return response, fmt.Errorf("Erro ao criar link de onboarding Stripe: %w", err)
	}

	// 4. Mapeie para o DTO de resposta
	response.URL = linkUrl
	return response, nil
}

func (s *nurseService) AddPrescription(nurseId string, visitId string, prescriptions []string) (dto.NurseVisitInfo, error) {
	visit, err := s.visitRepository.FindVisitById(visitId)
	if err != nil {
		return dto.NurseVisitInfo{}, fmt.Errorf("Erro ao buscar id da visita.")
	}

	fmt.Println("====")
	fmt.Println(visit.Status)
	fmt.Println("====")
	if visit.Status != "CONFIRMED" && visit.Status != "COMPLETED" {
		return dto.NurseVisitInfo{}, fmt.Errorf("A visita precisar estar confirmada ou completada para adicionar uma prescrição.")
	}

	if visit.NurseId != nurseId {
		return dto.NurseVisitInfo{}, fmt.Errorf("Essa visita é pertencente à outro enfermeiro.")
	}

	patient, err := s.userRepository.FindUserById(visit.PatientId)
	if err != nil {
		return dto.NurseVisitInfo{}, fmt.Errorf("Erro ao buscar id de enfermeiro(a).")
	}

	visitUpdates := bson.M{
		"prescriptions": prescriptions,
		"updated_at":    time.Now(),
	}

	_, err = s.visitRepository.UpdateVisitFields(visitId, visitUpdates)
	if err != nil {
		return dto.NurseVisitInfo{}, err
	}

	patientDto := dto.PatientInfoDto{
		ID:             patient.ID.Hex(),
		Name:           patient.Name,
		Email:          patient.Email,
		Phone:          patient.Phone,
		CEP:            patient.CEP,
		Street:         patient.Street,
		Number:         patient.Number,
		Complement:     patient.Complement,
		Neighborhood:   patient.Neighborhood,
		City:           patient.City,
		UF:             patient.UF,
		Latitude:       patient.Latitude,
		Longitude:      patient.Longitude,
		Cpf:            patient.Cpf,
		ProfileImageID: patient.ProfileImageID.Hex(),
	}

	visitDto := dto.VisitInfoDto{
		ID:            visit.ID.Hex(),
		Status:        visit.Status,
		PatientId:     visit.PatientId,
		PatientName:   visit.PatientName,
		Description:   visit.Description,
		Reason:        visit.Reason,
		CancelReason:  visit.CancelReason,
		VisitValue:    visit.VisitValue,
		VisitType:     visit.VisitType,
		VisitDate:     visit.VisitDate.Format("02/01/2006 15:04"),
		CreatedAt:     visit.CreatedAt.Format("02/01/2006 15:04"),
		UpdatedAt:     visit.UpdatedAt.Format("02/01/2006 15:04"),
		Prescriptions: prescriptions,
	}

	visitInfo := dto.NurseVisitInfo{
		Visit:   visitDto,
		Patient: patientDto,
	}

	return visitInfo, nil

}
