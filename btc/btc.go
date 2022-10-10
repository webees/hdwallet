/*
 * @Author: webees@qq.com
 * @Date: 2021-03-29 18:10:23
 * @Last Modified by:   webees@qq.com
 * @Last Modified time: 2021-03-29 18:10:23
 */
package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/webees/hdwallet/bip32"
	"github.com/webees/hdwallet/btc/txauthor"
	"github.com/webees/hdwallet/crypto/hash"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/gogf/gf/util/gconv"
)

const (
	purposeID = bip32.HardenedKeyStart + 49
	coinID    = bip32.HardenedKeyStart + 0
	addrID    = 0x05

	tCoinID = bip32.HardenedKeyStart + 1
	tAddrID = 0xC4

	// MinNondustOutput Any standard (ie P2PKH) output smaller than this value (in satoshis) will most likely be rejected by the network.
	// This is calculated by assuming a standard output will be 34 bytes
	minNondustOutput = 546        // satoshis
	omniHex          = "6f6d6e69" // Hex-encoded: "omni"
	usdtID           = 31
)

var (
	TEST               = false
	TestHDPrivateKeyID = [4]byte{0x04, 0x4A, 0x4E, 0x28}
	TestHDPublicKeyID  = [4]byte{0x04, 0x4A, 0x52, 0x62}
	HDPrivateKeyID     = [4]byte{0x04, 0x9D, 0x78, 0x78}
	HDPublicKeyID      = [4]byte{0x04, 0x9D, 0x7C, 0xB2}
)

type TxIn struct {
	Txid  string
	Vout  uint32
	Addr  string
	Value string
}

type TxOut struct {
	Addr  string
	Value string
}

func Xpub(pvk string, accountID uint32) (string, error) {
	accountID = bip32.HardenedKeyStart + accountID
	key, e := bip32.Xkey(pvk, purposeID, coinType(), accountID)
	if e != nil {
		return "", e
	}
	pbk, e := key.Neuter(pbkType())
	if e != nil {
		return "", e
	}
	return pbk.String(), nil
}

func Addr(xpub string, change uint32, index uint32) (string, error) {
	c, e := bip32.Xkey(xpub, change, index)
	if e != nil {
		return "", e
	}
	scriptHash := hash.Hash160(c.PubKeyBytes())
	var buf1 bytes.Buffer
	buf1.WriteByte(0x00)
	buf1.WriteByte(uint8(len(scriptHash)))
	buf1.Write(scriptHash)
	scriptSig := buf1.Bytes()
	addressHash := hash.Hash160(scriptSig)
	var buf2 bytes.Buffer
	buf2.WriteByte(addrType())
	buf2.Write(addressHash)
	addressBytes := buf2.Bytes()
	checksum := hash.Sha256d(addressBytes)[:4]
	var b bytes.Buffer
	b.Write(addressBytes)
	b.Write(checksum)
	return base58.Encode(b.Bytes()), nil
}

func UnsignedRaw(inputs []*TxIn, outputs []*TxOut, changeAddress string, relayFeePerKb int) (string, error) {
	wireOuts := make([]*wire.TxOut, 0, len(outputs)) // 准备输出
	for _, v := range outputs {
		pkScript, e := payScript(v.Addr)
		if e != nil {
			return "", e
		}
		to := wire.NewTxOut(gconv.Int64(v.Value), pkScript)
		wireOuts = append(wireOuts, to)
	}
	authoredTx, e := txauthor.NewUnsignedTransaction(wireOuts, btcutil.Amount(relayFeePerKb), fetchInputs(inputs), fetchChange(changeAddress)) // 创建交易数据
	if e != nil {
		return "", e
	}
	j, e := json.Marshal(authoredTx)
	if e != nil {
		return "", e
	}
	return base58.Encode(j), nil
}

func OminiUnsignedRaw(propertyID int, from string, to string, amount uint64, inputs []*TxIn, changeAddress string, relayFeePerKb int) (string, error) {
	wireOuts := make([]*wire.TxOut, 0, 3) // 准备输出，仅3个
	payload := fmt.Sprintf("%v%016x%016x", omniHex, propertyID, amount)
	b, e := hex.DecodeString(payload)
	if e != nil {
		return "", e
	}
	opreturnScript, e := txscript.NullDataScript(b)
	if e != nil {
		return "", e
	}
	wireOuts = append(wireOuts, wire.NewTxOut(0, opreturnScript)) // OP_RETURN
	pkScript, e := payScript(to)
	if e != nil {
		return "", e
	}
	wireOuts = append(wireOuts, wire.NewTxOut(minNondustOutput, pkScript))                                                                     // OMINI 接收地址
	authoredTx, e := txauthor.NewUnsignedTransaction(wireOuts, btcutil.Amount(relayFeePerKb), fetchInputs(inputs), fetchChange(changeAddress)) // 创建交易数据
	if e != nil {
		return "", e
	}
	j, e := json.Marshal(authoredTx)
	if e != nil {
		return "", e
	}
	return base58.Encode(j), nil
}

func SignedRaw(b58raw string, xpvk string) (string, error) {
	b := base58.Decode(b58raw)
	authoredTx := txauthor.AuthoredTx{}
	if e := json.Unmarshal(b, &authoredTx); e != nil {
		return "", e
	}
	if e := authoredTx.AddAllInputScripts(&secretSource{
		xpvk: xpvk,
	}); e != nil { // 签名
		return "", e
	}
	var buf bytes.Buffer
	if e := authoredTx.Tx.Serialize(&buf); e != nil { // 序列化
		return "", e
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func IsAddr(addr string) bool {
	if _, e := btcutil.DecodeAddress(addr, chainType()); e != nil {
		return false
	}
	return true
}
