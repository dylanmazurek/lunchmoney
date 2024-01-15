package handlers

import (
	"math"
	"strings"

	"github.com/Rhymond/go-money"
	"github.com/dylanmazurek/lunchmoney"
	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/shared"
	"github.com/rs/zerolog/log"
)

func AssetHandler(lma *lunchmoney.Client, asset *shared.Asset) {
	balanceFloat, err := asset.Balance.Float64()
	if err != nil {
		log.Error().
			Err(err).
			Str("ext-asset-id", asset.ExternalAssetID).
			Msg("unable to update asset")
	}

	if asset.AssetID == nil {
		log.Error().
			Str("ext-asset-id", asset.ExternalAssetID).
			Msg("unable to update asset, asset id not set")
	}

	assetId := *asset.AssetID

	currency := strings.ToUpper(asset.Currency)

	balance := money.NewFromFloat(math.Abs(balanceFloat), currency)

	lmAsset := &models.Asset{
		AssetID:     &assetId,
		Balance:     *balance,
		BalanceAsOf: asset.BalanceAsOf,
	}

	log.Info().
		Str("ext-asset-id", asset.ExternalAssetID).
		Int64("int-asset-id", assetId).
		Msg("updated asset")

	updatedAsset, err := lma.UpdateAsset(*asset.AssetID, lmAsset)
	if err != nil || updatedAsset.Error != nil {
		log.Error().
			Str("ext-asset-id", asset.ExternalAssetID).
			Int64("int-asset-id", assetId).
			Err(err).Msg("unable to update asset")
	}
}
