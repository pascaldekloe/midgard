package timeseries

import (
	"bytes"
	"fmt"
	"log"

	"github.com/pascaldekloe/metrics"

	"gitlab.com/thorchain/midgard/event"
)

// Double Depth Counting
var (
	addPerPoolAndAsset     = metrics.Must2LabelCounter("midgard_event_add_E8s_total", "pool", "asset")
	errataPerPoolAndAsset  = metrics.Must2LabelCounter("midgard_event_errata_E8s_total", "pool", "asset")
	gasPerPoolAndAsset     = metrics.Must2LabelCounter("midgard_event_gas_E8s_total", "pool", "asset")
	rewardsPerPoolAndAsset = metrics.Must2LabelCounter("midgard_event_rewards_E8s_total", "pool", "asset")
	stakePerPoolAndAsset   = metrics.Must2LabelCounter("midgard_event_stake_E8s_total", "pool", "asset")
	swapPerPoolAndAsset    = metrics.Must2LabelCounter("midgard_event_swap_E8s_total", "pool", "asset")
)

// Empty prevents the SQL driver from writing NULL values.
var empty = []byte{}

// EventListener is a singleton implementation which MUST be invoked seqentially
// in order of appearance.
var EventListener event.Listener = recorder

// Recorder gets initialised by Setup.
var recorder = &eventRecorder{
	runningTotals: *newRunningTotals(),
	linkedEvents:  *newLinkedEvents(),
}

type eventRecorder struct {
	runningTotals
	linkedEvents
}

func (r *eventRecorder) OnActiveVault(e *event.ActiveVault, meta *event.Metadata) {
	const q = `INSERT INTO active_vault_events (add_asgard_addr, block_timestamp)
VALUES ($1, $2)`
	_, err := DBExec(q, e.AddAsgardAddr, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("ActiveVault event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (r *eventRecorder) OnAdd(e *event.Add, meta *event.Metadata) {
	const q = `INSERT INTO add_events (tx, chain, from_addr, to_addr, asset, asset_E8, memo, rune_E8, pool, block_timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := DBExec(q, e.Tx, e.Chain, e.FromAddr, e.ToAddr, e.Asset, e.AssetE8, e.Memo, e.RuneE8, e.Pool, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("add event from height %d lost on %s", meta.BlockHeight, err)
		return
	}

	r.AddPoolAssetE8Depth(e.Pool, e.AssetE8)
	r.AddPoolRuneE8Depth(e.Pool, e.RuneE8)
	addPerPoolAndAsset(string(e.Pool), "pool").Add(uint64(e.AssetE8))
	addPerPoolAndAsset(string(e.Pool), "rune").Add(uint64(e.RuneE8))
}

func (r *eventRecorder) OnAsgardFundYggdrasil(e *event.AsgardFundYggdrasil, meta *event.Metadata) {
	const q = `INSERT INTO asgard_fund_yggdrasil_events (tx, asset, asset_E8, vault_key, block_timestamp)
VALUES ($1, $2, $3, $4, $5)`
	_, err := DBExec(q, e.Tx, e.Asset, e.AssetE8, e.VaultKey, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("asgard_fund_yggdrasil event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (_ *eventRecorder) OnBond(e *event.Bond, meta *event.Metadata) {
	const q = `INSERT INTO bond_events (tx, chain, from_addr, to_addr, asset, asset_E8, memo, bound_type, E8, block_timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := DBExec(q, e.Tx, e.Chain, e.FromAddr, e.ToAddr, e.Asset, e.AssetE8, e.Memo, e.BoundType, e.E8, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("bond event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (r *eventRecorder) OnErrata(e *event.Errata, meta *event.Metadata) {
	const q = `INSERT INTO errata_events (in_tx, asset, asset_E8, rune_E8, block_timestamp)
VALUES ($1, $2, $3, $4, $5)`
	_, err := DBExec(q, e.InTx, e.Asset, e.AssetE8, e.RuneE8, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("errata event from height %d lost on %s", meta.BlockHeight, err)
		return
	}

	r.AddPoolAssetE8Depth(e.Asset, e.AssetE8)
	r.AddPoolRuneE8Depth(e.Asset, e.RuneE8)
	errataPerPoolAndAsset(string(e.Asset), "pool").Add(uint64(e.AssetE8))
	errataPerPoolAndAsset(string(e.Asset), "rune").Add(uint64(e.RuneE8))
}

func (r *eventRecorder) OnFee(e *event.Fee, meta *event.Metadata) {
	const q = `INSERT INTO fee_events (tx, asset, asset_E8, pool_deduct, block_timestamp)
VALUES ($1, $2, $3, $4, $5)`
	_, err := DBExec(q, e.Tx, e.Asset, e.AssetE8, e.PoolDeduct, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("fee event from height %d lost on %s", meta.BlockHeight, err)
	}

	r.linkedEvents.OnFee(e, meta)
}

func (r *eventRecorder) OnGas(e *event.Gas, meta *event.Metadata) {
	const q = `INSERT INTO gas_events (asset, asset_E8, rune_E8, tx_count, block_timestamp)
VALUES ($1, $2, $3, $4, $5)`
	_, err := DBExec(q, e.Asset, e.AssetE8, e.RuneE8, e.TxCount, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("gas event from height %d lost on %s", meta.BlockHeight, err)
		return
	}

	r.AddPoolAssetE8Depth(e.Asset, -e.AssetE8)
	r.AddPoolRuneE8Depth(e.Asset, e.RuneE8)
	gasPerPoolAndAsset(string(e.Asset), "pool").Add(uint64(e.AssetE8))
	gasPerPoolAndAsset(string(e.Asset), "rune").Add(uint64(e.RuneE8))
}

func (r *eventRecorder) OnInactiveVault(e *event.InactiveVault, meta *event.Metadata) {
	const q = `INSERT INTO inactive_vault_events (add_asgard_addr, block_timestamp)
VALUES ($1, $2)`
	_, err := DBExec(q, e.AddAsgardAddr, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("InactiveVault event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (_ *eventRecorder) OnMessage(e *event.Message, meta *event.Metadata) {
	if e.FromAddr == nil {
		e.FromAddr = empty
	}
	if e.Action == nil {
		e.Action = empty
	}
	const q = `INSERT INTO message_events (from_addr, action, block_timestamp)
VALUES ($1, $2, $3)`
	_, err := DBExec(q, e.FromAddr, e.Action, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("message event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (_ *eventRecorder) OnNewNode(e *event.NewNode, meta *event.Metadata) {
	const q = `INSERT INTO new_node_events (node_addr, block_timestamp)
VALUES ($1, $2)`
	_, err := DBExec(q, e.NodeAddr, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("new_node event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (r *eventRecorder) OnOutbound(e *event.Outbound, meta *event.Metadata) {
	const q = `INSERT INTO outbound_events (tx, chain, from_addr, to_addr, asset, asset_E8, memo, in_tx, block_timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := DBExec(q, e.Tx, e.Chain, e.FromAddr, e.ToAddr, e.Asset, e.AssetE8, e.Memo, e.InTx, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("outound event from height %d lost on %s", meta.BlockHeight, err)
	}

	r.linkedEvents.OnOutbound(e, meta)
}

func (_ *eventRecorder) OnPool(e *event.Pool, meta *event.Metadata) {
	const q = `INSERT INTO pool_events (asset, status, block_timestamp)
VALUES ($1, $2, $3)`
	_, err := DBExec(q, e.Asset, e.Status, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("pool event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (r *eventRecorder) OnRefund(e *event.Refund, meta *event.Metadata) {
	const q = `INSERT INTO refund_events (tx, chain, from_addr, to_addr, asset, asset_E8, asset_2nd, asset_2nd_E8, memo, code, reason, block_timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := DBExec(q, e.Tx, e.Chain, e.FromAddr, e.ToAddr, e.Asset, e.AssetE8, e.Asset2nd, e.Asset2ndE8, e.Memo, e.Code, e.Reason, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("refund event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (_ *eventRecorder) OnReserve(e *event.Reserve, meta *event.Metadata) {
	const q = `INSERT INTO reserve_events (tx, chain, from_addr, to_addr, asset, asset_E8, memo, addr, E8, block_timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := DBExec(q, e.Tx, e.Chain, e.FromAddr, e.ToAddr, e.Asset, e.AssetE8, e.Memo, e.Addr, e.E8, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("reserve event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (r *eventRecorder) OnRewards(e *event.Rewards, meta *event.Metadata) {
	blockTimestamp := meta.BlockTimestamp.UnixNano()
	const q = "INSERT INTO rewards_events (bond_E8, block_timestamp) VALUES ($1, $2)"
	_, err := DBExec(q, e.BondE8, blockTimestamp)
	if err != nil {
		log.Printf("reserve event from height %d lost on %s", meta.BlockHeight, err)
		return
	}

	if len(e.PerPool) == 0 {
		return
	}

	// make batch insert work 🥴
	buf := bytes.NewBufferString("INSERT INTO rewards_event_entries (pool, rune_E8, block_timestamp) VALUES ")
	args := make([]interface{}, len(e.PerPool)*3)
	for i, p := range e.PerPool {
		fmt.Fprintf(buf, "($%d, $%d, $%d),", i*3+1, i*3+2, i*3+3)
		args[i*3], args[i*3+1], args[i*3+2] = p.Asset, p.E8, blockTimestamp
	}
	buf.Truncate(buf.Len() - 1) // last comma
	if _, err := DBExec(buf.String(), args...); err != nil {
		log.Printf("reserve event pools from height %d lost on %s", meta.BlockHeight, err)
		return
	}

	for _, a := range e.PerPool {
		r.AddPoolRuneE8Depth(a.Asset, a.E8)
		rewardsPerPoolAndAsset(string(a.Asset), "rune").Add(uint64(a.E8))
	}
}

func (_ *eventRecorder) OnSetIPAddress(e *event.SetIPAddress, meta *event.Metadata) {
	const q = `INSERT INTO set_ip_address_events (node_addr, ip_addr, block_timestamp)
VALUES ($1, $2, $3)`
	_, err := DBExec(q, e.NodeAddr, e.IPAddr, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("set_ip_address event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (_ *eventRecorder) OnSetMimir(e *event.SetMimir, meta *event.Metadata) {
	const q = `INSERT INTO set_mimir_events (key, value, block_timestamp)
VALUES ($1, $2, $3)`
	_, err := DBExec(q, e.Key, e.Value, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("set_mimir event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (_ *eventRecorder) OnSetNodeKeys(e *event.SetNodeKeys, meta *event.Metadata) {
	const q = `INSERT INTO set_node_keys_events (node_addr, secp256k1, ed25519, validator_consensus, block_timestamp)
VALUES ($1, $2, $3, $4, $5)`
	_, err := DBExec(q, e.NodeAddr, e.Secp256k1, e.Ed25519, e.ValidatorConsensus, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("set_node_keys event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (_ *eventRecorder) OnSetVersion(e *event.SetVersion, meta *event.Metadata) {
	const q = `INSERT INTO set_version_events (node_addr, version, block_timestamp)
VALUES ($1, $2, $3)`
	_, err := DBExec(q, e.NodeAddr, e.Version, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("set_version event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (_ *eventRecorder) OnSlash(e *event.Slash, meta *event.Metadata) {
	if len(e.Amounts) == 0 {
		log.Printf("slash event on pool %q ignored: zero amounts", e.Pool)
	}
	for _, a := range e.Amounts {
		const q = "INSERT INTO slash_amounts (pool, asset, asset_E8, block_timestamp) VALUES ($1, $2, $3, $4)"
		_, err := DBExec(q, e.Pool, a.Asset, a.E8, meta.BlockTimestamp.UnixNano())
		if err != nil {
			log.Printf("slash amount from height %d lost on %s", meta.BlockHeight, err)
		}
	}
}

func (r *eventRecorder) OnStake(e *event.Stake, meta *event.Metadata) {
	const q = `INSERT INTO stake_events (pool, asset_tx, asset_chain, asset_E8, rune_tx, rune_addr, rune_E8, stake_units, block_timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := DBExec(q, e.Pool, e.AssetTx, e.AssetChain, e.AssetE8, e.RuneTx, e.RuneAddr, e.RuneE8, e.StakeUnits, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("stake event from height %d lost on %s", meta.BlockHeight, err)
		return
	}

	r.AddPoolAssetE8Depth(e.Pool, e.AssetE8)
	r.AddPoolRuneE8Depth(e.Pool, e.RuneE8)
	stakePerPoolAndAsset(string(e.Pool), "pool").Add(uint64(e.AssetE8))
	stakePerPoolAndAsset(string(e.Pool), "rune").Add(uint64(e.RuneE8))
}

func (r *eventRecorder) OnSwap(e *event.Swap, meta *event.Metadata) {
	const q = `INSERT INTO swap_events (tx, chain, from_addr, to_addr, from_asset, from_E8, memo, pool, to_E8_min, trade_slip_BP, liq_fee_E8, liq_fee_in_rune_E8, block_timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := DBExec(q, e.Tx, e.Chain, e.FromAddr, e.ToAddr, e.FromAsset, e.FromE8, e.Memo, e.Pool, e.ToE8Min, e.TradeSlipBP, e.LiqFeeE8, e.LiqFeeInRuneE8, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("swap event from height %d lost on %s", meta.BlockHeight, err)
		return
	}

	if e.ToRune() {
		// Swap adds pool asset.
		r.AddPoolAssetE8Depth(e.Pool, e.FromE8)
		swapPerPoolAndAsset(string(e.Pool), "pool").Add(uint64(e.FromE8))
		// Swap deducts RUNE from pool with an event.Outbound.
	} else {
		// Swap adds RUNE to pool.
		r.AddPoolRuneE8Depth(e.Pool, e.FromE8)
		swapPerPoolAndAsset(string(e.Pool), "rune").Add(uint64(e.FromE8))
		// Swap deducts pool asset with an event.Outbound.
	}
}

func (_ *eventRecorder) OnTransfer(e *event.Transfer, meta *event.Metadata) {
	const q = `INSERT INTO transfer_events (from_addr, to_addr, rune_E8, block_timestamp)
VALUES ($1, $2, $3)`
	_, err := DBExec(q, e.FromAddr, e.ToAddr, e.RuneE8, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("transfer event from height %d lost on %s", meta.BlockHeight, err)
		return
	}
}

func (r *eventRecorder) OnUnstake(e *event.Unstake, meta *event.Metadata) {
	const q = `INSERT INTO unstake_events (tx, chain, from_addr, to_addr, asset, asset_E8, memo, pool, stake_units, basis_points, asymmetry, block_timestamp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := DBExec(q, e.Tx, e.Chain, e.FromAddr, e.ToAddr, e.Asset, e.AssetE8, e.Memo, e.Pool, e.StakeUnits, e.BasisPoints, e.Asymmetry, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("unstake event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (_ *eventRecorder) OnUpdateNodeAccountStatus(e *event.UpdateNodeAccountStatus, meta *event.Metadata) {
	if e.Former == nil {
		e.Former = empty
	}
	const q = `INSERT INTO update_node_account_status_events (node_addr, former, current, block_timestamp)
VALUES ($1, $2, $3, $4)`
	_, err := DBExec(q, e.NodeAddr, e.Former, e.Current, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("UpdateNodeAccountStatus event from height %d lost on %s", meta.BlockHeight, err)
	}
}

func (_ *eventRecorder) OnValidatorRequestLeave(e *event.ValidatorRequestLeave, meta *event.Metadata) {
	const q = `INSERT INTO validator_request_leave_events (tx, from_addr, node_addr, block_timestamp)
VALUES ($1, $2, $3, $4)`
	_, err := DBExec(q, e.Tx, e.FromAddr, e.NodeAddr, meta.BlockTimestamp.UnixNano())
	if err != nil {
		log.Printf("validator_request_leave event from height %d lost on %s", meta.BlockHeight, err)
	}
}
