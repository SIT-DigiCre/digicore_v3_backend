package api

type ReqPostEvent struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	Tags        []string `json:"tags"`
	FinishDate  string   `json:"finish_date"`
}

type ResPostEvent struct {
	EventID string `json:"event_id"`
}
