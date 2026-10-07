package api

const helpText = `numberkit — agentic-first number formatting and conversion service

The agent IS the interface. No UI, no SDK. Plain text API, JSON on demand.

AUTH: This service runs in no-auth mode by default. All endpoints are open.

RESPONSE FORMAT:
  Plain text by default (key=value pairs, one record per line).
  JSON on demand: send Accept: application/json or add ?format=json.

ERRORS:
  error: message | hint: what to do next

ENDPOINTS:

  GET/POST /format?num=1234567.89&decimals=2&separator=,&prefix=$&suffix=%
    Format a number with separators, decimal places, prefix/suffix.
    Params: num (required), decimals, separator, decimal_sep, prefix, suffix, trim (bool), group (bool)

  GET/POST /words?num=123.45
    Convert a number to English words. Supports negatives and decimals.

  GET/POST /ordinal?num=42
    Convert a number to ordinal form (1st, 2nd, 3rd, 4th...).

  GET/POST /roman?num=2024
    Convert an integer (1-3999) to Roman numerals.

  GET/POST /from-roman?roman=MMXXIV
    Convert Roman numerals to a number.

  GET/POST /base?num=255&from=10&to=16
    Convert a number between bases (2-36).

  GET/POST /scientific?num=1234567&precision=6
    Convert a number to scientific notation.

  GET/POST /from-scientific?sci=1.234e+5
    Convert scientific notation to a regular number.

  GET/POST /percentage?op=of&a=15&b=200
    Percentage calculations. op=of (a% of b), op=is (a is what % of b), op=change (% change from a to b).

  GET/POST /round?num=3.14159&decimals=2&mode=nearest
    Round a number. mode: nearest, up, down, toward_zero, away_zero.

  GET/POST /spell?words=one+hundred+twenty+three
    Convert English number words to digits.

  GET/POST /prime?num=17
    Check if a number is prime.

  GET/POST /factorize?num=360
    Get prime factorization of a number.

  GET/POST /gcd?a=48&b=18
    Greatest common divisor of two numbers.

  GET/POST /lcm?a=4&b=6
    Least common multiple of two numbers.

  GET/POST /fibonacci?count=10
    Generate the first n Fibonacci numbers.

  GET/POST /even?num=42
    Check if a number is even.

  GET/POST /odd?num=43
    Check if a number is odd.

  GET/POST /parity?num=42
    Get parity (even or odd) of a number.

  GET/POST /divmod?a=17&b=5
    Integer division with quotient and remainder.

  GET/POST /abs?num=-42
    Absolute value of a number.

  GET/POST /negate?num=42
    Negate a number.

  GET/POST /digit-sum?num=12345
    Sum of digits of a number.

  GET/POST /digit-count?num=12345
    Count of digits in a number.

  GET/POST /reverse?num=12345
    Reverse the digits of a number.

  GET/POST /collatz?num=27
    Generate the Collatz sequence from a number.

  GET/POST /perfect-square?num=144
    Check if a number is a perfect square.

  GET/POST /sqrt?num=144
    Integer square root of a number.

  GET/POST /power?base=2&exp=10
    Compute base^exponent for integers.

  GET/POST /factorial?num=20
    Compute n! (factorial).

  GET/POST /armstrong?num=153
    Check if a number is an Armstrong number.

  GET/POST /binary?num=255
    Convert a number to binary.

  GET/POST /hex?num=255
    Convert a number to hexadecimal.

  GET/POST /octal?num=255
    Convert a number to octal.

  GET/POST /from-binary?num=11111111
    Convert binary to decimal.

  GET/POST /from-hex?num=ff
    Convert hexadecimal to decimal.

  GET/POST /from-octal?num=377
    Convert octal to decimal.

  GET/POST /comma?num=1234567
    Add thousands separators to a number.

  GET/POST /uncomma?num=1,234,567
    Remove thousands separators from a number.

  GET/POST /fraction?num=0.375
    Convert a decimal to a simplified fraction.

  POST /mcp
    MCP (Model Context Protocol) JSON-RPC 2.0 endpoint.
    Supports: initialize, tools/list, tools/call.

  GET /help (or /.well-known/agent.md)
    This help text.

EXAMPLES:
  curl 'http://localhost:7780/words?num=12345'
  curl 'http://localhost:7780/roman?num=2024'
  curl 'http://localhost:7780/base?num=255&from=10&to=16'
  curl -H 'Accept: application/json' 'http://localhost:7780/prime?num=17'
  curl -d '1234567.89' 'http://localhost:7780/format?decimals=2&prefix=$'
`
