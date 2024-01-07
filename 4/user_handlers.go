package chiserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetUsers(dataStore *DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(2 * time.Second)
		select {
		default:
			users := make([]User, 0)
			dataStore.Users.Range(func(key, value interface{}) bool {
				users = append(users, value.(User))
				return true
			})

			respond(w, http.StatusOK, users)
		case <-r.Context().Done():
			log.Println("GetUsers: Request timeout")
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		}
	}
}

func GetUser(dataStore *DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(2 * time.Second)
		select {
		default:
			userIDStr := chi.URLParam(r, "id")

			if userIDStr == "" {
				respond(w, http.StatusBadRequest, "User ID is required")
				return
			}

			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				respond(w, http.StatusBadRequest, "Invalid User ID")
				return
			}

			if user, ok := dataStore.Users.Load(userID); ok {
				respond(w, http.StatusOK, user)
			} else {
				respond(w, http.StatusNotFound, nil)
			}
		case <-r.Context().Done():
			log.Println("GetUser: Request timeout")
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		}
	}
}

func CreateUser(dataStore *DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		default:
			var user User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				log.Println("CreateUser: Invalid JSON payload")
				respond(w, http.StatusBadRequest, "Invalid JSON payload")
				return
			}
			dataStore.Users.Store(user.ID, user)

			log.Printf("CreateUser: Creating new user - ID: %d, Name: %s\n", user.ID, user.Name)
			respond(w, http.StatusCreated, nil)
		case <-r.Context().Done():
			log.Println("CreateUser: Request timeout")
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		}
	}
}

func UpdateUser(dataStore *DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		default:
			userIDStr := chi.URLParam(r, "id")

			if userIDStr == "" {
				respond(w, http.StatusBadRequest, "User ID is required")
				return
			}

			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				respond(w, http.StatusBadRequest, "Invalid User ID")
				return
			}

			if _, ok := dataStore.Users.Load(userID); ok {
				var updatedUser User
				if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
					log.Println("UpdateUser: Invalid JSON payload")
					respond(w, http.StatusBadRequest, "Invalid JSON payload")
					return
				}
				dataStore.Users.Store(userID, updatedUser)

				log.Printf("UpdateUser: Updating user with ID %d\n", userID)
				respond(w, http.StatusOK, nil)
			} else {
				respond(w, http.StatusNotFound, nil)
			}
		case <-r.Context().Done():
			log.Println("UpdateUser: Request timeout")
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		}
	}
}

func DeleteUser(dataStore *DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		default:
			userIDStr := chi.URLParam(r, "id")

			if userIDStr == "" {
				respond(w, http.StatusBadRequest, "User ID is required")
				return
			}

			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				respond(w, http.StatusBadRequest, "Invalid User ID")
				return
			}

			if _, ok := dataStore.Users.Load(userID); ok {
				dataStore.Users.Delete(userID)

				log.Printf("DeleteUser: Deleting user with ID %d\n", userID)
				respond(w, http.StatusOK, nil)
			} else {
				respond(w, http.StatusNotFound, nil)
			}
		case <-r.Context().Done():
			log.Println("DeleteUser: Request timeout")
			http.Error(w, "Request timeout", http.StatusRequestTimeout)
		}
	}
}
