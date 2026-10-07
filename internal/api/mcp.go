package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/relentlessworks/numberkit/internal/model"
)

type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MCPTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type MCPToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

var mcpTools = []MCPTool{
	{Name: "format", Description: "Format a number with thousands separators, decimal places, prefix/suffix", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}, "decimals": map[string]interface{}{"type": "integer"}, "separator": map[string]interface{}{"type": "string"}, "prefix": map[string]interface{}{"type": "string"}, "suffix": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "words", Description: "Convert a number to English words", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "ordinal", Description: "Convert a number to ordinal form (1st, 2nd, 3rd)", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "roman", Description: "Convert an integer (1-3999) to Roman numerals", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "from_roman", Description: "Convert Roman numerals to a number", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"roman": map[string]interface{}{"type": "string"}}, "required": []string{"roman"}}},
	{Name: "base_convert", Description: "Convert a number between bases (2-36)", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}, "from": map[string]interface{}{"type": "integer"}, "to": map[string]interface{}{"type": "integer"}}, "required": []string{"num"}}},
	{Name: "scientific", Description: "Convert a number to scientific notation", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}, "precision": map[string]interface{}{"type": "integer"}}, "required": []string{"num"}}},
	{Name: "from_scientific", Description: "Convert scientific notation to a regular number", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"sci": map[string]interface{}{"type": "string"}}, "required": []string{"sci"}}},
	{Name: "percentage", Description: "Percentage calculations (of, is, change)", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"op": map[string]interface{}{"type": "string"}, "a": map[string]interface{}{"type": "number"}, "b": map[string]interface{}{"type": "number"}}, "required": []string{"a", "b"}}},
	{Name: "round", Description: "Round a number (nearest, up, down, toward_zero, away_zero)", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}, "decimals": map[string]interface{}{"type": "integer"}, "mode": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "spell", Description: "Convert English number words to digits", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"words": map[string]interface{}{"type": "string"}}, "required": []string{"words"}}},
	{Name: "prime", Description: "Check if a number is prime", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "factorize", Description: "Get prime factorization of a number", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "gcd", Description: "Greatest common divisor of two numbers", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"a": map[string]interface{}{"type": "string"}, "b": map[string]interface{}{"type": "string"}}, "required": []string{"a", "b"}}},
	{Name: "lcm", Description: "Least common multiple of two numbers", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"a": map[string]interface{}{"type": "string"}, "b": map[string]interface{}{"type": "string"}}, "required": []string{"a", "b"}}},
	{Name: "fibonacci", Description: "Generate the first n Fibonacci numbers", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"count": map[string]interface{}{"type": "integer"}}, "required": []string{"count"}}},
	{Name: "parity", Description: "Check if a number is even or odd", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "divmod", Description: "Integer division with quotient and remainder", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"a": map[string]interface{}{"type": "string"}, "b": map[string]interface{}{"type": "string"}}, "required": []string{"a", "b"}}},
	{Name: "abs", Description: "Absolute value of a number", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "digit_sum", Description: "Sum of digits of a number", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "digit_count", Description: "Count of digits in a number", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "reverse_digits", Description: "Reverse the digits of a number", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "collatz", Description: "Generate the Collatz sequence from a number", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "perfect_square", Description: "Check if a number is a perfect square", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "sqrt", Description: "Integer square root of a number", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "power", Description: "Compute base^exponent for integers", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"base": map[string]interface{}{"type": "string"}, "exp": map[string]interface{}{"type": "string"}}, "required": []string{"base", "exp"}}},
	{Name: "factorial", Description: "Compute n! (factorial)", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "armstrong", Description: "Check if a number is an Armstrong number", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "binary", Description: "Convert a number to binary", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "hex", Description: "Convert a number to hexadecimal", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "octal", Description: "Convert a number to octal", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "from_binary", Description: "Convert binary to decimal", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "from_hex", Description: "Convert hexadecimal to decimal", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "from_octal", Description: "Convert octal to decimal", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "comma", Description: "Add thousands separators to a number", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "uncomma", Description: "Remove thousands separators from a number", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
	{Name: "fraction", Description: "Convert a decimal to a simplified fraction", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"num": map[string]interface{}{"type": "string"}}, "required": []string{"num"}}},
}

func (h *Handler) mcp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed", "POST JSON-RPC 2.0 requests to /mcp")
		return
	}

	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"jsonrpc":"2.0","error":{"code":-32700,"message":"parse error"}}`)
		return
	}

	resp := MCPResponse{JSONRPC: "2.0", ID: req.ID}

	switch req.Method {
	case "initialize":
		resp.Result = map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "numberkit",
				"version": "0.1.0",
			},
		}

	case "tools/list":
		resp.Result = map[string]interface{}{"tools": mcpTools}

	case "tools/call":
		var params MCPToolCallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			resp.Error = &MCPError{Code: -32602, Message: "invalid params"}
			break
		}
		result, err := h.callMCPTool(params)
		if err != nil {
			resp.Error = &MCPError{Code: -32603, Message: err.Error()}
		} else {
			resp.Result = map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": result},
				},
			}
		}

	default:
		resp.Error = &MCPError{Code: -32601, Message: "method not found: " + req.Method}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) callMCPTool(params MCPToolCallParams) (string, error) {
	getStr := func(key string) string {
		if v, ok := params.Arguments[key]; ok {
			return fmt.Sprintf("%v", v)
		}
		return ""
	}
	getInt := func(key string, def int) int {
		if v, ok := params.Arguments[key]; ok {
			n := def
			fmt.Sscanf(fmt.Sprintf("%v", v), "%d", &n)
			return n
		}
		return def
	}
	getFloat := func(key string) float64 {
		if v, ok := params.Arguments[key]; ok {
			f := 0.0
			fmt.Sscanf(fmt.Sprintf("%v", v), "%g", &f)
			return f
		}
		return 0
	}

	switch params.Name {
	case "format":
		opts := model.FormatOptions{GroupDigits: true, Separator: ",", DecimalSep: "."}
		opts.Decimals = getInt("decimals", 0)
		if v := getStr("separator"); v != "" {
			opts.Separator = v
		}
		if v := getStr("prefix"); v != "" {
			opts.Prefix = v
		}
		if v := getStr("suffix"); v != "" {
			opts.Suffix = v
		}
		r, err := model.Format(getStr("num"), opts)
		if err != nil {
			return "", err
		}
		return r, nil
	case "words":
		return model.ToWords(getStr("num"))
	case "ordinal":
		return model.ToOrdinal(getStr("num"))
	case "roman":
		return model.ToRoman(getStr("num"))
	case "from_roman":
		return model.FromRoman(getStr("roman"))
	case "base_convert":
		return model.BaseConvert(getStr("num"), getInt("from", 10), getInt("to", 16))
	case "scientific":
		return model.ToScientific(getStr("num"), getInt("precision", 6))
	case "from_scientific":
		return model.FromScientific(getStr("sci"))
	case "percentage":
		return model.Percentage(getStr("op"), getFloat("a"), getFloat("b"))
	case "round":
		return model.RoundNumber(getStr("num"), getInt("decimals", 0), getStr("mode"))
	case "spell":
		return model.SpellToNumber(getStr("words"))
	case "prime":
		r, err := model.IsPrime(getStr("num"))
		if err != nil {
			return "", err
		}
		if r {
			return "true", nil
		}
		return "false", nil
	case "factorize":
		return model.Factorize(getStr("num"))
	case "gcd":
		return model.GCD(getStr("a"), getStr("b"))
	case "lcm":
		return model.LCM(getStr("a"), getStr("b"))
	case "fibonacci":
		r, err := model.Fibonacci(getInt("count", 10))
		if err != nil {
			return "", err
		}
		result := ""
		for i, v := range r {
			if i > 0 {
				result += "\n"
			}
			result += v
		}
		return result, nil
	case "parity":
		return model.Parity(getStr("num"))
	case "divmod":
		q, rem, err := model.Divmod(getStr("a"), getStr("b"))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("quotient=%s remainder=%s", q, rem), nil
	case "abs":
		return model.Abs(getStr("num"))
	case "digit_sum":
		return model.DigitSum(getStr("num"))
	case "digit_count":
		return model.DigitCount(getStr("num"))
	case "reverse_digits":
		return model.ReverseDigits(getStr("num"))
	case "collatz":
		r, err := model.Collatz(getStr("num"))
		if err != nil {
			return "", err
		}
		result := ""
		for i, v := range r {
			if i > 0 {
				result += "\n"
			}
			result += v
		}
		return result, nil
	case "perfect_square":
		r, err := model.IsPerfectSquare(getStr("num"))
		if err != nil {
			return "", err
		}
		if r {
			return "true", nil
		}
		return "false", nil
	case "sqrt":
		return model.Sqrt(getStr("num"))
	case "power":
		return model.Power(getStr("base"), getStr("exp"))
	case "factorial":
		return model.Factorial(getStr("num"))
	case "armstrong":
		r, err := model.IsArmstrong(getStr("num"))
		if err != nil {
			return "", err
		}
		if r {
			return "true", nil
		}
		return "false", nil
	case "binary":
		return model.ToBinary(getStr("num"))
	case "hex":
		return model.ToHex(getStr("num"))
	case "octal":
		return model.ToOctal(getStr("num"))
	case "from_binary":
		return model.FromBinary(getStr("num"))
	case "from_hex":
		return model.FromHex(getStr("num"))
	case "from_octal":
		return model.FromOctal(getStr("num"))
	case "comma":
		return model.CommaFormat(getStr("num"))
	case "uncomma":
		return model.Uncomma(getStr("num"))
	case "fraction":
		return model.ToFraction(getStr("num"))
	default:
		return "", fmt.Errorf("unknown tool: %s", params.Name)
	}
}
