package intersections

import (
	"fmt"
	"regexp"
)

type Intersection struct {
	Street1 string
	Street2 string
}

func Parse(input string) (*Intersection, error) {
	// Here we use a similar syntax to match the "corner"
	regex := `(?i)^\W*(.*?)\W*\s+(?:\band\b|\bat\b|&|\\|\/|\@)\s+(.*)(?:$|\W*)`
	re := regexp.MustCompile(regex)

	match := re.FindStringSubmatch(input)

	if len(match) < 3 {
		return nil, fmt.Errorf("failed to parse intersection from input string: %s", input)
	}

	return &Intersection{
		Street1: match[1],
		Street2: match[2],
	}, nil
}
