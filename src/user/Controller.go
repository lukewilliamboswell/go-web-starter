package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Controller struct {
	repo UserRepository
}

func NewController(repo UserRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) RegisterRoutes(router *chi.Mux) {
	router.Get("/users", c.GetUsers())
}

// Returns a function that handles GET requests to /users
func (c *Controller) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Fetch users from database
		users, err := c.repo.GetUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert users to JSON
		jsonData, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write JSON to response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

// Middleware which creates the users table if it doesn't exist, and then
// ensures the `X-Ms-Client-Principal-Id` header is present
// and check if the user is in the database
func (c *Controller) UserAuthenticationMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Only run this middleware if the request is to a route on /api/*
			// if chi.RouteContext(r.Context()).RoutePattern() == "/api/*" {

			// TODO log the request id i.e.
			// `X-Request-Id: 66ed202c-45ec-4ec9-b83a-6a4a4c540796`
			// when there is a server error

			log.Println("RUNNING IN HERE")

			// Get the user ID from the request header
			userID := r.Header.Get("X-Ms-Client-Principal-Id")
			if userID == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if the user ID exists in the `users` table
			user, err := c.repo.GetUser(userID)
			if err != nil {
				http.Error(w, "error checking for user", http.StatusInternalServerError)
				return
			}
			if user == (User{}) {
				// User not found, create a new user and register in database
				newUser := User{
					PrincipalId:       r.Header.Get("X-Ms-Client-Principal-Id"),
					PrincipalName:     r.Header.Get("X-Ms-Client-Principal-Name"),
					PrincipalProvider: r.Header.Get("X-Ms-Client-Principal-Idp"),
				}

				err = c.repo.InsertUser(newUser)
				if err != nil {
					log.Print("error inserting user: %w", err)
					http.Error(w, "error inserting user", http.StatusInternalServerError)
					return
				}

				// TODO
				// Redirect user to let them know they have been registered
				// http.Redirect(w, r, "/registered", http.StatusSeeOther)
				log.Printf("User registered: %+v", user)

			} else {

				// TODO
				// User found
				log.Printf("User ID exists: %+v", user.PrincipalId)

				// Call the next handler in the chain
				next.ServeHTTP(w, r)

			}

			// }

			// // Call the next handler in the chain
			// next.ServeHTTP(w, r)

		})
	}
}
