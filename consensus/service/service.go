package service

import (
	"BlockChainTest/database"
	"encoding/json"
	"github.com/AsynkronIT/protoactor-go/actor"
	accountService "github.com/drep-project/drep-chain/accounts/service"
	"github.com/drep-project/drep-chain/app"
	chainService "github.com/drep-project/drep-chain/chain/service"
	chainTypes "github.com/drep-project/drep-chain/chain/types"
	consensusTypes "github.com/drep-project/drep-chain/consensus/types"
	"github.com/drep-project/drep-chain/crypto/secp256k1"
	"github.com/drep-project/drep-chain/crypto/secp256k1/schnorr"
	"github.com/drep-project/drep-chain/crypto/sha3"
	"github.com/drep-project/drep-chain/log"
	p2pService "github.com/drep-project/drep-chain/network/service"
	p2pTypes "github.com/drep-project/drep-chain/network/types"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
	"reflect"
	"time"
)

const (
	blockInterval = time.Second*5
	minWaitTime = time.Millisecond * 500
)

type ConsensusService struct {
	p2pServer *p2pService.P2pService  `service:"p2p"`
	chainService *chainService.ChainService	`service:"chain"`
	walletService *accountService.AccountService `service:"accounts"`

	apis   []app.API

	pubkey *secp256k1.PublicKey
	privkey *secp256k1.PrivateKey
	producers []*consensusTypes.Produce
	consensusConfig *consensusTypes.ConsensusConfig

	pid *actor.PID

	currentHeight int64
	curMiner int
	leader *Leader
	member *Member
	consensusUsedTime time.Duration
	//

	quitRound chan struct{}  //Now no use
}

func (consensusService *ConsensusService) Name() string {
	return "consensus"
}

func (consensusService *ConsensusService) Api() []app.API {
	return consensusService.apis
}

func (consensusService *ConsensusService) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (consensusService *ConsensusService)  P2pMessages() map[int]interface{} {
	return map[int]interface{}{
		consensusTypes.MsgTypeSetUp:reflect.TypeOf(consensusTypes.Setup{}),
		consensusTypes.MsgTypeCommitment:reflect.TypeOf(consensusTypes.Commitment{}),
		consensusTypes.MsgTypeChallenge:reflect.TypeOf(consensusTypes.Challenge{}),
		consensusTypes.MsgTypeResponse:reflect.TypeOf(consensusTypes.Response{}),
		consensusTypes.MsgTypeFail:reflect.TypeOf(consensusTypes.Fail{}),
	}
}

func (consensusService *ConsensusService) Init(executeContext *app.ExecuteContext) error {
	phase := executeContext.GetConfig(consensusService.Name())
	consensusService.consensusConfig = &consensusTypes.ConsensusConfig{}
	err := json.Unmarshal(phase, consensusService.consensusConfig )
	if err != nil {
		return err
	}

	consensusService.currentHeight = 123 // TODO read from database
	consensusService.pubkey = consensusService.consensusConfig.MyPk
	consensusService.producers = consensusService.consensusConfig.Producers
	accountNode, _  := consensusService.walletService.Wallet.GetAccountByPubkey(consensusService.pubkey)
	consensusService.privkey = accountNode.PrivateKey

	props := actor.FromProducer(func() actor.Actor {
		return consensusService
	})
	pid, err := actor.SpawnNamed(props, "consensus_dbft")
	if err != nil {
		panic(err)
	}

	router :=  consensusService.p2pServer.Router
	chainP2pMessage := consensusService.P2pMessages()
	for msgType, _ := range chainP2pMessage {
		router.RegisterMsgHandler(msgType,pid)
	}
	consensusService.pid = pid
	consensusService.leader = NewLeader(consensusService.pubkey, consensusService.quitRound, consensusService.p2pServer)
	consensusService.member = NewMember(consensusService.privkey, consensusService.quitRound , consensusService.p2pServer)

	consensusService.apis = []app.API{
		app.API{
			Namespace: "consensus",
			Version:   "1.0",
			Service: &ConsensusApi{
				consensusService: consensusService,
			},
			Public: true,
		},
	}
	return nil
}

func (consensusService *ConsensusService) Start(executeContext *app.ExecuteContext) error {
	if !consensusService.isProduce() {
		return nil
	}
	minMember := len(consensusService.consensusConfig.Producers)
	for {
		log.Trace("node start", "Height", consensusService.currentHeight)
		var block *chainTypes.Block
		var err error
		if consensusService.consensusConfig.ConsensusMode == "solo" {
			block, err = consensusService.runAsSolo()
		} else {
			//TODO a more elegant implementation is needed: select live peer ,and Determine who is the leader
			participants := consensusService.CollectLiveMember()
			if len(participants) > 1 {
				isM, isL := consensusService.MoveToNextMiner(participants)
				if isL {
					consensusService.leader.UpdateStatus(participants, consensusService.curMiner, minMember, consensusService.currentHeight)
					block, err = consensusService.runAsLeader()
				}else if isM {
					consensusService.member.UpdateStatus(participants, consensusService.curMiner, minMember, consensusService.currentHeight)
					block, err = consensusService.runAsMember()
				}else{
					// backup node， return directly
					log.Debug("backup node")
					return nil
				}
			}else{
				err = errors.New("bft node not ready")
				time.Sleep(time.Second*10)
			}
		}
		if err != nil {
			log.Debug("Produce Block Fail", "reason", err.Error())
		}else{
			consensusService.chainService.ProcessBlock(block)
			consensusService.p2pServer.Broadcast(block)
			log.Info("Block Produced  ", "Height", database.GetMaxHeight())
		}

		time.Sleep(500*time.Millisecond)
		nextBlockTime, waitSpan :=  consensusService.GetWaitTime()
		log.Debug("Sleep", "nextBlockTime", nextBlockTime, "waitSpan", waitSpan)
		time.Sleep(waitSpan)
		consensusService.OnNewHeightUpdate(database.GetMaxHeight())
	}
	return nil
}

func (consensusService *ConsensusService) Stop(executeContext *app.ExecuteContext) error {
	return nil
}

// setLogConfig creates an log configuration from the set command line flags,
func (consensusService *ConsensusService) setLogConfig(ctx *cli.Context, homeDir string) {

}

func (consensusService *ConsensusService) runAsMember() (*chainTypes.Block, error) {
	consensusService.member.Reset()
	log.Trace("node member is going to process consensus for round 1")
	blockBytes, err := consensusService.member.ProcessConsensus()
	if err != nil {
		return nil, err
	}
	log.Trace("node member finishes consensus for round 1")

	block := &chainTypes.Block{}
	err = json.Unmarshal(blockBytes, block)
	if err != nil {
		return nil, err
	}
	pubkeys := consensusService.member.GetMembers()
	consensusService.member.Reset()
	log.Trace("node member is going to process consensus for round 2")
	multiSigBytes, err := consensusService.member.ProcessConsensus()
	if err != nil {
		return nil, err
	}
	multiSig := &chainTypes.MultiSignature{}
	err = json.Unmarshal(multiSigBytes, multiSig)
	if err != nil {
		return nil, err
	}
	block.MultiSig = multiSig
	//check multiSig

	sigmaPubKey := schnorr.CombinePubkeys(pubkeys)
	isValid :=  schnorr.Verify(sigmaPubKey, sha3.Hash256(blockBytes), multiSig.Sig.R, multiSig.Sig.S)
	if !isValid {
		return nil, errors.New("signature not correct")
	}
	consensusService.leader.Reset()
	log.Trace("node member finishes consensus for round 2")
	return block, nil
}

func (consensusService *ConsensusService) runAsLeader() (*chainTypes.Block, error) {
	consensusService.leader.Reset()

	membersPubkey := []*secp256k1.PublicKey{}
	for _, pub := range  consensusService.leader.members {
		membersPubkey = append(membersPubkey, pub.PubKey)
	}
	block, err := consensusService.chainService.GenerateBlock(membersPubkey)
	if err != nil {
		log.Error("generate block fail", "msg", err )
	}

	log.Trace("node leader is preparing process consensus for round 1", "Block",block)
	msg, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	log.Trace("node leader is going to process consensus for round 1")
	err, sig, bitmap := consensusService.leader.ProcessConsensus(msg)
	if err != nil {
		var str = err.Error()
		log.Error("Error occurs","msg", str)
		return nil, err
	}

	multiSig := &chainTypes.MultiSignature{Sig: *sig, Bitmap: bitmap}
	log.Trace("node leader is preparing process consensus for round 2")
	consensusService.leader.Reset()
	msg, err = json.Marshal(multiSig);
	if err != nil {
		return nil, err
	}
	log.Trace("node leader is going to process consensus for round 2")
	err, _, _ = consensusService.leader.ProcessConsensus(msg)
	if err != nil {
		return nil, err
	}
	log.Trace("node leader finishes process consensus for round 2")
	block.MultiSig = multiSig
	consensusService.leader.Reset()
	log.Trace("node leader finishes sending block")
	return block, nil
}

func (consensusService *ConsensusService) runAsSolo() (*chainTypes.Block, error){
	membersPubkey := []*secp256k1.PublicKey{}
	for _, produce := range  consensusService.producers {
		membersPubkey = append(membersPubkey, produce.Public)
	}
	block, _ := consensusService.chainService.GenerateBlock(membersPubkey)
	msg, err := json.Marshal(block)
	if err != nil {
		return block, nil
	}

	sig, err := consensusService.privkey.Sign(sha3.Hash256(msg))
	if err != nil {
		log.Error("sign block error")
		return nil, errors.New("sign block error")
	}
	multiSig := &chainTypes.MultiSignature{Sig: *sig, Bitmap: []byte{}}
	block.MultiSig = multiSig
	return block, nil
}

func (consensusService *ConsensusService) isProduce() bool {
	for _, produce := range  consensusService.producers {
		if produce.Public.IsEqual(consensusService.pubkey){
			return true
		}
	}
	return false
}

func (consensusService *ConsensusService) CollectLiveMember()[]*p2pTypes.Peer{
	liveMember := []*p2pTypes.Peer{}
	for _, produce := range consensusService.consensusConfig.Producers {
		if consensusService.pubkey.IsEqual(produce.Public) {
			liveMember = append(liveMember, nil)  // self
		}else{
			peer := consensusService.p2pServer.SelectPeer(produce.Public.Serialize(), produce.Ip)
			if peer != nil {
				liveMember = append(liveMember, peer)
			}
		}
	}
	return liveMember
}

func (consensusService *ConsensusService) MoveToNextMiner(liveMembers []*p2pTypes.Peer) (bool, bool) {
	consensusService.curMiner = int(consensusService.currentHeight%int64(len(liveMembers)))

	if liveMembers[consensusService.curMiner] == nil {
		return false, true
	} else{
		return true, false
	}
}

func (consensusService *ConsensusService) OnNewHeightUpdate(height int64) {
	if height > consensusService.currentHeight {
		consensusService.currentHeight = height
		log.Info("update new height","Height", height)
	}
}

func (consensusService *ConsensusService) GetMyPubkey() *secp256k1.PublicKey {
	return consensusService.pubkey
}

func (consensusService *ConsensusService) GetWaitTime() (time.Time,time.Duration){
	// max_delay_time +(min_block_interval)*windows = expected_block_interval*windows
	// 6h + 5s*windows = 10s*windows
	// windows = 4320
	lastBlockTime := time.Unix(database.GetHighestBlock().Header.Timestamp, 0)
	targetTime := lastBlockTime.Add(blockInterval)
	now := time.Now()
	if targetTime.Before(now) {
		return now.Add(time.Millisecond * 500 ), time.Millisecond * 500
	} else{
		return targetTime, targetTime.Sub(now)
	}
	/*
     window := int64(4320)
     endBlock := database.GetHighestBlock().Header
     if endBlock.Height < window {
		 lastBlockTime := time.Unix(database.GetHighestBlock().Header.Timestamp, 0)
		 span := time.Now().Sub(lastBlockTime)
		 if span > blockInterval {
			 span = 0
		 } else {
			 span = blockInterval - span
		 }
		 return span
	 }else{
	 	//wait for test
		 startHeight := endBlock.Height - window
		 if startHeight <0 {
			 startHeight = int64(0)
		 }
		 startBlock :=database.GetBlock(startHeight).Header

		 xx := window * 10 -(time.Unix(startBlock.Timestamp,0).Sub(time.Unix(endBlock.Timestamp,0))).Seconds()

		 span := time.Unix(startBlock.Timestamp,0).Sub(time.Unix(endBlock.Timestamp,0))  //window time
		 avgSpan := span.Nanoseconds()/window
		 return time.Duration(avgSpan) * time.Nanosecond
	 }
	*/
}