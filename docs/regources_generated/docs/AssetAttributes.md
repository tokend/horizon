# AssetAttributes

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AvailableForIssuance** | **Amount** | Asset volume authorized to be issued by an asset owner | 
**Details** | [**Details**](Details.md) |  | 
**Issued** | **Amount** | Asset volume that is currently in circulation | 
**MaxIssuanceAmount** | **Amount** | Max volume of an asset that can be in circulation | 
**PendingIssuance** | **Amount** | Asset volume to be distributed via [asset sale↪](https://tokend.gitbook.io/knowledge-base/platform-features/crowdfunding) but currently locked by the system | 
**Policies** | [**xdr.AssetPolicy**](Mask.md) |  | 
**PreIssuanceAssetSigner** | **string** | address of the signer responsible for pre-issuance. [Details↪](https://tokend.gitbook.io/knowledge-base/technical-details/key-entities/asset#pre-issued-asset-signer) | 
**TrailingDigits** | **int32** | Number of significant digits after the point | 
**Type** | **uint64** | Numeric type of asset | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


