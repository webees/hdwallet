/*
 * @Author: webees@qq.com
 * @Date: 2021-03-29 18:09:50
 * @Last Modified by: webees@qq.com
 * @Last Modified time: 2021-03-29 19:05:29
 */
package trx

import (
	"github.com/webees/hdwallet/bip32"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	purposeID = bip32.HardenedKeyStart + 44
	coinID    = bip32.HardenedKeyStart + 195
	// addrID is the byte prefix of the address used in TRON addresses.
	// It's supposed to be '0xa0' for testnet, and '0x41' for mainnet.
	// But the Shasta mainteiners don't use the testnet params. So the default value is 41.
	addrID = 0x41

	tCoinID = bip32.HardenedKeyStart + 1
)

var (
	TEST           = false
	HDPrivateKeyID = [4]byte{0x04, 0x88, 0xad, 0xe4}
	HDPublicKeyID  = [4]byte{0x04, 0x88, 0xb2, 0x1e}
)

func Xpub(pvk string, accountID uint32) (string, error) {
	key, e := bip32.Xkey(pvk, purposeID, coinType(), 2147483648, 0)
	if e != nil {
		return "", e
	}
	pbk, e := key.Neuter(HDPublicKeyID)
	if e != nil {
		return "", e
	}
	return pbk.String(), nil
}

func Addr(xpub string, index uint32) (string, error) {
	c, e := bip32.Xkey(xpub, index)
	if e != nil {
		return "", e
	}
	ec, e := c.ECPubKey()
	if e != nil {
		return "", e
	}
	ecdsa := ec.ToECDSA()
	return base58.CheckEncode(crypto.PubkeyToAddress(*ecdsa).Bytes(), addrID), nil
}
