package utils

import "bip32_threshold_wallet/node"

func CreateDevices(t, n uint32) []node.Device {
	pkShares, skShares, pk := node.GenSharedKey(t, n)
	chaincode, _ := node.NewMasterChainCode()
	index := uint32(0x0)

	devices := make([]node.Device, n)
	for i := uint32(0); i < n; i++ {
		device, _ := node.NewDevice(
			int(i),
			pkShares[i+1].Point,
			skShares[i+1].ShamirShare,
			pk,
			index,
			chaincode,
			nil,
		)
		devices[i] = device
	}
	return devices
}
