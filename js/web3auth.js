/**
 * Web3 authentication handler
 */
document.addEventListener('DOMContentLoaded', function() {
  // Initialize sWeb3App class when DOM is ready
  window.sWeb3AppClass = new sWeb3App();

  // Handle login button clicks
  document.addEventListener('click', function(event) {
    var $this = event.target;

    if ($this.classList.contains('LoginWeb3Button')) {
      event.preventDefault();
      actionLoginWeb3Button($this);
      return false;
    }
  });
});

/**
 * Handle Web3 login button action
 */
async function actionLoginWeb3Button($this) {
  const statusElement = document.getElementById('status-message');

  // Get target URLs from button attributes
  const userContentUrl = $this.getAttribute('data-user-content-url');
  const web3AuthUrl = $this.getAttribute('data-user-auth-url');

  if (!userContentUrl || !web3AuthUrl) {
    showError("Configuration error: Missing URL attributes");
    return;
  }

  // Display loading state
  $this.classList.add('loading');
  $this.textContent = 'Connecting...';
  statusElement.className = 'status-message';
  statusElement.textContent = 'Connecting to wallet...';

  try {
    // Initialize Web3 wallet connection
    const Web3WalletData = await window.sWeb3AppClass.setWeb3Wallet();
    if (Web3WalletData.error !== 0) {
      showError(Web3WalletData.errorMessage || 'Error connecting to wallet');
      resetButton($this);
      return;
    }

    statusElement.textContent = 'Wallet connected. Authenticating...';

    // Get account and signature data
    const currentAccount = Web3WalletData.RequestAccountsData.requestAccountsArr[0];
    const signMessage = Web3WalletData.checkSignData.SignData.signMessage;
    const signResult = Web3WalletData.checkSignData.SignData.signResult;

    // Submit authentication data to server
    const formData = new FormData();
    formData.append('web3-wallet', currentAccount);
    formData.append('signMessage', signMessage);
    formData.append('signResult', signResult);

    // Send authentication request
    fetch(web3AuthUrl, {
      method: 'POST',
      body: formData
    })
      .then(response => response.json())
      .then(data => {
        if (data.error === 0) {
          statusElement.className = 'status-message success-message';
          statusElement.textContent = 'Authentication successful! Redirecting...';
          setTimeout(() => {
            window.location.href = userContentUrl;
          }, 1000);
        } else {
          showError(data.errorMessage || 'Authentication failed');
          resetButton($this);
        }
      })
      .catch(error => {
        showError('Server error: ' + error.message);
        resetButton($this);
      });
  } catch (error) {
    showError('Error: ' + error.message);
    resetButton($this);
  }

  // Helper function to show errors
  function showError(message) {
    statusElement.className = 'status-message error-message';
    statusElement.textContent = message;
    console.error(message);
  }

  // Helper function to reset button state
  function resetButton(button) {
    button.classList.remove('loading');
    button.textContent = 'Login with Web3 Wallet';
  }
}