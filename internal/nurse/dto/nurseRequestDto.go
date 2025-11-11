package dto

type CancelReason struct {
	Reason string `json:"reason"`
}

type ReviewDTO struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type DeleteAccountPasswordDto struct {
	Password string `json:"password" binding:"required"`
}

type PrescriptionList struct{
	PrescriptionList []string `json:"prescription_list" binding:"required"`
}