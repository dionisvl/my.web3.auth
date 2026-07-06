package wallet

// Service resolves wallet/network metadata. It mirrors the PHP WalletService:
// balance is computed client-side (ethers.js), so the server only reports the
// configured network and its chain id.
type Service struct {
	network string
}

// networkIDs maps a network name to its Ethereum chain id.
var networkIDs = map[string]int{
	"mainnet": 1,
	"goerli":  5,
	"sepolia": 11155111,
	"avax":    43114, // Avalanche C-Chain
}

// New creates a Service for the given network. Empty falls back to mainnet,
// matching the original PHP default.
func New(network string) *Service {
	if network == "" {
		network = "mainnet"
	}
	return &Service{network: network}
}

// Details describes the wallet's network context for the templates/JSON API.
type Details struct {
	Address   string `json:"address"`
	Network   string `json:"network"`
	NetworkID int    `json:"networkId"`
}

// GetWalletDetails returns network details for an address.
func (s *Service) GetWalletDetails(address string) Details {
	id, ok := networkIDs[s.network]
	if !ok {
		id = 1
	}
	return Details{
		Address:   address,
		Network:   s.network,
		NetworkID: id,
	}
}

// APIConfig is the minimal client-facing config (the network name).
type APIConfig struct {
	Network string `json:"network"`
}

// GetAPIConfig returns the client API config.
func (s *Service) GetAPIConfig() APIConfig {
	return APIConfig{Network: s.network}
}
