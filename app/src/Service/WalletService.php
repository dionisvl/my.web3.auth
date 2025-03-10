<?php

declare(strict_types=1);

namespace Web3Auth\Service;

class WalletService
{
    private string $network;

    public function __construct(array $config = [])
    {
        $this->network = $config['network'] ?? 'mainnet';
    }

    public function getWalletDetails(string $address): array
    {
        $networks = [
            'mainnet' => 1,
            'goerli' => 5,
            'sepolia' => 11155111,
            'avax' => 43114, // Avalanche C-Chain
        ];

        return [
            'address' => $address,
            'network' => $this->network,
            'networkId' => $networks[$this->network] ?? 1
        ];
    }

    public function getApiConfig(): array
    {
        return [
            'network' => $this->network
        ];
    }
}
