package dto

import "medassist/internal/model"

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type AllNursesListDto struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	Specialization         string   `json:"specialization"`
	YearsExperience        int      `json:"years_experience"`
	PatientLocation        Location `json:"patient_location"`
	Price                  float32  `json:"price"`
	Shift                  string   `json:"shift"`
	Department             string   `json:"department"`
	Image                  string   `json:"image"`
	Available              bool     `json:"available"`
	Location               string   `json:"location"`
	City                   string   `json:"city"`
	UF                     string   `json:"uf"`
	Neighborhood           string   `json:"neighborhood"`
	Street                 string   `json:"street"`
	Latitude               float64  `json:"latitude"`
	Longitude              float64  `json:"longitude"`
	MaxPatientsPerDay      int      `json:"max_patients_per_day"`
	DaysAvailable          []string `json:"days_available"`
	Services               []string `json:"services"`
	AvailableNeighborhoods []string `json:"available_neighborhoods"`
}

type ReviewDTO struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type AvailabilityDTO struct {
	Day   string `json:"day"`
	Hours string `json:"hours"`
}

type NurseProfileResponseDTO struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Specialization string            `json:"specialization"`
	Experience     int               `json:"experience"`
	Rating         float64           `json:"rating"`
	Online         bool              `json:"online"`
	Price          float64           `json:"price"`
	Shift          string            `json:"shift"`
	Department     string            `json:"department"`
	Coren          string            `json:"coren"`
	Phone          string            `json:"phone"`
	Image          string            `json:"image"`
	Location       string            `json:"location"`
	Bio            string            `json:"bio"`
	Schedule       []model.Visit     `json:"schedules"`
	TotalPatients  int               `json:"total_patients"`
	Earnings       float64           `json:"earnings"`
	Qualifications []string          `json:"qualifications"`
	ProfileImageID string            `json:"profile_image_id"`
	Services       []string          `json:"services"`
	Reviews        []ReviewDTO       `json:"reviews"`
	Availability   []AvailabilityDTO `json:"availability"`
}

type NurseDto struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
	Image          string `json:"image"`
}

type AllVisitsDto struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Reason      string   `json:"reason"`
	VisitType   string   `json:"visit_type"`
	CreatedAt   string   `json:"created_at"`
	Date        string   `json:"date"`
	Status      string   `json:"status"`
	Nurse       NurseDto `json:"nurse"`
}

type VisitsResponseDto struct {
	AllVisits   []AllVisitsDto `json:"all_visits"`
	VisitsToday []AllVisitsDto `json:"visits_today"`
}

type PatientVisitInfo struct {
	Visit VisitInfoDto `json:"visit"`
	Nurse NurseInfoDto `json:"nurse"`
}

type VisitInfoDto struct {
	ID               string  `json:"id"`
	Status           string  `json:"status"`
	Description      string  `json:"description"`
	Reason           string  `json:"reason"`
	CancelReason     string  `json:"cancel_reason"`
	NurseId          string  `json:"nurse_id"`
	NurseName        string  `json:"nurse_name"`
	VisitValue       float64 `json:"visit_value"`
	VisitType        string  `json:"visit_type"`
	VisitDate        string  `json:"visit_date"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	ConfirmationCode string  `json:"confirmation_code"`
}

type NurseInfoDto struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Phone           string  `json:"phone"`
	Specialization  string  `json:"specialization"`
	YearsExperience int     `json:"years_experience"`
	Rating          float64 `json:"rating"`
	Coren           string  `json:"coren"`
	ProfileImageID  string  `json:"profile_image_id"`
}
