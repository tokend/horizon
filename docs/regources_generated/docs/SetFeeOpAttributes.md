# SetFeeOpAttributes

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccountAddress** | **string** |  | [optional] 
**AccountRole** | **uint64** |  | [optional] 
**AssetCode** | **string** | Unique identifier of the asset | 
**FeeType** | [**xdr.FeeType**](Enum.md) | * 0: \&quot;payment_fee\&quot; * 1: \&quot;offer_fee\&quot; * 2: \&quot;withdrawal_fee\&quot; * 3: \&quot;issuance_fee\&quot; * 4: \&quot;invest_fee\&quot; * 5: \&quot;capital_deployment_fee\&quot; * 6: \&quot;operation_fee\&quot; * 7: \&quot;payout_fee\&quot; * 8: \&quot;atomic_swap_sale_fee\&quot; * 9: \&quot;atomic_swap_purchase_fee\&quot;  | 
**FixedFee** | **Amount** | Fixed amount to pay | 
**IsDelete** | **bool** |  | 
**LowerBound** | **Amount** | Lower bound of fee applicability | 
**PercentFee** | **Amount** | Percent to pay | 
**Subtype** | **int32** | Subtype of the fee | 
**UpperBound** | **Amount** | Upper bound of fee applicability | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


