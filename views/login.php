<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title><?= $title ?? 'Web3 Wallet' ?></title>
    <meta name="description" content="Login with Web3 Wallets: MetaMask, Trust Wallet, Coinbase Wallet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://cdn.jsdelivr.net/npm/ethers@6.10.0/dist/ethers.umd.min.js"></script>
    <script src="/js/wallet.js"></script>
    <script src="/js/app.js"></script>
</head>
<body>
<div id="app">
    <div class="container text-center py-5">
        <div class="row justify-content-center">
            <div class="col-md-8 col-lg-6">
                <h1 class="mb-2">Login with Web3 Wallet</h1>
                <p class="lead mb-4">MetaMask, Trust Wallet, Coinbase Wallet</p>

                <button
                    data-user-content-url="/dashboard"
                    data-user-auth-url="/api/auth"
                    class="login-web3-btn btn btn-primary btn-lg"
                >
                    Login with Web3 Wallet
                </button>

                <div id="status-message" class="alert mt-4" style="display:none"></div>
            </div>
        </div>
    </div>
</div>
</body>
</html>
