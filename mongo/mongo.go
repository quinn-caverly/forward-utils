package mongo

import (
	"context"
	"fmt"

	"github.com/quinn-caverly/forward-utils/endpointstructs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateConnToBrand(brand string) (*mongo.Collection, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongo-service:27017"))
	if err != nil {
        return nil, fmt.Errorf("Was not able to connect to mongo via the service %w", err)
	}
	defer client.Disconnect(context.Background())
	collection := client.Database("products").Collection(brand)

	return collection, nil
}

// a) if the id already exists in the db, check if any of the colors are not currently added.
// b) if id not in db, write all the colors
// the return are the colors which need to be written to the filesystem.
func WriteUPE(pc endpointstructs.UniqueProductExpanded, coll *mongo.Collection) error {

}