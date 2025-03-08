/**
 * Main application JavaScript
 */
document.addEventListener('DOMContentLoaded', () => {
  window.walletApp = new WalletApp();

  setupEventListeners();

  if (document.getElementById('wallet-dashboard')) {
    initDashboard();
  }

  const statusElement = document.getElementById('status-message');
  if (statusElement) {
    statusElement.style.display = 'none';
  }
});

function setupEventListeners() {
  document.addEventListener('click', (event) => {
    if (event.target.classList.contains('login-web3-btn')) {
      event.preventDefault();
      handleWeb3Login(event.target);
    }
  });

  const refreshButton = document.getElementById('refresh-balance');
  if (refreshButton) {
    refreshButton.addEventListener('click', refreshWalletBalance);
  }
}

async function handleWeb3Login(button) {
  const statusElement = document.getElementById('status-message');
  const userContentUrl = button.dataset.userContentUrl;
  const web3AuthUrl = button.dataset.userAuthUrl;

  if (!userContentUrl || !web3AuthUrl) {
    showStatus("Configuration error: Missing URL attributes", "error");
    return;
  }

  // Display loading state
  button.classList.add('loading');
  button.textContent = 'Connecting...';
  showStatus('Connecting to wallet...', 'info');

  try {
    // Initialize Web3 wallet connection
    const web3Data = await window.walletApp.connect();
    if (web3Data.error !== 0) {
      showStatus(web3Data.errorMessage || 'Error connecting to wallet', "error");
      resetButton(button);
      return;
    }

    // Check if we have signature data
    if (!web3Data.signature) {
      showStatus('Failed to get signature from wallet', "error");
      resetButton(button);
      return;
    }

    showStatus('Wallet connected. Authenticating...', 'info');

    // Get account and signature data
    const currentAccount = web3Data.accounts[0];
    const signMessage = web3Data.signature.message;
    const signResult = web3Data.signature.result;

    const formData = new FormData();
    formData.append('wallet', currentAccount);
    formData.append('message', signMessage);
    formData.append('signature', signResult);

    const response = await fetch(web3AuthUrl, {
      method: 'POST',
      body: formData
    });

    const data = await response.json();

    if (data.error === 0) {
      showStatus('Authentication successful! Redirecting...', 'success');
      setTimeout(() => {
        window.location.href = userContentUrl;
      }, 1000);
    } else {
      showStatus(data.errorMessage || 'Authentication failed', 'error');
      resetButton(button);
    }
  } catch (error) {
    showStatus('Error: ' + error.message, 'error');
    resetButton(button);
  }
}

async function initDashboard() {
  const walletAddress = document.getElementById('wallet-address');
  const balanceElement = document.getElementById('wallet-balance');

  if (!walletAddress || !balanceElement) return;

  try {
    // Connect wallet
    const web3Data = await window.walletApp.connect();
    if (web3Data.error !== 0) {
      showStatus(web3Data.errorMessage || 'Error connecting to wallet', "error");
      return;
    }

    // Get balance
    await refreshWalletBalance();
  } catch (error) {
    showStatus('Error initializing dashboard: ' + error.message, 'error');
  }
}

async function refreshWalletBalance() {
  const balanceElement = document.getElementById('wallet-balance');
  if (!balanceElement) return;

  try {
    const balanceData = await window.walletApp.getBalance();
    if (balanceData.error === 0) {
      balanceElement.textContent = `${balanceData.balance} ETH`;
    } else {
      showStatus(balanceData.errorMessage || 'Error fetching balance', 'error');
    }
  } catch (error) {
    showStatus('Error refreshing balance: ' + error.message, 'error');
  }
}

function showStatus(message, type = 'info') {
  const statusElement = document.getElementById('status-message');
  if (!statusElement) return;

  statusElement.className = `alert alert-${type === 'error' ? 'danger' : type}`;
  statusElement.textContent = message;
  statusElement.style.display = 'block';

  if (type === 'error') {
    console.error(message);
  }
}

function resetButton(button) {
  button.classList.remove('loading');
  button.textContent = 'Login with Web3 Wallet';
}
