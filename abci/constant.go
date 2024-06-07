package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Bech32Prefix = "yymm"
)

func init() {
	SetConfig()
}

func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32Prefix, Bech32Prefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(Bech32Prefix+sdk.PrefixValidator+sdk.PrefixOperator, Bech32Prefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(Bech32Prefix+sdk.PrefixValidator+sdk.PrefixConsensus, Bech32Prefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic)
	sdk.SetCoinDenomRegex(func() string {
		return `[a-zA-Z][a-zA-Z0-9/:._-]{1,127}`
	})
}
