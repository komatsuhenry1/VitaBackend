package user

import (
	"context"
	"medassist/internal/auth/dto"
	"medassist/internal/model"
	"medassist/internal/repository"
	userDTO "medassist/internal/user/dto"
	"medassist/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	GetAllNurses() ([]userDTO.AllNursesListDto, error)
	GetFileByID(ctx context.Context, id primitive.ObjectID) (*dto.FileData, error)
	ContactUsMessage(contactUsDto userDTO.ContactUsDTO) error
	GetNurseProfile(nurseId string) (userDTO.NurseProfileResponseDTO, error)
	VisitSolicitation(userId string, createVisitDto userDTO.CreateVisitDto) error
	FindAllVisits(patientId string) ([]userDTO.AllVisitsDto, error)
}

type userService struct {
	userRepository  repository.UserRepository
	nurseRepository repository.NurseRepository
	visitRepository repository.VisitRepository
}

func NewUserService(userRepository repository.UserRepository, nurseRepository repository.NurseRepository, visitRepository repository.VisitRepository) UserService {
	return &userService{userRepository: userRepository, nurseRepository: nurseRepository, visitRepository: visitRepository}
}

func (s *userService) GetAllNurses() ([]userDTO.AllNursesListDto, error) {
	nurses, err := s.nurseRepository.GetAllNurses()
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

	qualifications := []string{"Pediatria", "Geriatria", "UTI"}
	services := []string{"Servico 1 ", "Servico 2", "Servico 3"}

	reviews := []userDTO.ReviewDTO{{ // funcao na repo que retorna uma lista de reviews
		Patient: "paciente name",
		Rating:  4.5,
		Comment: "Review comment",
		Date:    "Review date",
	}}

	availability := []userDTO.AvailabilityDTO{{ // funcao na repository que retorna lista de avalability
		Day:   "19/09/2010",
		Hours: "10:00",
	}}

	nurseProfile := userDTO.NurseProfileResponseDTO{
		ID:             nurse.ID.Hex(),
		Name:           nurse.Name,
		Specialization: nurse.Specialization,
		Experience:     nurse.YearsExperience,
		Rating:         nurse.Rating,
		Price:          nurse.Price,
		Shift:          nurse.Shift,
		Department:     nurse.Department,
		Image:          nurse.ProfileImageID.Hex(),
		Available:      nurse.Online,
		Location:       nurse.Address,
		Bio:            "NURSE BIO",
		Qualifications: qualifications,
		Services:       services,
		Reviews:        reviews,
		Availability:   availability,
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

	visit := model.Visit{
		ID:     primitive.NewObjectID(),
		Status: "PENDING",

		PatientId:   patientId,
		PatientName: patient.Name,
		PatientEmail: patient.Email,

		Description: createVisitDto.Description,
		Reason:      createVisitDto.Reason,

		NurseId:   createVisitDto.NurseId,
		NurseName: nurse.Name,

		VisitType: createVisitDto.VisitType,
		VisitDate: createVisitDto.VisitDate,
		VisitValue: createVisitDto.VisitValue,


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

func (h *userService) FindAllVisits(patientId string) ([]userDTO.AllVisitsDto, error) {
	visits, err := h.visitRepository.FindAllVisitsForPatient(patientId)
	if err != nil {
		return nil, err
	}

	var allVisitsDto []userDTO.AllVisitsDto

	for _, visit := range visits {
		if visit.Status =="CONFIRMED"{
			nurse, err := h.nurseRepository.FindNurseById(visit.NurseId)
			if err != nil {
				return nil, err
			}
	
			visit, err := h.visitRepository.FindVisitById(visit.ID.Hex())
			if err != nil {
				return nil, err
			}
	
			allVisitsDto = append(allVisitsDto, userDTO.AllVisitsDto{
				ID:          visit.ID.Hex(),
				Description: visit.Description,
				Reason:      visit.Reason,
				VisitType:   visit.VisitType,
				CreatedAt:   visit.CreatedAt.Format("02/01/2006 15:04"),
				Date:        visit.VisitDate.Format("02/01/2006 15:04"),
				Status:      visit.Status,
				Nurse: userDTO.NurseDto{
					ID:             nurse.ID.Hex(),
					Name:           nurse.Name,
					Specialization: nurse.Specialization,
					Image:          nurse.ProfileImageID.Hex(),
				},
			})
		}
	}

	return allVisitsDto, nil
}
