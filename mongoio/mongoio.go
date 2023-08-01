package mongoio

import (
	"context"
	"fmt"

	"github.com/quinn-caverly/forward-utils/endpointstructs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateConnToDatabase() (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo-service:27017"))
	if err != nil {
		return nil, nil, fmt.Errorf("Was not able to connect to mongo via the service %w", err)
	}

}

func CreateConnToBrand(brand string) (*mongo.Collection, *mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo-service:27017"))
	if err != nil {
		return nil, nil, fmt.Errorf("Was not able to connect to mongo via the service %w", err)
	}
	collection := client.Database("products").Collection(brand)

	return collection, client, nil
}

// a) if the id already exists in the db, check if any of the colors are not currently added.
// b) if id not in db, write all the colors
// returned are the colors which need to be written to the filesystem.
func WriteUPE(pc *endpointstructs.UniqueProductExpanded, coll *mongo.Collection) ([]endpointstructs.URLColorContainer, error) {

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

		return pc.URLColorContainers, nil
	} else if err != nil {
		return nil, fmt.Errorf("Error when attempting to find document of id in collection, %w", err)
	} else {

		toAdd := []endpointstructs.URLColorContainer{}
	outerLoop:
		for i := range pc.URLColorContainers {
			for j := range result.URLColorContainers {
				if pc.URLColorContainers[i].ColorAttr.ColorName == result.URLColorContainers[j].ColorAttr.ColorName {
					continue outerLoop
				}
			}
			toAdd = append(toAdd, pc.URLColorContainers[i])
		}

		result.URLColorContainers = append(result.URLColorContainers, toAdd...)

		//TODO it would be preferable to use Update here but inserting structs isnt well documented as far as I know
		_, err := coll.ReplaceOne(context.TODO(), bson.M{"_id": pc.Id}, result)
		if err != nil {
			return nil, fmt.Errorf("Error when replacing struct with updated ColorAttrs %w", err)
		}

		return toAdd, nil
	}
}

func ReadUPE(id string, coll *mongo.Collection) (endpointstructs.UniqueProductExpanded, error) {

	var result endpointstructs.UniqueProductExpanded
	err := coll.FindOne(
		context.TODO(),
		bson.D{{Key: "_id", Value: id}},
	).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return endpointstructs.UniqueProductExpanded{}, fmt.Errorf(fmt.Sprint("The given id: ", id, " is not in the database of the ", coll.Name(), " collection, "), err)
	} else if err != nil {
		return endpointstructs.UniqueProductExpanded{}, err
	}

	return result, nil
}
