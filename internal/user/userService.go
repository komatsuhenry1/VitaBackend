package user

import (
	"context"
	"encoding/json"
	"log"
	adminDTO "medassist/internal/admin/dto"
	"medassist/internal/auth/dto"
	"medassist/internal/model"
	"medassist/internal/repository"
	userDTO "medassist/internal/user/dto"
	"medassist/utils"
	"strconv"
	"time"

	"fmt"
	"medassist/internal/chat"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	GetAllNurses(patientId string) ([]userDTO.AllNursesListDto, error)
	GetFileByID(ctx context.Context, id primitive.ObjectID) (*dto.FileData, error)
	ContactUsMessage(contactUsDto userDTO.ContactUsDTO) error
	GetNurseProfile(nurseId string) (userDTO.NurseProfileResponseDTO, error)
	VisitSolicitation(userId string, createVisitDto userDTO.CreateVisitDto) error
	FindAllVisits(patientId string) (userDTO.VisitsResponseDto, error)
	UpdateUser(userId string, updates map[string]interface{}) (adminDTO.UserTypeResponse, error)
	DeleteUser(patientId string, deleteAccountPasswordDto userDTO.DeleteAccountPasswordDto) error
	ConfirmVisitService(visitId, patientId string) error
	GetOnlineNurses(userId string) ([]userDTO.AllNursesListDto, error)
	GetPatientVisitInfo(patientId, visitId string) (userDTO.PatientVisitInfo, error)
	AddReview(userId, visitId string, reviewDto userDTO.ReviewDTO) error
	ImmediateVisitSolicitation(patientId string, immediateVisitDto userDTO.ImmediateVisitDTO) (string, error)
	GetPatientProfile(patientId string) (userDTO.PatientProfileResponseDTO, error)
}

type userService struct {
	userRepository   repository.UserRepository
	nurseRepository  repository.NurseRepository
	visitRepository  repository.VisitRepository
	reviewRepository repository.ReviewRepository
	visitHub         *chat.Hub
}

func NewUserService(
	userRepository repository.UserRepository,
	nurseRepository repository.NurseRepository,
	visitRepository repository.VisitRepository,
	reviewRepository repository.ReviewRepository,
	visitHub *chat.Hub,
) UserService {
	return &userService{
		userRepository:   userRepository,
		nurseRepository:  nurseRepository,
		visitRepository:  visitRepository,
		reviewRepository: reviewRepository,
		visitHub:         visitHub,
	}
}

func (s *userService) GetAllNurses(patientId string) ([]userDTO.AllNursesListDto, error) {
	patient, err := s.userRepository.FindUserById(patientId)
	if err != nil {
		return []userDTO.AllNursesListDto{}, fmt.Errorf("Erro ao buscar id de paciente.")
	}
	nurses, err := s.nurseRepository.GetAllNurses(patient.City)
	if err != nil {
		return nil, err
	}

	return nurses, nil
}

func (s *userService) GetFileByID(ctx context.Context, id primitive.ObjectID) (*dto.FileData, error) {
	// Repassa os parâmetros corretamente para o repositório.
	return s.userRepository.FindFileByID(ctx, id)
}

func (h *userService) ContactUsMessage(contactUsDto userDTO.ContactUsDTO) error {
	err := utils.SendContactUsEmail(contactUsDto)
	if err != nil {
		return err
	}

	return nil
}

func (h *userService) GetNurseProfile(nurseId string) (userDTO.NurseProfileResponseDTO, error) {
	nurse, err := h.nurseRepository.FindNurseById(nurseId)
	if err != nil {
		return userDTO.NurseProfileResponseDTO{}, err
	}

	if nurse.MaxPatientsPerDay == 0 ||
		len(nurse.DaysAvailable) == 0 ||
		len(nurse.Services) == 0 ||
		len(nurse.AvailableNeighborhoods) == 0 {
		return userDTO.NurseProfileResponseDTO{}, fmt.Errorf("O enfermeiro ainda não preencheu os dados necessários para ser visto por pacientes.")
	}

	rating, err := h.reviewRepository.FindAverageRatingByNurseId(nurseId)
	if err != nil {
		return userDTO.NurseProfileResponseDTO{}, err
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
		Specialization: nurse.Specialization,
		Experience:     nurse.YearsExperience,
		Rating:         rating,
		Price:          nurse.Price,
		Shift:          nurse.Shift,
		Department:     nurse.Department,
		Image:          nurse.ProfileImageID.Hex(),
		Location:       nurse.Address,
		Neighborhood:   nurse.Neighborhood,
		Phone:          nurse.Phone,
		Online:         nurse.Online,
		Coren:          nurse.Coren,
		Bio:            nurse.Bio,
		Qualifications: nurse.Qualifications,
		Services:       nurse.Services,
		DaysAvailable:  nurse.DaysAvailable,
		StartTime:      nurse.StartTime,
		EndTime:        nurse.EndTime,
		ProfileImageID: nurse.ProfileImageID.Hex(),
		Reviews:        dtoReviews,
	}

	return nurseProfile, nil
}

func (h *userService) VisitSolicitation(patientId string, createVisitDto userDTO.CreateVisitDto) error {
	patient, err := h.userRepository.FindUserById(patientId)
	if err != nil {
		return err
	}

	nurse, err := h.nurseRepository.FindNurseById(createVisitDto.NurseId)
	if err != nil {
		return err
	}

	confirmationCode, err := utils.GenerateAuthCode()
	if err != nil {
		return fmt.Errorf("Erro ao gerar codigo de confirmação: %w", err)
	}

	visit := model.Visit{
		ID:               primitive.NewObjectID(),
		Status:           "PENDING",
		ConfirmationCode: strconv.Itoa(confirmationCode),

		PatientId:    patientId,
		PatientName:  patient.Name,
		PatientEmail: patient.Email,

		CEP:          createVisitDto.CEP,
		Street:       createVisitDto.Street,
		Number:       createVisitDto.Number,
		Complement:   createVisitDto.Complement,
		Neighborhood: createVisitDto.Neighborhood,

		Description: createVisitDto.Description,
		Reason:      createVisitDto.Reason,

		NurseId:   createVisitDto.NurseId,
		NurseName: nurse.Name,

		VisitType:        createVisitDto.VisitType,
		VisitDate:        createVisitDto.VisitDate,
		VisitValue:       createVisitDto.VisitValue,
		VisitRequestType: "SCHEDULED",

		PaymentIntentID: createVisitDto.PaymentIntentID,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = h.visitRepository.CreateVisit(visit)
	if err != nil {
		return err
	}

	//utils.SendEmailVisitSolicitation(nurse.Email, patient.Name, createVisitDto.VisitDate.String(), "100", patient.Address)
	utils.SendEmailVisitSolicitation(nurse.Email, patient.Name, createVisitDto.VisitDate.String(), visit.VisitValue, patient.Address)

	return nil
}

func (h *userService) FindAllVisits(patientId string) (userDTO.VisitsResponseDto, error) {

	responseDto := userDTO.VisitsResponseDto{
		AllVisits:   make([]userDTO.AllVisitsDto, 0),
		VisitsToday: make([]userDTO.AllVisitsDto, 0),
	}

	visits, err := h.visitRepository.FindAllVisitsForPatient(patientId)
	if err != nil {
		return userDTO.VisitsResponseDto{}, err
	}

	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		location = time.UTC
	}

	now := time.Now().In(location)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	tomorrowStart := todayStart.Add(24 * time.Hour)

	for _, visit := range visits {

		// MUDANÇA: Agora checamos o tipo do erro
		visitReview, err := h.reviewRepository.FindReviewByVisitId(visit.ID.Hex())
		if err != nil {
			// Se o erro NÃO FOR "documento não encontrado", é um erro real.
			if err != mongo.ErrNoDocuments {
				// É um erro inesperado (ex: DB offline), então paramos.
				return userDTO.VisitsResponseDto{}, err
			}
			// Se o erro FOR "mongo.ErrNoDocuments", está tudo bem.
			// 'visitReview' continuará sendo um 'model.Review{}' vazio (zero-value).
			// 'visitReview.Rating' será 0 (ou 0.0), que é o que queremos.
		}

		nurse, err := h.nurseRepository.FindNurseById(visit.NurseId)
		if err != nil {
			// Este erro (enfermeiro(a) não encontrado) deve ser fatal,
			// pois uma visita não pode existir sem um enfermeiro(a).
			return userDTO.VisitsResponseDto{}, err
		}

		visitDto := userDTO.AllVisitsDto{
			ID:          visit.ID.Hex(),
			Description: visit.Description,
			Reason:      visit.Reason,
			VisitType:   visit.VisitType,
			CreatedAt:   visit.CreatedAt.Format("02/01/2006 15:04"),
			Date:        visit.VisitDate.Format("02/01/2006 15:04"),
			Status:      visit.Status,
			// Se 'visitReview' estiver vazio (devido ao ErrNoDocuments),
			// 'visitReview.Rating' será o valor zero (ex: 0), que é o correto.
			Rating: visitReview.Rating,
			Nurse: userDTO.NurseDto{
				ID:             nurse.ID.Hex(),
				Name:           nurse.Name,
				Specialization: nurse.Specialization,
				Image:          nurse.ProfileImageID.Hex(),
			},
		}

		responseDto.AllVisits = append(responseDto.AllVisits, visitDto)

		visitDate := visit.VisitDate.In(location) // Garante que a data da visita está no mesmo fuso
		isConfirmed := visit.Status == "CONFIRMED"
		isToday := (visitDate.Equal(todayStart) || visitDate.After(todayStart)) && visitDate.Before(tomorrowStart)

		if isConfirmed && isToday {
			responseDto.VisitsToday = append(responseDto.VisitsToday, visitDto)
		}
	}

	return responseDto, nil
}

func (s *userService) UpdateUser(userId string, updates map[string]interface{}) (adminDTO.UserTypeResponse, error) {

	if emailRaw, ok := updates["email"]; ok {
		email, ok := emailRaw.(string)
		if ok {
			normalizedEmail := strings.ToLower(email)

			_, err := utils.EmailRegex(email)
			if err != nil {
				return adminDTO.UserTypeResponse{}, fmt.Errorf("Email no formato incorreto.")
			}

			existingUser, err := s.userRepository.FindUserByEmail(normalizedEmail)
			// se nao achar em user, busca em nurseRepositoru
			if err != nil {
				existingUser, err = s.nurseRepository.FindNurseByEmail(normalizedEmail)
			}

			if err == nil && existingUser.ID.Hex() != userId {
				return adminDTO.UserTypeResponse{}, fmt.Errorf("Email já está em uso por outro usuário")
			}

			updates["email"] = normalizedEmail
		}
	}

	if existingUser, err := s.userRepository.FindUserById(userId); err == nil && existingUser.Role == "PATIENT" {
		updated, err := s.userRepository.UpdateUser(userId, updates)
		if err != nil {
			return adminDTO.UserTypeResponse{}, fmt.Errorf("erro ao atualizar campos do usuario: %w", err)
		}
		return adminDTO.UserTypeResponse{
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

	return adminDTO.UserTypeResponse{}, fmt.Errorf("usuário não encontrado")
}

func (s *userService) DeleteUser(patientId string, deleteAccountPasswordDto userDTO.DeleteAccountPasswordDto) error {

	patient, err := s.userRepository.FindUserById(patientId)
	if err != nil {
		return fmt.Errorf("Erro ao buscar id de paciente.")
	}

	if !utils.ComparePassword(patient.Password, deleteAccountPasswordDto.Password) {
		return fmt.Errorf("Credenciais inválidas. Tente novamente.")
	}

	err = s.userRepository.DeleteUser(patientId)
	if err != nil {
		return fmt.Errorf("erro ao deletar visita: %w", err)
	}

	return nil
}

func (s *userService) ConfirmVisitService(visitId, patientId string) error {
	visit, err := s.visitRepository.FindVisitById(visitId)
	if err != nil {
		return fmt.Errorf("Erro ao buscar visita")
	}

	if visit.PatientId != patientId {
		return fmt.Errorf("Visita pertence à outro paciente.")
	}

	// logica de liberar o dinheiro retido para o enfermeiro

	if visit.Status != "CONFIRMED" {
		return fmt.Errorf("O status da visita deve estar com status confirmada para ser completada.")
	}

	visitUpdate := bson.M{
		"status":     "COMPLETED",
		"updated_at": time.Now(),
	}

	//salve user com status true/false
	visit, err = s.visitRepository.UpdateVisitFields(visitId, visitUpdate)
	if err != nil {
		return fmt.Errorf("Erro ao atualizar status de visita: %w", err)
	}

	return nil
}

func (s *userService) GetOnlineNurses(userId string) ([]userDTO.AllNursesListDto, error) {
	patient, err := s.userRepository.FindUserById(userId)
	if err != nil {
		return []userDTO.AllNursesListDto{}, nil
	}

	latitude := patient.Latitude
	longitude := patient.Longitude

	onlineNurses, err := s.nurseRepository.GetAllOnlineNurses(patient.City, latitude, longitude)
	if err != nil {
		return []userDTO.AllNursesListDto{}, nil
	}

	return onlineNurses, nil
}

func (s *userService) GetPatientVisitInfo(patientId, visitId string) (userDTO.PatientVisitInfo, error) {
	visit, err := s.visitRepository.FindVisitById(visitId)
	if err != nil {
		return userDTO.PatientVisitInfo{}, fmt.Errorf("Erro ao buscar id da visita.")
	}

	// if visit.Status == "PENDING" {
	// 	return userDTO.PatientVisitInfo{}, fmt.Errorf("O atendimento aindão não foi confirmado pelo enfermeiro(a).")
	// }

	if visit.PatientId != patientId {
		return userDTO.PatientVisitInfo{}, fmt.Errorf("Essa visita é pertencente à outro paciente.")
	}

	// validacao de ver os dados da visita apenas no dia de hoje
	// today := time.Now()
	// visitDate := visit.VisitDate

	// if today.Year() != visitDate.Year() || today.Month() != visitDate.Month() || today.Day() != visitDate.Day() {
	// 	return userDTO.PatientVisitInfo{}, fmt.Errorf("Esta visita não está agendada para hoje.")
	// }

	nurse, err := s.nurseRepository.FindNurseById(visit.NurseId)
	if err != nil {
		return userDTO.PatientVisitInfo{}, fmt.Errorf("Erro ao buscar id de enfermeiro(a).")
	}

	visitDto := userDTO.VisitInfoDto{
		ID:               visit.ID.Hex(),
		Status:           visit.Status,
		Description:      visit.Description,
		Reason:           visit.Reason,
		CancelReason:     visit.CancelReason,
		NurseId:          visit.NurseId,
		NurseName:        visit.NurseName,
		VisitValue:       visit.VisitValue,
		VisitType:        visit.VisitType,
		VisitDate:        visit.VisitDate.String(),
		CreatedAt:        visit.CreatedAt.String(),
		UpdatedAt:        visit.UpdatedAt.String(),
		ConfirmationCode: visit.ConfirmationCode,
	}

	nurseDto := userDTO.NurseInfoDto{
		ID:              nurse.ID.Hex(),
		Name:            nurse.Name,
		Email:           nurse.Email,
		Phone:           nurse.Phone,
		Specialization:  nurse.Specialization,
		YearsExperience: nurse.YearsExperience,
		Rating:          nurse.Rating,
		Coren:           nurse.Coren,
		ProfileImageID:  nurse.ProfileImageID.Hex(),
	}

	patientVisitInfo := userDTO.PatientVisitInfo{
		Visit: visitDto,
		Nurse: nurseDto,
	}

	return patientVisitInfo, nil
}

func (s *userService) AddReview(userId, visitId string, reviewDto userDTO.ReviewDTO) error {

	visit, err := s.visitRepository.FindVisitById(visitId)
	if err != nil {
		return fmt.Errorf("Erro ao buscar id da visita.")
	}

	//valida se a visita ja foi avaliada.

	if visit.Status != "COMPLETED" {
		return fmt.Errorf("A visita ainda não foi completada. Portanto não é possível deixar uma avaliação.")
	}

	if visit.PatientId != userId {
		return fmt.Errorf("Essa visita é pertencente à outro paciente.")
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
		ReviewType:  "PATIENT",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.reviewRepository.CreateReview(review)
	if err != nil {
		return fmt.Errorf("Erro ao criar review: %w", err)
	}

	return nil

}

func (s *userService) ImmediateVisitSolicitation(patientId string, immediateVisitDto userDTO.ImmediateVisitDTO) (string, error) {
	patient, err := s.userRepository.FindUserById(patientId)
	if err != nil {
		return "", fmt.Errorf("Erro ao buscar id de paciente.")
	}

	nurse, err := s.nurseRepository.FindNurseById(immediateVisitDto.NurseId)
	if err != nil {
		return "", fmt.Errorf("Erro ao buscar id de enfermeiro.")
	}
	if !nurse.Online {
		return "", fmt.Errorf("O(A) enfermeiro(a) %s não está online no momento e não pode receber solicitações imediatas", nurse.Name)
	}

	codeInt, _ := utils.GenerateAuthCode()
	code := strconv.Itoa(codeInt)

	visit := model.Visit{
		ID:               primitive.NewObjectID(),
		Status:           "PENDING",
		ConfirmationCode: code,

		PatientId:    patient.ID.Hex(),
		PatientName:  patient.Name,
		PatientEmail: patient.Email,

		CEP:          immediateVisitDto.CEP,
		Street:       immediateVisitDto.Street,
		Number:       immediateVisitDto.Number,
		Complement:   immediateVisitDto.Complement,
		Neighborhood: immediateVisitDto.Neighborhood,

		Description: immediateVisitDto.Description,
		Reason:      immediateVisitDto.Reason,

		NurseId:   nurse.ID.Hex(),
		NurseName: nurse.Name,

		PaymentIntentID: immediateVisitDto.PaymentIntentID,

		VisitValue: nurse.Price,

		VisitRequestType: "IMMEDIATE",
		VisitType:        immediateVisitDto.VisitType,
		VisitDate:        time.Now(),

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.visitRepository.CreateVisit(visit)
	if err != nil {
		return "", fmt.Errorf("Erro ao criar visita: %w", err)

	}

	// ===================================================================
	// 6. LÓGICA DE NOTIFICAÇÃO VIA WEBSOCKET (NOVA)
	// ===================================================================

	// 6a. Montar o payload (o que o app do enfermeiro vai receber)
	// (Você pode criar um DTO para isso)
	type NotificationPayload struct {
		Type        string  `json:"type"`     // Para o app saber que é uma nova chamada
		VisitID     string  `json:"visit_id"` // Para o enfermeiro aceitar/recusar
		PatientName string  `json:"patient_name"`
		PatientID   string  `json:"patient_id"`
		Reason      string  `json:"reason"`
		Value       float64 `json:"value"`
		Address     string  `json:"address"` // Um endereço formatado
	}

	address := fmt.Sprintf("%s, %s - %s", visit.Street, visit.Number, visit.Neighborhood)

	payload := NotificationPayload{
		Type:        "IMMEDIATE_VISIT_REQUEST",
		VisitID:     visit.ID.Hex(),
		PatientName: patient.Name,
		Reason:      visit.Reason,
		PatientID:   patient.ID.Hex(),
		Value:       visit.VisitValue,
		Address:     address,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Erro ao serializar payload do WebSocket para visita %s: %v", visit.ID.Hex(), err)
		// Não retorne erro aqui, o e-mail de fallback ainda pode funcionar
	} else {
		// 6b. Enviar a mensagem para o Hub
		sent := s.visitHub.SendToNurse(nurse.ID.Hex(), jsonPayload)

		if !sent {
			// Isso significa que o enfermeiro ficou offline no exato segundo
			// entre a verificação do 'if !nurse.Online' e agora.
			log.Printf("Alerta: Visita %s criada, mas enfermeiro %s ficou offline antes da notificação WS.", visit.ID.Hex(), nurse.ID.Hex())
		}
	}

	// 7. ENVIAR O E-MAIL (MANTIDO COMO FALLBACK)
	// Você pode manter isso, é uma boa garantia caso a notificação WS falhe.
	utils.SendEmailVisitSolicitation(nurse.Email, patient.Name, immediateVisitDto.VisitDate.String(), visit.VisitValue, patient.Address) // (patient.Address parece ser um campo antigo, talvez usar o 'address' formatado?)

	return patientId, nil
}

func (s *userService) GetPatientProfile(patientId string) (userDTO.PatientProfileResponseDTO, error) {

	patient, err := s.userRepository.FindUserById(patientId)
	if err != nil {
		return userDTO.PatientProfileResponseDTO{}, err
	}

	completedVisitsObjs, err := s.visitRepository.FindAllCompletedVisitsForPatient(patientId)
	if err != nil {
		return userDTO.PatientProfileResponseDTO{}, fmt.Errorf("Erro ao buscar visitas do paciente.")
	}

	// FAZER ISSO DEPOIS

	// rating, err := h.reviewRepository.FindAverageRatingByNurseId(nurseId)
	// if err != nil {
	// 	return userDTO.NurseProfileResponseDTO{}, err
	// }

	reviews, err := s.reviewRepository.FindAllPatientReviews(patientId)
	if err != nil {
		return userDTO.PatientProfileResponseDTO{}, fmt.Errorf("Erro ao buscar avaliações do paciente.")
	}

	var dtoReviews []userDTO.Reviews

	for _, review := range reviews {
		dtoReviews = append(dtoReviews, userDTO.Reviews{
			PatientName: review.NurseName,
			Rating:      review.Rating,
			Comment:     review.Comment,
		})
	}

	patientProfile := userDTO.PatientProfileResponseDTO{
		ID:             patient.ID,
		Name:           patient.Name,
		Email:          patient.Email,
		Phone:          patient.Phone,
		Address:        patient.Address,
		Cpf:            patient.Cpf,
		Rating:         5,
		TwoFactor:      patient.TwoFactor,
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
