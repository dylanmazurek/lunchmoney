package lunchmoney

import (
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()

	lmClient, _ := NewClient(os.Getenv("LUNCHMONEY_API_KEY"))

	assets, _ := lmClient.GetAssets(ctx)

	for _, asset := range *assets {
		fmt.Printf("asset: %s", asset.Name)
	}
}
