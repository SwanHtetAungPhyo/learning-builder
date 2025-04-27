package model

type Validator struct {
	Name  string `mapstructure:"name"`
	Stake int    `mapstructure:"stake"`
}

// ConsensusResult represents the result of consensus checking
type ConsensusResult struct {
	IsAchieved bool
	Validators []Validator
	TotalStake int64
}
