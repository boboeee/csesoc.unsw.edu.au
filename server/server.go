package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// H - interface for sending JSON
type H map[string]interface{}

// Post - struct to contain post data
type Post struct {
	PostID           int
	PostTitle        string
	PostSubtitle     string
	PostType         string
	PostCategory     int
	CreatedOn        int64
	LastEditedOn     int64
	PostContent      string
	PostLinkGithub   string
	PostLinkFacebook string
	ShowInMenu       bool
}

// Category - struct to contain category data
type Category struct {
	CategoryID   int
	CategoryName string
	Index        int
}

// User - struct to contain user data
type User struct {
	UserID    string //sha256 the zid
	UserToken string
	Role      string
}

// Sponsor - struct to contain sponsor data
type Sponsor struct {
	SponsorID   uuid.UUID
	SponsorName string
	SponsorLogo string
	SponsorTier string
	Expiry      int64
}

// Claims - struct to store jwt data
type Claims struct {
	HashedZID   [32]byte
	FirstName   string
	Permissions string
	jwt.StandardClaims
}

func main() {
	// Create new instance of echo
	e := echo.New()

	// Get flags
	dbURI := flag.String("dburi", "mongodb://127.0.0.1:27017", "The mongodb instance URI in the form: mongodb://ip:port")
	distPath := flag.String("dist", "../dist/",
		"Path to the built app distribution (e.g. using yarn build --mode production)")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	dbClient := setupDatabase(*dbURI)

	servePages(e, *distPath)
	serveAPI(e, dbClient)

	e.Logger.Fatal(e.Start(":1323"))
}

func servePages(e *echo.Echo, distPath string) {

	// Setup our assetHandler and point it to our static build location
	assetHandler := http.FileServer(http.Dir(distPath))

	// Setup a new echo route to load the build as our base path
	e.GET("/", echo.WrapHandler(assetHandler))

	e.GET("/favicon.ico", echo.WrapHandler(assetHandler))
	e.GET("/js/*", echo.WrapHandler(assetHandler))
	e.GET("/css/*", echo.WrapHandler(assetHandler))
	e.GET("/img/*", echo.WrapHandler(assetHandler))

	echo.NotFoundHandler = func(c echo.Context) error {
		// render your 404 page
		return c.String(http.StatusNotFound, "not found page")
	}
}

func setupDatabase(dbURI string) *mongo.Client {
	//Set client options
	clientOptions := options.Client().ApplyURI(dbURI)
	//Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	//Check connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func serveAPI(e *echo.Echo, dbClient *mongo.Client) {

	// Creating collections
	postsCollection := dbClient.Database("csesoc").Collection("posts")
	catCollection := dbClient.Database("csesoc").Collection("categories")
	sponsorCollection := dbClient.Database("csesoc").Collection("sponsors")
	// userCollection := dbClient.Database("csesoc").Collection("users")

	// e.POST("/login/", login(userCollection))

	// Routes for posts
	e.GET("/posts/", getPosts(postsCollection))
	e.POST("/post/", newPosts(postsCollection))
	e.PUT("/post/", updatePosts(postsCollection))
	e.DELETE("/post/", deletePosts(postsCollection))

	// Routes for categories
	e.GET("/category/:id/", getCategories(catCollection))
	e.GET("/category/", getAllCategories(catCollection))
	e.POST("/category/", newCategories(catCollection))
	e.PATCH("/category/", patchCategories(catCollection))
	e.DELETE("/category/", deleteCategories(catCollection))

	// Routes for sponsors
	e.POST("/sponsor/", newSponsors(sponsorCollection))
	e.DELETE("/sponsor/", deleteSponsors(sponsorCollection))
}

// func login(collection *mongo.Collection) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		zid := c.FormValue("zid")
// 		password := c.FormValue("password")
// 		permissions := c.FormValue("permissions")
// 		tokenString := Auth(collection, zid, password, permissions)
// 		return c.JSON(http.StatusOK, H{
// 			"token": tokenString,
// 		})
// 	}
// }

func getPosts(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.QueryParam("id")
		count, _ := strconv.Atoi(c.QueryParam("nPosts"))
		category := c.QueryParam("category")
		if id == "" {
			posts := GetAllPosts(collection, count, category)
			return c.JSON(http.StatusOK, H{
				"post": posts,
			})
		}

		idInt, _ := strconv.Atoi(id)
		result := GetPosts(collection, idInt, category)
		return c.JSON(http.StatusOK, H{
			"post": result,
		})
	}
}

func newPosts(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.FormValue("id"))
		category, _ := strconv.Atoi(c.FormValue("category"))
		showInMenu, _ := strconv.ParseBool(c.FormValue("showInMenu"))
		title := c.FormValue("title")
		subtitle := c.FormValue("subtitle")
		postType := c.FormValue("type")
		content := c.FormValue("content")
		github := c.FormValue("linkGithub")
		fb := c.FormValue("linkFacebook")
		NewPosts(collection, id, category, showInMenu, title, subtitle, postType, content, github, fb)
		return c.JSON(http.StatusOK, H{})
	}
}

func updatePosts(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.FormValue("id"))
		category, _ := strconv.Atoi(c.FormValue("category"))
		showInMenu, _ := strconv.ParseBool(c.FormValue("showInMenu"))
		title := c.FormValue("title")
		subtitle := c.FormValue("subtitle")
		postType := c.FormValue("type")
		content := c.FormValue("content")
		github := c.FormValue("linkGithub")
		fb := c.FormValue("linkFacebook")
		UpdatePosts(collection, id, category, showInMenu, title, subtitle, postType, content, github, fb)
		return c.JSON(http.StatusOK, H{})
	}
}

func deletePosts(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.FormValue("id"))
		DeletePosts(collection, id)
		return c.JSON(http.StatusOK, H{})
	}
}

func getAllCategories(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.FormValue("token")
		count, err := strconv.Atoi(c.QueryParam("count"))

		result, err := GetAllCategories(collection, count, token)

		if err != nil {
			return c.JSON(http.StatusBadRequest, H{
				"error": "Couldn't get all categories",
			})
		}

		return c.JSON(http.StatusOK, H{
			"categories": result,
		})
	}
}

func getCategories(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.FormValue("token")
		id, _ := strconv.Atoi(c.QueryParam("id"))

		result, err := GetCategories(collection, id, token)

		if err != nil {
			return c.JSON(http.StatusBadRequest, H{
				"error": "Couldn't get a category by ID",
			})
		}

		return c.JSON(http.StatusOK, H{
			"category": result,
		})

	}
}

func newCategories(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.FormValue("token")
		catID, _ := strconv.Atoi(c.FormValue("id"))
		index, _ := strconv.Atoi(c.FormValue("index"))
		name := c.FormValue("name")
		NewCategories(collection, catID, index, name, token)
		return c.JSON(http.StatusOK, H{})
	}
}

func patchCategories(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.FormValue("token")
		catID, _ := strconv.Atoi(c.FormValue("id"))
		name := c.FormValue("name")
		index, _ := strconv.Atoi(c.FormValue("index"))
		PatchCategories(collection, catID, name, index, token)
		return c.JSON(http.StatusOK, H{})
	}
}

func deleteCategories(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.FormValue("token")
		id, _ := strconv.Atoi(c.FormValue("id"))
		DeleteCategories(collection, id, token)
		return c.JSON(http.StatusOK, H{})
	}
}

func newSponsors(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.FormValue("token")
		expiryStr := c.FormValue("expiry")
		name := c.FormValue("name")
		logo := c.FormValue("logo")
		tier := c.FormValue("tier")
		NewSponsors(collection, expiryStr, name, logo, tier, token)
		return c.JSON(http.StatusOK, H{})
	}
}

func deleteSponsors(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.FormValue("token")
		id := c.FormValue("id")
		DeleteSponsors(collection, id, token)
		return c.JSON(http.StatusOK, H{})
	}
}
