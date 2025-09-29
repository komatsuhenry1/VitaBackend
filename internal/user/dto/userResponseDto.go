package dto

type AllNursesListDto struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Specialization  string  `json:"specialization"`
	YearsExperience int     `json:"years_experience"`
	Price           float32 `json:"price"`
	Shift           string  `json:"shift"`
	Department      string  `json:"department"`
	Image           string  `json:"image"`
	Available       bool    `json:"available"`
	Location        string  `json:"location"`
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

type NurseProfileResponseDTO struct{
	ID             string              `json:"id"`
	Name           string           `json:"name"`
	Specialization string           `json:"specialization"`
	Experience     int              `json:"experience"`
	Rating         float64          `json:"rating"`
	Price          float64          `json:"price"`
	Shift          string           `json:"shift"`
	Department     string           `json:"department"`
	Image          string           `json:"image"`
	Available      bool             `json:"available"`
	Location       string           `json:"location"`
	Bio            string           `json:"bio"`
	Qualifications []string         `json:"qualifications"`
	Services       []string         `json:"services"`
	Reviews        []ReviewDTO      `json:"reviews"`
	Availability   []AvailabilityDTO `json:"availability"`

}
