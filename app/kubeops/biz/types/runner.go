package types

type RunAdhocResult struct {
	Result *Result
}

type RunPlaybookResult struct {
	Result *Result
}

type Result struct {
	Id        string `json:"id,omitempty"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	Message   string `json:"message,omitempty"`
	Success   bool   `json:"success,omitempty"`
	Finished  bool   `json:"finished,omitempty"`
	Content   string `json:"content,omitempty"`
	Project   string `json:"project,omitempty"`
}

type ResultItems []*Result
