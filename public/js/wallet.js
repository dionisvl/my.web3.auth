/**
 * Web3 Wallet Integration using ethers.js
 */
class WalletApp {
  constructor() {
    this.provider = null;
    this.signer = null;
    this.wallet = null;
  }

  async connect() {
    try {
      if (!window.ethereum) {
        return {
          error: 1,
          errorMessage: "Web3 wallet not found. Please install MetaMask, Trust Wallet, or Coinbase Wallet."
        };
      }

      // Create ethers provider from Web3 provider
      this.provider = new ethers.BrowserProvider(window.ethereum);

      // Request account access
      const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });

      if (!accounts || accounts.length === 0) {
        return {
          error: 1,
          errorMessage: "No accounts found. Please unlock your wallet."
        };
      }

      this.wallet = accounts[0];

      // Get signer
      this.signer = await this.provider.getSigner();

      // Get signature for authentication
      const signature = await this.getSignature(this.wallet);
      if (signature.error !== 0) {
        return signature;
      }

      return {
        error: 0,
        accounts: accounts,
        signature: signature.data
      };
    } catch (error) {
      return {
        error: 1,
        errorMessage: "Error connecting to wallet: " + error.message
      };
    }
  }

  async getBalance() {
    try {
      if (!this.wallet || !this.provider) {
        return {
          error: 1,
          errorMessage: "Wallet not connected"
        };
      }

      const balance = await this.provider.getBalance(this.wallet);
      const formattedBalance = ethers.formatEther(balance);

      return {
        error: 0,
        balance: formattedBalance,
        rawBalance: balance.toString(),
        wallet: this.wallet
      };
    } catch (error) {
      return {
        error: 1,
        errorMessage: "Error getting balance: " + error.message
      };
    }
  }

  async getSignature(address) {
    try {
      const timestamp = new Date().getTime();
      const domain = window.location.hostname;
      const message = `Sign this message to authenticate on ${domain} at ${timestamp}`;

      // Use ethers to sign the message
      const signature = await this.signer.signMessage(message);

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
