# Web3 Wallet Authentication with PHP

A modern PHP application for Web3 wallet authentication and balance checking using Ethers.js.

## Features

- Authentication with Web3 wallets (MetaMask, Trust Wallet, Coinbase Wallet)
- Secure sign/verify process using Ethereum signatures
- Session management
- Wallet balance checking with Ethers.js
- Responsive UI

## Requirements

- PHP 8.0+
- Composer
- Web3 wallet (MetaMask, Trust Wallet, etc.)

## Installation

1. Clone the repository

2. Install dependencies
   ```
   composer install
   ```

3. Create and configure environment variables
   ```
   cp .env.example .env
   ```

4. Rum server
   ```
   make up
   ```

## Usage

1. Navigate to the homepage in your browser
2. Click "Login with Web3 Wallet"
3. Connect your wallet when prompted and sign the authentication message
4. View your wallet balance and network information on the dashboard

## License

MIT
