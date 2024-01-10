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
	balanceFloat, _ := asset.Balance.Float64()
	currency := strings.ToUpper(asset.Currency)

	balance := money.NewFromFloat(math.Abs(balanceFloat), currency)

	lmAsset := &models.Asset{
		AssetID:     asset.AssetID,
		Balance:     *balance,
		BalanceAsOf: asset.BalanceAsOf,
	}

	log.Info().
		Str("externalId", asset.ExternalAssetID).
		Int64("assetId", *asset.AssetID).
		Msg("updated asset")

	updatedAsset, err := lma.UpdateAsset(*asset.AssetID, lmAsset)
	if err != nil || updatedAsset.Error != nil {
		log.Error().
			Str("externalId", asset.ExternalAssetID).
			Int64("assetId", *asset.AssetID).
			Err(err).Msg("unable to update asset")
	}
}
