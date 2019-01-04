package model

import (
	"encoding/json"
	"gitlab.com/tokend/regources"
)

type Account struct {
	Id            string                `json:"id"`
	ResourceType  string                `json:"type"`
	Attributes    *AccountAttributes    `json:"attributes, omitempty"`
	Relationships *AccountRelationships `json:"relationships, omitempty"`
	Included      []*Model              `json:"included, omitempty"`
}

func (a *Account) MarshalSelf() ([]byte, error) {
	return json.MarshalIndent(a, "", "  ")
}

type AccountAttributes struct {
	AccountType  AccountType         `json:"account_type"`
	BlockReasons AccountBlockReasons `json:"block_reasons"`
	IsBlocked    bool                `json:"is_blocked"`
	Policies     AccountPolicies     `json:"policies"`
	Thresholds   AccountThresholds   `json:"thresholds"`
}

type AccountRelationships struct {
	// TODO
}

type AccountBlockReasons struct {
	Types []regources.Flag `json:"types"`
	TypeI int32            `json:"type_i"`
}

type AccountType struct {
	Type  string `json:"type"`
	TypeI int32  `json:"type_i"`
}

type AccountThresholds struct {
	LowThreshold  byte `json:"low_threshold"`
	MedThreshold  byte `json:"med_threshold"`
	HighThreshold byte `json:"high_threshold"`
}

type AccountPolicies struct {
	TypeI int32            `json:"type_i"`
	Types []regources.Flag `json:"types"`
}
