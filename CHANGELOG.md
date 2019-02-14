## Unreleased

### Added

Helper "test" for quick transaction envelopes unmarshal

### Fixed

500 error on /history endpoint when receiving history of payments

## 3.0.1-x.2

### Added

* added handler for `POST /v2/transactions`

### Changed

* (internal) Janus config

### Fixed

* Fixed rendering success operation state when corresponding request is rejected

## 3.0.1-x.1

## Fixed

* (internal) Fixed panic on ingest_v2 create account op trying to get referrer accountID which might not exist
* (internal) Fixed nil pointer exception on ingest_v2 on withdrawal participant effect handling

## 3.0.1-x.0

### Changed

* Updated XDR

## 3.0.0-x.2

### Added

* `/sales` endpoint
* `/order_book` endpoint

### Fixed

* changelog format

## 3.0.0-x.1

* `assets` owner is not rendering as `null` anymore
