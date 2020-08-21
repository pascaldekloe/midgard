package event

import (
	"errors"
	"log"
	"time"

	"github.com/pascaldekloe/metrics"
	// Tendermint is all about types? 🤔
	abci "github.com/tendermint/tendermint/abci/types"
	rpc "github.com/tendermint/tendermint/rpc/core/types"
	tendermint "github.com/tendermint/tendermint/types"
)

// Package Metrics
var (
	BlockHeight = metrics.MustInteger("midgard_chain_height", "The sequence identifier the last block read.")
	BlockTotal  = metrics.MustCounter("migdard_chain_blocks_total", "Read counter.")
	EventTotal  = metrics.MustCounter("midgard_chain_events_total", "Read counter.")
)

// Metadata has metadata for a block (from the chain).
type Metadata struct {
	BlockTimestamp time.Time // official acceptance moment
}

// Listener defines an event callback.
type Listener interface {
	OnAdd(*Add, *Metadata)
	OnBond(*Bond, *Metadata)
	OnErrata(*Errata, *Metadata)
	OnFee(*Fee, *Metadata)
	OnGas(*Gas, *Metadata)
	OnOutbound(*Outbound, *Metadata)
	OnPool(*Pool, *Metadata)
	OnRefund(*Refund, *Metadata)
	OnReserve(*Reserve, *Metadata)
	OnRewards(*Rewards, *Metadata)
	OnSlash(*Slash, *Metadata)
	OnStake(*Stake, *Metadata)
	OnSwap(*Swap, *Metadata)
	OnUnstake(*Unstake, *Metadata)
}

// Demux is a demultiplexer for events from the blockchain.
type Demux struct {
	Listener // destination

	// prevent memory allocation
	reuse struct {
		Add
		Bond
		Errata
		Fee
		Gas
		Outbound
		Pool
		Refund
		Reserve
		Rewards
		Stake
		Swap
		Unstake
	}
}

// Block invokes Listener for each transaction event in block.
func (d *Demux) Block(block *rpc.ResultBlockResults, meta *tendermint.BlockMeta) {
	BlockTotal.Add(1)
	BlockHeight.Set(meta.Header.Height)

	m := Metadata{BlockTimestamp: meta.Header.Time}

	// TODO(pascaldekloe): Find best way to ID an event.

	for txIndex, tx := range block.TxsResults {
		for eventIndex, event := range tx.Events {
			if err := d.event(event, &m); err != nil {
				log.Printf("block %s tx %d event %d type %q skipped: %s",
					meta.BlockID.String(), txIndex, eventIndex, event.Type, err)
			}
		}
	}
	for eventIndex, event := range block.BeginBlockEvents {
		if err := d.event(event, &m); err != nil {
			log.Printf("block %s begin event %d type %q skipped: %s",
				meta.BlockID.String(), eventIndex, event.Type, err)
		}
	}
	for eventIndex, event := range block.EndBlockEvents {
		if err := d.event(event, &m); err != nil {
			log.Printf("block %s end event %d type %q skipped: %s",
				meta.BlockID.String(), eventIndex, event.Type, err)
		}
	}
}

var errEventType = errors.New("unknown event type")

// Block notifies Listener for the transaction event.
// Errors do not include the event type in the message.
func (d *Demux) event(event abci.Event, meta *Metadata) error {
	EventTotal.Add(1)

	attrs := event.Attributes

	switch event.Type {
	case "add":
		if err := d.reuse.Add.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnAdd(&d.reuse.Add, meta)
	case "bond":
		if err := d.reuse.Bond.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnBond(&d.reuse.Bond, meta)
	case "errata":
		if err := d.reuse.Errata.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnErrata(&d.reuse.Errata, meta)
	case "fee":
		if err := d.reuse.Fee.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnFee(&d.reuse.Fee, meta)
	case "gas":
		if err := d.reuse.Gas.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnGas(&d.reuse.Gas, meta)
	case "outbound":
		if err := d.reuse.Outbound.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnOutbound(&d.reuse.Outbound, meta)
	case "pool":
		if err := d.reuse.Pool.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnPool(&d.reuse.Pool, meta)
	case "refund":
		if err := d.reuse.Refund.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnRefund(&d.reuse.Refund, meta)
	case "reserve":
		if err := d.reuse.Reserve.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnReserve(&d.reuse.Reserve, meta)
	case "rewards":
		if err := d.reuse.Rewards.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnRewards(&d.reuse.Rewards, meta)
	case "stake":
		if err := d.reuse.Stake.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnStake(&d.reuse.Stake, meta)
	case "swap":
		if err := d.reuse.Swap.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnSwap(&d.reuse.Swap, meta)
	case "unstake":
		if err := d.reuse.Unstake.LoadTendermint(attrs); err != nil {
			return err
		}
		d.Listener.OnUnstake(&d.reuse.Unstake, meta)
	case "ActiveVault", "InactiveVault", "asgard_fund_yggdrasil",
		"message", "transfer", "new_node", "UpdateNodeAccountStatus",
		"set_ip_address", "set_node_keys", "set_version",
		"set_mimir", "validator_request_leave":
		break // ignore
	default:
		return errEventType
	}

	return nil
}
