package avl

import (
	"github.com/SwanHtetAungPhyo/learning/mainNode/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConsensusCheck(t *testing.T) {
	tests := []struct {
		name          string
		validators    []model.Validator
		wantConsensus bool
		wantStake     int64
	}{
		{
			name: "Consensus achieved with more than 2/3 stake",
			validators: []model.Validator{
				{Name: "validator1", Stake: 1000},
				{Name: "validator2", Stake: 500},
				{Name: "validator3", Stake: 500},
			},
			wantConsensus: true,
			wantStake:     1500, // First two validators have enough stake
		},
		{
			name: "Consensus not achieved with less than 2/3 stake",
			validators: []model.Validator{
				{Name: "validator1", Stake: 300},
				{Name: "validator2", Stake: 300},
				{Name: "validator3", Stake: 300},
			},
			wantConsensus: false,
			wantStake:     900,
		},
		{
			name: "Consensus with single validator having more than 2/3 stake",
			validators: []model.Validator{
				{Name: "validator1", Stake: 1000},
				{Name: "validator2", Stake: 200},
				{Name: "validator3", Stake: 200},
			},
			wantConsensus: true,
			wantStake:     1000,
		},
		{
			name:          "Empty validator set",
			validators:    []model.Validator{},
			wantConsensus: false,
			wantStake:     0,
		},
		{
			name: "Complex scenario with 10 validators",
			validators: []model.Validator{
				{Name: "validator1", Stake: 1000},
				{Name: "validator2", Stake: 500},
				{Name: "validator3", Stake: 750},
				{Name: "validator4", Stake: 1200},
				{Name: "validator5", Stake: 900},
				{Name: "validator6", Stake: 800},
				{Name: "validator7", Stake: 1500},
				{Name: "validator8", Stake: 600},
				{Name: "validator9", Stake: 1100},
				{Name: "validator10", Stake: 950},
			},
			wantConsensus: true,
			wantStake:     6250, // Sum of top validators exceeding 2/3 of total
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create AVL tree and insert validators
			var root *Node
			for _, v := range tt.validators {
				root = root.Insert(v)
			}

			// Check consensus
			result := root.CheckConsensus()

			// Assert results
			assert.Equal(t, tt.wantConsensus, result.IsAchieved, "Consensus achievement mismatch")
			if tt.wantConsensus {
				assert.GreaterOrEqual(t, result.TotalStake, tt.wantStake,
					"Consensus stake should be greater than or equal to expected")
			}
		})
	}
}

func TestGetValidatorsInOrder(t *testing.T) {
	// Create test validators
	validators := []model.Validator{
		{Name: "validator1", Stake: 500},
		{Name: "validator2", Stake: 1000},
		{Name: "validator3", Stake: 750},
	}

	// Create AVL tree and insert validators
	var root *Node
	for _, v := range validators {
		root = root.Insert(v)
	}

	// Get validators in order
	orderedValidators := root.getValidatorsInOrder()

	// Assert the order (should be descending by stake)
	assert.Equal(t, 3, len(orderedValidators), "Should have all validators")
	assert.Equal(t, 1000, orderedValidators[0].Stake, "First validator should have highest stake")
	assert.Equal(t, 750, orderedValidators[1].Stake, "Second validator should have second highest stake")
	assert.Equal(t, 500, orderedValidators[2].Stake, "Third validator should have lowest stake")
}

func TestEdgeCases(t *testing.T) {
	t.Run("Nil root", func(t *testing.T) {
		var root *Node
		result := root.CheckConsensus()
		assert.False(t, result.IsAchieved, "Nil root should not achieve consensus")
		assert.Equal(t, int64(0), result.TotalStake, "Nil root should have zero stake")
		assert.Empty(t, result.Validators, "Nil root should have no validators")
	})

	t.Run("Single validator", func(t *testing.T) {
		var root *Node
		root = root.Insert(model.Validator{Name: "single", Stake: 1000})
		result := root.CheckConsensus()
		assert.True(t, result.IsAchieved, "Single validator should achieve consensus")
		assert.Equal(t, int64(1000), result.TotalStake, "Should have correct stake")
		assert.Len(t, result.Validators, 1, "Should have one validator")
	})

	t.Run("Zero stake validators", func(t *testing.T) {
		var root *Node
		root = root.Insert(model.Validator{Name: "v1", Stake: 0})
		root = root.Insert(model.Validator{Name: "v2", Stake: 0})
		result := root.CheckConsensus()
		assert.False(t, result.IsAchieved, "Zero stake validators should not achieve consensus")
		assert.Equal(t, int64(0), result.TotalStake, "Should have zero total stake")
	})
}

func TestConsensusThreshold(t *testing.T) {
	tests := []struct {
		name       string
		validators []model.Validator
		want       bool
	}{
		{
			name: "Exactly 2/3 stake",
			validators: []model.Validator{
				{Name: "v1", Stake: 667},
				{Name: "v2", Stake: 333},
			},
			want: true,
		},
		{
			name: "Just above 2/3 stake",
			validators: []model.Validator{
				{Name: "v1", Stake: 668},
				{Name: "v2", Stake: 332},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var root *Node
			for _, v := range tt.validators {
				root = root.Insert(v)
			}
			result := root.CheckConsensus()
			assert.Equal(t, tt.want, result.IsAchieved)
		})
	}
}
