package model

type BoardCell struct {
	Pos     int ` json:"pos" bson:"pos"`
	NextPos int ` json:"next_pos" bson:"next_pos"`
}

func (b *BoardCell) HasJumper() bool {
	return b.NextPos != b.Pos
}

func (b *BoardCell) GetNextPosition() int {
	return b.NextPos
}

type Board struct {
	Cells      []*BoardCell `json:"cells" bson:"cells"`
	TotalCells int          `json:"total_cells" bson:"total_cells"`
}

func (b *Board) GetCell(position int) *BoardCell {
	if position > b.TotalCells {
		return b.Cells[b.TotalCells] // Return the last cell if out of bounds
	} else if position <= 0 {
		return b.Cells[0] // Return the first cell if out of bounds
	}
	return b.Cells[position]
}

func NewBoard(totalCells int) *Board {
	// Initialize the board with cells
	cells := make([]*BoardCell, totalCells+1)
	for i := 0; i <= totalCells; i++ {
		cells[i] = &BoardCell{Pos: i, NextPos: i}
	}

	// Load a random preset
	preset := GetBoardPreset();

	// Apply snakes
	for _, s := range preset.Snakes {
		if s.Start <= totalCells && s.End <= totalCells {
			cells[s.Start].NextPos = s.End
		}
	}

	// Apply ladders
	for _, l := range preset.Ladders {
		if l.Start <= totalCells && l.End <= totalCells {
			cells[l.Start].NextPos = l.End
		}
	}

	return &Board{
		Cells:      cells,
		TotalCells: totalCells,
	}
}

func NewBoardFromCells(cells []*BoardCell) *Board {
	totalCells := len(cells) - 1 // Assuming cells[0] is a dummy cell

	return &Board{
		Cells:      cells,
		TotalCells: totalCells,
	}
}
