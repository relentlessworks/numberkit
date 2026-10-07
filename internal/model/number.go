package model

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// FormatOptions controls how a number is formatted.
type FormatOptions struct {
	Decimals    int    // number of decimal places
	Separator   string // thousands separator (default ",")
	DecimalSep  string // decimal separator (default ".")
	Prefix      string // prefix (e.g. "$")
	Suffix      string // suffix (e.g. "%")
	TrimZeros   bool   // trim trailing zeros after decimal point
	GroupDigits bool   // whether to group digits with separator
}

// Format formats a number string with the given options.
func Format(numStr string, opts FormatOptions) (string, error) {
	// Parse the number
	f, _, err := big.ParseFloat(numStr, 10, 256, big.ToNearestEven)
	if err != nil {
		return "", fmt.Errorf("invalid number: %s", numStr)
	}

	// Apply decimal places
	if opts.Decimals > 0 || opts.TrimZeros {
		f = setDecimals(f, opts.Decimals, opts.TrimZeros)
	}

	// Convert to string
	str := f.Text('f', -1)
	if opts.Decimals > 0 && !opts.TrimZeros {
		str = f.Text('f', opts.Decimals)
	}

	// Split into integer and fractional parts
	parts := strings.SplitN(str, ".", 2)
	intPart := parts[0]
	fracPart := ""
	if len(parts) > 1 {
		fracPart = parts[1]
	}

	// Handle negative sign
	negative := false
	if strings.HasPrefix(intPart, "-") {
		negative = true
		intPart = intPart[1:]
	}

	// Group digits
	if opts.GroupDigits && opts.Separator != "" {
		intPart = groupDigits(intPart, opts.Separator)
	}

	// Reassemble
	result := intPart
	if fracPart != "" {
		decSep := opts.DecimalSep
		if decSep == "" {
			decSep = "."
		}
		result += decSep + fracPart
	}

	if negative {
		result = "-" + result
	}

	return opts.Prefix + result + opts.Suffix, nil
}

func setDecimals(f *big.Float, decimals int, trim bool) *big.Float {
	if decimals > 0 {
		// Use strconv to format with the right number of decimals, then parse back
		// This ensures proper rounding
		s := f.Text('f', decimals)
		result, _, _ := big.ParseFloat(s, 10, 256, big.ToNearestEven)
		return result
	}
	if trim {
		return f
	}
	return f
}

func groupDigits(s, sep string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	var result strings.Builder
	first := n % 3
	if first > 0 {
		result.WriteString(s[:first])
		if n > first {
			result.WriteString(sep)
		}
	}
	for i := first; i < n; i += 3 {
		result.WriteString(s[i : i+3])
		if i+3 < n {
			result.WriteString(sep)
		}
	}
	return result.String()
}

// ToWords converts a number to its English word representation.
func ToWords(numStr string) (string, error) {
	f, _, err := big.ParseFloat(numStr, 10, 256, big.ToNearestEven)
	if err != nil {
		return "", fmt.Errorf("invalid number: %s", numStr)
	}

	// Check if negative
	negative := f.Sign() < 0
	if negative {
		f = new(big.Float).Abs(f)
	}

	// Check if zero
	if f.Sign() == 0 {
		return "zero", nil
	}

	// Split into integer and fractional parts
	str := f.Text('f', -1)
	parts := strings.SplitN(str, ".", 2)
	intPart := parts[0]
	fracPart := ""
	if len(parts) > 1 {
		fracPart = parts[1]
	}

	// Convert integer part
	intWords, err := integerToWords(intPart)
	if err != nil {
		return "", err
	}

	result := intWords

	// Convert fractional part as individual digits
	if fracPart != "" {
		var digitWords []string
		for _, c := range fracPart {
			d, _ := strconv.Atoi(string(c))
			digitWords = append(digitWords, ones[d])
		}
		result += " point " + strings.Join(digitWords, " ")
	}

	if negative {
		result = "negative " + result
	}

	return result, nil
}

var (
	ones = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	teens = []string{"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	tens = []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
	scales = []string{"", "thousand", "million", "billion", "trillion", "quadrillion", "quintillion", "sextillion", "septillion", "octillion", "nonillion", "decillion"}
)

func integerToWords(s string) (string, error) {
	// Remove leading zeros
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return "zero", nil
	}

	// Parse as big.Int
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return "", fmt.Errorf("invalid integer: %s", s)
	}

	if n.Sign() == 0 {
		return "zero", nil
	}

	// Process in groups of 3 digits
	var groups []string
	remainder := new(big.Int).Set(n)
	divisor := big.NewInt(1000)

	for remainder.Sign() > 0 {
		mod := new(big.Int)
		remainder, _ = remainder.DivMod(remainder, divisor, mod)
		groups = append(groups, mod.String())
	}

	var parts []string
	for i := len(groups) - 1; i >= 0; i-- {
		groupStr := threeDigitToWords(groups[i])
		if groupStr != "" {
			if i > 0 && i < len(scales) {
				groupStr += " " + scales[i]
			}
			parts = append(parts, groupStr)
		}
	}

	return strings.Join(parts, " "), nil
}

func threeDigitToWords(s string) string {
	n, _ := strconv.Atoi(s)
	if n == 0 {
		return ""
	}

	var parts []string

	hundreds := n / 100
	rest := n % 100

	if hundreds > 0 {
		parts = append(parts, ones[hundreds]+" hundred")
	}

	if rest > 0 {
		if rest < 10 {
			parts = append(parts, ones[rest])
		} else if rest < 20 {
			parts = append(parts, teens[rest-10])
		} else {
			t := rest / 10
			o := rest % 10
			if o > 0 {
				parts = append(parts, tens[t]+"-"+ones[o])
			} else {
				parts = append(parts, tens[t])
			}
		}
	}

	return strings.Join(parts, " and ")
}

func twoDigitToWords(s string) string {
	n, _ := strconv.Atoi(s)
	if n < 10 {
		return ones[n]
	}
	if n < 20 {
		return teens[n-10]
	}
	t := n / 10
	o := n % 10
	if o > 0 {
		return tens[t] + "-" + ones[o]
	}
	return tens[t]
}

// ToOrdinal converts a number to its ordinal form (1st, 2nd, 3rd, 4th...).
func ToOrdinal(numStr string) (string, error) {
	f, _, err := big.ParseFloat(numStr, 10, 256, big.ToNearestEven)
	if err != nil {
		return "", fmt.Errorf("invalid number: %s", numStr)
	}

	// Get the integer part
	intVal, _ := f.Int(nil)
	s := intVal.String()

	// Remove negative sign for suffix determination
	abs := s
	negative := false
	if strings.HasPrefix(abs, "-") {
		negative = true
		abs = abs[1:]
	}

	suffix := ordinalSuffix(abs)

	result := s + suffix
	if negative {
		// Already has the minus sign in s
	}
	return result, nil
}

func ordinalSuffix(s string) string {
	// Remove leading zeros
	s = strings.TrimLeft(s, "0")
	if s == "" {
		s = "0"
	}

	n, _ := strconv.Atoi(s)
	if n < 0 {
		n = -n
	}

	// Special cases: 11, 12, 13 use "th"
	if n%100 >= 11 && n%100 <= 13 {
		return "th"
	}

	switch n % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

// ToRoman converts a number to Roman numerals.
func ToRoman(numStr string) (string, error) {
	n, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid integer: %s", numStr)
	}

	if n < 1 || n > 3999 {
		return "", fmt.Errorf("number must be between 1 and 3999 for Roman numerals, got %d", n)
	}

	values := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	symbols := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}

	var result strings.Builder
	remaining := int(n)
	for i, v := range values {
		for remaining >= v {
			result.WriteString(symbols[i])
			remaining -= v
		}
	}

	return result.String(), nil
}

// FromRoman converts Roman numerals to a number.
func FromRoman(roman string) (string, error) {
	roman = strings.ToUpper(strings.TrimSpace(roman))
	if roman == "" {
		return "", fmt.Errorf("empty Roman numeral")
	}

	values := map[byte]int{'I': 1, 'V': 5, 'X': 10, 'L': 50, 'C': 100, 'D': 500, 'M': 1000}

	total := 0
	prev := 0

	for i := len(roman) - 1; i >= 0; i-- {
		v, ok := values[roman[i]]
		if !ok {
			return "", fmt.Errorf("invalid Roman numeral character: %c", roman[i])
		}
		if v < prev {
			total -= v
		} else {
			total += v
		}
		prev = v
	}

	// Verify by converting back
	verify, _ := ToRoman(strconv.Itoa(total))
	if verify != roman {
		return "", fmt.Errorf("invalid Roman numeral: %s", roman)
	}

	return strconv.Itoa(total), nil
}

// BaseConvert converts a number from one base to another.
func BaseConvert(numStr string, fromBase, toBase int) (string, error) {
	if fromBase < 2 || fromBase > 36 {
		return "", fmt.Errorf("from_base must be between 2 and 36, got %d", fromBase)
	}
	if toBase < 2 || toBase > 36 {
		return "", fmt.Errorf("to_base must be between 2 and 36, got %d", toBase)
	}

	numStr = strings.TrimSpace(numStr)
	negative := false
	if strings.HasPrefix(numStr, "-") {
		negative = true
		numStr = numStr[1:]
	}

	n, ok := new(big.Int).SetString(numStr, fromBase)
	if !ok {
		return "", fmt.Errorf("invalid number %q for base %d", numStr, fromBase)
	}

	result := n.Text(toBase)
	if negative {
		result = "-" + result
	}

	return result, nil
}

// ToScientific converts a number to scientific notation.
func ToScientific(numStr string, precision int) (string, error) {
	f, _, err := big.ParseFloat(numStr, 10, 256, big.ToNearestEven)
	if err != nil {
		return "", fmt.Errorf("invalid number: %s", numStr)
	}

	if f.Sign() == 0 {
		return "0e+0", nil
	}

	absF := new(big.Float).Abs(f)
	exp := 0

	// Find the exponent
	ten := big.NewFloat(10)

	for absF.Cmp(big.NewFloat(1)) >= 0 {
		absF.Quo(absF, ten)
		exp++
	}
	for absF.Cmp(big.NewFloat(1)) < 0 && absF.Sign() > 0 {
		absF.Mul(absF, ten)
		exp--
	}

	// Now absF is in [1, 10)
	// Multiply back to get the mantissa
	mantissa := new(big.Float).SetFloat64(0)
	mantissa.Set(f)
	if exp >= 0 {
		divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exp)), nil))
		mantissa.Quo(mantissa, divisor)
	} else {
		multiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-exp)), nil))
		mantissa.Mul(mantissa, multiplier)
	}

	if precision < 0 {
		precision = 6
	}

	mantStr := mantissa.Text('f', precision)

	sign := "+"
	if exp < 0 {
		sign = "-"
		exp = -exp
	}

	return fmt.Sprintf("%se%s%d", mantStr, sign, exp), nil
}

// FromScientific converts scientific notation to a regular number string.
func FromScientific(sciStr string) (string, error) {
	sciStr = strings.TrimSpace(sciStr)
	// Parse the mantissa and exponent
	parts := strings.SplitN(sciStr, "e", 2)
	if len(parts) != 2 {
		parts = strings.SplitN(sciStr, "E", 2)
	}
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid scientific notation: %s (expected format like 1.234e+5)", sciStr)
	}

	mantStr := parts[0]
	expStr := parts[1]

	exp, err := strconv.Atoi(expStr)
	if err != nil {
		return "", fmt.Errorf("invalid exponent: %s", expStr)
	}

	f, _, err := big.ParseFloat(mantStr, 10, 256, big.ToNearestEven)
	if err != nil {
		return "", fmt.Errorf("invalid mantissa: %s", mantStr)
	}

	if exp > 0 {
		multiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exp)), nil))
		f.Mul(f, multiplier)
	} else if exp < 0 {
		divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-exp)), nil))
		f.Quo(f, divisor)
	}

	return f.Text('f', -1), nil
}

// Percentage calculates various percentage operations.
func Percentage(operation string, a, b float64) (string, error) {
	switch operation {
	case "of":
		// What is a% of b?
		result := (a / 100) * b
		return formatFloat(result), nil
	case "is":
		// a is what % of b?
		if b == 0 {
			return "", fmt.Errorf("cannot calculate percentage of zero")
		}
		result := (a / b) * 100
		return formatFloat(result), nil
	case "change":
		// Percentage change from a to b
		if a == 0 {
			return "", fmt.Errorf("cannot calculate percentage change from zero")
		}
		result := ((b - a) / a) * 100
		return formatFloat(result), nil
	default:
		return "", fmt.Errorf("unknown operation: %s (use: of, is, change)", operation)
	}
}

func formatFloat(f float64) string {
	// Remove trailing zeros
	s := strconv.FormatFloat(f, 'f', -1, 64)
	return s
}

// RoundNumber rounds a number using the specified mode.
func RoundNumber(numStr string, decimals int, mode string) (string, error) {
	f, _, err := big.ParseFloat(numStr, 10, 256, big.ToNearestEven)
	if err != nil {
		return "", fmt.Errorf("invalid number: %s", numStr)
	}

	if decimals < 0 {
		decimals = 0
	}

	// For "nearest" mode, use big.Float's built-in formatting which does banker's rounding
	if mode == "nearest" {
		return f.Text('f', decimals), nil
	}

	// For other modes, we need manual rounding
	multiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	scaled := new(big.Float).Mul(f, multiplier)

	// Get the full string representation (no rounding)
	fullStr := scaled.Text('f', -1)

	// Split to get integer and fractional parts (truncated toward zero)
	intStr := fullStr
	fracStr := ""
	if idx := strings.Index(fullStr, "."); idx >= 0 {
		intStr = fullStr[:idx]
		fracStr = fullStr[idx+1:]
	}

	hasFrac := fracStr != "" && strings.TrimLeft(fracStr, "0") != ""

	// Parse the truncated integer
	intVal, ok := new(big.Int).SetString(intStr, 10)
	if !ok {
		intVal = big.NewInt(0)
	}

	negative := f.Sign() < 0

	switch mode {
	case "up":
		// Round toward positive infinity
		if hasFrac && !negative {
			intVal.Add(intVal, big.NewInt(1))
		}
	case "down":
		// Round toward negative infinity
		if hasFrac && negative {
			intVal.Sub(intVal, big.NewInt(1))
		}
	case "toward_zero":
		// Truncate toward zero - intVal already truncated
	case "away_zero":
		// Round away from zero
		if hasFrac {
			if negative {
				intVal.Sub(intVal, big.NewInt(1))
			} else {
				intVal.Add(intVal, big.NewInt(1))
			}
		}
	default:
		return "", fmt.Errorf("unknown rounding mode: %s (use: nearest, up, down, toward_zero, away_zero)", mode)
	}

	result := new(big.Float).SetInt(intVal)
	result.Quo(result, multiplier)

	return result.Text('f', decimals), nil
}

// SpellToNumber converts English number words to digits.
func SpellToNumber(words string) (string, error) {
	words = strings.ToLower(strings.TrimSpace(words))
	if words == "" {
		return "", fmt.Errorf("empty input")
	}

	// Handle "negative"
	negative := false
	if strings.HasPrefix(words, "negative ") {
		negative = true
		words = words[9:]
	} else if strings.HasPrefix(words, "minus ") {
		negative = true
		words = words[6:]
	}

	// Handle "point" for fractional part
	var fracWords []string
	hasFractional := false
	if idx := strings.Index(words, " point "); idx >= 0 {
		fracPart := words[idx+7:]
		words = words[:idx]
		hasFractional = true
		fracWords = strings.Fields(fracPart)
	}

	words = strings.TrimSpace(words)
	if words == "" {
		return "", fmt.Errorf("no number words found")
	}

	// Tokenize
	tokens := strings.Fields(words)

	// Build value maps
	onesMap := map[string]int64{
		"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4,
		"five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9,
	}
	teensMap := map[string]int64{
		"ten": 10, "eleven": 11, "twelve": 12, "thirteen": 13, "fourteen": 14,
		"fifteen": 15, "sixteen": 16, "seventeen": 17, "eighteen": 18, "nineteen": 19,
	}
	tensMap := map[string]int64{
		"twenty": 20, "thirty": 30, "forty": 40, "fifty": 50,
		"sixty": 60, "seventy": 70, "eighty": 80, "ninety": 90,
	}
	scalesMap := map[string]int64{
		"hundred": 100, "thousand": 1000, "million": 1000000,
		"billion": 1000000000, "trillion": 1000000000000,
		"quadrillion": 1000000000000000, "quintillion": 1000000000000000000,
	}

	var total int64 = 0
	var current int64 = 0

	for _, token := range tokens {
		// Handle hyphenated tens (e.g. "twenty-one")
		if strings.Contains(token, "-") {
			parts := strings.SplitN(token, "-", 2)
			if t, ok := tensMap[parts[0]]; ok {
				current += t
				if len(parts) > 1 {
					if o, ok := onesMap[parts[1]]; ok {
						current += o
					} else {
						return "", fmt.Errorf("unknown word: %s", parts[1])
					}
				}
				continue
			}
		}

		// Handle "and"
		if token == "and" {
			continue
		}

		if v, ok := onesMap[token]; ok {
			current += v
		} else if v, ok := teensMap[token]; ok {
			current += v
		} else if v, ok := tensMap[token]; ok {
			current += v
		} else if token == "hundred" {
			if current == 0 {
				current = 1
			}
			current *= 100
		} else if v, ok := scalesMap[token]; ok {
			if current == 0 {
				current = 1
			}
			if v >= 1000 {
				total += current * v
				current = 0
			} else {
				current *= v
			}
		} else {
			return "", fmt.Errorf("unknown word: %s", token)
		}
	}

	total += current

	result := strconv.FormatInt(total, 10)
	if negative {
		result = "-" + result
	}

	if hasFractional {
		var fracStr strings.Builder
		for _, w := range fracWords {
			if v, ok := onesMap[w]; ok {
				fracStr.WriteString(strconv.FormatInt(v, 10))
			} else {
				return "", fmt.Errorf("unknown fractional digit word: %s", w)
			}
		}
		if fracStr.Len() > 0 {
			result += "." + fracStr.String()
		}
	}

	return result, nil
}

// IsPrime checks if a number is prime.
func IsPrime(numStr string) (bool, error) {
	n, ok := new(big.Int).SetString(strings.TrimSpace(numStr), 10)
	if !ok {
		return false, fmt.Errorf("invalid number: %s", numStr)
	}

	if n.Sign() < 0 {
		n = new(big.Int).Abs(n)
	}

	if n.Cmp(big.NewInt(1)) <= 0 {
		return false, nil
	}

	return n.ProbablyPrime(20), nil
}

// Factorize returns the prime factorization of a number.
func Factorize(numStr string) (string, error) {
	n, ok := new(big.Int).SetString(strings.TrimSpace(numStr), 10)
	if !ok {
		return "", fmt.Errorf("invalid number: %s", numStr)
	}

	if n.Sign() < 0 {
		return "", fmt.Errorf("cannot factorize negative number")
	}

	if n.Cmp(big.NewInt(1)) <= 0 {
		if n.Sign() == 0 {
			return "0", nil
		}
		return "1", nil
	}

	var factors []string
	temp := new(big.Int).Set(n)
	divisor := big.NewInt(2)

	for {
		// Check if divisor^2 > temp (remaining)
		square := new(big.Int).Mul(divisor, divisor)
		if square.Cmp(temp) > 0 {
			break
		}
		count := 0
		mod := new(big.Int)
		for mod.Mod(temp, divisor).Sign() == 0 {
			temp.Quo(temp, divisor)
			count++
		}
		if count > 0 {
			if count == 1 {
				factors = append(factors, divisor.String())
			} else {
				factors = append(factors, fmt.Sprintf("%s^%d", divisor.String(), count))
			}
		}
		divisor.Add(divisor, big.NewInt(1))
	}

	if temp.Cmp(big.NewInt(1)) > 0 {
		factors = append(factors, temp.String())
	}

	return strings.Join(factors, " × "), nil
}

// GCD computes the greatest common divisor of two numbers.
func GCD(aStr, bStr string) (string, error) {
	a, ok := new(big.Int).SetString(strings.TrimSpace(aStr), 10)
	if !ok {
		return "", fmt.Errorf("invalid number: %s", aStr)
	}
	b, ok := new(big.Int).SetString(strings.TrimSpace(bStr), 10)
	if !ok {
		return "", fmt.Errorf("invalid number: %s", bStr)
	}

	if a.Sign() < 0 {
		a = new(big.Int).Abs(a)
	}
	if b.Sign() < 0 {
		b = new(big.Int).Abs(b)
	}

	return new(big.Int).GCD(nil, nil, a, b).String(), nil
}

// LCM computes the least common multiple of two numbers.
func LCM(aStr, bStr string) (string, error) {
	a, ok := new(big.Int).SetString(strings.TrimSpace(aStr), 10)
	if !ok {
		return "", fmt.Errorf("invalid number: %s", aStr)
	}
	b, ok := new(big.Int).SetString(strings.TrimSpace(bStr), 10)
	if !ok {
		return "", fmt.Errorf("invalid number: %s", bStr)
	}

	if a.Sign() == 0 || b.Sign() == 0 {
		return "0", nil
	}

	if a.Sign() < 0 {
		a = new(big.Int).Abs(a)
	}
	if b.Sign() < 0 {
		b = new(big.Int).Abs(b)
	}

	gcd := new(big.Int).GCD(nil, nil, a, b)
	product := new(big.Int).Mul(a, b)
	result := new(big.Int).Quo(product, gcd)

	return result.String(), nil
}

// Fibonacci generates the first n Fibonacci numbers.
func Fibonacci(n int) ([]string, error) {
	if n < 0 {
		return nil, fmt.Errorf("count must be non-negative, got %d", n)
	}
	if n == 0 {
		return []string{}, nil
	}
	if n > 1000 {
		return nil, fmt.Errorf("count too large (max 1000), got %d", n)
	}

	result := make([]string, 0, n)
	a := big.NewInt(0)
	b := big.NewInt(1)

	for i := 0; i < n; i++ {
		result = append(result, a.String())
		a, b = b, new(big.Int).Add(a, b)
	}

	return result, nil
}

// IsEven checks if a number is even.
func IsEven(numStr string) (bool, error) {
	n, ok := new(big.Int).SetString(strings.TrimSpace(numStr), 10)
	if !ok {
		return false, fmt.Errorf("invalid number: %s", numStr)
	}
	return new(big.Int).Mod(n, big.NewInt(2)).Sign() == 0, nil
}

// IsOdd checks if a number is odd.
func IsOdd(numStr string) (bool, error) {
	n, ok := new(big.Int).SetString(strings.TrimSpace(numStr), 10)
	if !ok {
		return false, fmt.Errorf("invalid number: %s", numStr)
	}
	return new(big.Int).Mod(n, big.NewInt(2)).Sign() != 0, nil
}

// Divmod returns the quotient and remainder of a divided by b.
func Divmod(aStr, bStr string) (string, string, error) {
	a, ok := new(big.Int).SetString(strings.TrimSpace(aStr), 10)
	if !ok {
		return "", "", fmt.Errorf("invalid number: %s", aStr)
	}
	b, ok := new(big.Int).SetString(strings.TrimSpace(bStr), 10)
	if !ok {
		return "", "", fmt.Errorf("invalid number: %s", bStr)
	}

	if b.Sign() == 0 {
		return "", "", fmt.Errorf("division by zero")
	}

	q := new(big.Int)
	r := new(big.Int)
	q.QuoRem(a, b, r)
	return q.String(), r.String(), nil
}

// Abs returns the absolute value of a number.
func Abs(numStr string) (string, error) {
	n, ok := new(big.Int).SetString(strings.TrimSpace(numStr), 10)
	if !ok {
		// Try as float
		f, _, err := big.ParseFloat(numStr, 10, 256, big.ToNearestEven)
		if err != nil {
			return "", fmt.Errorf("invalid number: %s", numStr)
		}
		return new(big.Float).Abs(f).Text('f', -1), nil
	}
	return new(big.Int).Abs(n).String(), nil
}

// Negate returns the negation of a number.
func Negate(numStr string) (string, error) {
	n, ok := new(big.Int).SetString(strings.TrimSpace(numStr), 10)
	if !ok {
		f, _, err := big.ParseFloat(numStr, 10, 256, big.ToNearestEven)
		if err != nil {
			return "", fmt.Errorf("invalid number: %s", numStr)
		}
		return new(big.Float).Neg(f).Text('f', -1), nil
	}
	return new(big.Int).Neg(n).String(), nil
}

// Parity returns "even" or "odd" for a number.
func Parity(numStr string) (string, error) {
	isEven, err := IsEven(numStr)
	if err != nil {
		return "", err
	}
	if isEven {
		return "even", nil
	}
	return "odd", nil
}

// DigitSum returns the sum of digits of a number.
func DigitSum(numStr string) (string, error) {
	s := strings.TrimSpace(numStr)
	// Remove sign and decimal point
	s = strings.TrimLeft(s, "-+")
	s = strings.ReplaceAll(s, ".", "")

	sum := big.NewInt(0)
	for _, c := range s {
		if c >= '0' && c <= '9' {
			sum.Add(sum, big.NewInt(int64(c-'0')))
		}
	}
	return sum.String(), nil
}

// DigitCount returns the number of digits in a number.
func DigitCount(numStr string) (string, error) {
	s := strings.TrimSpace(numStr)
	// Remove sign
	s = strings.TrimLeft(s, "-+")
	// Remove decimal point and everything after
	if idx := strings.Index(s, "."); idx >= 0 {
		s = s[:idx]
	}
	// Remove leading zeros
	s = strings.TrimLeft(s, "0")
	if s == "" {
		return "1", nil // zero has 1 digit
	}
	return strconv.Itoa(len(s)), nil
}

// ReverseDigits reverses the digits of a number.
func ReverseDigits(numStr string) (string, error) {
	s := strings.TrimSpace(numStr)
	negative := false
	if strings.HasPrefix(s, "-") {
		negative = true
		s = s[1:]
	}

	// Handle decimal
	intPart := s
	fracPart := ""
	if idx := strings.Index(s, "."); idx >= 0 {
		intPart = s[:idx]
		fracPart = s[idx+1:]
	}

	// Reverse integer part
	intRunes := []rune(intPart)
	for i, j := 0, len(intRunes)-1; i < j; i, j = i+1, j-1 {
		intRunes[i], intRunes[j] = intRunes[j], intRunes[i]
	}
	result := string(intRunes)

	if fracPart != "" {
		fracRunes := []rune(fracPart)
		for i, j := 0, len(fracRunes)-1; i < j; i, j = i+1, j-1 {
			fracRunes[i], fracRunes[j] = fracRunes[j], fracRunes[i]
		}
		result += "." + string(fracRunes)
	}

	if negative {
		result = "-" + result
	}

	return result, nil
}

// Collatz generates the Collatz sequence starting from n.
func Collatz(numStr string) ([]string, error) {
	n, ok := new(big.Int).SetString(strings.TrimSpace(numStr), 10)
	if !ok {
		return nil, fmt.Errorf("invalid number: %s", numStr)
	}

	if n.Sign() <= 0 {
		return nil, fmt.Errorf("number must be positive, got %s", numStr)
	}

	var result []string
	temp := new(big.Int).Set(n)
	maxSteps := 10000

	for temp.Cmp(big.NewInt(1)) != 0 && len(result) < maxSteps {
		result = append(result, temp.String())
		if new(big.Int).Mod(temp, big.NewInt(2)).Sign() == 0 {
			temp.Quo(temp, big.NewInt(2))
		} else {
			temp.Mul(temp, big.NewInt(3))
			temp.Add(temp, big.NewInt(1))
		}
	}
	result = append(result, "1")

	return result, nil
}

// IsPerfectSquare checks if a number is a perfect square.
func IsPerfectSquare(numStr string) (bool, error) {
	n, ok := new(big.Int).SetString(strings.TrimSpace(numStr), 10)
	if !ok {
		return false, fmt.Errorf("invalid number: %s", numStr)
	}

	if n.Sign() < 0 {
		return false, nil
	}

	sqrt := new(big.Int).Sqrt(n)
	square := new(big.Int).Mul(sqrt, sqrt)
	return square.Cmp(n) == 0, nil
}

// Sqrt returns the integer square root of a number.
func Sqrt(numStr string) (string, error) {
	n, ok := new(big.Int).SetString(strings.TrimSpace(numStr), 10)
	if !ok {
		return "", fmt.Errorf("invalid number: %s", numStr)
	}

	if n.Sign() < 0 {
		return "", fmt.Errorf("cannot compute square root of negative number")
	}

	return new(big.Int).Sqrt(n).String(), nil
}

// Power computes base^exponent for integers.
func Power(baseStr, expStr string) (string, error) {
	base, ok := new(big.Int).SetString(strings.TrimSpace(baseStr), 10)
	if !ok {
		return "", fmt.Errorf("invalid base: %s", baseStr)
	}
	exp, ok := new(big.Int).SetString(strings.TrimSpace(expStr), 10)
	if !ok {
		return "", fmt.Errorf("invalid exponent: %s", expStr)
	}

	if exp.Sign() < 0 {
		return "", fmt.Errorf("exponent must be non-negative for integer power")
	}

	if exp.Cmp(big.NewInt(10000)) > 0 {
		return "", fmt.Errorf("exponent too large (max 10000)")
	}

	return new(big.Int).Exp(base, exp, nil).String(), nil
}

// Factorial computes n! for non-negative integers.
func Factorial(numStr string) (string, error) {
	n, ok := new(big.Int).SetString(strings.TrimSpace(numStr), 10)
	if !ok {
		return "", fmt.Errorf("invalid number: %s", numStr)
	}

	if n.Sign() < 0 {
		return "", fmt.Errorf("factorial is not defined for negative numbers")
	}

	if n.Cmp(big.NewInt(10000)) > 0 {
		return "", fmt.Errorf("number too large for factorial (max 10000)")
	}

	result := big.NewInt(1)
	one := big.NewInt(1)
	i := big.NewInt(2)

	for i.Cmp(n) <= 0 {
		result.Mul(result, i)
		i.Add(i, one)
	}

	return result.String(), nil
}

// IsArmstrong checks if a number is an Armstrong number.
func IsArmstrong(numStr string) (bool, error) {
	s := strings.TrimSpace(numStr)
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return false, fmt.Errorf("invalid number: %s", numStr)
	}

	if n.Sign() < 0 {
		return false, nil
	}

	digits := s
	numDigits := len(digits)
	if numDigits == 0 {
		return false, nil
	}

	sum := big.NewInt(0)
	for _, c := range digits {
		d := int64(c - '0')
		// d^numDigits
		power := big.NewInt(d)
		for j := 1; j < numDigits; j++ {
			power.Mul(power, big.NewInt(d))
		}
		sum.Add(sum, power)
	}

	return sum.Cmp(n) == 0, nil
}

// ToBinary converts a number to its binary representation.
func ToBinary(numStr string) (string, error) {
	return BaseConvert(numStr, 10, 2)
}

// ToHex converts a number to hexadecimal.
func ToHex(numStr string) (string, error) {
	return BaseConvert(numStr, 10, 16)
}

// ToOctal converts a number to octal.
func ToOctal(numStr string) (string, error) {
	return BaseConvert(numStr, 10, 8)
}

// FromBinary converts a binary number to decimal.
func FromBinary(numStr string) (string, error) {
	return BaseConvert(numStr, 2, 10)
}

// FromHex converts a hexadecimal number to decimal.
func FromHex(numStr string) (string, error) {
	return BaseConvert(numStr, 16, 10)
}

// FromOctal converts an octal number to decimal.
func FromOctal(numStr string) (string, error) {
	return BaseConvert(numStr, 8, 10)
}

// CommaFormat adds thousands separators to a number.
func CommaFormat(numStr string) (string, error) {
	return Format(numStr, FormatOptions{
		GroupDigits: true,
		Separator:   ",",
	})
}

// Uncomma removes thousands separators from a number.
func Uncomma(numStr string) (string, error) {
	result := strings.ReplaceAll(numStr, ",", "")
	// Verify it's a valid number
	_, _, err := big.ParseFloat(result, 10, 256, big.ToNearestEven)
	if err != nil {
		return "", fmt.Errorf("invalid number after removing commas: %s", result)
	}
	return result, nil
}

// PercentToFraction converts a percentage to a fraction string.
func PercentToFraction(percentStr string) (string, error) {
	f, _, err := big.ParseFloat(strings.TrimSuffix(strings.TrimSpace(percentStr), "%"), 10, 256, big.ToNearestEven)
	if err != nil {
		return "", fmt.Errorf("invalid percentage: %s", percentStr)
	}

	// Divide by 100
	f.Quo(f, big.NewFloat(100))
	return f.Text('f', -1), nil
}

// ToFraction converts a decimal to a simplified fraction.
func ToFraction(numStr string) (string, error) {
	s := strings.TrimSpace(numStr)
	negative := false
	if strings.HasPrefix(s, "-") {
		negative = true
		s = s[1:]
	}

	parts := strings.SplitN(s, ".", 2)
	if len(parts) != 2 {
		// Integer - return as n/1
		n, ok := new(big.Int).SetString(s, 10)
		if !ok {
			return "", fmt.Errorf("invalid number: %s", numStr)
		}
		if negative {
			n = new(big.Int).Neg(n)
		}
		return n.String() + "/1", nil
	}

	intPart := parts[0]
	decPart := parts[1]

	if decPart == "" {
		decPart = "0"
	}

	// numerator = intPart * 10^len(decPart) + decPart
	// denominator = 10^len(decPart)
	denomExp := len(decPart)
	numerator := new(big.Int)
	if intPart != "" {
		numerator.SetString(intPart, 10)
	}
	decNum := new(big.Int)
	decNum.SetString(decPart, 10)
	denom := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(denomExp)), nil)
	numerator.Mul(numerator, denom)
	numerator.Add(numerator, decNum)

	if negative {
		numerator = new(big.Int).Neg(numerator)
	}

	// Simplify
	gcd := new(big.Int).GCD(nil, nil, numerator, denom)
	numerator.Quo(numerator, gcd)
	denom.Quo(denom, gcd)

	return numerator.String() + "/" + denom.String(), nil
}

// FormatFloat formats a float64 for display, removing trailing zeros.
func FormatFloatForDisplay(f float64) string {
	if math.IsInf(f, 0) {
		if f > 0 {
			return "infinity"
		}
		return "-infinity"
	}
	if math.IsNaN(f) {
		return "NaN"
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
