package collector

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EthBlockTimestamp struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

type blockResult struct {
	Timestamp hexutil.Uint64
}

func NewEthBlockTimestamp(rpc *rpc.Client) *EthBlockTimestamp {
	return &EthBlockTimestamp{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"eth_block_timestamp",
			"timestamp of the most recent block",
			nil,
			nil,
		),
	}
}

func (collector *EthBlockTimestamp) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EthBlockTimestamp) Collect(ch chan<- prometheus.Metric) {
	var bn hexutil.Uint64
	var result *blockResult

	if err := collector.rpc.Call(&bn, "eth_blockNumber"); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	if err := collector.rpc.Call(&result, "eth_getBlockByNumber", bn, 0); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	value := float64(result.Timestamp)

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
