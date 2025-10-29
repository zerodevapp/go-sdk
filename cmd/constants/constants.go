package constants

import "fmt"

// KernelVersion represents a supported kernel version
type KernelVersion string

const (
	KernelVersion031 KernelVersion = "0.3.1"
	KernelVersion032 KernelVersion = "0.3.2"
	KernelVersion033 KernelVersion = "0.3.3"
)

// KernelAddresses contains addresses for a specific kernel version
type KernelAddresses struct {
	AccountImplementationAddress string
	FactoryAddress               string
	MetaFactoryAddress           string
	InitCodeHash                 string
}

// KernelVersionToAddressesMap maps kernel versions to their respective addresses
var KernelVersionToAddressesMap = map[KernelVersion]KernelAddresses{
	KernelVersion031: {
		AccountImplementationAddress: "0xBAC849bB641841b44E965fB01A4Bf5F074f84b4D",
		FactoryAddress:               "0xaac5D4240AF87249B3f71BC8E4A2cae074A3E419",
		MetaFactoryAddress:           "0xd703aaE79538628d27099B8c4f621bE4CCd142d5",
		InitCodeHash:                 "0x85d96aa1c9a65886d094915d76ccae85f14027a02c1647dde659f869460f03e6",
	},
	KernelVersion032: {
		AccountImplementationAddress: "0xD830D15D3dc0C269F3dBAa0F3e8626d33CFdaBe1",
		FactoryAddress:               "0x7a1dBAB750f12a90EB1B60D2Ae3aD17D4D81EfFe",
		MetaFactoryAddress:           "0xd703aaE79538628d27099B8c4f621bE4CCd142d5",
		InitCodeHash:                 "0xc7c48c9dd12de68b8a4689b6f8c8c07b61d4d6fa4ddecdd86a6980d045fa67eb",
	},
	KernelVersion033: {
		AccountImplementationAddress: "0xd6CEDDe84be40893d153Be9d467CD6aD37875b28",
		FactoryAddress:               "0x6723b44Abeec4E71eBE3232BD5B455805baDD22f",
		MetaFactoryAddress:           "0xd703aaE79538628d27099B8c4f621bE4CCd142d5",
		InitCodeHash:                 "0xc452397f1e7518f8cea0566ac057e243bb1643f6298aba8eec8cdee78ee3b3dd",
	},
}

// GetAccountImplementationAddress returns the account implementation address for a given kernel version
func GetAccountImplementationAddress(version KernelVersion) (string, error) {
	addresses, ok := KernelVersionToAddressesMap[version]
	if !ok {
		return "", fmt.Errorf("unsupported kernel version: %s", version)
	}
	return addresses.AccountImplementationAddress, nil
}

// GetKernelAddresses returns all addresses for a given kernel version
func GetKernelAddresses(version KernelVersion) (KernelAddresses, error) {
	addresses, ok := KernelVersionToAddressesMap[version]
	if !ok {
		return KernelAddresses{}, fmt.Errorf("unsupported kernel version: %s", version)
	}
	return addresses, nil
}
