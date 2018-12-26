package chess

import (
	"unicode"
)

func (c *Chess) Move(step [4]int) (state [2]Piece) {
	oldPiece := c.Board[step[0]*9+step[1]]
	newPiece, ok := c.Board[step[2]*9+step[3]]
	if !ok {
		newPiece.Index = step[2]*9 + step[3]
		newPiece.NumberIndex = step[2]
		newPiece.AlphaIndex = step[3]
	}
	state[0], state[1] = oldPiece, newPiece
	delete(c.Board, oldPiece.Index)
	newPiece.Type = oldPiece.Type
	c.Board[newPiece.Index] = newPiece
	if c.Next == 'b' {
		c.Next = 'w'
	} else {
		c.Next = 'b'
	}
	if state[0].Type == 'k' {
		c.BlackKingIndex = state[1].Index
	} else if state[0].Type == 'K' {
		c.RedKingIndex = state[1].Index
	} else if state[1].Type == 'k' {
		c.BlackKingIndex = -1
	} else if state[1].Type == 'K' {
		c.RedKingIndex = -1
	}
	return
}

func (c *Chess) UnMove(state [2]Piece) {
	c.Board[state[0].Index] = state[0]
	if state[1].Type == 0 {
		delete(c.Board, state[1].Index)
	} else {
		c.Board[state[1].Index] = state[1]
	}
	if c.Next == 'b' {
		c.Next = 'w'
	} else {
		c.Next = 'b'
	}
	if state[0].Type == 'k' {
		c.BlackKingIndex = state[0].Index
	} else if state[0].Type == 'K' {
		c.RedKingIndex = state[0].Index
	} else if state[1].Type == 'k' {
		c.BlackKingIndex = state[1].Index
	} else if state[1].Type == 'K' {
		c.RedKingIndex = state[1].Index
	}
}

var KMoveRules = [4][2]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}
var AMoveRules = [4][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
var BMoveRules = [4][2]int{{-2, -2}, {-2, 2}, {2, -2}, {2, 2}}
var NMoveRules = [8][2]int{{-2, -1}, {-2, 1}, {-1, -2}, {-1, 2}, {1, -2}, {1, 2}, {2, -1}, {2, 1}}
var PMoveRules = [3][2]int{{0, -1}, {0, 1}, {1, 0}}

func checkStepAndSet(chess *Chess, steps *[][4]int, p *Piece, newN, newA int) {
	newP, ok := chess.Board[newN*9+newA]
	if !ok || (chess.Next == 'b' && unicode.IsUpper(rune(newP.Type))) || (chess.Next != 'b' && unicode.IsLower(rune(newP.Type))) {
		*steps = append(*steps, [4]int{p.NumberIndex, p.AlphaIndex, newN, newA})
	}
}

func (c *Chess) GetNextSteps() (steps [][4]int) {
	for _, p := range c.Board {
		if (c.Next == 'b' && unicode.IsUpper(rune(p.Type))) || (c.Next != 'b' && unicode.IsLower(rune(p.Type))) {
			continue
		}
		switch unicode.ToLower(rune(p.Type)) {
		case 'k': // 将、帅
			for _, r := range KMoveRules {
				if ((p.NumberIndex+r[0] >= 0 && p.NumberIndex+r[0] <= 2) || (p.NumberIndex+r[0] >= 7 && p.NumberIndex+r[0] <= 9)) &&
					(p.AlphaIndex+r[1] >= 'd'-'a' && p.AlphaIndex+r[1] <= 'f'-'a') {
					canMove := true
					if r[1] != 0 {
						for i := p.NumberIndex - 1; i >= 0; i-- {
							if e, ok := c.Board[i*9+p.AlphaIndex+r[1]]; ok && unicode.ToLower(rune(e.Type)) == 'k' {
								canMove = false
							} else if ok && unicode.ToLower(rune(e.Type)) != 'k' {
								break
							}
						}
						for i := p.NumberIndex + 1; i <= 9; i++ {
							if e, ok := c.Board[i*9+p.AlphaIndex+r[1]]; ok && unicode.ToLower(rune(e.Type)) == 'k' {
								canMove = false
							} else if ok && unicode.ToLower(rune(e.Type)) != 'k' {
								break
							}
						}
					}
					if canMove {
						checkStepAndSet(c, &steps, &p, p.NumberIndex+r[0], p.AlphaIndex+r[1])
					}
				}
			}
			for i := p.NumberIndex - 1; i >= 0; i-- {
				if e, ok := c.Board[i*9+p.AlphaIndex]; ok && unicode.ToLower(rune(e.Type)) == 'k' {
					checkStepAndSet(c, &steps, &p, i, p.AlphaIndex)
				} else if ok && unicode.ToLower(rune(e.Type)) != 'k' {
					break
				}
			}
			for i := p.NumberIndex + 1; i <= 9; i++ {
				if e, ok := c.Board[i*9+p.AlphaIndex]; ok && unicode.ToLower(rune(e.Type)) == 'k' {
					checkStepAndSet(c, &steps, &p, i, p.AlphaIndex)
				} else if ok && unicode.ToLower(rune(e.Type)) != 'k' {
					break
				}
			}
		case 'a': // 士、仕
			for _, r := range AMoveRules {
				if ((p.NumberIndex+r[0] >= 0 && p.NumberIndex+r[0] <= 2) || (p.NumberIndex+r[0] >= 7 && p.NumberIndex+r[0] <= 9)) &&
					(p.AlphaIndex+r[1] >= 'd'-'a' && p.AlphaIndex+r[1] <= 'f'-'a') {
					checkStepAndSet(c, &steps, &p, p.NumberIndex+r[0], p.AlphaIndex+r[1])
				}
			}
		case 'b': // 象、相
			for _, r := range BMoveRules {
				if c.Next == 'b' {
					if p.NumberIndex+r[0]/2 >= 5 && p.NumberIndex+r[0]/2 <= 9 && p.AlphaIndex+r[1]/2 >= 'a'-'a' && p.AlphaIndex+r[1]/2 <= 'i'-'a' {
						if _, ok := c.Board[(p.NumberIndex+r[0]/2)*9+p.AlphaIndex+r[1]/2]; !ok &&
							p.NumberIndex+r[0] >= 5 && p.NumberIndex+r[0] <= 9 && p.AlphaIndex+r[1] >= 'a'-'a' && p.AlphaIndex+r[1] <= 'i'-'a' {
							checkStepAndSet(c, &steps, &p, p.NumberIndex+r[0], p.AlphaIndex+r[1])
						}
					}
				} else {
					if p.NumberIndex+r[0]/2 >= 0 && p.NumberIndex+r[0]/2 <= 4 && p.AlphaIndex+r[1]/2 >= 'a'-'a' && p.AlphaIndex+r[1]/2 <= 'i'-'a' {
						if _, ok := c.Board[(p.NumberIndex+r[0]/2)*9+p.AlphaIndex+r[1]/2]; !ok &&
							p.NumberIndex+r[0] >= 0 && p.NumberIndex+r[0] <= 4 && p.AlphaIndex+r[1] >= 'a'-'a' && p.AlphaIndex+r[1] <= 'i'-'a' {
							checkStepAndSet(c, &steps, &p, p.NumberIndex+r[0], p.AlphaIndex+r[1])
						}
					}
				}
			}
		case 'n': // 马
			for _, r := range NMoveRules {
				if p.NumberIndex+r[0]/2 >= 0 && p.NumberIndex+r[0]/2 <= 9 && p.AlphaIndex+r[1]/2 >= 'a'-'a' && p.AlphaIndex+r[1]/2 <= 'i'-'a' {
					if _, ok := c.Board[(p.NumberIndex+r[0]/2)*9+p.AlphaIndex+r[1]/2]; !ok &&
						p.NumberIndex+r[0] >= 0 && p.NumberIndex+r[0] <= 9 && p.AlphaIndex+r[1] >= 'a'-'a' && p.AlphaIndex+r[1] <= 'i'-'a' {
						checkStepAndSet(c, &steps, &p, p.NumberIndex+r[0], p.AlphaIndex+r[1])
					}
				}
			}
		case 'r': // 车
			for i := 1; p.NumberIndex+i <= 9; i++ {
				checkStepAndSet(c, &steps, &p, p.NumberIndex+i, p.AlphaIndex)
				if _, ok := c.Board[(p.NumberIndex+i)*9+p.AlphaIndex]; ok {
					break
				}
			}
			for i := -1; p.NumberIndex+i >= 0; i-- {
				checkStepAndSet(c, &steps, &p, p.NumberIndex+i, p.AlphaIndex)
				if _, ok := c.Board[(p.NumberIndex+i)*9+p.AlphaIndex]; ok {
					break
				}
			}
			for i := 1; p.AlphaIndex+i <= 'i'-'a'; i++ {
				checkStepAndSet(c, &steps, &p, p.NumberIndex, p.AlphaIndex+i)
				if _, ok := c.Board[p.NumberIndex*9+p.AlphaIndex+i]; ok {
					break
				}
			}
			for i := -1; p.AlphaIndex+i >= 'a'-'a'; i-- {
				checkStepAndSet(c, &steps, &p, p.NumberIndex, p.AlphaIndex+i)
				if _, ok := c.Board[p.NumberIndex*9+p.AlphaIndex+i]; ok {
					break
				}
			}
		case 'c': // 炮
			var i int
			for i = 1; p.NumberIndex+i <= 9; i++ {
				if _, ok := c.Board[(p.NumberIndex+i)*9+p.AlphaIndex]; !ok {
					steps = append(steps, [4]int{p.NumberIndex, p.AlphaIndex, p.NumberIndex + i, p.AlphaIndex})
				} else {
					break
				}
			}
			for i++; p.NumberIndex+i <= 9; i++ {
				if _, ok := c.Board[(p.NumberIndex+i)*9+p.AlphaIndex]; ok {
					checkStepAndSet(c, &steps, &p, p.NumberIndex+i, p.AlphaIndex)
					break
				}
			}
			for i = -1; p.NumberIndex+i >= 0; i-- {
				if _, ok := c.Board[(p.NumberIndex+i)*9+p.AlphaIndex]; !ok {
					steps = append(steps, [4]int{p.NumberIndex, p.AlphaIndex, p.NumberIndex + i, p.AlphaIndex})
				} else {
					break
				}
			}
			for i--; p.NumberIndex+i >= 0; i-- {
				if _, ok := c.Board[(p.NumberIndex+i)*9+p.AlphaIndex]; ok {
					checkStepAndSet(c, &steps, &p, p.NumberIndex+i, p.AlphaIndex)
					break
				}
			}
			for i = 1; p.AlphaIndex+i <= 'i'-'a'; i++ {
				if _, ok := c.Board[p.NumberIndex*9+p.AlphaIndex+i]; !ok {
					steps = append(steps, [4]int{p.NumberIndex, p.AlphaIndex, p.NumberIndex, p.AlphaIndex + i})
				} else {
					break
				}
			}
			for i++; p.AlphaIndex+i <= 'i'-'a'; i++ {
				if _, ok := c.Board[p.NumberIndex*9+p.AlphaIndex+i]; ok {
					checkStepAndSet(c, &steps, &p, p.NumberIndex, p.AlphaIndex+i)
					break
				}
			}
			for i = -1; p.AlphaIndex+i >= 'a'-'a'; i-- {
				if _, ok := c.Board[p.NumberIndex*9+p.AlphaIndex+i]; !ok {
					steps = append(steps, [4]int{p.NumberIndex, p.AlphaIndex, p.NumberIndex, p.AlphaIndex + i})
				} else {
					break
				}
			}
			for i--; p.AlphaIndex+i >= 'a'-'a'; i-- {
				if _, ok := c.Board[p.NumberIndex*9+p.AlphaIndex+i]; ok {
					checkStepAndSet(c, &steps, &p, p.NumberIndex, p.AlphaIndex+i)
					break
				}
			}
		case 'p': // 卒、兵
			if c.Next == 'b' {
				if p.NumberIndex == 5 || p.NumberIndex == 6 {
					checkStepAndSet(c, &steps, &p, p.NumberIndex-1, p.AlphaIndex)
				} else {
					for _, r := range PMoveRules {
						if p.NumberIndex-r[0] >= 0 && p.NumberIndex-r[0] <= 4 && p.AlphaIndex+r[1] >= 'a'-'a' && p.AlphaIndex+r[1] <= 'i'-'a' {
							checkStepAndSet(c, &steps, &p, p.NumberIndex-r[0], p.AlphaIndex+r[1])
						}
					}
				}
			} else {
				if p.NumberIndex == 3 || p.NumberIndex == 4 {
					checkStepAndSet(c, &steps, &p, p.NumberIndex+1, p.AlphaIndex)
				} else {
					for _, r := range PMoveRules {
						if p.NumberIndex+r[0] >= 5 && p.NumberIndex+r[0] <= 9 && p.AlphaIndex+r[1] >= 'a'-'a' && p.AlphaIndex+r[1] <= 'i'-'a' {
							checkStepAndSet(c, &steps, &p, p.NumberIndex+r[0], p.AlphaIndex+r[1])
						}
					}
				}
			}
		}
	}
	return
}
