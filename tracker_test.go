package tracker

import (
	"fmt"
	"testing"
)

var gen1 = MockGenerator{}
var gen2 = MockGenerator{}
var gen3 = MockGenerator{}

var testPattern = Pattern{
	Track{
		Event{1, 1, gen1},
		Event{2, 1, gen1},
		Event{3, 1, gen1},
		Event{4, 1, gen1},
		Event{5, 1, gen1},
	},
	Track{
		Event{64, 127, gen2},
		Event{60, 127, gen2},
		Event{67, 127, gen2},
	},
	Track{
		Event{127, 127, gen3},
	},
}

func TestNothing(t *testing.T) {
	// TODO(aoeu): Remove after actually implementing tests.
	fmt.Println(testPattern)
}
