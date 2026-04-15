package virustotal

// Stats - статистика перевірок
type Stats struct {
	Malicious  int `json:"malicious"`
	Suspicious int `json:"suspicious"`
	Harmless   int `json:"harmless"`
	Undetected int `json:"undetected"`
}

// VTResponse - загальна структура для читання відповідей VT
type VTResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Status            string `json:"status"` // queued, in_progress, completed
			Stats             Stats  `json:"stats"`
			LastAnalysisStats Stats  `json:"last_analysis_stats"`
		} `json:"attributes"`
	} `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
