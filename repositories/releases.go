package repositories

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/mribichich/release-mgmt/entities"
)

// var currentId int
// var releases Releases
var session mgo.DialInfo
var releasesCollection *mgo.Collection

// Give us some seed data
func init() {
	// RepoCreateRelease(Release{Name: "Write presentation"})
	// RepoCreateRelease(Release{Name: "Host meetup"})

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	// defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	releasesCollection = session.DB("test").C("releases")

	// err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
	// 	&Person{"Cla", "+55 53 8402 8510"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// result := Person{}
	// err = c.Find(bson.M{"name": "Ale"}).One(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Phone:", result.Phone)
}

func FindAll() entities.Releases {
	result := entities.Releases{}

	err := releasesCollection.Find(bson.M{}).All(&result)

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("result:", result)

	return result
}

func RepoFindRelease(id bson.ObjectId) entities.Release {
	// for _, t := range releases {
	// 	if t.Id == id {
	// 		return t
	// 	}
	// }

	result := entities.Release{}

	err := releasesCollection.FindId(id).One(&result)

	if err != nil {
		log.Fatal(err)
	}

	// return empty Release if not found
	return entities.Release{}
}

func RepoCreateRelease(t entities.Release) entities.Release {
	// currentId += 1
	// t.Id = currentId
	// releases = append(releases, t)
	t.Id = bson.NewObjectId()
	t.Timestamp = time.Now()
	err := releasesCollection.Insert(t)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func RepoDestroyRelease(id bson.ObjectId) error {
	return releasesCollection.RemoveId(id)

	// for i, t := range releases {
	// 	if t.Id == id {
	// 		releases = append(releases[:i], releases[i+1:]...)
	// 		return nil
	// 	}
	// }

	// return fmt.Errorf("Could not find Release with id of %d to delete", id)
}
