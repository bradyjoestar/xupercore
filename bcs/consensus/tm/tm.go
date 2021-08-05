package tm

import (
	"github.com/xuperchain/xupercore/kernel/common/xcontext"
	"github.com/xuperchain/xupercore/kernel/consensus/base"
	cctx "github.com/xuperchain/xupercore/kernel/consensus/context"
	"time"
)

// TMConsensus 为TM出块的共识逻辑
type TMConsensus struct {
	ctx    cctx.ConsensusCtx
	status *TMStatus
	config *TMConfig
}

// CompeteMaster 返回是否为矿工以及是否需要进行SyncBlock
// 该函数返回两个bool，第一个表示是否当前应当出块，第二个表示是否当前需要向其他节点同步区块
func (s *TMConsensus) CompeteMaster(height int64) (bool, bool, error) {
	time.Sleep(time.Duration(s.config.Period) * time.Millisecond)

	if s.ctx.Address.Address == s.config.Miner {
		// single共识确定miner后只能通过共识升级改变miner，因此在单个single实例中miner是不可更改的
		// 此时一个miner从始至终都是自己在挖矿，故不需要向其他节点同步区块
		return true, false, nil
	}
	return false, false, nil
}

func (s *TMConsensus) CheckMinerMatch(ctx xcontext.XContext, block cctx.BlockInterface) (bool, error) {
	return true, nil
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
	return nil, nil
}

type TMConfig struct {
	Miner string `json:"miner"`
	// 单位为毫秒
	Period  int64 `json:"period"`
	Version int64 `json:"version"`
}
