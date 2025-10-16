package service

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_verifySignature(t *testing.T) {
	// 固定私钥（仅用于测试！不要在生产用）
	privKeyHex := "59c6995e998f97a5a0044978f6f2d7f6b9a6f9b7f8f4f6b9a8d7f6b5f4e3c2b1"
	privateKey, err := crypto.HexToECDSA(privKeyHex)
	assert.NoError(t, err)

	// 从私钥生成以太坊地址
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey).Hex()

	// 构造 SIWE 消息（简化版）
	// issuedAt := time.Now().UTC().Format(time.RFC3339)
	//	message := fmt.Sprintf(`example.com wants you to sign in with your Ethereum account:
	//%s
	//
	//Sign in with Ethereum to the app.
	//
	//URI: https://example.com
	//Version: 1
	//Chain ID: 1
	//Nonce: test-nonce-1234
	//Issued At: %s`, address, issuedAt)

	message := fmt.Sprintf(`account: %s`, address)

	fmt.Println("SIWE message:\n", message)

	// 按 EIP-191 规则加前缀
	prefixedMsg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	fmt.Println(prefixedMsg)
	// 对消息做 keccak256
	hash := crypto.Keccak256Hash([]byte(prefixedMsg))

	// 签名
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	assert.NoError(t, err)

	// 注意: go-ethereum 生成的 signature[64] 是 v 值(0或1)，ethers.js 需要 +27 才能用
	signature[64] += 27

	fmt.Println("\nsignature:", hexutil.Encode(signature))

	err = verifySignature(message, hexutil.Encode(signature))
	assert.NoError(t, err)
}
