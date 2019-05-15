## Unreleased

### Fixed

* Sale participation response rendering

### Added

* `GetAtomicSwapBidList` handler to handle `v3/atomic_swap_bids` request
* `GetAtomicSwapBid` handler to handle `v3/atomic_swap_bids/{id}` request

### Removed

* handling `/atomic_swap_bids` request

## 3.3.1-x.0

### Added

* `/v3/accounts/{id}/sales` endpoint to get list of sales account is allowed to participate in
* `/v3/sales/{id}/relationships/whitelist` endpoint to get list of accounts explicitly granted 
permission to participate in sale.
* `/v3/sales/{id}/relationships/participation` endpoint to get list of sale participations
* `/v3/account/{id}/converted_balances/{asset_code}` endpoint to get collection of converted balances
* Sale participation on sale close ingestion
* Account specific rules ingestion
* Sale version ingestion
* Bumped up ingest version

## 3.3.0

### Added

* Allow forcing reingest on application start up (see ingest2/main.go)
* `/v3/limits` endpoint to get limits list with filters
* Squash matches with the same price in history
* `/v3/movements` endpoint to get participants effect related to balance movements
* `/v3` endpoint to get horizon info
* Tx submission system v2 with retry
* `/v3/transactions` to submit txs
* Cancel poll and update poll end time ingestion

## Fixed
* `/v3/account` to include all limits imposed for account
* Issue with `/v3/assets` for assets with max_uint64 fields
* Issues with balances and limits having fields types for amount int64 instead of uint64. 
* Update doorman with SignerOf check fixed
* Problem with CPU usage on OS X
* Fixed `null` messages on failed transaction submit

## 3.2.0-x.9

## Fixed

* `/v3/transactions/` includes were never populated

## 3.2.0-x.8

## Fixed

* `/v3/sales` was setting the wrong caps if default quote asset was also in the list of quote assets

## 3.2.0-x.7

### Fixed

* `/v3/votes` ignoring page params and always returning all the votes

## 3.2.0-x.6

### Added

* `/v3/transactions/{id}` endpoint to get transaction by hash or ID

## 3.2.0-x.5

### Added

* `/v3/balances/{id}` endpoint
* docs for `/v3/balances/{id}` endpoint

### Fixed

* 500 for `GET /balances`  because of nil pointer on balance relationships
* Default quote asset include for `v3/sales`

## 3.2.0-x.4

### Fixed

* 500 for `GET /accounts?include=balances` because of nil pointer on balance relationships
* Panics on manage vote/manage poll/manage poll creation requests operation details

## 3.2.0-x.0

### Removed

* dependency from regources/v2

### Fixed

* saving equal values to transaction's `result_meta_xdr` and `result_xdr`
* all docs for proper regources generation
* `ManageAssetOpAttributes.PreissuanceSigner` -> `ManageAssetOpAttributes.PreIssuanceSigner`
* json: `preissued_signer` -> `pre_issuance_signer`
* missing poll id in create poll request response
* order book entries sorting by price in alphabetical order

### Added

* generated regources/generated in vendor
* new endpoint `/v3/order_books/{id}`

### Changed

* old endpoint `/v3/order_book/{id}` marked as deprecated

## 3.1.1

### Fixed

* psql error on poll ingesting

## 3.1.0-rc.0

### Added

* `GetTransactionList` handler to handler `v3/transactions` endpoint
* Messages for new error codes
* `GetBalanceList` handler to handle `v3/balances` request
* `GetPublicKeyEntry` handler to handle `v3/public_key_entries/{id}` request
* Receiver filter for create issuance requests
* Polls and votes ingestion
* `GetPollList` to handle `/v3/polls`
* `GetPoll` to handle `/v3/polls/{id}`
* `GetVoteList` to handle `/v3/polls/{id}/relationships/votes`
* `GetVote` to handle `/v3/polls/{id}/relationships/votes{voter}`
* `GetCreatePollRequests` to handle `/v3/create_poll_requests`, `/v3/create_poll_requests/{id}`
*  Docs on polls

### Fixed

* panic on handle set fee operation when account or account role does not exist
* panic on actions' `isAllowed` checking
* error on get reviewable request by reviewer
* order book is now sorted
* create change role ingestion causing 500 on operations
* Now Bad Request is correctly returned in case of invalid signature

### Removed

* charts, endpoint for charts etc.
* influx

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
