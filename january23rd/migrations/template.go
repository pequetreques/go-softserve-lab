package migrations

// migrate "github.com/eminetto/mongo-migrate"
// "github.com/globalsign/mgo"
// "github.com/globalsign/mgo/bson"
// "go.mongodb.org/mongo-driver/bson"

func init() {
	// migrate.Register(func(db *mgo.Database) error { // Up
	// 	return db.C("student").Update(
	// 			bson.M{"email": "peque@somemail.com"},
	// 			bson.M{"$set": bson.M{"type": 1}},
	// 		), func(db *mgo.Database) error { // Down
	// 			return db.C("student").Update(
	// 				bson.M{"email": "peque@somemail.com"},
	// 				bson.M{"$set": bson.M{"type": 0}},
	// 			)
	// 		}
	// })
}
