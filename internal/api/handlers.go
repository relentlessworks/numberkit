package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/relentlessworks/numberkit/internal/model"
)

type Handler struct {
	noAuth bool
}

func NewHandler(noAuth bool) *Handler {
	return &Handler{noAuth: noAuth}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/help", h.help)
	mux.HandleFunc("/.well-known/agent.md", h.help)
	mux.HandleFunc("/format", h.format)
	mux.HandleFunc("/words", h.words)
	mux.HandleFunc("/ordinal", h.ordinal)
	mux.HandleFunc("/roman", h.roman)
	mux.HandleFunc("/from-roman", h.fromRoman)
	mux.HandleFunc("/base", h.base)
	mux.HandleFunc("/scientific", h.scientific)
	mux.HandleFunc("/from-scientific", h.fromScientific)
	mux.HandleFunc("/percentage", h.percentage)
	mux.HandleFunc("/round", h.round)
	mux.HandleFunc("/spell", h.spell)
	mux.HandleFunc("/prime", h.prime)
	mux.HandleFunc("/factorize", h.factorize)
	mux.HandleFunc("/factors", h.factorize)
	mux.HandleFunc("/gcd", h.gcd)
	mux.HandleFunc("/lcm", h.lcm)
	mux.HandleFunc("/fibonacci", h.fibonacci)
	mux.HandleFunc("/fib", h.fibonacci)
	mux.HandleFunc("/even", h.even)
	mux.HandleFunc("/odd", h.odd)
	mux.HandleFunc("/parity", h.parity)
	mux.HandleFunc("/divmod", h.divmod)
	mux.HandleFunc("/abs", h.abs)
	mux.HandleFunc("/negate", h.negate)
	mux.HandleFunc("/digit-sum", h.digitSum)
	mux.HandleFunc("/digit-count", h.digitCount)
	mux.HandleFunc("/reverse", h.reverseDigits)
	mux.HandleFunc("/collatz", h.collatz)
	mux.HandleFunc("/perfect-square", h.perfectSquare)
	mux.HandleFunc("/sqrt", h.sqrt)
	mux.HandleFunc("/power", h.power)
	mux.HandleFunc("/factorial", h.factorial)
	mux.HandleFunc("/armstrong", h.armstrong)
	mux.HandleFunc("/binary", h.binary)
	mux.HandleFunc("/hex", h.hex)
	mux.HandleFunc("/octal", h.octal)
	mux.HandleFunc("/from-binary", h.fromBinary)
	mux.HandleFunc("/from-hex", h.fromHex)
	mux.HandleFunc("/from-octal", h.fromOctal)
	mux.HandleFunc("/comma", h.comma)
	mux.HandleFunc("/uncomma", h.uncomma)
	mux.HandleFunc("/fraction", h.fraction)
	mux.HandleFunc("/mcp", h.mcp)
	mux.HandleFunc("/", h.root)
}

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "numberkit — agentic-first number formatting and conversion service")
		fmt.Fprintln(w, "GET /help for the operating manual")
		return
	}
	writeError(w, r, http.StatusNotFound, "unknown endpoint: "+r.URL.Path, "GET /help to see available endpoints")
}

func wantsJSON(r *http.Request) bool {
	if r.Header.Get("Accept") == "application/json" {
		return true
	}
	if r.URL.Query().Get("format") == "json" {
		return true
	}
	return false
}

func writeError(w http.ResponseWriter, r *http.Request, status int, msg, hint string) {
	w.WriteHeader(status)
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error":"%s","hint":"%s"}`, escapeJSON(msg), escapeJSON(hint))
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "error: %s | hint: %s\n", msg, hint)
	}
}

func writeRecord(w http.ResponseWriter, r *http.Request, fields ...string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "{")
		for i := 0; i < len(fields); i += 2 {
			if i > 0 {
				fmt.Fprint(w, ",")
			}
			fmt.Fprintf(w, `"%s":"%s"`, fields[i], escapeJSON(fields[i+1]))
		}
		fmt.Fprint(w, "}")
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		for i := 0; i < len(fields); i += 2 {
			if i > 0 {
				fmt.Fprint(w, " ")
			}
			fmt.Fprintf(w, "%s=%s", fields[i], fields[i+1])
		}
		fmt.Fprintln(w)
	}
}

func writeLines(w http.ResponseWriter, r *http.Request, lines []string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "[")
		for i, line := range lines {
			if i > 0 {
				fmt.Fprint(w, ",")
			}
			fmt.Fprintf(w, `"%s"`, escapeJSON(line))
		}
		fmt.Fprint(w, "]")
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		for _, line := range lines {
			fmt.Fprintln(w, line)
		}
	}
}

func writeRaw(w http.ResponseWriter, r *http.Request, content string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"result":"%s"}`, escapeJSON(content))
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, content)
	}
}

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}

func getParam(r *http.Request, name string) string {
	val := r.URL.Query().Get(name)
	if val != "" {
		return val
	}
	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		val = r.PostFormValue(name)
	}
	return val
}

func getBody(r *http.Request) string {
	if r.Method == http.MethodPost {
		body, _ := io.ReadAll(r.Body)
		return strings.TrimSpace(string(body))
	}
	return ""
}

func getInput(r *http.Request, names ...string) string {
	for _, name := range names {
		val := getParam(r, name)
		if val != "" {
			return val
		}
	}
	body := getBody(r)
	if body != "" {
		return body
	}
	return ""
}

func (h *Handler) help(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, helpText)
}

func (h *Handler) format(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=1234567 or POST body")
		return
	}
	opts := model.FormatOptions{
		GroupDigits: true,
		Separator:   ",",
		DecimalSep:  ".",
	}
	if v := getParam(r, "decimals"); v != "" {
		fmt.Sscanf(v, "%d", &opts.Decimals)
	}
	if v := getParam(r, "separator"); v != "" {
		opts.Separator = v
	}
	if v := getParam(r, "decimal_sep"); v != "" {
		opts.DecimalSep = v
	}
	if v := getParam(r, "prefix"); v != "" {
		opts.Prefix = v
	}
	if v := getParam(r, "suffix"); v != "" {
		opts.Suffix = v
	}
	if v := getParam(r, "trim"); v == "true" || v == "1" {
		opts.TrimZeros = true
	}
	if v := getParam(r, "group"); v == "false" || v == "0" {
		opts.GroupDigits = false
	}
	result, err := model.Format(num, opts)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "ensure the number is valid, e.g. ?num=1234567.89")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) words(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=123 or POST body")
		return
	}
	result, err := model.ToWords(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "ensure the number is valid, e.g. ?num=123.45")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) ordinal(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=42 or POST body")
		return
	}
	result, err := model.ToOrdinal(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "ensure the number is a valid integer, e.g. ?num=42")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) roman(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide an integer 1-3999 via ?num=2024 or POST body")
		return
	}
	result, err := model.ToRoman(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide an integer between 1 and 3999, e.g. ?num=2024")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) fromRoman(w http.ResponseWriter, r *http.Request) {
	roman := getInput(r, "roman", "value", "r")
	if roman == "" {
		writeError(w, r, http.StatusBadRequest, "missing roman numeral", "provide a Roman numeral via ?roman=MMXXIV or POST body")
		return
	}
	result, err := model.FromRoman(roman)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid Roman numeral like MMXXIV, XIV, or XLII")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) base(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=255&from=10&to=16 or POST body")
		return
	}
	fromBase := 10
	toBase := 16
	if v := getParam(r, "from"); v != "" {
		fmt.Sscanf(v, "%d", &fromBase)
	}
	if v := getParam(r, "to"); v != "" {
		fmt.Sscanf(v, "%d", &toBase)
	}
	result, err := model.BaseConvert(num, fromBase, toBase)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide valid number and bases (2-36), e.g. ?num=255&from=10&to=16")
		return
	}
	writeRecord(w, r, "input", num, "from_base", fmt.Sprintf("%d", fromBase), "to_base", fmt.Sprintf("%d", toBase), "result", result)
}

func (h *Handler) scientific(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=1234567 or POST body")
		return
	}
	precision := 6
	if v := getParam(r, "precision"); v != "" {
		fmt.Sscanf(v, "%d", &precision)
	}
	result, err := model.ToScientific(num, precision)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "ensure the number is valid, e.g. ?num=1234567")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) fromScientific(w http.ResponseWriter, r *http.Request) {
	sci := getInput(r, "sci", "value", "v", "notation")
	if sci == "" {
		writeError(w, r, http.StatusBadRequest, "missing scientific notation", "provide notation via ?sci=1.234e+5 or POST body")
		return
	}
	result, err := model.FromScientific(sci)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide valid scientific notation like 1.234e+5 or 6.022E23")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) percentage(w http.ResponseWriter, r *http.Request) {
	operation := getParam(r, "op")
	if operation == "" {
		operation = "of"
	}
	aStr := getParam(r, "a")
	bStr := getParam(r, "b")
	if aStr == "" || bStr == "" {
		writeError(w, r, http.StatusBadRequest, "missing parameters", "provide ?op=of&a=15&b=200 (what is 15% of 200), ?op=is&a=30&b=150 (30 is what % of 150), or ?op=change&a=100&b=150 (% change from 100 to 150)")
		return
	}
	var a, b float64
	fmt.Sscanf(aStr, "%g", &a)
	fmt.Sscanf(bStr, "%g", &b)
	result, err := model.Percentage(operation, a, b)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "use op=of (a% of b), op=is (a is what % of b), or op=change (% change from a to b)")
		return
	}
	writeRecord(w, r, "operation", operation, "a", aStr, "b", bStr, "result", result)
}

func (h *Handler) round(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=3.14159&decimals=2 or POST body")
		return
	}
	decimals := 0
	if v := getParam(r, "decimals"); v != "" {
		fmt.Sscanf(v, "%d", &decimals)
	}
	mode := getParam(r, "mode")
	if mode == "" {
		mode = "nearest"
	}
	result, err := model.RoundNumber(num, decimals, mode)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "use mode=nearest, up, down, toward_zero, or away_zero. e.g. ?num=3.14159&decimals=2&mode=nearest")
		return
	}
	writeRecord(w, r, "input", num, "decimals", fmt.Sprintf("%d", decimals), "mode", mode, "result", result)
}

func (h *Handler) spell(w http.ResponseWriter, r *http.Request) {
	words := getInput(r, "words", "text", "value")
	if words == "" {
		writeError(w, r, http.StatusBadRequest, "missing words", "provide number words via ?words=one+hundred+twenty+three or POST body")
		return
	}
	result, err := model.SpellToNumber(words)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide valid English number words like 'one hundred twenty three' or 'negative five point two five'")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) prime(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=17 or POST body")
		return
	}
	result, err := model.IsPrime(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid integer, e.g. ?num=17")
		return
	}
	isPrimeStr := "false"
	if result {
		isPrimeStr = "true"
	}
	writeRecord(w, r, "number", num, "is_prime", isPrimeStr)
}

func (h *Handler) factorize(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=360 or POST body")
		return
	}
	result, err := model.Factorize(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a positive integer, e.g. ?num=360")
		return
	}
	writeRecord(w, r, "number", num, "factors", result)
}

func (h *Handler) gcd(w http.ResponseWriter, r *http.Request) {
	a := getParam(r, "a")
	b := getParam(r, "b")
	if a == "" || b == "" {
		writeError(w, r, http.StatusBadRequest, "missing parameters", "provide two numbers via ?a=48&b=18")
		return
	}
	result, err := model.GCD(a, b)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide two valid integers, e.g. ?a=48&b=18")
		return
	}
	writeRecord(w, r, "a", a, "b", b, "gcd", result)
}

func (h *Handler) lcm(w http.ResponseWriter, r *http.Request) {
	a := getParam(r, "a")
	b := getParam(r, "b")
	if a == "" || b == "" {
		writeError(w, r, http.StatusBadRequest, "missing parameters", "provide two numbers via ?a=4&b=6")
		return
	}
	result, err := model.LCM(a, b)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide two valid integers, e.g. ?a=4&b=6")
		return
	}
	writeRecord(w, r, "a", a, "b", b, "lcm", result)
}

func (h *Handler) fibonacci(w http.ResponseWriter, r *http.Request) {
	count := 10
	if v := getParam(r, "count"); v != "" {
		fmt.Sscanf(v, "%d", &count)
	}
	if v := getParam(r, "n"); v != "" {
		fmt.Sscanf(v, "%d", &count)
	}
	result, err := model.Fibonacci(count)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a count between 0 and 1000, e.g. ?count=10")
		return
	}
	writeLines(w, r, result)
}

func (h *Handler) even(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=42 or POST body")
		return
	}
	result, err := model.IsEven(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid integer, e.g. ?num=42")
		return
	}
	isEvenStr := "false"
	if result {
		isEvenStr = "true"
	}
	writeRecord(w, r, "number", num, "is_even", isEvenStr)
}

func (h *Handler) odd(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=43 or POST body")
		return
	}
	result, err := model.IsOdd(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid integer, e.g. ?num=43")
		return
	}
	isOddStr := "false"
	if result {
		isOddStr = "true"
	}
	writeRecord(w, r, "number", num, "is_odd", isOddStr)
}

func (h *Handler) parity(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=42 or POST body")
		return
	}
	result, err := model.Parity(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid integer, e.g. ?num=42")
		return
	}
	writeRecord(w, r, "number", num, "parity", result)
}

func (h *Handler) divmod(w http.ResponseWriter, r *http.Request) {
	a := getParam(r, "a")
	b := getParam(r, "b")
	if a == "" || b == "" {
		writeError(w, r, http.StatusBadRequest, "missing parameters", "provide two numbers via ?a=17&b=5")
		return
	}
	q, rem, err := model.Divmod(a, b)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide two valid integers with non-zero divisor, e.g. ?a=17&b=5")
		return
	}
	writeRecord(w, r, "a", a, "b", b, "quotient", q, "remainder", rem)
}

func (h *Handler) abs(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=-42 or POST body")
		return
	}
	result, err := model.Abs(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid number, e.g. ?num=-42")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) negate(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=42 or POST body")
		return
	}
	result, err := model.Negate(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid number, e.g. ?num=42")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) digitSum(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=12345 or POST body")
		return
	}
	result, err := model.DigitSum(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid number, e.g. ?num=12345")
		return
	}
	writeRecord(w, r, "number", num, "digit_sum", result)
}

func (h *Handler) digitCount(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=12345 or POST body")
		return
	}
	result, err := model.DigitCount(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid number, e.g. ?num=12345")
		return
	}
	writeRecord(w, r, "number", num, "digit_count", result)
}

func (h *Handler) reverseDigits(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=12345 or POST body")
		return
	}
	result, err := model.ReverseDigits(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid number, e.g. ?num=12345")
		return
	}
	writeRecord(w, r, "input", num, "reversed", result)
}

func (h *Handler) collatz(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a positive integer via ?num=27 or POST body")
		return
	}
	result, err := model.Collatz(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a positive integer, e.g. ?num=27")
		return
	}
	writeLines(w, r, result)
}

func (h *Handler) perfectSquare(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=144 or POST body")
		return
	}
	result, err := model.IsPerfectSquare(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid integer, e.g. ?num=144")
		return
	}
	isPS := "false"
	if result {
		isPS = "true"
	}
	writeRecord(w, r, "number", num, "is_perfect_square", isPS)
}

func (h *Handler) sqrt(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=144 or POST body")
		return
	}
	result, err := model.Sqrt(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a non-negative integer, e.g. ?num=144")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) power(w http.ResponseWriter, r *http.Request) {
	base := getParam(r, "base")
	exp := getParam(r, "exp")
	if base == "" || exp == "" {
		writeError(w, r, http.StatusBadRequest, "missing parameters", "provide base and exponent via ?base=2&exp=10")
		return
	}
	result, err := model.Power(base, exp)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide valid integers, e.g. ?base=2&exp=10")
		return
	}
	writeRecord(w, r, "base", base, "exponent", exp, "result", result)
}

func (h *Handler) factorial(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a non-negative integer via ?num=20 or POST body")
		return
	}
	result, err := model.Factorial(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a non-negative integer (max 10000), e.g. ?num=20")
		return
	}
	writeRecord(w, r, "number", num, "factorial", result)
}

func (h *Handler) armstrong(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=153 or POST body")
		return
	}
	result, err := model.IsArmstrong(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid integer, e.g. ?num=153")
		return
	}
	isArm := "false"
	if result {
		isArm = "true"
	}
	writeRecord(w, r, "number", num, "is_armstrong", isArm)
}

func (h *Handler) binary(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=255 or POST body")
		return
	}
	result, err := model.ToBinary(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid integer, e.g. ?num=255")
		return
	}
	writeRecord(w, r, "number", num, "binary", result)
}

func (h *Handler) hex(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=255 or POST body")
		return
	}
	result, err := model.ToHex(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid integer, e.g. ?num=255")
		return
	}
	writeRecord(w, r, "number", num, "hex", result)
}

func (h *Handler) octal(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=255 or POST body")
		return
	}
	result, err := model.ToOctal(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid integer, e.g. ?num=255")
		return
	}
	writeRecord(w, r, "number", num, "octal", result)
}

func (h *Handler) fromBinary(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "binary", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing binary number", "provide a binary number via ?num=11111111 or POST body")
		return
	}
	result, err := model.FromBinary(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid binary number (0s and 1s), e.g. ?num=11111111")
		return
	}
	writeRecord(w, r, "binary", num, "decimal", result)
}

func (h *Handler) fromHex(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "hex", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing hex number", "provide a hex number via ?num=ff or POST body")
		return
	}
	result, err := model.FromHex(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid hex number (0-9, a-f), e.g. ?num=ff")
		return
	}
	writeRecord(w, r, "hex", num, "decimal", result)
}

func (h *Handler) fromOctal(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "octal", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing octal number", "provide an octal number via ?num=377 or POST body")
		return
	}
	result, err := model.FromOctal(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid octal number (0-7), e.g. ?num=377")
		return
	}
	writeRecord(w, r, "octal", num, "decimal", result)
}

func (h *Handler) comma(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number via ?num=1234567 or POST body")
		return
	}
	result, err := model.CommaFormat(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid number, e.g. ?num=1234567")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) uncomma(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a number with commas via ?num=1,234,567 or POST body")
		return
	}
	result, err := model.Uncomma(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a number with commas, e.g. ?num=1,234,567")
		return
	}
	writeRaw(w, r, result)
}

func (h *Handler) fraction(w http.ResponseWriter, r *http.Request) {
	num := getInput(r, "num", "number", "n", "value")
	if num == "" {
		writeError(w, r, http.StatusBadRequest, "missing number", "provide a decimal number via ?num=0.375 or POST body")
		return
	}
	result, err := model.ToFraction(num)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error(), "provide a valid decimal number, e.g. ?num=0.375")
		return
	}
	writeRecord(w, r, "number", num, "fraction", result)
}
