# Web3 Wallet Auth (Go)

A small Go application for Web3 wallet authentication and balance checking.
The frontend uses ethers.js; the backend verifies Ethereum signatures
(EIP-191 `personal_sign`) via [go-ethereum](https://github.com/ethereum/go-ethereum).

Rewritten from the original PHP/Slim implementation, which now lives under
[`legacy-platforms/php`](./legacy-platforms/php) (kept for reference).

## Features

- Authentication with Web3 wallets (MetaMask, Trust Wallet, Coinbase Wallet)
- Signature verify via secp256k1 public-key recovery (go-ethereum)
- Cookie-based sessions (`gorilla/sessions`)
- Wallet balance checking with ethers.js (client-side)
- Single self-contained binary — templates and JS embedded via `//go:embed`

## Requirements

- Go 1.23+
- A Web3 wallet (MetaMask, Trust Wallet, etc.)

## Run

```
cp .env.example .env      # set SESSION_KEY for persistent sessions
make run                  # or: go run ./cmd/server
```

Then open http://localhost:8080.

With Docker + Traefik (as the original):

```
make up
```

## Test

```
make test
```

Unit tests cover signature verification (valid/invalid/tampered/malformed cases).

## Configuration

| Env           | Default   | Description                                  |
|---------------|-----------|----------------------------------------------|
| `APP_PORT`    | `8080`    | HTTP listen port                             |
| `ETH_NETWORK` | `sepolia` | Network name (mainnet/goerli/sepolia/avax)   |
| `SESSION_KEY` | *(empty)* | Cookie signing key; empty = ephemeral random |

## Layout

```
cmd/server         entrypoint
internal/auth      signature verification + sessions
internal/wallet    network → chain id
internal/handlers  HTTP handlers (net/http, Go 1.22+ mux)
internal/config    env config
web/               embedded templates + static JS
legacy-platforms/  original PHP implementation
```

## Usage

1. Open the homepage
2. Click "Login with Web3 Wallet"
3. Connect your wallet and sign the authentication message
4. View your wallet balance and network on the dashboard

## License

MIT
