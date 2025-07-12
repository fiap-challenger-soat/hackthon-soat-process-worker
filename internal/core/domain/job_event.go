package domain

type JobErrorEvent struct {
	JobID string `json:"job_id"`
	Email string `json:"email"`
}

type JobMessageEvent struct {
	JobID string `json:"job_id"`
}
