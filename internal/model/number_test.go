package model

import (
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		num    string
		opts   FormatOptions
		expect string
	}{
		{"1234567", FormatOptions{GroupDigits: true, Separator: ","}, "1,234,567"},
		{"1234567.89", FormatOptions{GroupDigits: true, Separator: ",", Decimals: 2}, "1,234,567.89"},
		{"-1234567", FormatOptions{GroupDigits: true, Separator: ","}, "-1,234,567"},
		{"1000000", FormatOptions{GroupDigits: true, Separator: ".", DecimalSep: ","}, "1.000.000"},
		{"42", FormatOptions{GroupDigits: true, Separator: ","}, "42"},
		{"1234567", FormatOptions{GroupDigits: false}, "1234567"},
		{"1234.5678", FormatOptions{GroupDigits: true, Separator: ",", Decimals: 2}, "1,234.57"},
	}
	for _, tt := range tests {
		got, err := Format(tt.num, tt.opts)
		if err != nil {
			t.Errorf("Format(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("Format(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestToWords(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"0", "zero"},
		{"1", "one"},
		{"10", "ten"},
		{"15", "fifteen"},
		{"21", "twenty-one"},
		{"100", "one hundred"},
		{"101", "one hundred and one"},
		{"1000", "one thousand"},
		{"12345", "twelve thousand three hundred and forty-five"},
		{"1000000", "one million"},
		{"-5", "negative five"},
		{"123.45", "one hundred and twenty-three point four five"},
	}
	for _, tt := range tests {
		got, err := ToWords(tt.num)
		if err != nil {
			t.Errorf("ToWords(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("ToWords(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestToOrdinal(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"1", "1st"},
		{"2", "2nd"},
		{"3", "3rd"},
		{"4", "4th"},
		{"11", "11th"},
		{"12", "12th"},
		{"13", "13th"},
		{"21", "21st"},
		{"22", "22nd"},
		{"23", "23rd"},
		{"100", "100th"},
		{"101", "101st"},
		{"111", "111th"},
		{"112", "112th"},
		{"113", "113th"},
	}
	for _, tt := range tests {
		got, err := ToOrdinal(tt.num)
		if err != nil {
			t.Errorf("ToOrdinal(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("ToOrdinal(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestToRoman(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"1", "I"},
		{"4", "IV"},
		{"9", "IX"},
		{"14", "XIV"},
		{"40", "XL"},
		{"90", "XC"},
		{"400", "CD"},
		{"900", "CM"},
		{"2024", "MMXXIV"},
		{"3999", "MMMCMXCIX"},
	}
	for _, tt := range tests {
		got, err := ToRoman(tt.num)
		if err != nil {
			t.Errorf("ToRoman(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("ToRoman(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestFromRoman(t *testing.T) {
	tests := []struct {
		roman  string
		expect string
	}{
		{"I", "1"},
		{"IV", "4"},
		{"IX", "9"},
		{"XIV", "14"},
		{"XL", "40"},
		{"XC", "90"},
		{"CD", "400"},
		{"CM", "900"},
		{"MMXXIV", "2024"},
		{"MMMCMXCIX", "3999"},
	}
	for _, tt := range tests {
		got, err := FromRoman(tt.roman)
		if err != nil {
			t.Errorf("FromRoman(%s) error: %v", tt.roman, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("FromRoman(%s) = %q, want %q", tt.roman, got, tt.expect)
		}
	}
}

func TestBaseConvert(t *testing.T) {
	tests := []struct {
		num      string
		from     int
		to       int
		expect   string
	}{
		{"255", 10, 16, "ff"},
		{"255", 10, 2, "11111111"},
		{"255", 10, 8, "377"},
		{"ff", 16, 10, "255"},
		{"11111111", 2, 10, "255"},
		{"377", 8, 10, "255"},
		{"z", 36, 10, "35"},
		{"10", 10, 36, "a"},
	}
	for _, tt := range tests {
		got, err := BaseConvert(tt.num, tt.from, tt.to)
		if err != nil {
			t.Errorf("BaseConvert(%s, %d, %d) error: %v", tt.num, tt.from, tt.to, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("BaseConvert(%s, %d, %d) = %q, want %q", tt.num, tt.from, tt.to, got, tt.expect)
		}
	}
}

func TestIsPrime(t *testing.T) {
	primes := []string{"2", "3", "5", "7", "11", "13", "17", "19", "23", "29", "31", "97", "101"}
	notPrimes := []string{"0", "1", "4", "6", "8", "9", "10", "15", "21", "25", "100"}

	for _, n := range primes {
		result, err := IsPrime(n)
		if err != nil {
			t.Errorf("IsPrime(%s) error: %v", n, err)
			continue
		}
		if !result {
			t.Errorf("IsPrime(%s) = false, want true", n)
		}
	}

	for _, n := range notPrimes {
		result, err := IsPrime(n)
		if err != nil {
			t.Errorf("IsPrime(%s) error: %v", n, err)
			continue
		}
		if result {
			t.Errorf("IsPrime(%s) = true, want false", n)
		}
	}
}

func TestFactorize(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"2", "2"},
		{"12", "2^2 × 3"},
		{"360", "2^3 × 3^2 × 5"},
		{"100", "2^2 × 5^2"},
		{"17", "17"},
		{"1", "1"},
	}
	for _, tt := range tests {
		got, err := Factorize(tt.num)
		if err != nil {
			t.Errorf("Factorize(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("Factorize(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		a, b   string
		expect string
	}{
		{"48", "18", "6"},
		{"17", "5", "1"},
		{"100", "75", "25"},
		{"0", "5", "5"},
		{"7", "7", "7"},
	}
	for _, tt := range tests {
		got, err := GCD(tt.a, tt.b)
		if err != nil {
			t.Errorf("GCD(%s, %s) error: %v", tt.a, tt.b, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("GCD(%s, %s) = %q, want %q", tt.a, tt.b, got, tt.expect)
		}
	}
}

func TestLCM(t *testing.T) {
	tests := []struct {
		a, b   string
		expect string
	}{
		{"4", "6", "12"},
		{"3", "5", "15"},
		{"12", "18", "36"},
		{"7", "7", "7"},
		{"0", "5", "0"},
	}
	for _, tt := range tests {
		got, err := LCM(tt.a, tt.b)
		if err != nil {
			t.Errorf("LCM(%s, %s) error: %v", tt.a, tt.b, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("LCM(%s, %s) = %q, want %q", tt.a, tt.b, got, tt.expect)
		}
	}
}

func TestFibonacci(t *testing.T) {
	result, err := Fibonacci(10)
	if err != nil {
		t.Fatalf("Fibonacci(10) error: %v", err)
	}
	expected := []string{"0", "1", "1", "2", "3", "5", "8", "13", "21", "34"}
	if len(result) != len(expected) {
		t.Fatalf("Fibonacci(10) returned %d items, want %d", len(result), len(expected))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Fibonacci(10)[%d] = %q, want %q", i, result[i], v)
		}
	}
}

func TestIsEven(t *testing.T) {
	evens := []string{"0", "2", "4", "10", "100", "-4"}
	odds := []string{"1", "3", "7", "11", "-5"}
	for _, n := range evens {
		result, _ := IsEven(n)
		if !result {
			t.Errorf("IsEven(%s) = false, want true", n)
		}
	}
	for _, n := range odds {
		result, _ := IsEven(n)
		if result {
			t.Errorf("IsEven(%s) = true, want false", n)
		}
	}
}

func TestIsOdd(t *testing.T) {
	odds := []string{"1", "3", "7", "11", "-5"}
	evens := []string{"0", "2", "4", "10", "100", "-4"}
	for _, n := range odds {
		result, _ := IsOdd(n)
		if !result {
			t.Errorf("IsOdd(%s) = false, want true", n)
		}
	}
	for _, n := range evens {
		result, _ := IsOdd(n)
		if result {
			t.Errorf("IsOdd(%s) = true, want false", n)
		}
	}
}

func TestDivmod(t *testing.T) {
	tests := []struct {
		a, b, eq, er string
	}{
		{"17", "5", "3", "2"},
		{"10", "3", "3", "1"},
		{"100", "7", "14", "2"},
		{"0", "5", "0", "0"},
	}
	for _, tt := range tests {
		q, r, err := Divmod(tt.a, tt.b)
		if err != nil {
			t.Errorf("Divmod(%s, %s) error: %v", tt.a, tt.b, err)
			continue
		}
		if q != tt.eq {
			t.Errorf("Divmod(%s, %s) quotient = %q, want %q", tt.a, tt.b, q, tt.eq)
		}
		if r != tt.er {
			t.Errorf("Divmod(%s, %s) remainder = %q, want %q", tt.a, tt.b, r, tt.er)
		}
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"-5", "5"},
		{"5", "5"},
		{"0", "0"},
		{"-42", "42"},
		{"42", "42"},
	}
	for _, tt := range tests {
		got, err := Abs(tt.num)
		if err != nil {
			t.Errorf("Abs(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("Abs(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestDigitSum(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"12345", "15"},
		{"0", "0"},
		{"999", "27"},
		{"100", "1"},
	}
	for _, tt := range tests {
		got, err := DigitSum(tt.num)
		if err != nil {
			t.Errorf("DigitSum(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("DigitSum(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestDigitCount(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"0", "1"},
		{"5", "1"},
		{"42", "2"},
		{"12345", "5"},
		{"1000000", "7"},
	}
	for _, tt := range tests {
		got, err := DigitCount(tt.num)
		if err != nil {
			t.Errorf("DigitCount(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("DigitCount(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestReverseDigits(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"12345", "54321"},
		{"100", "001"},
		{"-123", "-321"},
	}
	for _, tt := range tests {
		got, err := ReverseDigits(tt.num)
		if err != nil {
			t.Errorf("ReverseDigits(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("ReverseDigits(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestCollatz(t *testing.T) {
	result, err := Collatz("6")
	if err != nil {
		t.Fatalf("Collatz(6) error: %v", err)
	}
	expected := []string{"6", "3", "10", "5", "16", "8", "4", "2", "1"}
	if len(result) != len(expected) {
		t.Fatalf("Collatz(6) returned %d items, want %d", len(result), len(expected))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Collatz(6)[%d] = %q, want %q", i, result[i], v)
		}
	}
}

func TestIsPerfectSquare(t *testing.T) {
	squares := []string{"0", "1", "4", "9", "16", "25", "100", "144"}
	notSquares := []string{"2", "3", "5", "7", "10", "15", "99", "145"}
	for _, n := range squares {
		result, _ := IsPerfectSquare(n)
		if !result {
			t.Errorf("IsPerfectSquare(%s) = false, want true", n)
		}
	}
	for _, n := range notSquares {
		result, _ := IsPerfectSquare(n)
		if result {
			t.Errorf("IsPerfectSquare(%s) = true, want false", n)
		}
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"0", "0"},
		{"1", "1"},
		{"4", "2"},
		{"9", "3"},
		{"10", "3"},
		{"144", "12"},
		{"10000", "100"},
	}
	for _, tt := range tests {
		got, err := Sqrt(tt.num)
		if err != nil {
			t.Errorf("Sqrt(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("Sqrt(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestPower(t *testing.T) {
	tests := []struct {
		base, exp, expect string
	}{
		{"2", "10", "1024"},
		{"2", "0", "1"},
		{"10", "3", "1000"},
		{"5", "2", "25"},
		{"1", "100", "1"},
		{"0", "5", "0"},
	}
	for _, tt := range tests {
		got, err := Power(tt.base, tt.exp)
		if err != nil {
			t.Errorf("Power(%s, %s) error: %v", tt.base, tt.exp, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("Power(%s, %s) = %q, want %q", tt.base, tt.exp, got, tt.expect)
		}
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"0", "1"},
		{"1", "1"},
		{"5", "120"},
		{"10", "3628800"},
		{"20", "2432902008176640000"},
	}
	for _, tt := range tests {
		got, err := Factorial(tt.num)
		if err != nil {
			t.Errorf("Factorial(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("Factorial(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestIsArmstrong(t *testing.T) {
	armstrong := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "153", "370", "371", "407", "1634"}
	notArmstrong := []string{"10", "12", "100", "154", "200", "1000"}
	for _, n := range armstrong {
		result, _ := IsArmstrong(n)
		if !result {
			t.Errorf("IsArmstrong(%s) = false, want true", n)
		}
	}
	for _, n := range notArmstrong {
		result, _ := IsArmstrong(n)
		if result {
			t.Errorf("IsArmstrong(%s) = true, want false", n)
		}
	}
}

func TestSpellToNumber(t *testing.T) {
	tests := []struct {
		words  string
		expect string
	}{
		{"zero", "0"},
		{"one", "1"},
		{"ten", "10"},
		{"twenty-one", "21"},
		{"one hundred", "100"},
		{"one hundred and one", "101"},
		{"one thousand", "1000"},
		{"twelve thousand three hundred and forty-five", "12345"},
		{"negative five", "-5"},
		{"one hundred twenty-three point four five", "123.45"},
	}
	for _, tt := range tests {
		got, err := SpellToNumber(tt.words)
		if err != nil {
			t.Errorf("SpellToNumber(%s) error: %v", tt.words, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("SpellToNumber(%s) = %q, want %q", tt.words, got, tt.expect)
		}
	}
}

func TestToFraction(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"0.5", "1/2"},
		{"0.25", "1/4"},
		{"0.75", "3/4"},
		{"0.375", "3/8"},
		{"1.5", "3/2"},
		{"3", "3/1"},
		{"0.1", "1/10"},
	}
	for _, tt := range tests {
		got, err := ToFraction(tt.num)
		if err != nil {
			t.Errorf("ToFraction(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("ToFraction(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestCommaFormat(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"1234567", "1,234,567"},
		{"1000000", "1,000,000"},
		{"42", "42"},
		{"0", "0"},
	}
	for _, tt := range tests {
		got, err := CommaFormat(tt.num)
		if err != nil {
			t.Errorf("CommaFormat(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("CommaFormat(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestUncomma(t *testing.T) {
	tests := []struct {
		num    string
		expect string
	}{
		{"1,234,567", "1234567"},
		{"1,000,000", "1000000"},
		{"42", "42"},
	}
	for _, tt := range tests {
		got, err := Uncomma(tt.num)
		if err != nil {
			t.Errorf("Uncomma(%s) error: %v", tt.num, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("Uncomma(%s) = %q, want %q", tt.num, got, tt.expect)
		}
	}
}

func TestPercentage(t *testing.T) {
	// 15% of 200 = 30
	result, err := Percentage("of", 15, 200)
	if err != nil {
		t.Errorf("Percentage(of, 15, 200) error: %v", err)
	} else if result != "30" {
		t.Errorf("Percentage(of, 15, 200) = %q, want %q", result, "30")
	}

	// 30 is what % of 150 = 20
	result, err = Percentage("is", 30, 150)
	if err != nil {
		t.Errorf("Percentage(is, 30, 150) error: %v", err)
	} else if result != "20" {
		t.Errorf("Percentage(is, 30, 150) = %q, want %q", result, "20")
	}

	// % change from 100 to 150 = 50
	result, err = Percentage("change", 100, 150)
	if err != nil {
		t.Errorf("Percentage(change, 100, 150) error: %v", err)
	} else if result != "50" {
		t.Errorf("Percentage(change, 100, 150) = %q, want %q", result, "50")
	}
}

func TestRoundNumber(t *testing.T) {
	tests := []struct {
		num      string
		decimals int
		mode     string
		expect   string
	}{
		{"3.14159", 2, "nearest", "3.14"},
		{"3.14159", 3, "nearest", "3.142"},
		{"3.5", 0, "nearest", "4"},
		{"2.5", 0, "nearest", "2"}, // banker's rounding
		{"3.1", 0, "up", "4"},
		{"3.9", 0, "down", "3"},
		{"-3.5", 0, "up", "-3"},
		{"-3.5", 0, "down", "-4"},
	}
	for _, tt := range tests {
		got, err := RoundNumber(tt.num, tt.decimals, tt.mode)
		if err != nil {
			t.Errorf("RoundNumber(%s, %d, %s) error: %v", tt.num, tt.decimals, tt.mode, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("RoundNumber(%s, %d, %s) = %q, want %q", tt.num, tt.decimals, tt.mode, got, tt.expect)
		}
	}
}

func TestToScientific(t *testing.T) {
	result, err := ToScientific("1234567", 6)
	if err != nil {
		t.Fatalf("ToScientific error: %v", err)
	}
	if !strings.Contains(result, "1.234567") || !strings.Contains(result, "e+6") {
		t.Errorf("ToScientific(1234567) = %q, expected 1.234567e+6", result)
	}
}

func TestFromScientific(t *testing.T) {
	result, err := FromScientific("1.5e+3")
	if err != nil {
		t.Fatalf("FromScientific error: %v", err)
	}
	if result != "1500" {
		t.Errorf("FromScientific(1.5e+3) = %q, want %q", result, "1500")
	}
}
