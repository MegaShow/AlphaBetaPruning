package chess

func GetBestMove(fen string, steps []string) string {
	if len(steps) >= 1 && steps[0] == "" {
		steps = steps[1:]
	}
	chess := NewChess(fen, steps)
	step := Search(chess)
	return step
}
