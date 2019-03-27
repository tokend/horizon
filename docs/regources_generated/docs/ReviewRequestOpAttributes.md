# ReviewRequestOpAttributes

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Action** | [**xdr.ReviewRequestOpAction**](Enum.md) | * 1: \&quot;approve\&quot; * 2: \&quot;reject\&quot; * 3: \&quot;permanent_reject\&quot;  | 
**AddedTasks** | **int32** | Tasks that were added on the request review | 
**ExternalDetails** | [**Details**](Details.md) |  | 
**IsFulfilled** | **bool** | Whether request being reviewed was fulfilled | 
**Reason** | **string** | Reject reason | 
**RemovedTasks** | **int32** | Tasks that were removed on the request review | 
**RequestHash** | **string** | Hash of the request being reviewed | 
**RequestId** | **int32** | ID of the request being reviewed | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


