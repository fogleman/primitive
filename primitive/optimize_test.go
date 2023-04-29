package primitive

import (
	"fmt"
	"testing"
)

type hillClimbTestState struct {
	matchValues []int
	testValue   int
	unmoveCalls int
}

func (htState *hillClimbTestState) Copy() Annealable {
	return &hillClimbTestState{
		htState.matchValues, htState.testValue, htState.unmoveCalls}
}

// When we mutate, always increment the next value to be tested
// but save the current value in case we need to return to it
func (htState *hillClimbTestState) DoMove() interface{} {

	htState.testValue++
	return 0
}

func (htState *hillClimbTestState) UndoMove(undo interface{}) {
	//fmt.Println(fmt.Sprintf("Called undoMove. Unmove calls: %d", htState.unmoveCalls))
	htState.unmoveCalls += 1
}

// So the optimized value is going to be 1/ however many of the
// match values can be divided into the test value. So we can simply
// incriment as we mutate, and the behavior should be predictable.
func (htState *hillClimbTestState) Energy() float64 {

	// Continue to iterate until we run out of values or until
	// the current value is greater than the test value
	// i <= htState.matchValues[i] &&

	num_divisible := 0
	for i := 0; i < len(htState.matchValues) &&
		htState.testValue >= htState.matchValues[i]; i++ {
		quotient := htState.testValue % htState.matchValues[i]
		if quotient == 0 {
			num_divisible++
		}

	}

	if num_divisible == 0 {
		return 10.0
	}

	return 1 / float64(num_divisible)
}

func createHillClimbTestState() *hillClimbTestState {
	state := new(hillClimbTestState)
	state.testValue = 1
	state.matchValues = []int{2, 3, 5, 7, 11, 13, 17, 23, 29, 31, 37, 41}
	state.unmoveCalls = 0
	return state
}

type testHillClimbCase struct {
	age              int
	expected_energy  float64
	best_value       int
	exp_unmove_calls int
}

func TestHillClimb(t *testing.T) {

	cases := []testHillClimbCase{
		{1, 1, 2, 0}, {2, 1, 2, 0}, {3, 1, 2, 0}, {4, .5, 6, 3}, {5, .5, 6, 3},
		{6, .5, 6, 3}, {7, .5, 6, 3}, {8, .5, 6, 3}, {9, .5, 6, 3}, {10, .5, 6, 3},
		{11, .5, 6, 3}, {12, .5, 6, 3}, {13, .5, 6, 3}, {14, .5, 6, 3},
		{15, .5, 6, 3}, {16, .5, 6, 3}, {17, .5, 6, 3}, {18, .5, 6, 3},
		{19, .5, 6, 3}, {20, .5, 6, 3}, {21, .5, 6, 3}, {22, .5, 6, 3},
		{23, .5, 6, 3}, {24, .333333, 30, 26}, {25, .333333, 30, 26},
		{26, .3333333, 30, 26}, {27, .333333, 30, 26}, {28, .333333, 30, 26},
		{29, .3333333, 30, 26}, {30, .333333, 30, 26},
	}

	for _, c := range cases {

		//fmt.Println(fmt.Sprintf("\n Starting test case for age %d", c.age))
		testingState := createHillClimbTestState()
		resultState := HillClimb(testingState, c.age)

		if fmt.Sprintf("%f", resultState.Energy()) != fmt.Sprintf("%f", c.expected_energy) {
			t.Error(fmt.Sprintf(
				"Unexpected energy value in HillClimb for test value %d: %f. Expecting %f",
				resultState.(*hillClimbTestState).testValue,
				resultState.(*hillClimbTestState).Energy(),
				c.expected_energy))
			break
		}

		if resultState.(*hillClimbTestState).testValue != c.best_value {
			t.Error(fmt.Sprintf(
				"Unexpected best value in HillClimb: %d",
				resultState.(*hillClimbTestState).testValue))
			break
		}

		if resultState.(*hillClimbTestState).unmoveCalls != c.exp_unmove_calls {
			t.Error(fmt.Sprintf("Unexpected number of unmove calls: %d for age value: %d",
				resultState.(*hillClimbTestState).unmoveCalls, c.age))
			break
		}
	}
}
