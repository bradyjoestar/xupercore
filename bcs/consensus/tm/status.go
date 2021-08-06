package tm

import (
	"encoding/json"
	"sync"
)

type ValidatorsInfo struct {
	Validators []string `json:"validators"`
}

type TMStatus struct {
	startHeight int64
	mutex       sync.RWMutex
	newHeight   int64
	index       int
	config      *TMConfig
}

// GetVersion 返回pow所在共识version
func (s *TMStatus) GetVersion() int64 {
	return s.config.Version
}

// GetConsensusBeginInfo 返回该实例初始高度
func (s *TMStatus) GetConsensusBeginInfo() int64 {
	return s.startHeight
}

// GetStepConsensusIndex 获取共识item所在consensus slice中的index
func (s *TMStatus) GetStepConsensusIndex() int {
	return s.index
}

// GetConsensusName 获取共识类型
func (s *TMStatus) GetConsensusName() string {
	return "tm"
}

// GetCurrentTerm 获取当前状态机term
func (s *TMStatus) GetCurrentTerm() int64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.newHeight
}

// GetCurrentValidatorsInfo 获取当前矿工信息
func (s *TMStatus) GetCurrentValidatorsInfo() []byte {
	miner := ValidatorsInfo{
		Validators: []string{s.config.Miner},
	}
	m, _ := json.Marshal(miner)
	return m
}
