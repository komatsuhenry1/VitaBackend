package dto

import "time"

type DocumentInfoResponse struct {
	Name        string `json:"name"`         // Um nome amigável, ex: "Documento de Licença (COREN)"
	Type        string `json:"type"`         // Um identificador, ex: "license_document"
	DownloadURL string `json:"download_url"` // A URL para baixar o arquivo
	ImageID     string `json:"image_id"`
}

type DashboardAdminDataResponse struct {
	// Métricas Principais (KPIs)
	TotalNurses         int     `json:"total_nurses"`
	TotalPatients       int     `json:"total_patients"`
	NumberVisits        int64   `json:"number_visits"` // Mudei para int64 para consistência
	VisitsToday         int64   `json:"visits_today"`  // Mudei para int64 para consistência
	AverageNurseRating  float64 `json:"average_nurse_rating"`
	TotalRevenueLast30Days float64 `json:"total_revenue_last_30_days"`

	// Métricas de Atividade e Crescimento
	NursesOnline              int64 `json:"nurses_online"`
	NewNursesLast30Days     int64 `json:"new_nurses_last_30_days"`
	NewPatientsLast30Days   int64 `json:"new_patients_last_30_days"`
	CompletedVisitsLast30Days int64 `json:"completed_visits_last_30_days"`

	// Métricas de Gestão e Alertas
	PendentApprovations int64                                `json:"pendent_approvations"`
	PendingNurses       []NursesFieldsForDashboardResponse `json:"nurses_ids_pendent_approvations"` // Renomeei aqui
	NursesInactive      int64                                `json:"nurses_inactive"`
	PatientsInactive    int64                                `json:"patients_inactive"`
	
	// Métricas Estratégicas
	MostCommonSpecialization string `json:"most_common_specialization"`
}

type NursesFieldsForDashboardResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserTypeResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	// PixKey      string `json:"pix_key"`
	Cpf         string `json:"cpf"`
	Password    string `json:"password"`
	Hidden      bool   `json:"hidden"`
	Role        string `json:"role"`
	FirstAccess bool   `json:"first_access"`
}

type NurseTypeResponse struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	Phone            string  `json:"phone"`
	Address          string  `json:"address"`
	Cpf              string  `json:"cpf"`
	PixKey           string  `json:"pix_key"`
	Password         string  `json:"password"`
	VerificationSeal bool    `json:"verification_seal"`
	Coren            string  `json:"coren"`          // registro profissional
	Specialization   string  `json:"specialization"` // área (ex: pediatrics, geriatrics, ER)
	Shift            string  `json:"shift"`          // manhã, tarde, noite
	Department       string  `json:"department"`     // setor/hospital onde trabalha
	YearsExperience  int     `json:"years_experience"`
	Price            float64 `json:"price"`
	Bio              string  `json:"bio"`
	Hidden           bool    `json:"hidden"`
	Role             string  `json:"role"`
	Online           bool    `json:"online"`
	FirstAccess      bool    `json:"first_access"`
	StartTime        string  `json:"start_time"`
	EndTime          string  `json:"end_time"`
}

type VisitTypeResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`

	PatientId    string `json:"patient_id"`
	PatientName  string `json:"patient_name"`
	PatientEmail string `json:"patient_email"`

	Description  string `json:"description"`
	Reason       string `json:"reason"`
	CancelReason string `json:"cancel_reason"`

	NurseId   string `json:"nurse_id"`
	NurseName string `json:"nurse_name"`

	VisitValue float64   `json:"value"`
	VisitType  string    `json:"visit_type"`
	VisitDate  time.Time `json:"visit_date"`
}

type UserListsResponse struct {
	Users  []UserTypeResponse  `json:"users"`
	Nurses []NurseTypeResponse `json:"nurses"`
	Visits []VisitTypeResponse `json:"visits"`
}
