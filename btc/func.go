/*
 * @Author: webees@qq.com
 * @Date: 2021-03-29 18:10:26
 * @Last Modified by: webees@qq.com
 * @Last Modified time: 2021-03-29 18:18:31
 */
package btc

import (
	"github.com/webees/hdwallet/bip32"
	"github.com/webees/hdwallet/btc/txauthor"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/gogf/gf/util/gconv"
)

func addrScript(address string) (btcutil.Address, error) {
	addr, e := btcutil.DecodeAddress(address, chainType())
	if e != nil {
		return nil, e
	}
	return addr, nil
}

func payScript(address string) ([]byte, error) {
	addr, e := addrScript(address)
	if e != nil {
		return nil, e
	}
	pkScript, e := txscript.PayToAddrScript(addr)
	if e != nil {
		return nil, e
	}
	return pkScript, nil
}

// 准备找零地址脚本
func fetchChange(address string) txauthor.ChangeSource {
	f := func() ([]byte, error) {
		script, e := payScript(address)
		if e != nil {
			return nil, e
		}
		return script, nil
	}
	return txauthor.ChangeSource(f)
}

// Current inputs and their total value.  These are closed over by the returned input source and reused across multiple calls.
func fetchInputs(inputs []*TxIn) txauthor.InputSource {
	currentTotal := btcutil.Amount(0)
	currentInputs := make([]*wire.TxIn, 0, len(inputs))
	currentScripts := make([][]byte, 0, len(inputs))
	currentInputValues := make([]btcutil.Amount, 0, len(inputs))
	f := func(target btcutil.Amount) (btcutil.Amount, []*wire.TxIn, []btcutil.Amount, [][]byte, error) {
		for currentTotal < target && len(inputs) != 0 {
			u := inputs[0]
			inputs = inputs[1:]
			hash, e := chainhash.NewHashFromStr(u.Txid)
			if e != nil {
				return 0, nil, nil, nil, e
			}
			index := u.Vout
			prevOut := wire.NewOutPoint(hash, index)
			pkScript, e := payScript(u.Addr)
			if e != nil {
				return 0, nil, nil, nil, e
			}
			nextInput := wire.NewTxIn(prevOut, pkScript, nil)
			currentTotal += btcutil.Amount(gconv.Int64(u.Value))
			currentInputs = append(currentInputs, nextInput)
			currentScripts = append(currentScripts, pkScript)
			currentInputValues = append(currentInputValues, btcutil.Amount(gconv.Int64(u.Value)))
		}
		return currentTotal, currentInputs, currentInputValues, currentScripts, nil
	}
	return txauthor.InputSource(f)
}

type secretSource struct {
	xpvk string // 扩展私钥
}

// 通过地址遍历获取私钥
func (t *secretSource) GetKey(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
	var i uint32 = 0
	pbk, e := Xpub(t.xpvk, 0) // 从私钥扩展公钥
	if e != nil {
		return nil, false, e
	}
	for ; i < 9999; i++ { // TODO: 地址索引，需定期聚币，控制循环次数
		s, e := Addr(pbk, 0, i)
		if e != nil {
			return nil, false, e
		}
		if addr.EncodeAddress() == s {
			break
		}
	}
	var cid uint32
	if TEST {
		cid = tCoinID
	} else {
		cid = coinID
	}
	xkey, e := bip32.Xkey(t.xpvk, purposeID, cid, 2147483648, 0, i) // 默认账户0
	if e != nil {
		return nil, false, e
	}
	privKey, e := xkey.ECPrivKey()
	if e != nil {
		return nil, false, e
	}
	return privKey, true, nil
}

func (t *secretSource) GetScript(addr btcutil.Address) ([]byte, error) {
	return nil, nil // 暂未使用
}

func (t *secretSource) ChainParams() *chaincfg.Params {
	return chainType()
}

func chainType() *chaincfg.Params {
	if TEST {
		return &chaincfg.TestNet3Params
	} else {
		return &chaincfg.MainNetParams
	}
}

func coinType() uint32 {
	if TEST {
		return tCoinID
	} else {
		return coinID
	}
}

func addrType() uint8 {
	if TEST {
		return tAddrID
	} else {
		return addrID
	}
}

func pbkType() [4]byte {
	if TEST {
		return TestHDPublicKeyID
	} else {
		return HDPublicKeyID
	}
}
