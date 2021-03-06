package conway

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromSlice(t *testing.T) {
	require := require.New(t)
	var grid Grid
	var err error

	// Valid input.
	grid, err = FromSlice([][]int{
		{1, 0, 0},
		{0, 0, 0},
		{0, 1, 1},
	})
	require.Nil(err)
	expected := []Cell{{0, 0}, {1, 2}, {2, 2}}
	actual := aliveCells(grid)
	require.Subset(expected, actual)
	require.Subset(actual, expected)

	// Invalid input.
	_, err = FromSlice([][]int{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	})
	require.Nil(err)
}

func TestFromString(t *testing.T) {
	require := require.New(t)
	var grid Grid
	var err error

	// Valid input, newline-delimited.
	grid, err = FromString(`
		x..
		...
		.xx
    `)
	require.Nil(err)
	expected := []Cell{{0, 0}, {1, 2}, {2, 2}}
	actual := aliveCells(grid)
	require.Subset(expected, actual)
	require.Subset(actual, expected)

	// Valid input, semicolon-delimited.
	grid, err = FromString(`
		x..;...;.xx     `)
	require.Nil(err)
	expected = []Cell{{0, 0}, {1, 2}, {2, 2}}
	actual = aliveCells(grid)
	require.Subset(expected, actual)
	require.Subset(actual, expected)

	// Invalid input.
	grid, err = FromString(`
        ...
        ...
        ...
    `)
	require.Nil(err)
}

func TestRandomGrid(t *testing.T) {
	require := require.New(t)
	var grid, _expected Grid
	var err error

	// All living Cells.
	grid, err = RandomGrid(3, 3, 1.0)
	require.Nil(err)
	_expected, _ = FromSlice([][]int{
		{1, 1, 1},
		{1, 1, 1},
		{1, 1, 1},
	})
	expected := aliveCells(_expected)
	actual := aliveCells(grid)
	require.Subset(expected, actual)
	require.Subset(actual, expected)

	// Some living Cells.
	grid, err = RandomGrid(3, 3, 0.5)
	require.Nil(err)

	// No living Cells.
	grid, err = RandomGrid(0, 3, 0.5)
	require.NotNil(err)
	require.Nil(grid)
	grid, err = RandomGrid(3, 0, 0.5)
	require.NotNil(err)
	require.Nil(grid)
	grid, err = RandomGrid(3, 3, 0)
	require.NotNil(err)
	require.Nil(grid)
}

func TestCell_neighbors(t *testing.T) {
	require := require.New(t)
	cell := Cell{0, 2}
	actual := cell.neighbors()
	expected := []Cell{
		Cell{-1, 1}, Cell{0, 1}, Cell{1, 1},
		Cell{-1, 2}, Cell{1, 2},
		Cell{-1, 3}, Cell{0, 3}, Cell{1, 3},
	}
	require.Equal(expected, actual)
}

func TestGrid_liveNeighbors(t *testing.T) {
	require := require.New(t)
	grid := mkGrid([][]int{
		{1, 0, 0, 1},
		{1, 1, 1, 0},
		{0, 0, 1, 1},
		{0, 0, 0, 0},
	})
	require.Equal(2, grid.liveNeighbors(Cell{0, 0}))
	require.Equal(4, grid.liveNeighbors(Cell{2, 1}))
	require.Equal(3, grid.liveNeighbors(Cell{2, 2}))
}

func TestGrid_nextCell(t *testing.T) {
	assert := assert.New(t)
	grid := mkGrid([][]int{
		{1, 0, 0, 1, 0},
		{0, 1, 0, 0, 0},
		{1, 0, 0, 1, 0},
		{1, 1, 0, 1, 1},
		{0, 0, 0, 1, 1},
	})
	// 0 live neighbors dies
	assert.Equal(Dead, grid.nextCell(Cell{3, 0}))
	// 1 live neighbor dies
	assert.Equal(Dead, grid.nextCell(Cell{0, 0}))
	// 2 live neighbors lives
	assert.Equal(Alive, grid.nextCell(Cell{1, 1}))
	assert.Equal(Alive, grid.nextCell(Cell{3, 2}))
	assert.Equal(Alive, grid.nextCell(Cell{1, 3}))
	// 3 live neighbors lives
	assert.Equal(Alive, grid.nextCell(Cell{0, 2}))
	assert.Equal(Alive, grid.nextCell(Cell{0, 3}))
	assert.Equal(Alive, grid.nextCell(Cell{3, 4}))
	assert.Equal(Alive, grid.nextCell(Cell{4, 4}))
	// 4+ live neighbors dies
	assert.Equal(Dead, grid.nextCell(Cell{3, 3}))
	assert.Equal(Dead, grid.nextCell(Cell{4, 3}))

	// 0-2 live neighbors stays dead
	assert.Equal(Dead, grid.nextCell(Cell{1, 0}))
	assert.Equal(Dead, grid.nextCell(Cell{2, 0}))
	assert.Equal(Dead, grid.nextCell(Cell{4, 0}))
	assert.Equal(Dead, grid.nextCell(Cell{3, 1}))
	assert.Equal(Dead, grid.nextCell(Cell{4, 1}))
	assert.Equal(Dead, grid.nextCell(Cell{0, 4}))
	assert.Equal(Dead, grid.nextCell(Cell{1, 4}))
	// 3 live neighbors is revived
	assert.Equal(Alive, grid.nextCell(Cell{0, 1}))
	assert.Equal(Alive, grid.nextCell(Cell{2, 1}))
	assert.Equal(Alive, grid.nextCell(Cell{4, 2}))
	assert.Equal(Alive, grid.nextCell(Cell{2, 4}))
	// 4+ live neighbors stays dead
	assert.Equal(Dead, grid.nextCell(Cell{1, 2}))
	assert.Equal(Dead, grid.nextCell(Cell{2, 2}))
	assert.Equal(Dead, grid.nextCell(Cell{2, 3}))
}

func TestGrid_withNeighbors(t *testing.T) {
	require := require.New(t)
	grid := mkGrid([][]int{
		{1, 0},
		{0, 1},
	})
	actual := grid.withNeighbors()
	expected := []Cell{
		Cell{-1, -1}, Cell{0, -1}, Cell{1, -1},
		Cell{-1, 0}, Cell{0, 0}, Cell{1, 0}, Cell{2, 0},
		Cell{-1, 1}, Cell{0, 1}, Cell{1, 1}, Cell{2, 1},
		Cell{0, 2}, Cell{1, 2}, Cell{2, 2},
	}
	actual_ := allCells(actual)
	require.Subset(expected, actual_)
	require.Subset(actual_, expected)
}

func TestGrid_Next(t *testing.T) {
	require := require.New(t)
	type gridPair struct {
		start Grid
		next  Grid
	}
	pairs := []gridPair{
		{
			start: mkGrid([][]int{
				{0, 1, 0},
				{0, 1, 0},
				{0, 1, 0},
			}),
			next: mkGrid([][]int{
				{0, 0, 0},
				{1, 1, 1},
				{0, 0, 0},
			}),
		},
		{
			start: mkGrid([][]int{
				{0, 0, 0, 0},
				{0, 1, 1, 1},
				{1, 1, 1, 0},
				{0, 0, 0, 0},
			}),
			next: mkGrid([][]int{
				{0, 0, 1, 0},
				{1, 0, 0, 1},
				{1, 0, 0, 1},
				{0, 1, 0, 0},
			}),
		},
		{
			start: mkGrid([][]int{
				{0, 0, 0, 0},
				{0, 1, 1, 0},
				{0, 1, 1, 0},
				{0, 0, 0, 0},
			}),
			next: mkGrid([][]int{
				{0, 0, 0, 0},
				{0, 1, 1, 0},
				{0, 1, 1, 0},
				{0, 0, 0, 0},
			}),
		},
	}

	for _, pair := range pairs {
		next, _ := pair.start.Next()
		expected := aliveCells(pair.next)
		actual := aliveCells(next)
		require.Subset(expected, actual)
		require.Subset(actual, expected)
	}
}

func TestGrid_xyBounds(t *testing.T) {
	require := require.New(t)
	grid := mkGrid([][]int{
		{1, 0, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 1, 1, 0},
	})
	grid.Set(Cell{-2, 0}, Alive)
	min, max := grid.xyBounds()
	require.Equal(Cell{-2, 0}, min)
	require.Equal(Cell{3, 3}, max)
}

func TestGrid_Show(t *testing.T) {
	require := require.New(t)
	grid := mkGrid([][]int{
		{1, 0, 0},
		{0, 0, 0},
		{0, 1, 1},
	})
	actual := grid.Show()
	x := LiveCellRepr
	o := DeadCellRepr
	expected := strings.Join(
		[]string{
			x + o + o,
			o + o + o,
			o + x + x,
		},
		"\n",
	)
	require.Equal(strings.TrimSpace(expected), strings.TrimSpace(actual))
}

func mkGrid(rows [][]int) Grid {
	grid, _ := FromSlice(rows)
	return grid
}

func aliveCells(grid Grid) []Cell {
	cells := []Cell{}
	for cell, state := range grid {
		if state == Alive {
			cells = append(cells, cell)
		}
	}
	return cells
}

func allCells(grid Grid) []Cell {
	cells := make([]Cell, len(grid))
	for cell := range grid {
		cells = append(cells, cell)
	}
	return cells
}
