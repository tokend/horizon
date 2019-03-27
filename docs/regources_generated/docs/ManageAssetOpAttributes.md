# ManageAssetOpAttributes

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Action** | [**xdr.ManageAssetAction**](Enum.md) | * 0: \&quot;create_asset_creation_request\&quot; * 1: \&quot;create_asset_update_request\&quot; * 2: \&quot;cancel_asset_request\&quot; * 3: \&quot;change_preissued_asset_signer\&quot; * 4: \&quot;update_max_issuance\&quot;  | 
**AssetCode** | **string** | Asset to manage | 
**CreatorDetails** | [**Details**](Details.md) |  | 
**MaxIssuanceAmount** | **Amount** |  | 
**Policies** | Pointer to **xdr.AssetPolicy** | Bit mask. * 1:  \&quot;transferable\&quot; * 2:  \&quot;base_asset\&quot; * 4:  \&quot;stats_quote_asset\&quot; * 8:  \&quot;withdrawable\&quot; * 16: \&quot;issuance_manual_review_required\&quot; * 32: \&quot;can_be_base_in_atomic_swap\&quot; * 64: \&quot;can_be_quote_in_atomic_swap\&quot;  | [optional] 
**PreIssuanceSigner** | **string** | Address of preissuance signer | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


