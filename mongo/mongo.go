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
// returned are the colors which need to be written to the filesystem.
func WriteUPE(pc *endpointstructs.UniqueProductExpanded, coll *mongo.Collection) ([]endpointstructs.ColorAttr, error) {

	var result endpointstructs.UniqueProductExpanded
	err := coll.FindOne(
		context.TODO(),
		bson.D{{Key: "_id", Value: pc.Id}},
	).Decode(&result)

	if err == mongo.ErrNoDocuments {
		_, err = coll.InsertOne(context.TODO(), pc)
		if err != nil {
			return nil, fmt.Errorf("Error when inserting into mongo %w", err)
		}

		return pc.ColorAttrs, nil
	} else {

		toAdd := []endpointstructs.ColorAttr{}
    	outerLoop:
    		for i := range pc.ColorAttrs {
    			for j := range result.ColorAttrs {
    				if pc.ColorAttrs[i].ColorName == result.ColorAttrs[j].ColorName {
    					continue outerLoop
    				}
    			}
    			toAdd = append(toAdd, pc.ColorAttrs[i])
    		}

		result.ColorAttrs = append(result.ColorAttrs, toAdd...)

		//TODO it would be preferable to use Update here but inserting structs isnt well documented as far as I know
		_, err := coll.ReplaceOne(context.Background(), bson.D{{Key: "_id", Value: pc.Id}}, result)
		if err != nil {
			return nil, fmt.Errorf("Error when replacing struct with updated ColorAttrs %w", err)
		}

		return toAdd, nil
	}
}
