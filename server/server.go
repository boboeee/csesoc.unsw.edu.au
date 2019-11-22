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
	PostID        int
	PostTitle     string
	PostSubtitle  string
	PostType      string
	PostCategory  int
	CreatedOn     int64
	LastEditedOn  int64
	PostContent   string
	CanonicalLink string
	ImageLink     string
	ResourceLink  string
	ShowInMenu    bool
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
	SponsorLink string
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

	e.Logger.Fatal(e.Start(":80"))
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

	// Create a new API subroute
	apiRoute := e.Group("/api/v1/")

	// apiRoute.POST("/login/", login(userCollection))

	// Routes for posts
	apiRoute.GET("posts", getPosts(postsCollection))
	apiRoute.POST("post", newPosts(postsCollection))
	apiRoute.PUT("post", updatePosts(postsCollection))
	apiRoute.DELETE("post", deletePosts(postsCollection))

	// Routes for categories
	apiRoute.GET("category/:id", getCategories(catCollection))
	apiRoute.GET("category", getAllCategories(catCollection))
	apiRoute.POST("category", newCategories(catCollection))
	apiRoute.PATCH("category", patchCategories(catCollection))
	apiRoute.DELETE("category", deleteCategories(catCollection))

	// Routes for sponsors
	apiRoute.GET("sponsor", getAllSponsors(sponsorCollection))
	apiRoute.POST("sponsor", newSponsors(sponsorCollection))
	apiRoute.DELETE("sponsor", deleteSponsors(sponsorCollection))

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
		category, _ := strconv.Atoi(c.QueryParam("category"))

		if id == "" {
			posts, err := GetAllPosts(collection, count, category)
			if err != nil {
				return c.JSON(http.StatusBadRequest, H{
					"error": "Couldn't get all posts",
				})
			}

			return c.JSON(http.StatusOK, H{
				"posts": posts,
			})

		}

		idInt, _ := strconv.Atoi(id)
		result, err := GetPosts(collection, idInt, category)

		if err != nil {
			return c.JSON(http.StatusBadRequest, H{
				"error": "Couldn't get post with that ID",
			})
		}

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
		imageLink := c.FormValue("imageLink")
		resourceLink := c.FormValue("resourceLink")
		canonicalLink := c.FormValue("canonicalLink")
		err := NewPosts(collection, id, category, showInMenu, title, subtitle, postType, content, imageLink, resourceLink, canonicalLink)
		if err != nil {
			return c.JSON(http.StatusBadRequest, H{
				"error": "Couldn't create that post with the given data",
			})
		}

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
		imageLink := c.FormValue("imageLink")
		resourceLink := c.FormValue("resourceLink")
		canonicalLink := c.FormValue("canonicalLink")
		err := UpdatePosts(collection, id, category, showInMenu, title, subtitle, postType, content, imageLink, resourceLink, canonicalLink)
		if err != nil {
			return c.JSON(http.StatusBadRequest, H{
				"error": "Couldn't update that post with the given data",
			})
		}

		return c.JSON(http.StatusOK, H{})
	}
}

func deletePosts(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.FormValue("id"))
		err := DeletePosts(collection, id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, H{
				"error": "Couldn't delete that post",
			})
		}

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
		err := NewCategories(collection, catID, index, name, token)

		if err != nil {
			return c.JSON(http.StatusBadRequest, H{
				"error": "Couldn't create that category with the given data",
			})
		}

		return c.JSON(http.StatusOK, H{})
	}
}

func patchCategories(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.FormValue("token")
		catID, _ := strconv.Atoi(c.FormValue("id"))
		name := c.FormValue("name")
		index, _ := strconv.Atoi(c.FormValue("index"))
		err := PatchCategories(collection, catID, name, index, token)

		if err != nil {
			return c.JSON(http.StatusBadRequest, H{
				"error": "Couldn't edit that category with the given data",
			})
		}

		return c.JSON(http.StatusOK, H{})
	}
}

func deleteCategories(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.FormValue("token")
		id, _ := strconv.Atoi(c.FormValue("id"))
		err := DeleteCategories(collection, id, token)

		if err != nil {
			return c.JSON(http.StatusBadRequest, H{
				"error": "Couldn't delete that category",
			})
		}

		return c.JSON(http.StatusOK, H{})
	}
}

func getAllSponsors(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		count, _ := strconv.Atoi(c.QueryParam("count"))
		result, err := GetAllSponsors(collection, count)

		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, H{
				"error": "Couldn't get all sponsors",
			})
		}

		return c.JSON(http.StatusOK, H{
			"sponsors": result,
		})
	}
}

func newSponsors(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.FormValue("token")
		expiryStr := c.FormValue("expiry")
		name := c.FormValue("name")
		logo := c.FormValue("logo")
		tier := c.FormValue("tier")
		link := c.FormValue("link")

		err := NewSponsors(collection, expiryStr, name, logo, tier, link, token)

		if err != nil {
			return c.JSON(http.StatusBadRequest, H{
				"error": "Couldn't create a sponsor with those details",
			})
		}

		return c.JSON(http.StatusOK, H{})
	}
}

func deleteSponsors(collection *mongo.Collection) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.FormValue("token")
		id := c.FormValue("id")
		err := DeleteSponsors(collection, id, token)
		if err != nil {
			return c.JSON(http.StatusBadRequest, H{
				"error": "Couldn't delete that sponsor",
			})
		}

		return c.JSON(http.StatusOK, H{})
	}
}
