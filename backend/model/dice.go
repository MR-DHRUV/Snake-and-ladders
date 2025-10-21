package model

import (
	"math/rand"
	"time"
)


type DiceManager struct {
	DiceCount int `json:"dice_count" bson:"dice_count"`
	Sides     int `json:"sides" bson:"sides"`
}

func NewDiceManager(diceCount, sides int) *DiceManager {
	rand.Seed(time.Now().UnixNano())

	return &DiceManager{
		DiceCount: diceCount,
		Sides:     sides,
	}
}

func (dm *DiceManager) RollAll() (int, []int) {
	sum := 0
	rolls := make([]int, dm.DiceCount)
	for i := 0; i < dm.DiceCount; i++ {
		roll := rand.Intn(dm.Sides) + 1 // Generate a random number between 1 and Sides
		sum += roll
		rolls[i] = roll
	}

	return sum, rolls
}
