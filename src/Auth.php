<?php

declare(strict_types=1);

namespace Web3Auth;

class Auth
{
    /**
     * Authenticate user with Web3 wallet
     */
    public function authenticate(array $data): array
    {
        // Validate required fields
        if (empty($data['wallet']) || empty($data['message']) || empty($data['signature'])) {
            return ['error' => 1, 'errorMessage' => 'Missing required parameters'];
        }

        // Validate wallet address format
        if (!preg_match('/^0x[a-fA-F0-9]{40}$/', $data['wallet'])) {
            return ['error' => 1, 'errorMessage' => 'Invalid wallet address format'];
        }

        /*
         * TODO: Server-side signature verification
         * For production, implement verification with:
         * - web3.php library (https://github.com/web3p/web3.php)
         * - or kornrunner/ethereum-sign-verify package
         */

        // Store wallet address in session
        $_SESSION['wallet'] = $data['wallet'];
        $_SESSION['login_time'] = time();

        return ['error' => 0];
    }

    /**
     * Check if user is authenticated
     */
    public function isAuthenticated(): bool
    {
        return !empty($_SESSION['wallet']);
    }

    /**
     * Get authenticated wallet address
     */
    public function getWallet(): ?string
    {
        return $_SESSION['wallet'] ?? null;
    }

    /**
     * Logout user
     */
    public function logout(): void
    {
        if (!empty($_SESSION['wallet'])) {
            unset($_SESSION['wallet']);
            unset($_SESSION['login_time']);
        }
    }
}
