package resource

import (
	"fmt"

	"bullioncoin.githost.io/development/go/amount"
	"bullioncoin.githost.io/development/horizon/db2/core"
)

type AssetAmountInfo struct {
	AvailableEmission  string            `json:"available_emission"`
	CoinsInCirculation string            `json:"coins_in_circulation"`
	Stats              map[string]string `json:"stats"`
}

func NewAssetAmountInfo(availableEmission, coinsInCirculation int64) (aai AssetAmountInfo) {
	aai.AvailableEmission = amount.String(availableEmission)
	aai.CoinsInCirculation = amount.String(coinsInCirculation)
	return aai
}

type CoinsAmountInfo map[string]AssetAmountInfo

func NewCoinsAmountInfo(availableEmissions, coinsInCirculation map[string]int64, stats []core.AssetStat) CoinsAmountInfo {
	cai := make(CoinsAmountInfo)

	for asset, aeAmount := range availableEmissions {
		ccAmount := int64(0)
		if amountInt, ok := coinsInCirculation[asset]; ok {
			ccAmount = amountInt
		}
		cai[asset] = NewAssetAmountInfo(aeAmount, ccAmount)
	}

	for asset, ccAmount := range coinsInCirculation {
		if _, ok := cai[asset]; ok {
			continue
		}
		cai[asset] = NewAssetAmountInfo(int64(0), ccAmount)
	}

	for _, stat := range stats {
		asset := cai[stat.Asset]
		asset.Stats = map[string]string{
			"hundreds":  fmt.Sprintf("%d", stat.Hundreds),
			"ones":      fmt.Sprintf("%d", stat.Ones),
			"remainder": amount.String(stat.Remainder),
		}
		cai[stat.Asset] = asset
	}

	return cai
}
