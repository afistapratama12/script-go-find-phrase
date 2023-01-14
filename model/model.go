package model

type Phrase struct {
	Id             int
	Phrase         string
	WalletEtherBSC string
	WalletLitecoin string
	WalletBitcoin  string
}

type Phrases []Phrase

type Result struct {
	Id              int
	PhraseID        int
	BalanceEther    string
	BalanceBSC      string
	BalanceTokenBSC string
	BalanceLitecoin string
	BalanceBitcoin  string
}

type Results []Result
