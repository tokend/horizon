# SignerRuleAttributes

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Action** | **xdr.SignerRuleAction** | defines an action to be performed over the specified resource * 1:  \&quot;any\&quot; * 2:  \&quot;create\&quot; * 3:  \&quot;create_for_other\&quot; * 4:  \&quot;update\&quot; * 5:  \&quot;manage\&quot; * 6:  \&quot;send\&quot; * 7:  \&quot;remove\&quot; * 8:  \&quot;cancel\&quot; * 9:  \&quot;review\&quot; * 10: \&quot;receive_atomic_swap\&quot; * 11: \&quot;participate\&quot; * 12: \&quot;bind\&quot; * 13: \&quot;update_max_issuance\&quot; * 14: \&quot;check\&quot;  | 
**Details** | [**Details**](Details.md) |  | 
**IsDefault** | **bool** | defines whether this rule should be included into all new roles | 
**IsForbid** | **bool** | defines whether the specified action is forbidden | 
**Resource** | [**xdr.SignerRuleResource**](map[string]interface{}.md) | defines a resource to which the rule is applied. TODO: add link to XDR | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


