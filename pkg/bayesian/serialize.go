package bayesian

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

// NewClassifierFromFile loads an existing classifier from
// file. The classifier was previously saved with a call
// to c.WriteToFile(string).
func NewClassifierFromFile(name string) (c *Classifier, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", name, err)
	}
	defer file.Close()

	return NewClassifierFromReader(file)
}

// NewClassifierFromReader deserializes of a Gob encoded classifier.
func NewClassifierFromReader(r io.Reader) (c *Classifier, err error) {
	dec := gob.NewDecoder(r)
	w := new(Classifier)
	err = dec.Decode(w)
	if err != nil {
		return nil, fmt.Errorf("failed to decode as Classifier: %w", err)
	}

	return w, nil
}

// WriteToFile serializes this classifier to a file.
func (c *Classifier) WriteToFile(name string) (err error) {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open a file: %w", err)
	}
	defer file.Close()

	return c.Write(file)
}

// WriteTo serializes this classifier to GOB and write to Writer.
func (c *Classifier) Write(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	err = enc.Encode(c)
	if err != nil {
		return fmt.Errorf("failed to encode Classifier: %w", err)
	}

	return nil
}
