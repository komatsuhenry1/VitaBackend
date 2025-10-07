package dto

type DocumentInfoResponse struct {
	Name        string `json:"name"`         // Um nome amigável, ex: "Documento de Licença (COREN)"
	Type        string `json:"type"`         // Um identificador, ex: "license_document"
	DownloadURL string `json:"download_url"` // A URL para baixar o arquivo
	ImageID     string `json:"image_id"`
}

type DashboardAdminDataResponse struct {
	TotalNurses         int                                `json:"total_nurses"`
	TotalPatients       int                                `json:"total_patients"`
	VisitsToday         int                                `json:"visits_today"`
	NumberVisits        int                                `json:"number_visits"`
	PendentApprovations int                                `json:"pendent_approvations"`
	NursesFields        []NursesFieldsForDashboardResponse `json:"nurses_ids_pendent_approvations"`
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
	LicenseNumber    string  `json:"license_number"` // registro profissional
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

type UserListsResponse struct {
	Users  []UserTypeResponse  `json:"users"`
	Nurses []NurseTypeResponse `json:"nurses"`
}
