package main

import (
	accountTypes"github.com/drep-project/drep-chain/accounts/types"
	"github.com/drep-project/drep-chain/log"
	p2pTypes "github.com/drep-project/drep-chain/network/types"
	chainTypes  "github.com/drep-project/drep-chain/chain/types"
consensusTypes  "github.com/drep-project/drep-chain/consensus/types"
rpcTypes  "github.com/drep-project/drep-chain/rpc/types"
	"github.com/drep-project/drep-chain/common"
	"github.com/drep-project/drep-chain/crypto"
	"github.com/drep-project/drep-chain/crypto/secp256k1"
	"github.com/drep-project/drep-chain/crypto/sha3"
	"BlockChainTest/util/flags"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	path2 "path"
)

var (
	parentNode = accountTypes.NewNode(nil,common.ChainIdType{})
	pathFlag = flags.DirectoryFlag{
		Name:  "path",
		Usage: "keystore save to",
	}
)

func main() {

	app := flags.NewApp("", "the drep command line interface")
	app.Flags = []cli.Flag{
		pathFlag,
	}
	app.Action = gen
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func gen(ctx *cli.Context) error {
	appPath := getCurPath()
	cfgPath := path2.Join(appPath,"config.json")
	nodeItems, err := parserConfig(cfgPath)
	if err != nil {
		return err
	}
	path := ""
	if ctx.GlobalIsSet(pathFlag.Name) {
		path = ctx.GlobalString(pathFlag.Name)
	}else{
		path = appPath
	}
	bootsNodes := []config.BootNode{}
	standbyKey := []*secp256k1.PrivateKey{}
	nodes := []*accountTypes.Node{}
	produces := []*config.Produce{}
	for i:=0; i< len(nodeItems); i++{
		aNode := getAccount(nodeItems[i].Name)
		nodes = append(nodes, aNode)
		bootsNodes = append(bootsNodes,p2pTypes.BootNode{
			PubKey:(*secp256k1.PublicKey)(&aNode.PrivateKey.PublicKey),
			IP :nodeItems[i].Ip,
			Port:nodeItems[i].Port,
		})
		standbyKey = append(standbyKey, aNode.PrivateKey)
		producer := &consensusTypes.Produce{
			Ip:nodeItems[i].Ip,
			Port:nodeItems[i].Port,
			Public: (*secp256k1.PublicKey)(&aNode.PrivateKey.PublicKey),
		}
		produces = append(produces, producer)
	}
	cfg :=config.NodeConfig{}
	err = json.Unmarshal([]byte(ConfTemplate),&cfg)
	if err != nil {
		return  err
	}


//	accounConfig := accountTypes.Config{}
	logConfig := log.LogConfig{}
	rpcConfig := rpcTypes.RpcConfig{}

	p2pConfig := p2pTypes.P2pConfig{}
	p2pConfig.ListerAddr = "0.0.0.0"
	p2pConfig.BootNodes = bootsNodes

	consensusConfig := consensusTypes.ConsensusConfig{}
	consensusConfig.ConsensusMode = "bft"
	consensusConfig.Producers = produces

	chainConfig := chainTypes.ChainConfig{}
	for i:=0; i<len(nodeItems); i++{
		consensusConfig.MyPk = (*secp256k1.PublicKey)(&standbyKey[i].PublicKey)
		consensusConfig.PrvKey = standbyKey[i]
		userDir :=  path2.Join(path,nodeItems[i].Name)
		os.MkdirAll(userDir, os.ModeDir|os.ModePerm)
		keyStorePath := path2.Join(userDir, "keystore")

		store := accounts.NewFileStore(keyStorePath)
		password := string(sha3.Hash256([]byte("123")))
		store.StoreKey(nodes[i],password)
		cfgPath := path2.Join(userDir, "config.json")
		saveConfig(&cfg,cfgPath)
	}

	return nil
}

func getAccount(name string) *accountTypes.Node {
	node := RandomNode([]byte(name))
	return node
}

func saveConfig(cfg *config.NodeConfig, path string) {
   content, _ := json.MarshalIndent(cfg,"","\t")
   ioutil.WriteFile(path,content,0644)
}

func RandomNode(seed []byte) *accountTypes.Node {
	var (
		prvKey *secp256k1.PrivateKey
		chainCode []byte
	)

	h := hmAC(seed, accountTypes.DrepMark)
	prvKey, _ = secp256k1.PrivKeyFromBytes(h[:accountTypes.KeyBitSize])
	chainCode = h[accountTypes.KeyBitSize:]
	addr :=  crypto.PubKey2Address(prvKey.PubKey())
	return &accountTypes.Node{
		PrivateKey: prvKey,
		Address: &addr,
		ChainId: common.ChainIdType{},
		ChainCode: chainCode,
	}
}

func hmAC(message, key []byte) []byte {
	h := hmac.New(sha512.New, key)
	h.Write(message)
	return h.Sum(nil)
}

func getCurPath() string {
	dir, _ := os.Getwd()
	return dir
}

func parserConfig(cfgPath string) ([]*NodeItem, error) {
	content, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil,err
	}
	cfg := []*NodeItem{}
	err = json.Unmarshal([]byte(content), &cfg)
	if err != nil {
		return nil,err
	}
	return cfg, nil
}


type NodeItem struct {
	Name string
	Ip string
	Port int
}