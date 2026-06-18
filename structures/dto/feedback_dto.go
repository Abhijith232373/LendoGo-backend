package dto

type SubmitFeedbackReq struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type UpdateFeedbackStatusReq struct {
	Status string `json:"status"` 
}