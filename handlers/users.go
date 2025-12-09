package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/coddicted/go-mongo-crud/db"
	"github.com/coddicted/go-mongo-crud/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// helper to get collection
func getCollection() *mongo.Collection {
	collName := "users" // or read from env if you want
	return db.Client.Database("testdb").Collection(collName)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = primitive.NewObjectID()
	coll := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := coll.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	coll := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	idParam := params["id"]
	oid, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	coll := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	idParam := params["id"]
	oid, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var update models.User
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coll := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// create update document (only fields provided will be updated)
	updateDoc := bson.M{}
	if update.FirstName != "" {
		updateDoc["first_name"] = update.FirstName
	}
	if update.LastName != "" {
		updateDoc["last_name"] = update.LastName
	}
	if update.Email != "" {
		updateDoc["email"] = update.Email
	}
	if update.Age != 0 {
		updateDoc["age"] = update.Age
	}

	if len(updateDoc) == 0 {
		http.Error(w, "no fields provided to update", http.StatusBadRequest)
		return
	}

	res, err := coll.UpdateByID(ctx, oid, bson.M{"$set": updateDoc}, options.Update().SetUpsert(false))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.MatchedCount == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// return updated doc
	var updated models.User
	if err := coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updated)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	idParam := params["id"]
	oid, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	coll := getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.DeletedCount == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// w.WriteHeader(http.StatusNoContent) // 204
	json.NewEncoder(w).Encode(map[string]string{"message": "user deleted"})
}
