<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Login with Web3 Wallet</title>
    <meta name="description" content="Login with Web3 Wallets: MetaMask, Trust Wallet, Coinbase Wallet ...">
    <!-- Load Web3 first -->
    <script src="https://cdn.jsdelivr.net/npm/web3@4.16.0/dist/web3.min.js"></script>
    <!-- Then load our custom scripts -->
    <script src="/web3-auth/js/sWeb3App.js"></script>
    <script src="/web3-auth/js/web3auth.js"></script>
    <link rel="stylesheet" href="/web3-auth/css/web3auth.css" type="text/css">
</head>
<body>
<div class="pageAccount">
    <h1>Login with Web3 Wallet</h1>
    <h2>Wallets: MetaMask, Trust Wallet, Coinbase Wallet ...</h2>

    <a data-user-content-url="/web3-auth/user-content.php" data-user-auth-url="/web3-auth/web3-auth.php" class="LoginWeb3Button" href="#">Login with Web3 Wallet</a>

    <div id="status-message" class="status-message"></div>
</div>
</body>
</html>