package model

import (
	"crypto/rand"
	"math/big"
)

type Jumper struct {
	Start int
	End   int
}

type BoardPreset struct {
	Snakes  []Jumper
	Ladders []Jumper
}

var boardPresets = []BoardPreset{
	// Preset 1: The Great Fall - A massive snake at 99 sends you back to the start.
	{
		Snakes: []Jumper{
			{99, 2}, {85, 54}, {73, 51}, {68, 47}, {43, 22}, {31, 10},
		},
		Ladders: []Jumper{
			{4, 25}, {13, 46}, {33, 49}, {42, 62}, {50, 69}, {63, 81}, {74, 92},
		},
	},
	// Preset 2: The Final Gauntlet - The 90s are a treacherous path.
	{
		Snakes: []Jumper{
			{98, 77}, {96, 84}, {92, 72}, {88, 67}, {65, 44}, {48, 12},
		},
		Ladders: []Jumper{
			{3, 21}, {9, 30}, {23, 41}, {35, 56}, {52, 70}, {61, 79}, {80, 99},
		},
	},
	// Preset 3: The Rollercoaster - Many ups and downs throughout the board.
	{
		Snakes: []Jumper{
			{95, 75}, {89, 53}, {76, 58}, {62, 42}, {49, 11}, {36, 6},
		},
		Ladders: []Jumper{
			{2, 23}, {8, 31}, {20, 38}, {28, 77}, {40, 59}, {51, 67}, {71, 91},
		},
	},
	// Preset 4: Minimalist - Fewer snakes and ladders make each one more critical.
	{
		Snakes: []Jumper{
			{87, 24}, {71, 40}, {54, 34}, {32, 10}, {99, 63},
		},
		Ladders: []Jumper{
			{5, 27}, {15, 45}, {28, 84}, {41, 62}, {70, 89},
		},
	},
	// Preset 5: The Long Climb - Fewer but longer ladders and snakes.
	{
		Snakes: []Jumper{
			{97, 25}, {83, 43}, {64, 18}, {55, 35},
		},
		Ladders: []Jumper{
			{4, 56}, {13, 47}, {37, 79}, {48, 68}, {61, 98},
		},
	},
	// Preset 6: Second Chance - A big drop from 94, but a ladder at 19 offers recovery.
	{
		Snakes: []Jumper{
			{94, 16}, {78, 57}, {66, 45}, {59, 38}, {33, 12},
		},
		Ladders: []Jumper{
			{7, 26}, {19, 60}, {32, 51}, {48, 69}, {63, 82}, {72, 90},
		},
	},
	// Preset 7: Crowded Middle - The board is most dangerous in the 40-70 range.
	{
		Snakes: []Jumper{
			{99, 80}, {70, 50}, {64, 41}, {56, 44}, {49, 29}, {46, 25},
		},
		Ladders: []Jumper{
			{3, 22}, {9, 31}, {21, 42}, {28, 77}, {51, 67}, {71, 91}, {80, 98},
		},
	},
	// Preset 8: Fast Start - Early ladders give a boost, but the top half is guarded.
	{
		Snakes: []Jumper{
			{93, 73}, {86, 52}, {69, 34}, {53, 30}, {98, 4},
		},
		Ladders: []Jumper{
			{6, 17}, {11, 49}, {26, 47}, {36, 57}, {61, 82}, {72, 88},
		},
	},
	// Preset 9: The Two Drops - Two significant snakes define the game.
	{
		Snakes: []Jumper{
			{99, 39}, {62, 19}, {46, 25}, {34, 1}, {88, 54},
		},
		Ladders: []Jumper{
			{2, 38}, {7, 14}, {15, 26}, {21, 42}, {28, 84}, {51, 67}, {72, 91},
		},
	},
	// Preset 10: Balanced Classic - A standard, well-rounded board.
	{
		Snakes: []Jumper{
			{97, 78}, {95, 75}, {92, 88}, {89, 68}, {74, 53}, {62, 19}, {48, 26},
		},
		Ladders: []Jumper{
			{1, 38}, {4, 14}, {9, 31}, {21, 42}, {28, 84}, {36, 44}, {51, 67}, {71, 91},
		},
	},
	{
		// Volatile peak — a brutal 99->3 snake
		Snakes: []Jumper{
			{99, 3}, {92, 88}, {75, 54}, {62, 19}, {47, 16},
		},
		Ladders: []Jumper{
			{2, 25}, {9, 31}, {20, 38}, {28, 84}, {70, 90},
		},
	},
	{
		// Early boost but a hidden 99->2 trap
		Snakes: []Jumper{
			{99, 2}, {95, 72}, {87, 65}, {54, 34},
		},
		Ladders: []Jumper{
			{3, 22}, {8, 30}, {13, 46}, {41, 59},
		},
	},
	{
		// Balanced, mid-board climbs
		Snakes: []Jumper{
			{96, 53}, {88, 69}, {73, 50},
		},
		Ladders: []Jumper{
			{4, 18}, {11, 29}, {27, 45}, {33, 61}, {76, 97},
		},
	},
	{
		// Slow-and-steady ladders, few high snakes
		Snakes: []Jumper{
			{98, 77}, {84, 58}, {66, 44},
		},
		Ladders: []Jumper{
			{1, 10}, {5, 14}, {21, 41}, {35, 53},
		},
	},
	{
		// Two huge swings including 99->4
		Snakes: []Jumper{
			{99, 4}, {67, 24},
		},
		Ladders: []Jumper{
			{2, 15}, {6, 23}, {16, 39}, {43, 69},
		},
	},
	{
		// High ladder reward around the middle
		Snakes: []Jumper{
			{94, 71}, {80, 59}, {47, 20},
		},
		Ladders: []Jumper{
			{7, 32}, {12, 36}, {30, 52}, {56, 78},
		},
	},
	{
		// Fast lanes and a few punishing drops
		Snakes: []Jumper{
			{90, 48}, {76, 55}, {69, 33},
		},
		Ladders: []Jumper{
			{10, 27}, {22, 42}, {34, 60}, {50, 85},
		},
	},
	{
		// Small set, high drama with 99->3 again
		Snakes: []Jumper{
			{99, 3}, {85, 63},
		},
		Ladders: []Jumper{
			{2, 19}, {15, 37}, {28, 50}, {45, 66},
		},
	},
	{
		// Short snakes but helpful clustered ladders
		Snakes: []Jumper{
			{93, 14}, {58, 40},
		},
		Ladders: []Jumper{
			{4, 9}, {13, 31}, {24, 44}, {35, 68},
		},
	},
	{
		// A couple of deep drops and steady progressions
		Snakes: []Jumper{
			{97, 5}, {81, 68}, {61, 38},
		},
		Ladders: []Jumper{
			{6, 17}, {18, 29}, {25, 46}, {49, 72},
		},
	},
}


func GetBoardPreset() BoardPreset {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(boardPresets))))
	if err != nil {
		// Fallback to the first preset if randomness fails (very unlikely)
		return boardPresets[0]
	}
	return boardPresets[n.Int64()]
}