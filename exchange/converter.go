package exchange

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/swarmfund/horizon/db2/core"
)

type assetProvider interface {
	GetAssetsForPolicy(policy uint32) ([]core.Asset, error)
	GetAssetPairsForCodes(baseAssets []string, quoteAssets []string) ([]core.AssetPair, error)
}

type Converter struct {
	assetProvider assetProvider
	baseAssets    []string
}

type assetProviderImpl struct {
	q core.QInterface
}

func (p assetProviderImpl) GetAssetsForPolicy(policy uint32) ([]core.Asset, error) {
	return p.q.Assets().ForPolicy(policy).Select()
}
func (p assetProviderImpl) GetAssetPairsForCodes(baseAssets []string, quoteAssets []string) ([]core.AssetPair, error) {
	return p.q.AssetPairs().ForAssets(baseAssets, quoteAssets).Select()
}

func NewConverter(q core.QInterface) (*Converter, error) {
	return newConverter(assetProviderImpl{q: q})
}

func newConverter(assetProvider assetProvider) (*Converter, error) {
	baseAssets, err := assetProvider.GetAssetsForPolicy(uint32(xdr.AssetPolicyBaseAsset))
	if err != nil {
		return nil, errors.Wrap(err, "failed to load base assets")
	}

	result := &Converter{
		assetProvider: assetProvider,
	}

	for i := range baseAssets {
		result.baseAssets = append(result.baseAssets, baseAssets[i].Code)
	}

	return result, nil
}

func (c *Converter) loadPairsWithBaseAssets(asset string) ([]core.AssetPair, error) {

	direct, err := c.assetProvider.GetAssetPairsForCodes([]string{asset}, c.baseAssets)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load direct asset pairs")
	}

	reverse, err := c.assetProvider.GetAssetPairsForCodes(c.baseAssets, []string{asset})
	if err != nil {
		return nil, errors.Wrap(err, "failed to load reverse asset pairs")
	}

	return append(direct, reverse...), nil
}

func (c *Converter) tryLoadDirect(fromAsset, toAsset string) (*core.AssetPair, error) {
	pairs, err := c.assetProvider.GetAssetPairsForCodes([]string{fromAsset, toAsset}, []string{fromAsset, toAsset})
	if err != nil {
		return nil, errors.Wrap(err, "failed to load direct asset pairs")
	}

	for i, pair := range pairs {
		if (pair.BaseAsset == fromAsset && pair.QuoteAsset == toAsset) ||
			(pair.BaseAsset == toAsset && pair.QuoteAsset == fromAsset) {
			return &pairs[i], nil
		}
	}

	return nil, nil
}

func (c *Converter) convertWithMaxPath(amount int64, fromAsset, toAsset string, fromPairs, toPairs []core.AssetPair) (*int64, error) {
	converted := false
	var result int64
	for _, fromPair := range fromPairs {
		hopAmount, isConverted, err := fromPair.ConvertFromSourceAsset(fromAsset, amount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert from asset to hop asset")
		}

		if !isConverted {
			continue
		}

		for _, toPair := range toPairs {
			if !fromPair.IsOverlaps(toPair) {
				continue
			}

			destAmount, isConverted, err := fromPair.ConvertToDestAsset(toAsset, hopAmount)
			if err != nil {
				return nil, errors.Wrap(err, "failed to convert to toAsset")
			}

			if !isConverted {
				continue
			}

			converted = true
			if destAmount > result {
				result = destAmount
			}
		}
	}

	if converted {
		return &result, nil
	}

	return nil, nil
}

func (c *Converter) TryToConvertWithOneHop(amount int64, fromAsset, toAsset string) (*int64, error) {
	if fromAsset == toAsset {
		return &amount, nil
	}

	directPair, err := c.tryLoadDirect(fromAsset, toAsset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load direct asset pair")
	}

	if directPair != nil {
		result, isConverted, err := directPair.ConvertToDestAsset(toAsset, amount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert using direct pair")
		}

		if !isConverted {
			return nil, nil
		}

		return &result, nil
	}

	fromAssetPairs, err := c.loadPairsWithBaseAssets(fromAsset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load pairs with base asset for from asset")
	}

	toAssetPairs, err := c.loadPairsWithBaseAssets(toAsset)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load pairs with base asset for to asset")
	}

	return c.convertWithMaxPath(amount, fromAsset, toAsset, fromAssetPairs, toAssetPairs)
}
