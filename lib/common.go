package lib

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/crypto/sha3"
	"golang.org/x/xerrors"
)

const USERADDRESS = "0x7A2994d6e1b8E5889Cc96664cF9873BFf33F8deD"
const ENDPOINT = "https://devchain.metamemo.one:8501"

var (
	AccessKey string
	SecretKey string
	EndPoint  string
)

func New() (*miniogo.Client, error) {
	optionsStaticCreds := &miniogo.Options{
		Creds:        credentials.NewStaticV4(AccessKey, SecretKey, ""),
		Secure:       false,
		BucketLookup: miniogo.BucketLookupAuto,
	}

	client, err := miniogo.New(EndPoint, optionsStaticCreds)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func Signmsg(ctx context.Context, sk string, value big.Int, destaddr string, filemd5 string) (*types.Transaction, error) {
	client, err := ethclient.Dial(ENDPOINT)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(sk)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, xerrors.New("key type not right")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println("sender address: ", fromAddress)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, err
	}

	sendvalue := big.NewInt(0)

	log.Println("getting gasPrice..")
	gasPrice := big.NewInt(1000)
	log.Println("gasprice: ", gasPrice)

	toAddress := common.HexToAddress(destaddr)
	log.Println("to address:", toAddress)
	tokenAddress := common.HexToAddress("0xa678920287Eac5b8e81578268469AA2BaaD2eC87")
	// tokenAddress := common.HexToAddress("0xf3783070Ffe8eDd3C7F89bc136ba7c0512F18627")
	log.Println("erc20 token addr:", tokenAddress)

	log.Println("getting methodID from fnSig..")
	transferFnSignature := []byte("transfer(address,uint256)")

	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	log.Println(hexutil.Encode(methodID))

	log.Println("padding address..")
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	log.Println(hexutil.Encode(paddedAddress))

	log.Println("padding amount..")
	amount := &value

	// pad amount to 32 bytes
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	log.Println(hexutil.Encode(paddedAmount))

	fmt.Println("constructing tx data..")
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	data = append(data, []byte(filemd5)...)
	log.Println("md5 hex", hexutil.Encode([]byte(filemd5)))
	log.Println("tx data: ", hexutil.Encode(data))

	gasLimit := uint64(300000)

	// construct tx with all parts
	log.Println("constructing tx..")
	tx := types.NewTransaction(nonce, tokenAddress, sendvalue, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, err
	}
	log.Println("chainID: ", chainID)

	log.Println("signing tx..")
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}

	// log.Println("sending tx..")
	// err = client.SendTransaction(ctx, signedTx)
	// if err != nil {
	// 	return "", err
	// }

	// log.Printf("tx sent: %s\n", signedTx.Hash().Hex())

	// WaitTx(endpoint, signedTx.Hash())
	return signedTx, nil
}

func Approve(ctx context.Context, sk string, tokenaddr, destaddr common.Address, value *big.Int) (*types.Transaction, error) {
	client, err := ethclient.Dial(ENDPOINT)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(sk)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, xerrors.New("key type not right")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, err
	}

	sendvalue := big.NewInt(0)
	gasPrice := big.NewInt(1000)

	log.Println("erc20 token addr:", tokenaddr)

	approveFnSignature := []byte("approve(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(approveFnSignature)
	methodID := hash.Sum(nil)[:4] 

	log.Println("padding address..")
	paddedAddress := common.LeftPadBytes(destaddr.Bytes(), 32)

	paddedAmount := common.LeftPadBytes(value.Bytes(), 32)

	log.Println("constructing tx data..")
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit := uint64(300000)
	tx := types.NewTransaction(nonce, tokenaddr, sendvalue, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, err
	}
	log.Println("chainID: ", chainID)

	log.Println("signing tx..")
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// func getTransactionReceipt(endPoint string, hash common.Hash) *types.Receipt {
// 	client, err := ethclient.Dial(endPoint)
// 	if err != nil {
// 		return nil
// 	}
// 	defer client.Close()
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
// 	defer cancel()
// 	receipt, err := client.TransactionReceipt(ctx, hash)
// 	if err != nil {
// 		log.Printf("get transaction %s receipt fail: %s", hash, err)
// 	}
// 	tx, flag, err := client.TransactionByHash(context.Background(), hash)
// 	if err != nil {
// 		return nil
// 	}
// 	log.Printf("tx data:%x\n", tx.Data())
// 	log.Println(flag)
// 	return receipt
// }

// func WaitTx(ep string, txHash common.Hash) {
// 	log.Println("tx hash:", txHash)

// 	log.Println("waiting tx complete...")

// 	time.Sleep(30 * time.Second)
// 	receipt := getTransactionReceipt(ep, txHash)
// 	log.Println("tx status:", receipt.Status)
// 	log.Println(receipt.Logs[0].Index, receipt.Logs[0].BlockNumber)
// }
