# Changelog

## Unreleased

### Added

* Asset type to asset, manage asset, create asset request resources

### Fixed

* Use xdr.AccountRuleAction and xdr.SignerRuleAction
* (internal) proper types for reviewable requests amounts
* (internal) removed KYCData from ChangeRoleRequest

## Unreleased

### Changed

* Rename is_forbid to forbids
* Rename details to creator_details in ops:
    create_aml_alert_request
    create_change_role_request
    create_issuance_request
    create_manage_limits_request
    create_pressuance_request
    create_sale_request
    create_withdraw_request
    manage_asset

### Added

* Relations `Limits` and `ExternalSystemIDs` to Account

## 3.0.1-x.0

### Changed

* Bumped tokend/go to `3.0.2-x.0` (switched to XDR `> 3.0.1`)

## 3.0.0-x.1

### Changed

* `PreIssuanceRequestAttrs.Amount` now has `Amount` type
* `WithdrawalRequestAttrs.Fee` is now type `Fee` instead of `FeeStr`

### Removed

* `FeeStr` type
