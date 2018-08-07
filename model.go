package main

type KraRet struct {
	Error  []interface{} `json:"error"`
	Result KraRes        `json:"result"`
}

type KraRes struct {
	XXBTZCAD TradeData `json:"xxbtzcad"`
}

type TradeData struct {
	Asks [][]interface{} `json:"asks"`
}

type Price struct {
	Vol   float64
	Price float64
}

type QuadRet struct {
	TS   string     `json:"timestamp"`
	Bids [][]string `json:"bids"`
}
