package tm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuperchain/xupercore/kernel/common/xcontext"
	"github.com/xuperchain/xupercore/kernel/consensus"
	"github.com/xuperchain/xupercore/kernel/consensus/base"
	cctx "github.com/xuperchain/xupercore/kernel/consensus/context"
	"github.com/xuperchain/xupercore/kernel/consensus/def"
	"strconv"
	"time"
)

// 本次TM改造支持TM的升级，即Miner地址可变
var (
	MinerAddressErr = errors.New("Block's proposer must be equal to its address.")
	BeginBlockChan  = make(chan int, 1)
)

func init() {
	consensus.Register("tm", NewTMConsensus)
}

// TMConsensus TM为单点出块的共识逻辑
type TMConsensus struct {
	ctx    cctx.ConsensusCtx
	status *TMStatus
	config *TMConfig
}

// NewTMConsensus 初始化实例
func NewTMConsensus(cCtx cctx.ConsensusCtx, cCfg def.ConsensusConfig) base.ConsensusImplInterface {
	// 解析config中需要的字段
	if cCtx.XLog == nil {
		return nil
	}
	// TODO:cCtx.BcName需要注册表吗？
	if cCtx.Crypto == nil || cCtx.Address == nil {
		cCtx.XLog.Error("TM::NewTMConsensus::CryptoClient in context is nil")
		return nil
	}
	if cCtx.Ledger == nil {
		cCtx.XLog.Error("TM::NewTMConsensus::Ledger in context is nil")
		return nil
	}
	if cCfg.ConsensusName != "tm" {
		cCtx.XLog.Error("TM::NewTMConsensus::consensus name in config is wrong", "name", cCfg.ConsensusName)
		return nil
	}

	config, err := buildConfigs([]byte(cCfg.Config))
	if err != nil {
		cCtx.XLog.Error("TM::NewTMConsensus::TM parse config", "error", err)
		return nil
	}

	// newHeight取上一共识的最高值，因为此时BeginHeight也许并为生产出来
	status := &TMStatus{
		startHeight: cCfg.StartHeight,
		newHeight:   cCfg.StartHeight - 1,
		index:       cCfg.Index,
		config:      config,
	}
	tm := &TMConsensus{
		ctx:    cCtx,
		config: config,
		status: status,
	}
	return tm
}

// CompeteMaster 返回是否为矿工以及是否需要进行SyncBlock
// 该函数返回两个bool，第一个表示是否当前应当出块，第二个表示是否当前需要向其他节点同步区块
func (s *TMConsensus) CompeteMaster(height int64) (bool, bool, error) {
	//time.Sleep(time.Duration(s.config.Period) * time.Millisecond)

	//tmp block here
	v := <-BeginBlockChan
	fmt.Println(time.Now().UTC().String())
	fmt.Println("BeginBlockChain:" + strconv.Itoa(v))

	if s.ctx.Address.Address == s.config.Miner {
		// TM共识确定miner后只能通过共识升级改变miner，因此在单个TM实例中miner是不可更改的
		// 此时一个miner从始至终都是自己在挖矿，故不需要向其他节点同步区块
		return true, false, nil
	}
	return false, false, nil
}

// CheckMinerMatch 查看block是否合法
// ATTENTION: TODO: 上层需要先检查VerifyMerkle(block)
func (s *TMConsensus) CheckMinerMatch(ctx xcontext.XContext, block cctx.BlockInterface) (bool, error) {
	// 检查区块的区块头是否hash正确
	bid, err := block.MakeBlockId()
	if err != nil {
		return false, err
	}
	if !bytes.Equal(bid, block.GetBlockid()) {
		ctx.GetLog().Warn("TM::CheckMinerMatch::equal blockid error")
		return false, err
	}
	// 检查矿工地址是否合法
	if string(block.GetProposer()) != s.config.Miner {
		ctx.GetLog().Warn("TM::CheckMinerMatch::miner check error", "blockid", block.GetBlockid(),
			"proposer", string(block.GetProposer()), "local proposer", s.config.Miner)
		return false, err
	}
	//验证签名
	//1 验证一下签名和公钥是不是能对上
	k, err := s.ctx.Crypto.GetEcdsaPublicKeyFromJsonStr(block.GetPublicKey())
	if err != nil {
		ctx.GetLog().Warn("TM::CheckMinerMatch::get ecdsa from block error", "error", err)
		return false, err
	}
	chkResult, _ := s.ctx.Crypto.VerifyAddressUsingPublicKey(string(block.GetProposer()), k)
	if chkResult == false {
		ctx.GetLog().Warn("TM::CheckMinerMatch::address is not match publickey")
		return false, err
	}
	//2 验证地址
	addr, err := s.ctx.Crypto.GetAddressFromPublicKey(k)
	if err != nil {
		return false, err
	}
	if addr != string(block.GetProposer()) {
		return false, MinerAddressErr
	}
	//3 验证一下签名是否正确
	valid, err := s.ctx.Crypto.VerifyECDSA(k, block.GetSign(), block.GetBlockid())
	if err != nil {
		ctx.GetLog().Warn("TM::CheckMinerMatch::verifyECDSA error",
			"error", err, "sign", block.GetSign())
	}
	return valid, err
}

// ProcessBeforeMiner 开始挖矿前进行相应的处理, 返回是否需要truncate, 返回写consensusStorage, 返回err
func (s *TMConsensus) ProcessBeforeMiner(timestamp int64) ([]byte, []byte, error) {
	return nil, nil, nil
}

// CalculateBlock 矿工挖矿时共识需要做的工作, 如PoW时共识需要完成存在性证明
func (s *TMConsensus) CalculateBlock(block cctx.BlockInterface) error {
	return nil
}

// ProcessConfirmBlock 用于确认块后进行相应的处理
func (s *TMConsensus) ProcessConfirmBlock(block cctx.BlockInterface) error {
	return nil
}

// GetStatus 获取区块链共识信息
func (s *TMConsensus) GetConsensusStatus() (base.ConsensusStatus, error) {
	return s.status, nil
}

// 共识实例的挂起逻辑, 另: 若共识实例发现绑定block结构有误，会直接停掉当前共识实例并panic
func (s *TMConsensus) Stop() error {
	return nil
}

// 共识实例的启动逻辑
func (s *TMConsensus) Start() error {
	return nil
}

// ParseConsensusStorage 共识占用blockinterface的专有存储，特定共识需要提供parse接口，在此作为接口高亮
// TM共识没有用到区块存储信息, 故返回空
func (s *TMConsensus) ParseConsensusStorage(block cctx.BlockInterface) (interface{}, error) {
	return nil, nil
}

type TMConfig struct {
	Miner string `json:"miner"`
	// 单位为毫秒
	Period  int64 `json:"period"`
	Version int64 `json:"version"`
}

func buildConfigs(input []byte) (*TMConfig, error) {
	v := make(map[string]string)
	err := json.Unmarshal(input, &v)
	if err != nil {
		return nil, fmt.Errorf("unmarshal TM config error")
	}

	config := &TMConfig{
		Miner: v["miner"],
	}

	if v["version"] != "" {
		config.Version, err = strconv.ParseInt(v["version"], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse version error: %v, %v", err, v["version"])
		}
	}

	if v["period"] != "" {
		config.Period, err = strconv.ParseInt(v["period"], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse period error: %v, %v", err, v["period"])
		}
	}

	return config, nil
}
