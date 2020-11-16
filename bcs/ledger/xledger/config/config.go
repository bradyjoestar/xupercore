package config

import (
	"fmt"

	"github.com/xuperchain/xupercore/lib/utils"

	"github.com/spf13/viper"
)

type XLedgerConf struct {
	// kv storage type
	KVStorage string `yaml:"kvStorage,omitempty"`
}

func LoadLedgerConf(cfgFile string) (*XLedgerConf, error) {
	cfg := GetDefLedgerConf()
	err := cfg.loadConf(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("load ledger config failed.err:%s", err)
	}

	return cfg, nil
}

func GetDefLedgerConf() *XLedgerConf {
	return &XLedgerConf{
		KVStorage: "leveldb",
	}
}

func (t *XLedgerConf) loadConf(cfgFile string) error {
	if cfgFile == "" || !utils.FileIsExist(cfgFile) {
		return fmt.Errorf("config file set error.path:%s", cfgFile)
	}

	viperObj := viper.New()
	viperObj.SetConfigFile(cfgFile)
	err := viperObj.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config failed.path:%s,err:%v", cfgFile, err)
	}

	if err = viperObj.Unmarshal(t); err != nil {
		return fmt.Errorf("unmatshal config failed.path:%s,err:%v", cfgFile, err)
	}

	return nil
}
