/**
 * sWeb3App - Web3 wallet connection helper class
 */
class sWeb3App {
  constructor() {
    this.web3 = null;
    this.isWeb3Available = false;
  }

  async setWeb3Wallet() {
    try {
      // Check if Web3 provider exists
      if (window.ethereum) {
        this.web3 = new Web3(window.ethereum);
        this.isWeb3Available = true;
      } else if (window.web3) {
        this.web3 = new Web3(window.web3.currentProvider);
        this.isWeb3Available = true;
      } else {
        return {
          error: 1,
          errorMessage: "Web3 wallet not found. Please install MetaMask, Trust Wallet, or Coinbase Wallet."
        };
      }

      // Request account access
      const RequestAccountsData = await this.requestAccounts();
      if (RequestAccountsData.error !== 0) {
        return RequestAccountsData;
      }

      // Get signature for authentication
      const checkSignData = await this.checkSign(RequestAccountsData.requestAccountsArr[0]);
      if (checkSignData.error !== 0) {
        return checkSignData;
      }

      return {
        error: 0,
        RequestAccountsData: RequestAccountsData,
        checkSignData: checkSignData
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
        requestAccountsArr: accounts
      };
    } catch (error) {
      return {
        error: 1,
        errorMessage: "Error requesting accounts: " + error.message
      };
    }
  }

  async checkSign(address) {
    try {
      // Create a unique message for signing
      const timestamp = new Date().getTime();
      const domain = window.location.hostname;
      const signMessage = `Sign this message to authenticate on ${domain} at ${timestamp}`;

      // Request signature from user
      const signResult = await window.ethereum.request({
        method: 'personal_sign',
        params: [signMessage, address]
      });

      return {
        error: 0,
        SignData: {
          signMessage: signMessage,
          signResult: signResult
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