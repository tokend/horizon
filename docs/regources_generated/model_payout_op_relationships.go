/*
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package regources

type PayoutOpRelationships struct {
	Asset         *Relation `json:"asset,omitempty"`
	SourceAccount *Relation `json:"source_account,omitempty"`
	SourceBalance *Relation `json:"source_balance,omitempty"`
}
