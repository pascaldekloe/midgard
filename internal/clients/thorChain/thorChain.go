package thorChain

import (
	"encoding/json"
	"fmt"
	"gitlab.com/thorchain/midgard/internal/common"
	"net/http"
	"sort"
	"strings"

	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/midgard/internal/clients/blockchains/binance"
	"gitlab.com/thorchain/midgard/internal/clients/thorChain/types"
	"gitlab.com/thorchain/midgard/internal/config"
	"gitlab.com/thorchain/midgard/internal/models"
	"gitlab.com/thorchain/midgard/internal/store/timescale"
)

// API to talk to thorchain
type API struct {
	logger        zerolog.Logger
	cfg           config.ThorChainConfiguration
	baseUrl       string
	baseRPCUrl    string
	netClient     *http.Client
	wg            *sync.WaitGroup
	stopChan      chan struct{}
	store         *timescale.Client
	binanceClient *binance.Client
}

// NewBinanceClient create a new instance of API which can talk to thorChain
func NewAPIClient(cfg config.ThorChainConfiguration, binanceClient *binance.Client, timescale *timescale.Client) (*API, error) {
	if len(cfg.Host) == 0 {
		return nil, errors.New("thorchain host is empty")
	}
	return &API{
		cfg:    cfg,
		logger: log.With().Str("module", "thorchain").Logger(),
		netClient: &http.Client{
			Timeout: cfg.ReadTimeout,
		},
		baseUrl:       fmt.Sprintf("%s://%s/thorchain", cfg.Scheme, cfg.Host),
		baseRPCUrl:    fmt.Sprintf("%s://%s", cfg.Scheme, cfg.RPCHost),
		stopChan:      make(chan struct{}),
		wg:            &sync.WaitGroup{},
		store:         timescale,
		binanceClient: binanceClient,
	}, nil
}

func (api *API) getGenesis() (types.Genesis, error) {
	uri := fmt.Sprintf("%s/genesis", api.baseRPCUrl)
	api.logger.Debug().Msg(uri)
	resp, err := api.netClient.Get(uri)
	if err != nil {
		return types.Genesis{}, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			api.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var genesis types.Genesis
	if err := json.NewDecoder(resp.Body).Decode(&genesis); nil != err {
		return types.Genesis{}, errors.Wrap(err, "failed to unmarshal events")
	}

	return genesis, nil
}

func (api *API) processGenesis(genesisTime types.Genesis) error {
	api.logger.Debug().Msg("processGenesisTime")

	record := models.NewGenesis(genesisTime)
	_, err := api.store.CreateGenesis(record)
	if err != nil {
		return errors.Wrap(err, "failed to create genesis record")
	}
	return nil
}

func (api *API) getEvents(id int64) ([]types.Event, error) {
	uri := fmt.Sprintf("%s/events/%d", api.baseUrl, id)
	api.logger.Debug().Msg(uri)
	resp, err := api.netClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			api.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var events []types.Event
	if err := json.NewDecoder(resp.Body).Decode(&events); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal events")
	}
	return events, nil
}

// returns (maxID, events, err)
func (api *API) processEvents(id int64) (int64, int, error) {
	events, err := api.getEvents(id)
	if err != nil {
		return id, 0, errors.Wrap(err, "failed to get events")
	}

	// sort events lowest ID first. Ensures we don't process an event out of order
	sort.Slice(events[:], func(i, j int) bool {
		return events[i].ID < events[j].ID
	})

	maxID := id
	// pts := make([]client.Point, 0)
	for _, evt := range events {
		if maxID < evt.ID {
			maxID = evt.ID
			api.logger.Info().Int64("maxID", maxID).Msg("new maxID")
		}
		if evt.OutTxs == nil {
			outTx, err := api.GetOutTx(evt)
			if err != nil {
				api.logger.Err(err).Msg("GetOutTx failed")
			} else {
				evt.OutTxs = outTx
			}
		}
		switch strings.ToLower(evt.Type) {
		case "swap":
			err = api.processSwapEvent(evt)
			if err != nil {
				api.logger.Err(err).Msg("processSwapEvent failed")
				continue
			}
		case "stake":
			err = api.processStakingEvent(evt)
			if err != nil {
				api.logger.Err(err).Msg("processStakingEvent failed")
				continue
			}
		case "unstake":
			err = api.processUnstakeEvent(evt)
			if err != nil {
				api.logger.Err(err).Msg("processUnstakeEvent failed")
				continue
			}
		case "rewards":
			err = api.processRewardEvent(evt)
			if err != nil {
				api.logger.Err(err).Msg("processRewardEvent failed")
				continue
			}
		case "add":
			err = api.processAddEvent(evt)
			if err != nil {
				api.logger.Err(err).Msg("processAddEvent failed")
				continue
			}
		case "pool":
			err = api.processPoolEvent(evt)
			if err != nil {
				api.logger.Err(err).Msg("processPoolEvent failed")
				continue
			}
		case "gas":
			err = api.processGasEvent(evt)
			if err != nil {
				api.logger.Err(err).Msg("processGasEvent failed")
				continue
			}
		case "refund":
			err = api.processRefundEvent(evt)
			if err != nil {
				api.logger.Err(err).Msg("processRefundEvent failed")
				continue
			}
		default:
			api.logger.Info().Str("evt.Type", evt.Type).Msg("Unknown event type")
			continue
		}
	}
	return maxID, len(events), nil
}

func (api *API) processSwapEvent(evt types.Event) error {
	api.logger.Debug().Msg("processSwapEvent")
	var swap types.EventSwap
	err := json.Unmarshal(evt.Event, &swap)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal swap event")
	}
	record := models.NewSwapEvent(swap, evt)
	err = api.store.CreateSwapRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create swap record")
	}
	return nil
}

func (api *API) processStakingEvent(evt types.Event) error {
	api.logger.Debug().Msg("processStakingEvent")
	var stake types.EventStake
	err := json.Unmarshal(evt.Event, &stake)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal stake event")
	}
	record := models.NewStakeEvent(stake, evt)
	err = api.store.CreateStakeRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create stake record")
	}
	return nil
}

func (api *API) processUnstakeEvent(evt types.Event) error {
	api.logger.Debug().Msg("processUnstakeEvent")
	var unstake types.EventUnstake
	err := json.Unmarshal(evt.Event, &unstake)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal unstake event")
	}
	record := models.NewUnstakeEvent(unstake, evt)
	err = api.store.CreateUnStakesRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create unstake record")
	}
	return nil
}

func (api *API) processRewardEvent(evt types.Event) error {
	api.logger.Debug().Msg("processRewardEvent")
	var rewards types.EventRewards
	err := json.Unmarshal(evt.Event, &rewards)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal rewards event")
	}
	record := models.NewRewardEvent(rewards, evt)
	err = api.store.CreateRewardRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create rewards record")
	}
	return nil
}

func (api *API) processAddEvent(evt types.Event) error {
	api.logger.Debug().Msg("processAddEvent")
	var add types.EventAdd
	err := json.Unmarshal(evt.Event, &add)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal add event")
	}
	record := models.NewAddEvent(add, evt)
	err = api.store.CreateAddRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create add record")
	}
	return nil
}

func (api *API) processPoolEvent(evt types.Event) error {
	api.logger.Debug().Msg("processPoolEvent")
	var pool types.EventPool
	err := json.Unmarshal(evt.Event, &pool)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal pool event")
	}
	record := models.NewPoolEvent(pool, evt)
	err = api.store.CreatePoolRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create pool record")
	}
	return nil
}

func (api *API) processGasEvent(evt types.Event) error {
	api.logger.Debug().Msg("processGasEvent")
	var gas types.EventGas
	err := json.Unmarshal(evt.Event, &gas)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal gas event")
	}
	record := models.NewGasEvent(gas, evt)
	err = api.store.CreateGasRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create gas record")
	}
	return nil
}
func (api *API) processRefundEvent(evt types.Event) error {
	api.logger.Debug().Msg("processRefundEvent")
	var refund types.EventRefund
	err := json.Unmarshal(evt.Event, &refund)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal refund event")
	}
	record := models.NewRefundEvent(refund, evt)
	err = api.store.CreateRefundRecord(record)
	if err != nil {
		return errors.Wrap(err, "failed to create refund record")
	}
	return nil
}

// StartScan start to scan
func (api *API) StartScan() error {
	api.logger.Info().Msg("start thorchain event scanning")
	if !api.cfg.EnableScan {
		api.logger.Debug().Msg("Scan not enabled.")
		return nil
	}
	api.wg.Add(1)
	go api.scan()
	return nil
}

func (api *API) scan() {
	api.logger.Info().Msg("getting thorchain genesis")
	genesisTime, err := api.getGenesis()
	if err != nil {
		api.logger.Error().Err(err).Msg("failed to get genesis from thorchain")
	}

	err = api.processGenesis(genesisTime)
	if err != nil {
		api.logger.Error().Err(err).Msg("failed to set genesis in db")
	}
	api.logger.Info().Msg("processed thorchain genesis")

	defer api.wg.Done()

	api.logger.Info().Msg("start thorchain event scanning")
	defer api.logger.Info().Msg("thorchain event scanning stopped")
	currentPos := int64(1) // we start from 1
	maxID, err := api.store.GetMaxID()
	if err != nil {
		api.logger.Error().Err(err).Msg("failed to get currentPos from data store")
	} else {
		api.logger.Info().Int64("previous pos", maxID).Msg("find previous maxID")
		currentPos = maxID + 1
	}
	for {
		api.logger.Debug().Msg("sleeping thorchain scan")
		time.Sleep(time.Second * 1)
		select {
		case <-api.stopChan:
			return
		default:
			api.logger.Debug().Int64("currentPos", currentPos).Msg("request events")
			maxID, events, err := api.processEvents(currentPos)
			if err != nil {
				api.logger.Error().Err(err).Msg("failed to get events from thorchain")
				continue // we will retry a bit later
			}
			if events == 0 { // nothing in it
				select {
				case <-api.stopChan:
				case <-time.After(api.cfg.NoEventsBackoff):
					api.logger.Debug().Str("NoEventsBackoff", api.cfg.NoEventsBackoff.String()).Msg("Finished executing NoEventsBackoff")
				}
				continue
			}
			currentPos = maxID + 1
		}
	}
}

func (api *API) StopScan() error {
	api.logger.Info().Msg("stop scan request received")
	close(api.stopChan)
	api.wg.Wait()

	return nil
}

//Query output transaction for a given event from THORNode
func (api *API) GetOutTx(event types.Event) (common.Txs, error) {
	if event.InTx.ID.IsEmpty() {
		return nil, nil
	}
	uri := fmt.Sprintf("%s/keysign/%d", api.baseUrl, event.Height)
	api.logger.Debug().Msg(uri)
	resp, err := api.netClient.Get(uri)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); nil != err {
			api.logger.Error().Err(err).Msg("failed to close response body")
		}
	}()

	var chainTxout types.QueryResTxOut
	if err := json.NewDecoder(resp.Body).Decode(&chainTxout); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal chainTxout")
	}
	var outTxs common.Txs
	for _, chain := range chainTxout.Chains {
		for _, tx := range chain.TxArray {
			if tx.InHash == event.InTx.ID {
				outTx := common.Tx{
					ID:        tx.OutHash,
					ToAddress: tx.ToAddress,
					Memo:      tx.Memo,
					Chain:     tx.Chain,
					Coins: common.Coins{
						tx.Coin,
					},
				}
				if outTx.ID.IsEmpty() {
					outTx.ID = common.UnknownTxID
				}
				outTxs = append(outTxs, outTx)
			}
		}
	}
	return outTxs, nil
}
