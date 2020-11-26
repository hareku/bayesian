package bayesian_test

import (
	"os"
	"testing"

	"github.com/hareku/bayesian/pkg/bayesian"
)

func TestNewClassifierFromFile(t *testing.T) {
	good := bayesian.Class("good")
	bad := bayesian.Class("bad")

	c, err := bayesian.NewClassifier(good, bad)
	if err != nil {
		t.Fatalf("NewClassifier failed: %s", err)
	}
	c.Learn([]string{"tall", "handsome", "rich"}, good)
	c.Learn([]string{"bald", "poor", "ugly"}, bad)

	name := "./test-learned"
	t.Cleanup(func() {
		os.Remove(name)
	})

	err = c.WriteToFile(name)
	if err != nil {
		t.Fatalf("WriteToFile failed: %s", err)
	}

	c2, err := bayesian.NewClassifierFromFile(name)
	if err != nil {
		t.Fatalf("NewClassifierFromFile failed: %s", err)
	}

	if len(c2.Classes) != 2 {
		t.Errorf("classes len is not 2: %v", c2.Classes)
	}

	word := c2.Words.Get("tall")
	if word == nil {
		t.Fatal("serialized word is nill")
	}
	if word.Total != 1 {
		t.Errorf("word count is not 1, it's %d", word.Total)
	}
	if word.Count[good] != 1 {
		t.Errorf("good word count is not 1, it's %d", word.Count[good])
	}
	if word.Count[bad] != 0 {
		t.Errorf("bad word count is not 0, it's %d", word.Count[bad])
	}
}
