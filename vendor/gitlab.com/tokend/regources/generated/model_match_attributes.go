/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type MatchAttributes struct {
	// defines quote amount of the match
	BaseAmount Amount `json:"base_amount"`
	// defines whether match is a secondary market match or match after sale end. `0` - secondary market, `saleId` - (for specific sale) or `-1`
	OrderBookId string `json:"order_book_id"`
	// defines price of the match
	Price Amount `json:"price"`
	// defines base amount of the match
	QuoteAmount Amount `json:"quote_amount"`
}
