// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// AssetDetail defines model for AssetDetail.
type AssetDetail struct {
	Asset       *Asset   `json:"asset,omitempty"`
	DateCreated *int64   `json:"dateCreated,omitempty"`
	Logo        *string  `json:"logo,omitempty"`
	Name        *string  `json:"name,omitempty"`
	PriceRune   *float64 `json:"priceRune,omitempty"`
}

// Error defines model for Error.
type Error struct {
	Error string `json:"error"`
}

// PoolDetail defines model for PoolDetail.
type PoolDetail struct {
	Asset            *Asset   `json:"asset,omitempty"`
	AssetDepth       *int64   `json:"assetDepth,omitempty"`
	AssetROI         *float64 `json:"assetROI,omitempty"`
	AssetStakedTotal *int64   `json:"assetStakedTotal,omitempty"`
	BuyAssetCount    *int64   `json:"buyAssetCount,omitempty"`
	BuyFeeAverage    *int64   `json:"buyFeeAverage,omitempty"`
	BuyFeesTotal     *int64   `json:"buyFeesTotal,omitempty"`
	BuySlipAverage   *float64 `json:"buySlipAverage,omitempty"`
	BuyTxAverage     *int64   `json:"buyTxAverage,omitempty"`
	BuyVolume        *int64   `json:"buyVolume,omitempty"`
	PoolDepth        *int64   `json:"poolDepth,omitempty"`
	PoolFeeAverage   *int64   `json:"poolFeeAverage,omitempty"`
	PoolFeesTotal    *int64   `json:"poolFeesTotal,omitempty"`
	PoolROI          *float64 `json:"poolROI,omitempty"`
	PoolROI12        *float64 `json:"poolROI12,omitempty"`
	PoolSlipAverage  *float64 `json:"poolSlipAverage,omitempty"`
	PoolStakedTotal  *int64   `json:"poolStakedTotal,omitempty"`
	PoolTxAverage    *int64   `json:"poolTxAverage,omitempty"`
	PoolUnits        *int64   `json:"poolUnits,omitempty"`
	PoolVolume       *int64   `json:"poolVolume,omitempty"`
	PoolVolume24hr   *int64   `json:"poolVolume24hr,omitempty"`
	Price            *float64 `json:"price,omitempty"`
	RuneDepth        *int64   `json:"runeDepth,omitempty"`
	RuneROI          *float64 `json:"runeROI,omitempty"`
	RuneStakedTotal  *int64   `json:"runeStakedTotal,omitempty"`
	SellAssetCount   *int64   `json:"sellAssetCount,omitempty"`
	SellFeeAverage   *int64   `json:"sellFeeAverage,omitempty"`
	SellFeesTotal    *int64   `json:"sellFeesTotal,omitempty"`
	SellSlipAverage  *float64 `json:"sellSlipAverage,omitempty"`
	SellTxAverage    *int64   `json:"sellTxAverage,omitempty"`
	SellVolume       *int64   `json:"sellVolume,omitempty"`
	StakeTxCount     *int64   `json:"stakeTxCount,omitempty"`
	StakersCount     *int64   `json:"stakersCount,omitempty"`
	StakingTxCount   *int64   `json:"stakingTxCount,omitempty"`
	Status           *string  `json:"status,omitempty"`
	SwappersCount    *int64   `json:"swappersCount,omitempty"`
	SwappingTxCount  *int64   `json:"swappingTxCount,omitempty"`
	WithdrawTxCount  *int64   `json:"withdrawTxCount,omitempty"`
}

// Stakers defines model for Stakers.
type Stakers string

// StakersAddressData defines model for StakersAddressData.
type StakersAddressData struct {
	PoolsArray  *[]Asset `json:"poolsArray,omitempty"`
	TotalEarned *int64   `json:"totalEarned,omitempty"`
	TotalROI    *float64 `json:"totalROI,omitempty"`
	TotalStaked *int64   `json:"totalStaked,omitempty"`
}

// StakersAssetData defines model for StakersAssetData.
type StakersAssetData struct {
	Asset           *Asset   `json:"asset,omitempty"`
	AssetEarned     *int64   `json:"assetEarned,omitempty"`
	AssetROI        *float64 `json:"assetROI,omitempty"`
	AssetStaked     *int64   `json:"assetStaked,omitempty"`
	DateFirstStaked *int64   `json:"dateFirstStaked,omitempty"`
	PoolEarned      *int64   `json:"poolEarned,omitempty"`
	PoolROI         *float64 `json:"poolROI,omitempty"`
	PoolStaked      *int64   `json:"poolStaked,omitempty"`
	RuneEarned      *int64   `json:"runeEarned,omitempty"`
	RuneROI         *float64 `json:"runeROI,omitempty"`
	RuneStaked      *int64   `json:"runeStaked,omitempty"`
	StakeUnits      *int64   `json:"stakeUnits,omitempty"`
}

// StatsData defines model for StatsData.
type StatsData struct {
	DailyActiveUsers   *int64 `json:"dailyActiveUsers,omitempty"`
	DailyTx            *int64 `json:"dailyTx,omitempty"`
	MonthlyActiveUsers *int64 `json:"monthlyActiveUsers,omitempty"`
	MonthlyTx          *int64 `json:"monthlyTx,omitempty"`
	PoolCount          *int64 `json:"poolCount,omitempty"`
	TotalAssetBuys     *int64 `json:"totalAssetBuys,omitempty"`
	TotalAssetSells    *int64 `json:"totalAssetSells,omitempty"`
	TotalDepth         *int64 `json:"totalDepth,omitempty"`
	TotalEarned        *int64 `json:"totalEarned,omitempty"`
	TotalStakeTx       *int64 `json:"totalStakeTx,omitempty"`
	TotalStaked        *int64 `json:"totalStaked,omitempty"`
	TotalTx            *int64 `json:"totalTx,omitempty"`
	TotalUsers         *int64 `json:"totalUsers,omitempty"`
	TotalVolume        *int64 `json:"totalVolume,omitempty"`
	TotalVolume24hr    *int64 `json:"totalVolume24hr,omitempty"`
	TotalWithdrawTx    *int64 `json:"totalWithdrawTx,omitempty"`
}

// TxDetails defines model for TxDetails.
type TxDetails struct {
	Date    *int64  `json:"date,omitempty"`
	Events  *Event  `json:"events,omitempty"`
	Gas     *Gas    `json:"gas,omitempty"`
	Height  *int64  `json:"height,omitempty"`
	In      *Tx     `json:"in,omitempty"`
	Options *Option `json:"options,omitempty"`
	Out     *Tx     `json:"out,omitempty"`
	Pool    *Asset  `json:"pool,omitempty"`
	Status  *string `json:"status,omitempty"`
	Type    *string `json:"type,omitempty"`
}

// Asset defines model for asset.
type Asset struct {
	Chain  *string `json:"chain,omitempty"`
	Symbol *string `json:"symbol,omitempty"`
	Ticker *string `json:"ticker,omitempty"`
}

// Coin defines model for coin.
type Coin struct {
	Amount *int64 `json:"amount,omitempty"`
	Asset  *Asset `json:"asset,omitempty"`
}

// Coins defines model for coins.
type Coins []Coin

// Event defines model for event.
type Event struct {
	Fee        *int64   `json:"fee,omitempty"`
	Slip       *float64 `json:"slip,omitempty"`
	StakeUnits *int64   `json:"stakeUnits,omitempty"`
}

// Gas defines model for gas.
type Gas struct {
	Amount *int64 `json:"amount,omitempty"`
	Asset  *Asset `json:"asset,omitempty"`
}

// Option defines model for option.
type Option struct {
	Asymmetry           *float64 `json:"asymmetry,omitempty"`
	PriceTarget         *int64   `json:"priceTarget,omitempty"`
	WithdrawBasisPoints *int64   `json:"withdrawBasisPoints,omitempty"`
}

// Tx defines model for tx.
type Tx struct {
	Address *string `json:"address,omitempty"`
	Coins   *Coins  `json:"coins,omitempty"`
	Memo    *string `json:"memo,omitempty"`
	TxID    *string `json:"txID,omitempty"`
}

// AssetsDetailedResponse defines model for AssetsDetailedResponse.
type AssetsDetailedResponse AssetDetail

// GeneralErrorResponse defines model for GeneralErrorResponse.
type GeneralErrorResponse Error

// PoolsDetailedResponse defines model for PoolsDetailedResponse.
type PoolsDetailedResponse PoolDetail

// PoolsResponse defines model for PoolsResponse.
type PoolsResponse []Asset

// StakersAddressDataResponse defines model for StakersAddressDataResponse.
type StakersAddressDataResponse StakersAddressData

// StakersAssetDataResponse defines model for StakersAssetDataResponse.
type StakersAssetDataResponse StakersAssetData

// StakersResponse defines model for StakersResponse.
type StakersResponse []Stakers

// StatsResponse defines model for StatsResponse.
type StatsResponse StatsData

// TxDetailedResponse defines model for TxDetailedResponse.
type TxDetailedResponse []TxDetails

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get Asset Information// (GET /v1/assets/{asset})
	GetAssetInfo(ctx echo.Context, asset string) error
	// Get Documents// (GET /v1/doc)
	GetDocs(ctx echo.Context) error
	// Get Health// (GET /v1/health)
	GetHealth(ctx echo.Context) error
	// Get Asset Pools// (GET /v1/pools)
	GetPools(ctx echo.Context) error
	// Get Pools Data// (GET /v1/pools/{asset})
	GetPoolsData(ctx echo.Context, asset string) error
	// Get Stakers// (GET /v1/stakers)
	GetStakersData(ctx echo.Context) error
	// Get Staker Data// (GET /v1/stakers/{address})
	GetStakersAddressData(ctx echo.Context, address string) error
	// Get Staker Pool Data// (GET /v1/stakers/{address}/{asset})
	GetStakersAddressAndAssetData(ctx echo.Context, address string, asset string) error
	// Get Global Stats// (GET /v1/stats)
	GetStats(ctx echo.Context) error
	// Get Swagger// (GET /v1/swagger.json)
	GetSwagger(ctx echo.Context) error
	// Get the Proxied Pool Addresses// (GET /v1/thorchain/pool_addresses)
	GetThorchainProxiedEndpoints(ctx echo.Context) error
	// Get transaction// (GET /v1/tx/asset/{asset})
	GetTxDetailsByAsset(ctx echo.Context, asset string) error
	// Get transaction// (GET /v1/tx/{address})
	GetTxDetails(ctx echo.Context, address string) error
	// Get transaction// (GET /v1/tx/{address}/asset/{asset})
	GetTxDetailsByAddressAsset(ctx echo.Context, address string, asset string) error
	// Get transaction// (GET /v1/tx/{address}/txid/{txid})
	GetTxDetailsByAddressTxId(ctx echo.Context, address string, txid string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetAssetInfo converts echo context to params.
func (w *ServerInterfaceWrapper) GetAssetInfo(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "asset" -------------
	var asset string

	err = runtime.BindStyledParameter("simple", false, "asset", ctx.Param("asset"), &asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetAssetInfo(ctx, asset)
	return err
}

// GetDocs converts echo context to params.
func (w *ServerInterfaceWrapper) GetDocs(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetDocs(ctx)
	return err
}

// GetHealth converts echo context to params.
func (w *ServerInterfaceWrapper) GetHealth(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetHealth(ctx)
	return err
}

// GetPools converts echo context to params.
func (w *ServerInterfaceWrapper) GetPools(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPools(ctx)
	return err
}

// GetPoolsData converts echo context to params.
func (w *ServerInterfaceWrapper) GetPoolsData(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "asset" -------------
	var asset string

	err = runtime.BindStyledParameter("simple", false, "asset", ctx.Param("asset"), &asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPoolsData(ctx, asset)
	return err
}

// GetStakersData converts echo context to params.
func (w *ServerInterfaceWrapper) GetStakersData(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStakersData(ctx)
	return err
}

// GetStakersAddressData converts echo context to params.
func (w *ServerInterfaceWrapper) GetStakersAddressData(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStakersAddressData(ctx, address)
	return err
}

// GetStakersAddressAndAssetData converts echo context to params.
func (w *ServerInterfaceWrapper) GetStakersAddressAndAssetData(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// ------------- Path parameter "asset" -------------
	var asset string

	err = runtime.BindStyledParameter("simple", false, "asset", ctx.Param("asset"), &asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStakersAddressAndAssetData(ctx, address, asset)
	return err
}

// GetStats converts echo context to params.
func (w *ServerInterfaceWrapper) GetStats(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStats(ctx)
	return err
}

// GetSwagger converts echo context to params.
func (w *ServerInterfaceWrapper) GetSwagger(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetSwagger(ctx)
	return err
}

// GetThorchainProxiedEndpoints converts echo context to params.
func (w *ServerInterfaceWrapper) GetThorchainProxiedEndpoints(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetThorchainProxiedEndpoints(ctx)
	return err
}

// GetTxDetailsByAsset converts echo context to params.
func (w *ServerInterfaceWrapper) GetTxDetailsByAsset(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "asset" -------------
	var asset string

	err = runtime.BindStyledParameter("simple", false, "asset", ctx.Param("asset"), &asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTxDetailsByAsset(ctx, asset)
	return err
}

// GetTxDetails converts echo context to params.
func (w *ServerInterfaceWrapper) GetTxDetails(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTxDetails(ctx, address)
	return err
}

// GetTxDetailsByAddressAsset converts echo context to params.
func (w *ServerInterfaceWrapper) GetTxDetailsByAddressAsset(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// ------------- Path parameter "asset" -------------
	var asset string

	err = runtime.BindStyledParameter("simple", false, "asset", ctx.Param("asset"), &asset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter asset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTxDetailsByAddressAsset(ctx, address, asset)
	return err
}

// GetTxDetailsByAddressTxId converts echo context to params.
func (w *ServerInterfaceWrapper) GetTxDetailsByAddressTxId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "address" -------------
	var address string

	err = runtime.BindStyledParameter("simple", false, "address", ctx.Param("address"), &address)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter address: %s", err))
	}

	// ------------- Path parameter "txid" -------------
	var txid string

	err = runtime.BindStyledParameter("simple", false, "txid", ctx.Param("txid"), &txid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter txid: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTxDetailsByAddressTxId(ctx, address, txid)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router runtime.EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET("/v1/assets/:asset", wrapper.GetAssetInfo)
	router.GET("/v1/doc", wrapper.GetDocs)
	router.GET("/v1/health", wrapper.GetHealth)
	router.GET("/v1/pools", wrapper.GetPools)
	router.GET("/v1/pools/:asset", wrapper.GetPoolsData)
	router.GET("/v1/stakers", wrapper.GetStakersData)
	router.GET("/v1/stakers/:address", wrapper.GetStakersAddressData)
	router.GET("/v1/stakers/:address/:asset", wrapper.GetStakersAddressAndAssetData)
	router.GET("/v1/stats", wrapper.GetStats)
	router.GET("/v1/swagger.json", wrapper.GetSwagger)
	router.GET("/v1/thorchain/pool_addresses", wrapper.GetThorchainProxiedEndpoints)
	router.GET("/v1/tx/asset/:asset", wrapper.GetTxDetailsByAsset)
	router.GET("/v1/tx/:address", wrapper.GetTxDetails)
	router.GET("/v1/tx/:address/asset/:asset", wrapper.GetTxDetailsByAddressAsset)
	router.GET("/v1/tx/:address/txid/:txid", wrapper.GetTxDetailsByAddressTxId)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"xPK5aVnw6fnfAQAA//+H/Mz6PkIAAA==",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
