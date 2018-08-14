package reviewablerequest2

// Represents Reviewable request
type ReviewableRequest struct {
	ID           string   `json:"id"`
	PT           string   `json:"paging_token"`
	Requestor    string   `json:"requestor"`
	Reviewer     string   `json:"reviewer"`
	Reference    *string  `json:"reference"`
	RejectReason string   `json:"reject_reason"`
	Hash         string   `json:"hash"`
	Details      *Details `json:"details"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	RequestState
}

func (r *ReviewableRequest) PagingToken() string {
	return r.PT
}
