<?php

declare(strict_types=1);

namespace Web3Auth\Service;

class AuthService
{
    public function authenticate(array $data): array
    {
        if (empty($data['wallet']) || empty($data['message']) || empty($data['signature'])) {
            return ['error' => 1, 'errorMessage' => 'Missing required parameters'];
        }

        if (!preg_match('/^0x[a-fA-F0-9]{40}$/', $data['wallet'])) {
            return ['error' => 1, 'errorMessage' => 'Invalid wallet address format'];
        }

        // TODO: Realize sign checking with correct library

        $_SESSION['wallet'] = $data['wallet'];
        $_SESSION['login_time'] = time();

        return ['error' => 0];
    }

    public function isAuthenticated(): bool
    {
        return !empty($_SESSION['wallet']);
    }

    public function getWallet(): ?string
    {
        return $_SESSION['wallet'] ?? null;
    }

    public function logout(): void
    {
        unset($_SESSION['wallet'], $_SESSION['login_time']);
    }
}
