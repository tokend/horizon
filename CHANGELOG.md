## Unreleased

### Added

* `GetBalanceList` handler to handle `v3/balances` request
* `GetPublicKeyEntry` handler to handle `v3/public_key_enties/{id}` request

### Fixed

* panic on actions' `isAllowed` checking

### Removed

* charts, endpoint for charts etc.

## 3.0.1-x.21
## Fixed
* Invalid fees joined for account

## 3.0.1-x.20
## Fixed
* Not valid percent amount for manage offer op

## 3.0.1-x.19
## Added
* Error codes for manageSigner and skiped op
## Fixed
* Information disclosure without permission for reviewable requests v1 

## 3.0.1-x.18
## Fixed

* Ingest v2 stop on sale participation delete

## 3.0.1-x.17
* Fixed issue with 0 sale state

## 3.0.1-x.16

* Added license op and stamp op to history

## 3.0.1-x.15

### Fixed
* Ingest v2 fail on cancel sale request

## 3.0.1-x.14

### Added

* Asset type to asset and create asset response
* Ingest of details of LicenseOp and StampOp
* Error for v3 history, when trying to request filtration by non existing account/balance

### Changed

* xdr revision
* fix limits select
* Messages for error codes

### Fixed
* Change role request creator details null
* Fixed state of the reviewable requests
* Added xdr type to reviewable request
* Fees list when account type filter is not requested
* 401 on v1 operations when skipSigCheck is on
* Ingest v2 stop on check_sale_state op


# 3.0.1-x.13

### Fixed

* Returning account and signer rule action using xdr type in account and signer rules responses
* (internal) proper types for reviewable requests amounts
* Removed KYCData from ChangeRoleRequest

# 3.0.1-x.11

### Fixed

* 500 on `/history` when there are update signer role operations
* pagination in docs
* all endpoints in docs has padlocks and 401 errcode where needed
* required & non-required fields

# 3.0.1-x.10

### Added

* Added xdr revision to root response

### Fixed
* Fixed issue with empty filter present in url params handled as requested filter

# 3.0.1-x.9

### Changed

* Rename `Details`/`Reason`/`NewDetails` field names in operation which create reqeusts unified to `creatorDetails`
* Rename isForbid to forbids
* Fixed issues with fee
* Fixed issues with filters by flags or dest account for change role requests in v1

### Fixed

* Error codes
* `request_details` format to satisfy JSON API spec

# 3.0.1-x.7

### Added 

* Signers endpoint
* Account/Signer Role/Rules
* Endpoints for reviewable requests
* Docs on reviewable requests and sale

### Changed

* `Details`/`Reason`/`NewDetails` field names in reviewable request types unified to `creatorDetails`
* Corresponding messages in `messages` map
* Calculated fee response

### Removed

* `UpdateSaleEndTimeRequest` type

### Fixed

* Fee bounds
* Fee calculation for account

# 3.0.1-x.6

### Fixed
* 500 on fees

# 3.0.1-x.5

### Added

* Helper "test" for quick transaction envelopes unmarshal

### Fixed

* 500 error on /history endpoint when receiving history of payments
* Fixed issue with balanceID not been set for balance effects
* Fixed Fee 500 if asset does not exists


# 3.0.1-x.4

## Fixed

* Switched endpoint version to only track major version
* Switched config back to default

# 3.0.1-x.3

## Added

Added back proxy to API

## Fixed

* Allow to specify filter for primary market offers via orderBookID = -1
* (internal) Fixed ingest v2. Participant effect has not included asset.
* (internal) Fixed ingest v2. Correctly handle `fulfilled` on review of request

# 3.0.1-x.2

## Added

* added handler for `POST /v3.0/transactions`

## Changed

* (internal) Janus config
* `v2`-prefixed endpoints updated to `v3.0`-prefixed

## Fixed

* Fixed rendering success operation state when corresponding request is rejected

# 3.0.1-x.1

## Fixed

* (internal) Fixed panic on ingest_v2 create account op trying to get referrer accountID which might not exist
* (internal) Fixed nil pointer exception on ingest_v2 on withdrawal participant effect handling

# 3.0.1-x.0

## Changed

* Updated XDR

# 3.0.0-x.2

## Added

* `/sales` endpoint
* `/order_book` endpoint

## Fixed

* changelog format

# 3.0.0-x.1

* `assets` owner is not rendering as `null` anymore
