# OfferAttributes

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BaseAmount** | **Amount** | defines the amount of offer in the base asset | 
**CreatedAt** | [**time.Time**](time.Time.md) | defines the time when the offer was created | 
**Fee** | [**Fee**](Fee.md) |  | 
**IsBuy** | **bool** | defines whether an offer created is on buying or selling the base_asset, or both | 
**OrderBookId** | **int64** | defines whether an offer created is on selling or trading. Could be either &#x60;0&#x60; (secondary market) or some &#x60;saleId&#x60; (for specific sale) or &#x60;-1&#x60; | 
**Price** | **Amount** | defines the price of an offer | 
**QuoteAmount** | **Amount** | defines the amount of offer in the quote asset | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


