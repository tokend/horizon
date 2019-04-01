package rgenerated

type AssetPairAttributes struct {
	Policies xdr.AssetPairPolicy `json:"policies"`
	// defines an asset pair price as quote asset divided by base asset (i.e., amount of quote asset per 1 base asset)
	Price Amount `json:"price"`
}
