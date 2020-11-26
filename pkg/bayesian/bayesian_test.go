package bayesian_test

import (
	"testing"

	"github.com/hareku/bayesian/pkg/bayesian"
)

func TestClassifier(t *testing.T) {
	good := bayesian.Class("good")
	bad := bayesian.Class("bad")

	c, err := bayesian.NewClassifier(good, bad)
	if err != nil {
		t.Fatalf("NewClassifier failed: %s", err)
	}

	c.Learn([]string{"tall", "handsome", "rich"}, good)
	c.Learn([]string{"bald", "poor", "ugly"}, bad)

	score := c.ProbScores([]string{"tall"})
	if score[good] <= score[bad] {
		t.Errorf("tall should be classed good: %v", score)
	}

	score = c.ProbScores([]string{"bald"})
	if score[good] >= score[bad] {
		t.Errorf("bald should be classed bad: %v", score)
	}

	score = c.ProbScores([]string{"tall", "unknown"})
	if score[good] <= score[bad] {
		t.Errorf("tall and unknown should be classed good: %v", score)
	}
}
