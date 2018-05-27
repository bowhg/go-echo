package main

import (
	"net/http"

	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	mongo_host := "172.17.0.2:27017"
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	//query all
	e.GET("/posts", func(c echo.Context) (err error) {
		session, err := mgo.Dial(mongo_host)
		defer session.Close()
		s := session.DB("facebook_feed_me").C("streams")
		var posts []FacebookPost
		err = s.Find(bson.M{}).All(&posts)
		if err != nil {
			panic(err)
		}
		return c.JSON(http.StatusOK, posts)
	})
	//query between date
	e.GET("/post", func(c echo.Context) (err error) {
		session, err := mgo.Dial(mongo_host)
		defer session.Close()
		s := session.DB("facebook_feed_me").C("streams")
		fromDate := c.QueryParam("from")
		toDate := c.QueryParam("to")
		var posts []FacebookPost
		err = s.Find(bson.M{"created_time": bson.M{"$gte": fromDate, "$lt": toDate}}).All(&posts)
		if err != nil {
			panic(err)
		}
		return c.JSON(http.StatusOK, posts)
	})
	//query from
	e.GET("/post/:id", func(c echo.Context) (err error) {
		session, err := mgo.Dial(mongo_host)
		defer session.Close()
		s := session.DB("facebook_feed_me").C("streams")
		id := c.Param("id")
		var post FacebookPost
		err = s.Find(bson.M{"id": id}).One(&post)
		if err != nil {
			panic(err)
		}
		return c.JSON(http.StatusOK, post)
	})
	e.Logger.Fatal(e.Start(":1323"))
	//try to print query
	/*session, err := mgo.Dial(mongo_host)
	defer session.Close()
	c := session.DB("facebook_feed_me").C("streams")
	var posts []FacebookPost
	err = c.Find(bson.M{}).All(&posts)

	if err != nil {
		panic(err)
	}
	fmt.Println("Results All: ", posts)
	*/
}

type FacebookPost struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Message      string        `json:"message" bson:"message,omitempty"`
	Created_time string        `json:"created_time" bson:"created_time,omitempty"`
	Post_id      string        `json:"post_id" bson:"id,omitempty"`
}
