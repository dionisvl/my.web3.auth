<?php

declare(strict_types=1);

namespace Web3Auth\Service;

use Elliptic\EC;
use kornrunner\Keccak;

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

        // Verify the signature matches the wallet address
        try {
            $isValid = $this->verifySignature(
                $data['message'],
                $data['signature'],
                $data['wallet']
            );

            if (!$isValid) {
                return ['error' => 1, 'errorMessage' => 'Invalid signature'];
            }
        } catch (\Exception $e) {
            return ['error' => 1, 'errorMessage' => 'Signature verification error: ' . $e->getMessage()];
        }

        $_SESSION['wallet'] = $data['wallet'];
        $_SESSION['login_time'] = time();
        $_SESSION['auth_method'] = 'web3_signature';

        $_SESSION['last_signature'] = [
            'message' => $data['message'],
            'signature' => $data['signature'],
            'timestamp' => time()
        ];

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

    public function getLastSignature(): ?array
    {
        return $_SESSION['last_signature'] ?? null;
    }

    public function logout(): void
    {
        unset($_SESSION['wallet'], $_SESSION['login_time'], $_SESSION['auth_method'], $_SESSION['last_signature']);
    }

    private function verifySignature(string $message, string $signature, string $address): bool
    {
        $address = strtolower(str_replace('0x', '', $address));
        $signature = str_replace('0x', '', $signature);

        // Check signature length
        if (strlen($signature) !== 130) {
            return false;
        }

        // Split signature parts
        $r = substr($signature, 0, 64);
        $s = substr($signature, 64, 64);
        $v = ord(hex2bin(substr($signature, 128, 2)));

        // Convert recovery param
        if ($v < 27) {
            $v += 27;
        }
        $v -= 27;

        // Create message hash
        $messagePrefix = "\x19Ethereum Signed Message:\n" . strlen($message);
        $messageHash = Keccak::hash($messagePrefix . $message, 256);

        // Recover the public key
        $ec = new EC('secp256k1');
        $publicKey = $ec->recoverPubKey($messageHash, [
            'r' => $r,
            's' => $s
        ], $v);

        // Convert public key to address
        $publicKeyHex = $publicKey->encode('hex');
        $publicKeyBytes = substr(hex2bin($publicKeyHex), 1);
        $recoveredAddress = strtolower(substr(Keccak::hash($publicKeyBytes, 256), -40));

        // Compare addresses
        return $address === $recoveredAddress;
    }
}
