package main

import (
	"fmt"
	"strings"
)

func formatFloat(f float64) string {
	if f < 10 {
		// Print with up to 7 non-zero decimal numbers
		return fmt.Sprintf("%.*g", 7, f)
	} else if f < 1000 {
		// Print with 2 decimal numbers
		return fmt.Sprintf("%.2f", f)
	}
	return fmt.Sprintf("%s", comma(fmt.Sprintf("%.2f", f)))
}

func comma(s string) string {
	// Split the string into two parts: before and after the decimal point
	parts := strings.Split(s, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = "." + parts[1]
	}

	// Insert commas every three digits from the end
	n := len(intPart)
	for i := n - 3; i > 0; i -= 3 {
		intPart = intPart[:i] + "," + intPart[i:]
	}

	return intPart + decPart
}
