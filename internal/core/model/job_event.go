package model

type JobErrorEvent struct {
	JobID    string `json:"job_id"`
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	FailedAt string `json:"failed_at"`
}

// JobMessage define a estrutura de dados esperada no corpo da mensagem SQS.
// Esta struct poderia vir do repositório de 'contratos' para ser compartilhada.
type JobMessageEvent struct {
	JobID     string `json:"job_id"`
	VideoPath string `json:"video_path"`
}
