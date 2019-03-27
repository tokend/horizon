# ReviewableRequestAttributes

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllTasks** | **int32** | All tasks that have been set for a request | 
**CreatedAt** | [**time.Time**](time.Time.md) | Time when a request has been submitted | 
**ExternalDetails** | [**Details**](Details.md) |  | [optional] 
**Hash** | **string** | Hash of a particular request | 
**PendingTasks** | **int32** | Tasks that have not been removed yet | 
**Reference** | **string** | Reference for the request | [optional] 
**RejectReason** | **string** | Details on why a request has been rejected | 
**State** | **string** | String representation of the request&#39;s state | 
**StateI** | **int32** | Integer representation of the request&#39;s state | 
**UpdatedAt** | [**time.Time**](time.Time.md) | Last time when a request has been updated | 
**XdrType** | **xdr.ReviewableRequestType** |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


