package functions

import (
	"math"

	"github.com/Rhymond/go-money"
	"github.com/dylanmazurek/lunchmoney"
	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/shared"
	"github.com/rs/zerolog/log"
)

func AssetHandler(lma *lunchmoney.Client, asset *shared.Asset) {
	balanceFloat, _ := asset.Balance.Float64()
	balance := money.NewFromFloat(math.Abs(balanceFloat), *asset.Currency)

	lmAsset := &models.Asset{
		AssetID:     &asset.AssetID,
		Balance:     *balance,
		BalanceAsOf: asset.BalanceAsOf,
	}

	log.Debug().Msg("updated asset")

	updatedAsset, err := lma.UpdateAsset(asset.AssetID, lmAsset)
	if err != nil || updatedAsset.Error != nil {
		log.Error().Err(err).Msg("unable to update asset")
	}
}
