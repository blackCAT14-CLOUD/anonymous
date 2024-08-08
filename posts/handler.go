package posts

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"strconv"
)

func CreatePostHandler(service PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var payload PostPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		post, err := service.CreatePost(token, &payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse(w, http.StatusCreated, post)
	}
}

func GetAllPostsHandler(service PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Récupérer les paramètres de pagination depuis la requête
		offsetStr := r.URL.Query().Get("offset")
		limitStr := r.URL.Query().Get("limit")

		// Convertir les paramètres en entiers
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0 // Valeur par défaut
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10 // Valeur par défaut
		}

		// Récupérer les posts avec pagination
		posts, err := service.GetAllPosts(offset, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse(w, http.StatusOK, posts)
	}
}


func GetPostsByUserHandler(service PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		posts, err := service.GetPostsByUser(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, http.StatusOK, posts)
	}
}

func UpdatePostHandler(service PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		postID := chi.URLParam(r, "postID")

		var payload PostPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		updatedPost, err := service.UpdatePost(token, postID, &payload)
		if err != nil {
			if err.Error() == "unauthorized" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse(w, http.StatusOK, updatedPost)
	}
}

func DeletePostHandler(service PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		postID := chi.URLParam(r, "postID")

		err := service.DeletePost(token, postID)
		if err != nil {
			if err.Error() == "unauthorized" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if err.Error() == "post not found" {
				http.Error(w, "Post not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}


func LikePostHandler(service PostService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        postID := chi.URLParam(r, "postID")
        if err := service.LikePost(token, postID); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func UnlikePostHandler(service PostService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        postID := chi.URLParam(r, "postID")
        if err := service.UnlikePost(token, postID); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func AddReactionHandler(service PostService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        postID := chi.URLParam(r, "postID")
        var req struct {
            ReactionType string `json:"reaction_type"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        if err := service.AddReaction(token, postID, req.ReactionType); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}

func RemoveReactionHandler(service PostService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        postID := chi.URLParam(r, "postID")
        if err := service.RemoveReaction(token, postID); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }
}


func GetLikesCountHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        postID := chi.URLParam(r, "postID")
        query := `SELECT COUNT(*) FROM post_likes WHERE post_id = $1`
        var likesCount int
        err := db.QueryRow(query, postID).Scan(&likesCount)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        response := struct {
            PostID      string `json:"post_id"`
            LikesCount  int    `json:"likes_count"`
        }{
            PostID:     postID,
            LikesCount: likesCount,
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}


func GetReactionsCountHandler(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        postID := chi.URLParam(r, "postID")
        query := `SELECT reaction_type, COUNT(*) FROM post_reactions WHERE post_id = $1 GROUP BY reaction_type`
        rows, err := db.Query(query, postID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        reactionsCount := make(map[string]int)
        for rows.Next() {
            var reactionType string
            var count int
            if err := rows.Scan(&reactionType, &count); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            reactionsCount[reactionType] = count
        }
        response := struct {
            PostID         string            `json:"post_id"`
            ReactionsCount map[string]int    `json:"reactions_count"`
        }{
            PostID:         postID,
            ReactionsCount: reactionsCount,
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}