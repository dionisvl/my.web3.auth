<?php

declare(strict_types=1);

use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Slim\Factory\AppFactory;
use Slim\Views\PhpRenderer;
use Web3Auth\Service\AuthService;
use Web3Auth\Service\WalletService;

require_once __DIR__ . '/../vendor/autoload.php';

// Load environment variables
$dotenv = Dotenv\Dotenv::createImmutable(__DIR__ . '/..');
$dotenv->safeLoad();

// Start session
session_start();

// Initialize Slim app
$app = AppFactory::create();

// Add routing middleware
$app->addRoutingMiddleware();
$app->addErrorMiddleware(true, true, true);

// Set up view renderer
$renderer = new PhpRenderer(__DIR__ . '/../views');

// Set up services
$authService = new AuthService();

$walletService = new WalletService([
        'network' => $_ENV['ETH_NETWORK'] ?? 'sepolia',
]);

// Define routes
$app->get('/', function (Request $request, Response $response) use ($renderer, $authService) {
    if ($authService->isAuthenticated()) {
        return $response->withHeader('Location', '/dashboard')->withStatus(302);
    }

    return $renderer->render($response, 'login.php', [
            'title' => 'Login with Web3 Wallet',
            'network' => $_ENV['ETH_NETWORK'] ?? 'sepolia'
    ]);
});

$app->get('/dashboard', function (Request $request, Response $response) use ($renderer, $authService, $walletService) {
    if (!$authService->isAuthenticated()) {
        return $response->withHeader('Location', '/')->withStatus(302);
    }

    $wallet = $authService->getWallet();
    $walletDetails = $walletService->getWalletDetails($wallet);

    return $renderer->render($response, 'dashboard.php', [
            'title' => 'Web3 Wallet Dashboard',
            'wallet' => $wallet,
            'walletDetails' => $walletDetails,
            'apiConfig' => $walletService->getApiConfig()
    ]);
});

$app->post('/api/auth', function (Request $request, Response $response) use ($authService) {
    $params = (array)$request->getParsedBody();

    $result = $authService->authenticate([
            'wallet' => $params['wallet'] ?? '',
            'message' => $params['message'] ?? '',
            'signature' => $params['signature'] ?? ''
    ]);

    $response->getBody()->write(json_encode($result));
    return $response->withHeader('Content-Type', 'application/json');
});

$app->get('/api/wallet', function (Request $request, Response $response) use ($authService, $walletService) {
    if (!$authService->isAuthenticated()) {
        $response->getBody()->write(json_encode([
                'error' => 1,
                'errorMessage' => 'Not authenticated'
        ]));
        return $response->withHeader('Content-Type', 'application/json')->withStatus(401);
    }

    $wallet = $authService->getWallet();
    $walletDetails = $walletService->getWalletDetails($wallet);

    $response->getBody()->write(json_encode([
            'error' => 0,
            'wallet' => $wallet,
            'walletDetails' => $walletDetails,
            'apiConfig' => $walletService->getApiConfig()
    ]));

    return $response->withHeader('Content-Type', 'application/json');
});

$app->get('/api/logout', function (Request $request, Response $response) use ($authService) {
    $authService->logout();
    return $response->withHeader('Location', '/')->withStatus(302);
});

// Run the app
$app->run();
