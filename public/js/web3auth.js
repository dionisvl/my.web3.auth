/**
 * Web3 authentication handler
 */
document.addEventListener('DOMContentLoaded', () => {
  window.web3App = new Web3App();

  document.addEventListener('click', (event) => {
    if (event.target.classList.contains('login-web3-btn')) {
      event.preventDefault();
      handleWeb3Login(event.target);
    }
  });
});

async function handleWeb3Login(button) {
  const statusElement = document.getElementById('status-message');
  const userContentUrl = button.getAttribute('data-user-content-url');
  const web3AuthUrl = button.getAttribute('data-user-auth-url');

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
    const web3Data = await window.web3App.connect();
    if (web3Data.error !== 0) {
      showStatus(web3Data.errorMessage || 'Error connecting to wallet', "error");
      resetButton(button);
      return;
    }

    showStatus('Wallet connected. Authenticating...', 'info');

    // Get account and signature data
    const currentAccount = web3Data.accounts[0];
    const signMessage = web3Data.signature.message;
    const signResult = web3Data.signature.result;

    // Submit authentication data to server
    const formData = new FormData();
    formData.append('wallet', currentAccount);
    formData.append('message', signMessage);
    formData.append('signature', signResult);

    // Send authentication request
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

function showStatus(message, type = 'info') {
  const statusElement = document.getElementById('status-message');
  statusElement.className = `alert alert-${type === 'error' ? 'danger' : type}`;
  statusElement.textContent = message;

  if (type === 'error') {
    console.error(message);
  }
}

function resetButton(button) {
  button.classList.remove('loading');
  button.textContent = 'Login with Web3 Wallet';
}