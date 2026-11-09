package auction

import (
	"context"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateAuction_ShouldCloseAuctionAfterDuration(t *testing.T) {
	ctx := context.Background()

	if os.Getenv("MONGODB_URL") == "" {
		t.Setenv("MONGODB_URL", "mongodb://admin:admin@localhost:27017/?authSource=admin")
	}
	t.Setenv("MONGODB_DB", "auctions_test")
	t.Setenv("AUCTION_DURATION", "2s")

	database, err := mongodb.NewMongoDBConnection(ctx)
	assert.Nil(t, err)
	defer database.Drop(ctx)

	repository := NewAuctionRepository(database)

	auction, internalErr := auction_entity.CreateAuction(
		"Notebook", "Eletronicos", "Notebook usado em bom estado", auction_entity.Used)
	assert.Nil(t, internalErr)

	internalErr = repository.CreateAuction(ctx, auction)
	assert.Nil(t, internalErr)

	auctionFound, internalErr := repository.FindAuctionById(ctx, auction.Id)
	assert.Nil(t, internalErr)
	assert.Equal(t, auction_entity.Active, auctionFound.Status)

	time.Sleep(3 * time.Second)

	auctionFound, internalErr = repository.FindAuctionById(ctx, auction.Id)
	assert.Nil(t, internalErr)
	assert.Equal(t, auction_entity.Completed, auctionFound.Status)
}
