# Unreleased

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

* `Limits` and `ExternalSystemIDs` to `/accounts` endpoint
* `/sales` endpoint
* `/order_book` endpoint

## Fixed

* changelog format

# 3.0.0-x.1

* `assets` owner is not rendering as `null` anymore
