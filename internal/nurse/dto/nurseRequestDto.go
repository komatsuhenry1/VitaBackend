package dto

type CancelReason struct {
	Reason string `json:"reason"`
}

type ReviewDTO struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}
