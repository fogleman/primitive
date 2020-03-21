package primitive

// REQ 2.4.0 2.4.1/2.4.2/2.4.3/2.4.4
type Difficulty int

const (
	easy   Difficulty = 0
	medium Difficulty = 2
	hard   Difficulty = 3
	expert Difficulty = 4
)

func setDifficulty() int {
	//Todo easy=50, medium=100, hard=150, expert=300
	return 0
}
