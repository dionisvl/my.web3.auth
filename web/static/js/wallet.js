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
      if (window.ethereum) {
        this.provider = new ethers.BrowserProvider(window.ethereum);
      } else if (window.trustwallet) {
        this.provider = new ethers.BrowserProvider(window.trustwallet);
      } else if (window.web3) {
        this.provider = new ethers.BrowserProvider(window.web3.currentProvider);
      } else {
        return {
          error: 1,
          errorMessage: "Web3 wallet not found. Please install MetaMask, Trust Wallet, or Coinbase Wallet."
        };
      }

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
      let provider;
      if (window.ethereum) {
        provider = window.ethereum;
      } else if (window.trustwallet) {
        provider = window.trustwallet;
      } else if (window.web3) {
        provider = window.web3.currentProvider;
      } else {
        return {
          error: 1,
          errorMessage: "Web3 wallet not found. Please install MetaMask, Trust Wallet, or Coinbase Wallet."
        };
      }

      this.provider = new ethers.BrowserProvider(provider);

      const accounts = await provider.request({ method: 'eth_requestAccounts' });

      if (!accounts || accounts.length === 0) {
        return {
          error: 1,
          errorMessage: "No accounts found. Please unlock your wallet."
        };
      }

      this.wallet = accounts[0];

      this.signer = await this.provider.getSigner();

      const signature = await this.getSignature();
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

  async getSignature() {
    try {
      // Fetch a single-use, server-issued challenge (nonce + expiry) so a
      // captured signature cannot be replayed. The wallet signs it verbatim.
      const resp = await fetch('/api/nonce', {
        credentials: 'same-origin'
      });
      if (!resp.ok) {
        return {
          error: 1,
          errorMessage: `Failed to obtain authentication challenge (HTTP ${resp.status})`
        };
      }
      const challenge = await resp.json();
      if (!challenge || challenge.error !== 0 || !challenge.message) {
        return {
          error: 1,
          errorMessage: challenge.errorMessage || "Failed to obtain authentication challenge"
        };
      }

      const message = challenge.message;
      const signature = await this.signer.signMessage(message);

      return {
        error: 0,
        data: {
          message: message,
          result: signature,
          csrfToken: challenge.csrfToken
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
