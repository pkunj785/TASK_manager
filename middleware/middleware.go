package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"com.serve_volt/models"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	loadTheEnv()
	crateDBIntance()
}

func loadTheEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading the .env file")
	}
}

func crateDBIntance() {
	connectionURI := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("DB_COLLECTION_NAME")

	clientOpt := options.Client().ApplyURI(connectionURI)
	client, err := mongo.Connect(context.TODO(), clientOpt)

	if err != nil {
		log.Fatalln(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("connected to mongodb!")

	collection = client.Database(dbName).Collection(collectionName)
	fmt.Println("collection instance created")
}


func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	payload := getAllTasks()
	json.NewEncoder(w).Encode(payload)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	var task models.TaskList
	json.NewDecoder(r.Body).Decode(&task)
	insertOneTask(task)
	json.NewEncoder(w).Encode(task)
}

func TaskComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","PUT")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")

	params := chi.URLParam(r, "id")
	taskComplete(params)
	json.NewEncoder(w).Encode(params)
}

func UndoTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","PUT")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")

	params := chi.URLParam(r, "id")
	undoTask(params)
	json.NewEncoder(w).Encode(params)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","DELETE")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	
	params := chi.URLParam(r, "id")
	deleteOneTask(params)

}

func DeleteAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	count := deleteAllTasks()
	json.NewEncoder(w).Encode(count)
}


///////////////////////////////////////////

func getAllTasks() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatalln(err)
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatalln(err)
	}

	cur.Close(context.Background())
	return results
}

func taskComplete(task string)  {
	id,_ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status": true }}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("modified count:", result.ModifiedCount)
}

func insertOneTask(task models.TaskList){
	insertResult, err := collection.InsertOne(context.Background(), task)
	
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Insertes a single record", insertResult.InsertedID)
}

func undoTask(task string) {
	id,_ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status": false }}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("modified count:", result.ModifiedCount)
}

func deleteOneTask(task string)  {
	id,_ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	del, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatalln(err)
	}
	
	fmt.Println("Deleted Document", del.DeletedCount)
}

func deleteAllTasks() int64{
	del, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Deleted Document", del.DeletedCount)
	return del.DeletedCount
}