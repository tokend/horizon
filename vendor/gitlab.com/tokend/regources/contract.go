package regources

type Contract struct {
	ID            string   `json:"id"`
	Contractor    string   `json:"contractor"`
	Customer      string   `json:"customer"`
	Escrow        string   `json:"escrow"`
	Disputer      string   `json:"disputer,omitempty"`
	StartTime     string   `json:"start_time"`
	EndTime       string   `json:"end_time"`
	Details       []string `json:"details"`
	Invoices      []string `json:"invoices,omitempty"`
	DisputeReason string   `json:"dispute_reason,omitempty"`
	State         string   `json:"state"`
}
