<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title><?= $title ?? 'Web3 Dashboard' ?></title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://cdn.jsdelivr.net/npm/ethers@6.10.0/dist/ethers.umd.min.js"></script>
    <script src="/js/wallet.js"></script>
    <script src="/js/app.js"></script>
</head>
<body>
<div id="app"
     data-network="<?= $apiConfig['network'] ?? 'mainnet' ?>">

    <div class="container py-5" id="wallet-dashboard">
        <div class="row">
            <div class="col-md-12 mb-4">
                <div class="card">
                    <div class="card-header d-flex justify-content-between align-items-center">
                        <h3 class="m-0">Web3 Wallet Dashboard</h3>
                        <a href="<?= $config['routes']['logout'] ?>" class="btn btn-outline-secondary btn-sm">Logout</a>
                    </div>
                    <div class="card-body">
                        <div class="row">
                            <div class="col-md-6">
                                <div class="mb-3">
                                    <h5>Wallet Address</h5>
                                    <code id="wallet-address" class="d-block p-2 bg-light"><?= htmlspecialchars($wallet) ?></code>
                                </div>
                                <div class="mb-3">
                                    <h5>Network</h5>
                                    <p id="network-name"><?= htmlspecialchars($walletDetails['network']) ?></p>
                                </div>
                                <div class="mb-3">
                                    <h5>Balance</h5>
                                    <div class="d-flex align-items-center">
                                        <p id="wallet-balance" class="me-3">Loading...</p>
                                        <button id="refresh-balance" class="btn btn-sm btn-primary">Refresh</button>
                                    </div>
                                </div>
                            </div>
                            <div class="col-md-6">
                                <div class="card">
                                    <div class="card-header">Transaction History</div>
                                    <div class="card-body">
                                        <p class="text-muted">Connect your wallet to view transactions</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div id="status-message" class="alert mt-4" style="display:none"></div>
    </div>
</div>
</body>
</html>
