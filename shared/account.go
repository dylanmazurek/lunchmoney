package shared

type Account struct {
	AssetSource     string `bson:"assetSource"`
	ExternalAssetID string `bson:"externalAssetId"`
	ExternalUserID  string `bson:"externalUserId"`
	AssetName       string `bson:"assetName"`
	Currency        string `bson:"currency"`

	AssetID int64 `bson:"assetId"`
}
