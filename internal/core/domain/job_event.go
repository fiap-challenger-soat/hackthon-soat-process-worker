package domain

type JobErrorEvent struct {
	JobID string `json:"job_id"`
}

type JobMessageEvent struct {
	JobID     string `json:"job_id"`
}
