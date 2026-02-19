package jupiter

// TokenStats represents trading statistics for a time interval.
type TokenStats struct {
	PriceChange       *float64 `json:"priceChange,omitempty"`
	LiquidityChange   *float64 `json:"liquidityChange,omitempty"`
	VolumeChange      *float64 `json:"volumeChange,omitempty"`
	BuyVolume         *float64 `json:"buyVolume,omitempty"`
	SellVolume        *float64 `json:"sellVolume,omitempty"`
	BuyOrganicVolume  *float64 `json:"buyOrganicVolume,omitempty"`
	SellOrganicVolume *float64 `json:"sellOrganicVolume,omitempty"`
	NumBuys           *int     `json:"numBuys,omitempty"`
	NumSells          *int     `json:"numSells,omitempty"`
	NumTraders        *int     `json:"numTraders,omitempty"`
	NumOrganicBuyers  *int     `json:"numOrganicBuyers,omitempty"`
	NumNetBuyers      *int     `json:"numNetBuyers,omitempty"`
}

// Audit contains token security audit information.
type Audit struct {
	MintAuthorityDisabled   *bool    `json:"mintAuthorityDisabled,omitempty"`
	FreezeAuthorityDisabled *bool    `json:"freezeAuthorityDisabled,omitempty"`
	TopHoldersPercentage    *float64 `json:"topHoldersPercentage,omitempty"`
}

// FirstPool contains information about the token's first liquidity pool.
type FirstPool struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

// SwapInfo represents a single swap step in a route plan.
type SwapInfo struct {
	AmmKey    string `json:"ammKey"`
	Label     string `json:"label,omitempty"`
	InputMint string `json:"inputMint"`
	OutputMint string `json:"outputMint"`
	InAmount  string `json:"inAmount"`
	OutAmount string `json:"outAmount"`
	FeeAmount string `json:"feeAmount,omitempty"`
	FeeMint   string `json:"feeMint,omitempty"`
}

// RoutePlanStep represents a single step in a swap route plan.
type RoutePlanStep struct {
	SwapInfo SwapInfo `json:"swapInfo"`
	Percent  *int     `json:"percent,omitempty"`
	Bps      *int     `json:"bps,omitempty"`
}

// PlatformFee represents platform fee information.
type PlatformFee struct {
	Amount string `json:"amount"`
	FeeBps int    `json:"feeBps"`
}

// ExecuteRequest is the common request body for execute endpoints.
type ExecuteRequest struct {
	SignedTransaction string `json:"signedTransaction"`
	RequestID         string `json:"requestId"`
}

// ExecuteResponse is the common response from execute endpoints.
type ExecuteResponse struct {
	Status    string `json:"status"`
	Signature string `json:"signature,omitempty"`
	Error     string `json:"error,omitempty"`
	Code      *int   `json:"code,omitempty"`
}
