/*
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package regources

type PaymentOpRelationships struct {
	AccountFrom *Relation `json:"account_from,omitempty"`
	AccountTo   *Relation `json:"account_to,omitempty"`
	Asset       *Relation `json:"asset,omitempty"`
	BalanceFrom *Relation `json:"balance_from,omitempty"`
	BalanceTo   *Relation `json:"balance_to,omitempty"`
}
