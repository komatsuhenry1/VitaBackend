package dto

type NurseVisitsListsDto struct {
	Pending   []VisitDto `json:"pending"`
	Confirmed []VisitDto `json:"confirmed"`
	Completed []VisitDto `json:"completed"`
}

type VisitDto struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Reason      string `json:"reason"`
	VisitType   string `json:"visit_type"`
	CreatedAt   string `json:"created_at"`
	Date        string `json:"date"`
	Status      string `json:"status"`
	PatientName string `json:"patient_name"`
	NurseName   string `json:"nurse_name"`
}
