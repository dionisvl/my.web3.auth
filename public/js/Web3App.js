/**
 * Web3App - Web3 wallet connection helper class
 */
class Web3App {
  constructor() {
    this.web3 = null;
  }

  async connect() {
    try {
      if (!window.ethereum) {
        return {
          error: 1,
          errorMessage: "Web3 wallet not found. Please install MetaMask, Trust Wallet, or Coinbase Wallet."
        };
      }

      this.web3 = new Web3(window.ethereum);

      // Request account access
      const accounts = await this.requestAccounts();
      if (accounts.error !== 0) {
        return accounts;
      }

      // Get signature for authentication
      const signature = await this.getSignature(accounts.accounts[0]);
      if (signature.error !== 0) {
        return signature;
      }

      return {
        error: 0,
        accounts: accounts.accounts,
        signature: signature.data
      };
    } catch (error) {
      return {
        error: 1,
        errorMessage: "Error connecting to wallet: " + error.message
      };
    }
  }

  async requestAccounts() {
    try {
      const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });

      if (!accounts || accounts.length === 0) {
        return {
          error: 1,
          errorMessage: "No accounts found. Please unlock your wallet."
        };
      }

      return {
        error: 0,
        accounts: accounts
      };
    } catch (error) {
      return {
        error: 1,
        errorMessage: "Error requesting accounts: " + error.message
      };
    }
  }

  async getSignature(address) {
    try {
      const timestamp = new Date().getTime();
      const domain = window.location.hostname;
      const message = `Sign this message to authenticate on ${domain} at ${timestamp}`;

      const signature = await window.ethereum.request({
        method: 'personal_sign',
        params: [message, address]
      });

      return {
        error: 0,
        data: {
          message: message,
          result: signature
        }
      };
    } catch (error) {
      return {
        error: 1,
        errorMessage: "Error signing message: " + error.message
      };
    }
  }
}