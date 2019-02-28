package main

import (
	"github.com/drep-project/drep-chain/app"
	"os"
	"fmt"
	path2 "path"
	"io/ioutil"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"gopkg.in/urfave/cli.v1"
	"github.com/drep-project/drep-chain/common"
	"github.com/drep-project/drep-chain/crypto"
	"github.com/drep-project/drep-chain/crypto/sha3"
	"github.com/drep-project/drep-chain/crypto/secp256k1"
	"github.com/drep-project/drep-chain/pkgs/log"

	p2pTypes "github.com/drep-project/drep-chain/network/types"
	"github.com/drep-project/drep-chain/rpc"
	chainTypes "github.com/drep-project/drep-chain/chain/types"
	accountTypes "github.com/drep-project/drep-chain/pkgs/accounts/types"
	consensusTypes "github.com/drep-project/drep-chain/pkgs/consensus/types"
	accountComponent "github.com/drep-project/drep-chain/pkgs/accounts/component"
)

var (
	pasword = "123"
	parentNode = chainTypes.NewNode(nil,app.ChainIdType{})
	pathFlag = common.DirectoryFlag{
		Name:  "path",
		Usage: "keystore save to",
	}
)

func main() {
	app := cli.NewApp()
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
	bootsNodes := []p2pTypes.BootNode{}
	standbyKey := []*secp256k1.PrivateKey{}
	nodes := []*chainTypes.Node{}
	produces := []*consensusTypes.Producer{}
	for i:=0; i< len(nodeItems); i++{
		aNode := getAccount(nodeItems[i].Name)
		nodes = append(nodes, aNode)
		bootsNodes = append(bootsNodes,p2pTypes.BootNode{
			//PubKey:(*secp256k1.PublicKey)(&aNode.PrivateKey.PublicKey),
			IP :nodeItems[i].Ip,
			Port:nodeItems[i].Port,
		})
		standbyKey = append(standbyKey, aNode.PrivateKey)
		producer := &consensusTypes.Producer{
			Ip:nodeItems[i].Ip,
			Port:nodeItems[i].Port,
			Public: (*secp256k1.PublicKey)(&aNode.PrivateKey.PublicKey),
		}
		produces = append(produces, producer)
	}


	logConfig := log.LogConfig{}
	logConfig.LogLevel = 3

	rpcConfig := rpc.RpcConfig{}
	rpcConfig.IPCEnabled = true
	rpcConfig.HTTPEnabled = true
	p2pConfig := p2pTypes.P2pConfig{}
	p2pConfig.ListerAddr = "0.0.0.0"
	p2pConfig.Port = 55555
	p2pConfig.BootNodes = bootsNodes

	consensusConfig := consensusTypes.ConsensusConfig{}
	consensusConfig.EnableConsensus = true
	consensusConfig.ConsensusMode = "bft"
	consensusConfig.Producers = produces

	chainConfig := chainTypes.ChainConfig{}
	chainConfig.RemotePort = 55555
	chainConfig.ChainId = app.ChainIdType{}
	chainConfig.GenesisPK = "0x03177b8e4ef31f4f801ce00260db1b04cc501287e828692a404fdbc46c7ad6ff26"
	
	walletConfig := accountTypes.Config{}
	walletConfig.WalletPassword = pasword
	for i:=0; i<len(nodeItems); i++{
		consensusConfig.MyPk = (*secp256k1.PublicKey)(&standbyKey[i].PublicKey)
		p2pConfig.PrvKey = standbyKey[i]
		userDir :=  path2.Join(path,nodeItems[i].Name)
		os.MkdirAll(userDir, os.ModeDir|os.ModePerm)
		keyStorePath := path2.Join(userDir, "keystore")

		store := accountComponent.NewFileStore(keyStorePath)
		password := string(sha3.Hash256([]byte(pasword)))
		store.StoreKey(nodes[i],password)

		cfgPath := path2.Join(userDir, "config.json")
		fs, _ := os.Create(cfgPath)
		offset := int64(0)
		fs.WriteAt([]byte("{\n"),offset)
		offset = int64(2)

		offset = writePhase(fs, "log", logConfig, offset)
		offset = writePhase(fs, "rpc",rpcConfig, offset)
		offset = writePhase(fs, "consensus",consensusConfig, offset)
		offset = writePhase(fs, "p2p",p2pConfig, offset)
		offset = writePhase(fs, "chain",chainConfig, offset)
		offset = writePhase(fs, "accounts",walletConfig, offset)

		fs.Truncate(offset - 2)
		fs.WriteAt([]byte("\n}"), offset-2)

	}
	return nil
}

func writePhase(fs *os.File, name string, config interface{},  offset int64) int64 {
	bytes, _ := json.MarshalIndent(config, "	", "      ")
	bytes = append([]byte("	\""+name+"\" : "), bytes...)
	fs.WriteAt(bytes, offset)
	offset += int64(len(bytes))

	fs.WriteAt([]byte(",\n"),offset)
	offset += 2
	return offset
}

func getAccount(name string) *chainTypes.Node {
	node := RandomNode([]byte(name))
	return node
}

func RandomNode(seed []byte) *chainTypes.Node {
	var (
		prvKey *secp256k1.PrivateKey
		chainCode []byte
	)

	h := hmAC(seed, chainTypes.DrepMark)
	prvKey, _ = secp256k1.PrivKeyFromBytes(h[:chainTypes.KeyBitSize])
	chainCode = h[chainTypes.KeyBitSize:]
	addr :=  crypto.PubKey2Address(prvKey.PubKey())
	return &chainTypes.Node{
		PrivateKey: prvKey,
		Address: &addr,
		ChainId: app.ChainIdType{},
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