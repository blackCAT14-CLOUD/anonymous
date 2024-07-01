package searchalgorithm

import (
    "anonymous/models"
    "fmt"
    "github.com/jmoiron/sqlx"
)

type SearchService interface {
    Search(query string) (*SearchResults, error)
}

type searchService struct {
    db *sqlx.DB
}

func NewSearchService(db *sqlx.DB) SearchService {
    return &searchService{db: db}
}

type SearchResults struct {
    Users []*models.User `json:"users"`
    Posts []*models.Post `json:"posts"`
}

func (s *searchService) Search(query string) (*SearchResults, error) {
    var results SearchResults

    // Recherche des utilisateurs
    usersQuery := `SELECT id, username, email, profile_picture FROM users WHERE username ILIKE $1`
    usersRows, err := s.db.Queryx(usersQuery, "%"+query+"%")
    if err != nil {
        return nil, fmt.Errorf("error searching users: %w", err)
    }
    defer usersRows.Close()

    for usersRows.Next() {
        var user models.User
        if err := usersRows.StructScan(&user); err != nil {
            return nil, fmt.Errorf("error scanning user: %w", err)
        }
        results.Users = append(results.Users, &user)
    }

    // Recherche des posts
    postsQuery := `SELECT id, user_id, content, created_at FROM posts WHERE content ILIKE $1`
    postsRows, err := s.db.Queryx(postsQuery, "%"+query+"%")
    if err != nil {
        return nil, fmt.Errorf("error searching posts: %w", err)
    }
    defer postsRows.Close()

    for postsRows.Next() {
        var post models.Post
        if err := postsRows.StructScan(&post); err != nil {
            return nil, fmt.Errorf("error scanning post: %w", err)
        }
        results.Posts = append(results.Posts, &post)
    }

    return &results, nil
}
