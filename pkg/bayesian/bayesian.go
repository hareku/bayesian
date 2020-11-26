package bayesian

import (
	"fmt"
)

// defaultProb is the tiny non-zero probability that a word
// we have not seen before appears in the class.
const defaultProb = 1e-11

// Class defines a class that the classifier will filter:
// C = {C_1, ..., C_n}. You should define your classes as a
// set of constants, for example as follows:
//
//    const (
//        Good Class = "Good"
//        Bad Class = "Bad
//    )
//
// Class values should be unique.
type Class string

// Classifier implements the Naive Bayesian Classifier.
type Classifier struct {
	Classes  []Class
	DocCount map[Class]int64
	Words    WordsTrie
}

// NewClassifier returns a new classifier. The classes provided
// should be at least 2 in number and unique.
func NewClassifier(classes ...Class) (*Classifier, error) {
	n := len(classes)

	if n < 2 {
		return nil, fmt.Errorf("require at least two classes, given %d classes", n)
	}

	// check uniqueness
	seen := make(map[Class]bool, n)
	for _, class := range classes {
		if seen[class] {
			return nil, fmt.Errorf("classes must be unique, %q is duplicated", class)
		}
		seen[class] = true
	}

	c := &Classifier{
		Classes:  classes,
		DocCount: make(map[Class]int64, n),
		Words:    NewWordsTrie(),
	}
	return c, nil
}

// WordData contains the total word count for each classes and total word count.
type WordData struct {
	Total int64
	Count map[Class]int64
}

// getWordProb returns P(W|C_j) -- the probability of seeing
// a particular word W in a document of this class.
func (c *Classifier) getWordProb(class Class, word string) float64 {
	wordData := c.Words.Get(word)
	if wordData == nil {
		return defaultProb
	}

	count, ok := wordData.Count[class]
	if !ok {
		return defaultProb
	}

	return float64(count) / float64(wordData.Total)
}

// getPriors returns the prior probabilities for the
// classes provided -- P(C_j).
//
// TODO: There is a way to smooth priors, currently
// not implemented here.
func (c *Classifier) getPriors() map[Class]float64 {
	n := len(c.Classes)
	sum := int64(0)
	priors := make(map[Class]float64, n)

	for _, class := range c.Classes {
		total := c.DocCount[class]
		priors[class] = float64(total)
		sum += total
	}

	if sum != 0 {
		for i, p := range priors {
			priors[i] = p / float64(sum)
		}
	}
	return priors
}

// Learn supervised learns the given document.
func (c *Classifier) Learn(document []string, class Class) {
	c.DocCount[class]++

	for _, word := range document {
		wordData := c.Words.GetOrNew(word)
		wordData.Count[class]++
		wordData.Total++
	}
}

// ProbScores calculates probabilities.
func (c *Classifier) ProbScores(doc []string) map[Class]float64 {
	n := len(c.Classes)
	scores := make(map[Class]float64, n)
	priors := c.getPriors()
	sum := float64(0)

	for _, class := range c.Classes {
		score := priors[class]
		for _, word := range doc {
			score *= c.getWordProb(class, word)
		}
		scores[class] = score
		sum += score
	}

	for i, s := range scores {
		scores[i] = s / sum
	}

	return scores
}
