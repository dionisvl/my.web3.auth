<?php
declare(strict_types=1);

require_once __DIR__ . '/../../vendor/autoload.php';

// Start session
if (session_status() === PHP_SESSION_NONE) {
    session_start();
}

use Web3Auth\Auth;

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    exit(json_encode(['error' => 1, 'errorMessage' => 'Invalid request method']));
}

// Check referer to prevent CSRF
$referer = $_SERVER['HTTP_REFERER'] ?? '';
if ($referer && strpos($referer, $_SERVER['HTTP_HOST']) === false) {
    exit(json_encode(['error' => 1, 'errorMessage' => 'Invalid referer']));
}

$auth = new Auth();
$result = $auth->authenticate([
    'wallet' => $_POST['wallet'] ?? '',
    'message' => $_POST['message'] ?? '',
    'signature' => $_POST['signature'] ?? ''
]);

header('Content-Type: application/json');
echo json_encode($result);
