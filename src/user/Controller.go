package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type contextKey string

const userKey contextKey = "user"

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
func (c *Controller) UserAuthenticationMiddleware(version string, exceptions []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Check if the current request path matches any of the exceptions
			for _, path := range exceptions {
				if r.URL.Path == path {
					// Skip the middleware and call the next handler directly
					next.ServeHTTP(w, r)
					return
				}
			}

			// Get the user ID from the request header
			userID := r.Header.Get("X-Ms-Client-Principal-Id")
			if userID == "" {
				http.Error(w, "Expected X-Ms-Client-Principal-Id header to be set", http.StatusUnauthorized)
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

				http.Redirect(w, r, "/registered.html", http.StatusSeeOther)

			} else if user.AccessLevel == DENY_ACCESS {

				// Redirect user to let them know they have been registered
				if r.URL.Path == "/api/*" {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
				} else {
					http.Redirect(w, r, "/registered.html", http.StatusSeeOther)
				}

			} else {

				// Add the User object to the request context
				ctx := context.WithValue(r.Context(), userKey, user)

				// Create a new request with the updated context
				r = r.WithContext(ctx)

				// Call the next handler in the chain
				next.ServeHTTP(w, r)

			}
		})
	}
}
