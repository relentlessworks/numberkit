# numberkit

Agentic-first number formatting and conversion service. Format numbers with separators, convert to words, ordinals, Roman numerals, number bases, scientific notation, percentages, fractions, and more. Plain text API, agent-driven, single Go binary.

## Quick Start

```bash
make build
./numberkit
```

Then:

```bash
curl 'http://localhost:7780/words?num=12345'
# one hundred twenty three thousand four hundred and five

curl 'http://localhost:7780/roman?num=2024'
# MMXXIV

curl 'http://localhost:7780/base?num=255&from=10&to=16'
# input=255 from_base=10 to_base=16 result=ff

curl 'http://localhost:7780/factorize?num=360'
# number=360 factors=2^3 × 3^2 × 5

curl -H 'Accept: application/json' 'http://localhost:7780/prime?num=17'
# {"number":"17","is_prime":"true"}
```

## Principles

- **The agent IS the interface** — No UI, no SDK. The API is the product.
- **Plain text by default** — One labeled, grepable line per record. JSON on demand via `Accept: application/json` or `?format=json`.
- **Instructive errors** — Every 4xx includes a hint telling the agent what to do next.
- **Self-documenting** — `GET /help` returns a one-page operating manual.
- **Single static binary** — Go, CGO_ENABLED=0, zero external dependencies.
- **Zero config defaults** — Runs out of the box.
- **MCP connector** — Speaks Model Context Protocol at `POST /mcp`.

## Endpoints

| Endpoint | Description |
|----------|-------------|
| `/format` | Format number with separators, decimals, prefix/suffix |
| `/words` | Convert number to English words |
| `/ordinal` | Convert to ordinal (1st, 2nd, 3rd) |
| `/roman` | Convert to Roman numerals (1-3999) |
| `/from-roman` | Convert Roman numerals to number |
| `/base` | Convert between bases (2-36) |
| `/scientific` | Convert to scientific notation |
| `/from-scientific` | Convert scientific notation to number |
| `/percentage` | Percentage calculations (of, is, change) |
| `/round` | Round numbers (nearest, up, down, toward_zero, away_zero) |
| `/spell` | Convert English words to digits |
| `/prime` | Check if prime |
| `/factorize` | Prime factorization |
| `/gcd` | Greatest common divisor |
| `/lcm` | Least common multiple |
| `/fibonacci` | Generate Fibonacci sequence |
| `/even` | Check if even |
| `/odd` | Check if odd |
| `/parity` | Get parity (even/odd) |
| `/divmod` | Integer division with remainder |
| `/abs` | Absolute value |
| `/negate` | Negate number |
| `/digit-sum` | Sum of digits |
| `/digit-count` | Count of digits |
| `/reverse` | Reverse digits |
| `/collatz` | Collatz sequence |
| `/perfect-square` | Check if perfect square |
| `/sqrt` | Integer square root |
| `/power` | Integer power (base^exp) |
| `/factorial` | Factorial (n!) |
| `/armstrong` | Check if Armstrong number |
| `/binary` | Convert to binary |
| `/hex` | Convert to hexadecimal |
| `/octal` | Convert to octal |
| `/from-binary` | Convert binary to decimal |
| `/from-hex` | Convert hex to decimal |
| `/from-octal` | Convert octal to decimal |
| `/comma` | Add thousands separators |
| `/uncomma` | Remove thousands separators |
| `/fraction` | Convert decimal to simplified fraction |
| `/mcp` | MCP JSON-RPC 2.0 endpoint |
| `/help` | Operating manual |

## Configuration

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `-addr` | `NUMBERKIT_ADDR` | `:7780` | Listen address |
| `-no-auth` | `NUMBERKIT_NO_AUTH` | `true` | Disable authentication |

## Build

```bash
make build    # CGO_ENABLED=0 go build -trimpath ./cmd/numberkit
make test     # go test -race ./...
make vet      # go vet ./...
```

## License

MIT
