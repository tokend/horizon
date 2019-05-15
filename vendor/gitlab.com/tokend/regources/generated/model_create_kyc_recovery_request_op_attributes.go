/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type CreateKycRecoveryRequestOpAttributes struct {
	// tasks set on request creation
	AllTasks       *int64             `json:"all_tasks,omitempty"`
	CreatorDetails Details            `json:"creator_details"`
	Identity       *uint32            `json:"identity,omitempty"`
	SignersData    []UpdateSignerData `json:"signers_data"`
	// Weight of the signer fo the account
	Weight *uint32 `json:"weight,omitempty"`
}
