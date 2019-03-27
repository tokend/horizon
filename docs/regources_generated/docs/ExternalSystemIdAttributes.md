# ExternalSystemIdAttributes

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BindedAt** | [**time.Time**](time.Time.md) | the time when the external system ID was binded | 
**Data** | **string** | identifier of an account in the external system. | 
**ExpiresAt** | [**time.Time**](time.Time.md) | this ID can be binded to another account in the system after the expiration time | 
**ExternalSystemType** | **int32** | type of the external system | 
**IsDeleted** | **bool** | if true, this external system ID will not be released back to bool after the expiration but will rather be removed | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


