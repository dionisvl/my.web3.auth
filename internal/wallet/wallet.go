package wallet

// Service reports network metadata. Balance is computed client-side (ethers.js).
type Service struct {
	network string
}

var networkIDs = map[string]int{
	"mainnet": 1,
	"goerli":  5,
	"sepolia": 11155111,
	"avax":    43114, // Avalanche C-Chain
}

// New creates a Service for the given network; empty falls back to mainnet.
func New(network string) *Service {
	if network == "" {
		network = "mainnet"
	}
	return &Service{network: network}
}

type Details struct {
	Address   string `json:"address"`
	Network   string `json:"network"`
	NetworkID int    `json:"networkId"`
}

func (s *Service) GetWalletDetails(address string) Details {
	id, ok := networkIDs[s.network]
	if !ok {
		id = 1
	}
	return Details{Address: address, Network: s.network, NetworkID: id}
}

type APIConfig struct {
	Network string `json:"network"`
}

func (s *Service) GetAPIConfig() APIConfig {
	return APIConfig{Network: s.network}
}
