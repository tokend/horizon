package lorem

import "testing"

func TestRandomNameIter(t *testing.T) {
	for i := 0; i < 512; i++ {
		sofar := map[string]struct{}{}
		generator := RandomNameIter()
		for got := range generator {
			// there are should be no duplicates
			_, ok := sofar[got]
			if ok {
				t.Fatalf("got duplicate, %s", got)
			}

			// we expect it to generate at least 2k random names consistently
			if len(sofar) == 2000 {
				close(generator)
				break
			}

			sofar[got] = struct{}{}
		}
	}
}
