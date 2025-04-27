package avl

import (
	"github.com/SwanHtetAungPhyo/learning/mainNode/internal/model"
	"github.com/gofiber/fiber/v2/log"
	"sort"
)

type Node struct {
	Val    model.Validator
	Left   *Node
	Right  *Node
	Height int16
}

func NewNode(val model.Validator) *Node {
	return &Node{Val: val,
		Left:   nil,
		Right:  nil,
		Height: 1}
}

func (n *Node) Insert(val model.Validator) *Node {
	if n == nil {
		return NewNode(val)
	}
	if val.Stake < n.Val.Stake {
		n.Left = n.Left.Insert(val)
	} else {
		n.Right = n.Right.Insert(val)
	}
	n.Height = max(n.GetHeight(n.Left), n.GetHeight(n.Right)) + 1
	balanceFactor := n.GetBalanceFactor(n)

	if balanceFactor > 1 && n.GetBalanceFactor(n.Left) >= 0 {
		return n.RightRotate(n)
	}
	if balanceFactor < -1 && n.GetBalanceFactor(n.Right) <= 0 {
		return n.LeftRotate(n)
	}
	if balanceFactor > 1 && n.GetBalanceFactor(n.Left) < 0 {
		n.Left.Left = n.Left.Left.LeftRotate(n.Left.Left)
		return n.RightRotate(n)
	}
	if balanceFactor < -1 && n.GetBalanceFactor(n.Right) > 0 {
		n.Right.Right = n.Right.Right.RightRotate(n.Right.Right)
		return n.LeftRotate(n)
	}
	return n
}

func (n *Node) Delete(val model.Validator) {

}

func (n *Node) Search(val model.Validator) *Node {
	if n == nil {
		return nil
	}
	if val.Stake < n.Val.Stake {
		return n.Left.Search(val)
	} else {
		return n.Right.Search(val)
	}
}

func (n *Node) GetHighestValidator() *Node {
	if n == nil {
		return nil
	}
	if n.Right == nil {
		return n
	}
	return n.Right.GetHighestValidator()
}

func (n *Node) GetHeight(node *Node) int16 {
	if n == nil {
		return 0
	}
	return n.Height
}
func (n *Node) GetBalanceFactor(node *Node) int16 {
	if n == nil {
		return 0
	}
	return n.GetHeight(n.Left) - n.GetHeight(n.Right)
}

func (n *Node) LeftRotate(x *Node) *Node {
	y := x.Right
	T2 := y.Left
	y.Left = x
	x.Right = T2
	x.Height = max(x.Left.Height, x.Right.Height) + 1
	y.Height = max(y.Left.Height, y.Right.Height) + 1
	return y
}

func (n *Node) RightRotate(x *Node) *Node {
	y := x.Left
	T2 := y.Right
	y.Right = x
	x.Left = T2
	x.Height = max(x.Left.Height, x.Right.Height) + 1
	y.Height = max(y.Left.Height, y.Right.Height) + 1
	return y
}

func (n *Node) CheckConsensus() *model.ConsensusResult {
	validators := n.getValidatorsInOrder()
	if len(validators) == 0 {
		return &model.ConsensusResult{
			IsAchieved: false,
			Validators: nil,
			TotalStake: 0,
		}
	}

	var totalStake int64
	for _, v := range validators {
		totalStake += int64(v.Stake)
	}

	// If total stake is 0, consensus cannot be achieved
	if totalStake == 0 {
		return &model.ConsensusResult{
			IsAchieved: false,
			Validators: validators,
			TotalStake: 0,
		}
	}

	var currentStake int64
	var consensusValidators []model.Validator

	// Calculate 2/3 threshold
	threshold := (totalStake*2 + 2) / 3 // Using ceiling division for precise threshold

	for _, v := range validators {
		consensusValidators = append(consensusValidators, v)
		currentStake += int64(v.Stake)

		if currentStake > threshold {
			return &model.ConsensusResult{
				IsAchieved: true,
				Validators: consensusValidators,
				TotalStake: currentStake,
			}
		}
	}

	return &model.ConsensusResult{
		IsAchieved: currentStake >= threshold,
		Validators: consensusValidators,
		TotalStake: currentStake,
	}
}
func (n *Node) getValidatorsInOrder() []model.Validator {
	var validators []model.Validator
	n.InorderTraversal(&validators)

	sort.Slice(validators, func(i, j int) bool {
		return validators[i].Stake > validators[j].Stake
	})

	return validators
}

func (n *Node) InorderTraversal(validators *[]model.Validator) {
	if n == nil {
		return
	}

	n.Left.InorderTraversal(validators)
	*validators = append(*validators, n.Val)
	log.Info(n.Val)
	n.Right.InorderTraversal(validators)
}
