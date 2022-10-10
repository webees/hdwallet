/*
 * @Author: webees@qq.com
 * @Date: 2021-03-29 18:10:18
 * @Last Modified by:   webees@qq.com
 * @Last Modified time: 2021-03-29 18:10:18
 */
package btc

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/webees/hdwallet/bip32"
	"github.com/webees/hdwallet/bip39"

	"github.com/btcsuite/btcd/txscript"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/glog"
)

func TestBtc(t *testing.T) {
	fmt.Println("\n######################################## BTC ########################################")
	mnemonic := "owner mosquito uphold squirrel utility fat warrior wheat vital chapter shoulder horn"
	seed, _ := bip39.NewSeed(mnemonic, "")
	pvk, _ := bip32.NewMaster(seed, HDPrivateKeyID)
	pbk, _ := Xpub(pvk.String(), 0)
	fmt.Println("\n助记词：  ", mnemonic)
	fmt.Println("扩展私钥: ", pvk.String())
	fmt.Println("账户公钥：", pbk)
	addr := gmap.New(true)
	start := time.Now()
	for i := 0; i < 9; i++ {
		s, _ := Addr(pbk, 0, uint32(i))
		addr.Set(i, s)
		fmt.Println("地址", i, "=", s)
	}
	elapsed := time.Since(start)
	fmt.Println("\n耗时：", elapsed)
}

func TestTbtc(t *testing.T) {
	TEST = true
	fmt.Println("\n######################################## tBTC ########################################")
	mnemonic := "owner mosquito uphold squirrel utility fat warrior wheat vital chapter shoulder horn"
	seed, _ := bip39.NewSeed(mnemonic, "")
	pvk, _ := bip32.NewMaster(seed, TestHDPrivateKeyID)
	pbk, _ := Xpub(pvk.String(), 0)
	fmt.Println("\n助记词：  ", mnemonic)
	fmt.Println("扩展私钥: ", pvk.String())
	fmt.Println("账户公钥：", pbk)
	addr := gmap.New(true)
	start := time.Now()
	for i := 0; i < 9; i++ {
		s, _ := Addr(pbk, 0, uint32(i))
		addr.Set(i, s)
		fmt.Println("地址", i, "=", s)
	}
	elapsed := time.Since(start)
	fmt.Println("\n耗时：", elapsed)
}

func TestUnsignedRaw(t *testing.T) {
	TEST = true
	fmt.Println("\n######################################## UnsignedRaw ########################################")
	txIn := make([]*TxIn, 0, 1)
	txOut := make([]*TxOut, 0, 1)
	txIn = append(txIn, &TxIn{
		Txid:  "4c554ad676a5a68a735f3badf1a01a5fd6271f561b59537c67cef9232b683d40",
		Vout:  0,
		Addr:  "2NCNqcThPY7ntaqrGR48YfgWb9u47tspL2b",
		Value: "4000000",
	})
	txOut = append(txOut, &TxOut{
		Addr:  "2NCNqcThPY7ntaqrGR48YfgWb9u47tspL2b",
		Value: "2000000",
	})
	txOut = append(txOut, &TxOut{
		Addr:  "2N9icdqaUQgZgeiV36P81h4HB7mHy3PP79F",
		Value: "1000000",
	})
	unsigned, e := UnsignedRaw(txIn, txOut, "2NCNqcThPY7ntaqrGR48YfgWb9u47tspL2b", 1400)
	if e != nil {
		glog.Error(e)
	}
	fmt.Println("待签名数据：", unsigned)
}

func TestSignedRaw(t *testing.T) {
	TEST = true
	fmt.Println("\n######################################## SignedRaw ########################################")
	mnemonic := "owner mosquito uphold squirrel utility fat warrior wheat vital chapter shoulder horn"
	seed, _ := bip39.NewSeed(mnemonic, "")
	pvk, _ := bip32.NewMaster(seed, HDPrivateKeyID)
	b58raw := "oaje7NM9BKwUzQkLenDxsWbVD7hKdhyT3amJok1Mf5VmTiSngsnzRfnbLUmfscvCUkaN8fm52W3oE59n7Rt2VR2Ufqj6x4PZWjW2BPBRXuL7WfG7Vv5WqiwcMNRrBuc9ZJY8fqcU6C3FhDPYpUz3FtKVvFJPffFCN5kzMWAuqGeZLvgeX64MbV51VT6G22mjGKksg6Soc6UwDyA7XCGokTXVhno6Q5r4Ka6qaDqKc6dNABuLe1mws9jCNaeBNrdfunh4MgopwhJscrF4qRfAvwtAnZys4kP29ZsGKpxbFS4UwzuS8pAsWSqeeKx2dLAMPQDUadopxtbrkBxNurBpgVdznZTRaDU6zESh5mH3PN1vHEEHSMDEdkyKtUVqcdshSPg2h52Z9a46GvhEDUaRK1n1JGefEiAH7X5YDc98b45xdA2B8FpD3v3Y2Q3CLJvFbREvwR9Z7sQYtnUm2HBQfD7z9x7dT7wkMB9MV2AFFq7k7hf2yDEG9sNphtxv473XPeyEZZdgEdBwGkSrzHW9CNC8sbbYJPKM9X1unUN2L3tvNPep6aFxkwL24Yn4bVbMACzL7xCfZkLeLoCny13T5ntn4NBmB3bpgozj8ZKLDsYVcACkVufKeQwTeSiLkAAyWQCnUfH9iD3nPL9pJnJ2F34xFoAYega8VMHFfypDFgmf2Q3KtitFtjZK7oy2wJTVJjKUiELYGARCUALBKPWNeKe4VHq9FqAk85EkBi3sPoSHArCbXxhnJxYverPiC1yt1qZYY7SmjxcaxRi3Yi4Qyoc14mrhRia1RuJ8N7uqRan6BKZEFvQSmS"
	signed, e := SignedRaw(b58raw, pvk.String())
	if e != nil {
		glog.Error(e)
	}
	fmt.Println("已签名数据：", signed)
}

func TestOmini(t *testing.T) {
	fmt.Println("\n######################################## Omini ########################################")
	payload := fmt.Sprintf("%v%016x%016x", omniHex, usdtID, int64(999))
	b, e := hex.DecodeString(payload)
	if e != nil {
		return
	}
	opreturnScript, e := txscript.NullDataScript(b)
	if e != nil {
		return
	}
	fmt.Println(hex.EncodeToString(opreturnScript))
}

func TestIsAddr(t *testing.T) {
	fmt.Println("\n######################################## IsAddr ########################################")
	TEST = false
	fmt.Println("123456789                                   =", IsAddr("123456789"))
	fmt.Println("bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq  =", IsAddr("bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq"))
	fmt.Println("3HanbdYe9B2uLeQwYC8L2eENmm7CKM8dQG          =", IsAddr("3HanbdYe9B2uLeQwYC8L2eENmm7CKM8dQG"))
	fmt.Println("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa          =", IsAddr("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"))
	TEST = true
	fmt.Println("2MyH38gFs4zToCsjaSWi1TcHS5d8HQddwfX         =", IsAddr("2MyH38gFs4zToCsjaSWi1TcHS5d8HQddwfX"))
}
