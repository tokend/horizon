# Horizon

## Changelog

All notable changes to this project will be documented in [this file](./CHANGELOG.md). This project adheres to Semantic Versioning.
Please be careful during upgrading to 3.12.5 - Ingest from scratch is recommended. 

## Startup

- Entry point package: `cmd/horizon`
- Copy config file from config.yaml
- Specify config file in command line by option `--config <path_to_config>`
- Build and run Horizon

## Contribution notes

* When rendering the response that contains slice fields, **NEVER** render `null` if the slice is empty.


## Ingest Performance 
Data: catchup from start for first 400 ledgers with 16590 transaction in total

|Version of Horizon| v1| v2|
|:---:|:---:|:----:|
|Batch load of 100 ledgers| 2m13.0366828s | 1m49.039161s|
|Fixed accounts map for v2| -|55.0308887s|