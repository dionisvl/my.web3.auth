/**
 * Web3 Wallet Integration using ethers.js
 */
class WalletApp {
  constructor() {
    this.provider = null;
    this.signer = null;
    this.wallet = null;
  }

  async initProvider() {
    try {
      if (!window.ethereum) {
        return {
          error: 1,
          errorMessage: "Web3 wallet not found. Please install MetaMask, Trust Wallet, or Coinbase Wallet."
        };
      }

      this.provider = new ethers.BrowserProvider(window.ethereum);
      return { error: 0 };
    } catch (error) {
      return {
        error: 1,
        errorMessage: "Error initializing provider: " + error.message
      };
    }
  }

  async connect() {
    try {
      if (!window.ethereum) {
        return {
          error: 1,
          errorMessage: "Web3 wallet not found. Please install MetaMask, Trust Wallet, or Coinbase Wallet."
        };
      }

      this.provider = new ethers.BrowserProvider(window.ethereum);

      const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });

      if (!accounts || accounts.length === 0) {
        return {
          error: 1,
          errorMessage: "No accounts found. Please unlock your wallet."
        };
      }

      this.wallet = accounts[0];

      this.signer = await this.provider.getSigner();

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
      if (!this.wallet) {
        return {
          error: 1,
          errorMessage: "Wallet not connected"
        };
      }

      if (!this.provider) {
        const result = await this.initProvider();
        if (result.error !== 0) {
          return result;
        }
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
