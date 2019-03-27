# AccountRuleAttributes

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Action** | **xdr.AccountRuleAction** | defines an action to be performed over the specified resource * 1:  \&quot;any\&quot; * 2:  \&quot;create\&quot; * 3:  \&quot;create_for_other\&quot; * 4:  \&quot;create_with_tasks\&quot; * 5:  \&quot;manage\&quot; * 6:  \&quot;send\&quot; * 7:  \&quot;withdraw\&quot; * 8:  \&quot;receive_issuance\&quot; * 9:  \&quot;receive_payment\&quot; * 10: \&quot;receive_atomic_swap\&quot; * 11: \&quot;participate\&quot; * 12: \&quot;bind\&quot; * 13: \&quot;update_max_issuance\&quot; * 14: \&quot;check\&quot; * 15: \&quot;cancel\&quot;  | 
**Details** | [**Details**](Details.md) |  | 
**IsForbid** | **bool** | defines whether or not the specified action is forbidden | 
**Resource** | [**xdr.AccountRuleResource**](.md) | defines resource to which the rule is applied. TODO: add link to XDR | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


