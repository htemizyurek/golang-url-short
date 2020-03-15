package findURL

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*MONGODBNAME mongo database name */
/*MONGOHOST mongo database host address */
/*MONGODBCOLLECTINNAME mongo database url collenction name */
const (
	MONGODBNAME          string = "url-shortener"
	MONGOHOST            string = "localhost:27017"
	MONGODBCOLLECTINNAME string = "urls"
)

/*MongoData database fields*/
type MongoData struct {
	URL string
}

/*Find mongoDB url find*/
func Find(value string) string {
	clientOptions := options.Client().ApplyURI("mongodb://local-mongo-user:1234567@" + MONGOHOST + "/" + MONGODBNAME + "?retryWrites=true&w=majority&authMechanism=SCRAM-SHA-256")
	clientOptions.SetConnectTimeout(20 * time.Second)
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal("client error: ", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("ping error: ", err)
	}

	log.Println("Connected to MongoDB!")

	var result MongoData
	filter := bson.M{"value": value}
	collection := client.Database(MONGODBNAME).Collection(MONGODBCOLLECTINNAME)
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("founded url:", result.URL)

	return result.URL
}
