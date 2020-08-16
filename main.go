package main

// Import our dependencies. We'll use the standard HTTP library as well as the gorilla router for this app
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/auth0-community/go-auth0"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	jose "gopkg.in/square/go-jose.v2"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("failed to load .env; here's why: ", err)
	}

	// Here we are instantiating the gorilla/mux router
	r := mux.NewRouter()

	// On the default page we will simply serve our static index page.
	r.Handle("/", http.FileServer(http.Dir("./views/")))
	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Handle("/status", StatusHandler).Methods("GET")
	/* We will add the middleware to our products and feedback routes. The status route will be publicly accessible */
	r.Handle("/products", authMiddleware(ProductsHandler)).Methods("GET")
	r.Handle("/products/{slug}/feedback", authMiddleware(AddFeedbackHandler)).Methods("POST")

	// Our application will run on port 3000. Here we declare the port and pass in our router.
	if err := http.ListenAndServe(":3020", handlers.LoggingHandler(os.Stdout, r)); err != nil {
		log.Println("error upon serving:", err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretProvider := auth0.NewKeyProvider([]byte(os.Getenv("AUTH0SECRET")))
		configuration := auth0.NewConfiguration(
			secretProvider,
			[]string{os.Getenv("AUTH0AUDIENCE")},
			os.Getenv("AUTH0DOMAIN"),
			jose.HS256,
		)
		validator := auth0.NewValidator(configuration, nil)

		token, err := validator.ValidateRequest(r)

		if err == nil {
			next.ServeHTTP(w, r)
			return
		}
		fmt.Println(err)
		fmt.Println("Token is not valid:", token)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized"))
	})
}

// Product will be the first type to be created
//   This type will contain information about VR experiences
type Product struct {
	ID          int
	Name        string
	Slug        string
	Description string
}

/* We will create our catalog of VR experiences and store them in a slice. */
var products = []Product{
	{ID: 1, Name: "Hover Shooters", Slug: "hover-shooters", Description: "Shoot your way to the top on 14 different hoverboards"},
	{ID: 2, Name: "Ocean Explorer", Slug: "ocean-explorer", Description: "Explore the depths of the sea in this one of a kind underwater experience"},
	{ID: 3, Name: "Dinosaur Park", Slug: "dinosaur-park", Description: "Go back 65 million years in the past and ride a T-Rex"},
	{ID: 4, Name: "Cars VR", Slug: "cars-vr", Description: "Get behind the wheel of the fastest cars in the world."},
	{ID: 5, Name: "Robin Hood", Slug: "robin-hood", Description: "Pick up the bow and arrow and master the art of archery"},
	{ID: 6, Name: "Real World VR", Slug: "real-world-vr", Description: "Explore the seven wonders of the world in VR"},
}

// StatusHandler will be invoked when the user calls the /status route
//   It will simply return a string with the message "API is up and running"
var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("API is up and running"))
})

// ProductsHandler will be called when the user makes a GET request to the /products endpoint.
//    This handler will return a list of products available for users to review
var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// Here we are converting the slice of products to JSON
	payload, _ := json.Marshal(products)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
	if _, err := w.Write([]byte(payload)); err != nil {
		log.Println("error writing payload:", err)
	}
})

// AddFeedbackHandler will add either positive or negative feedback to the product
//   We would normally save this data to the database - but for this demo, we'll fake it
//   so that as long as the request is successful and we can match a product to our catalog of products
//   we'll return an OK status.
var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var product Product
	vars := mux.Vars(r)
	slug := vars["slug"]

	for _, p := range products {
		if p.Slug == slug {
			product = p
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if product.Slug != "" {
		payload, _ := json.Marshal(product)
		if _, err := w.Write([]byte(payload)); err != nil {
			log.Println("error writing payload:", err)
		}
	} else {
		_, _ = w.Write([]byte("Product Not Found"))
	}
})
