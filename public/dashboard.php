<?php
declare(strict_types=1);

require_once __DIR__ . '/../vendor/autoload.php';

// Start session
if (session_status() === PHP_SESSION_NONE) {
    session_start();
}

use Web3Auth\Auth;

$auth = new Auth();

if (!$auth->isAuthenticated()) {
    header('Location: /');
    exit;
}

$wallet = $auth->getWallet();
?>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Web3 Dashboard</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
<div class="container py-5">
    <div class="card">
        <div class="card-header">
            <h3>Web3 Dashboard</h3>
        </div>
        <div class="card-body">
            <p class="card-text">Wallet Address: <code><?= htmlspecialchars($wallet) ?></code></p>
            <p class="card-text">User Content Area</p>
        </div>
        <div class="card-footer">
            <a href="/api/logout.php" class="btn btn-secondary">Logout</a>
        </div>
    </div>
</div>
</body>
</html>
