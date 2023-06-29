package model

type Exam struct {
	PrepodName     string `json:"prepodName"`
	ExamDate       string `json:"examDate"`
	DisciplNum     string `json:"disciplNum"`
	PrepodNameEnc  string `json:"prepodNameEnc"`
	PrepodLogin    string `json:"prepodLogin"`
	AudNum         string `json:"audNum"`
	BuildNum       string `json:"buildNum"`
	DisciplNameEnc string `json:"disciplNameEnc"`
	DisciplName    string `json:"disciplName"`
	ExamTime       string `json:"examTime"`
}
