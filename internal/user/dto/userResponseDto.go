package dto

import "medassist/internal/model"

type AllNursesListDto struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	Specialization         string   `json:"specialization"`
	YearsExperience        int      `json:"years_experience"`
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
	Patient string  `json:"patient"`
	Rating  float64 `json:"rating"`
	Comment string  `json:"comment"`
	Date    string  `json:"date"`
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
