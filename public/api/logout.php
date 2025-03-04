<?php
declare(strict_types=1);

require_once __DIR__ . '/../../vendor/autoload.php';

// Start session
if (session_status() === PHP_SESSION_NONE) {
    session_start();
}

use Web3Auth\Auth;

$auth = new Auth();
$auth->logout();

header('Location: /');
