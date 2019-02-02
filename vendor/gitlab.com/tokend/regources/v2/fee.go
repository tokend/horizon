package regources

// Fee - describes fee happened on balance. Direction of fee depends on the operation (depending on effect might be
// charged, locked, unlocked, for all incoming effects but unlocked it's always charged)
type Fee struct {
	Fixed             Amount `json:"fixed"`
	CalculatedPercent Amount `json:"calculated_percent"`
}

type FeeStr struct {
	Fixed             string `json:"fixed"`
	CalculatedPercent string `json:"calculated_percent"`
}
