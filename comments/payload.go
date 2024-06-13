package comments

import (
    "github.com/google/uuid"
)

type CommentPayload struct {
    PostID      string `json:"post_id"`
    UserID      string `json:"user_id"`
    ContentType string `json:"content_type"`
    Content     string `json:"content"`
}

func (p *CommentPayload) Validate() map[string]string {
    errors := make(map[string]string)

    // Vérifier la présence et la validité de PostID
    if p.PostID == "" {
        errors["post_id"] = "Post ID is required"
    } else if !isValidUUID(p.PostID) {
        errors["post_id"] = "Post ID must be a valid UUID"
    }

    // Vérifier la présence et la validité de UserID
    if p.UserID == "" {
        errors["user_id"] = "User ID is required"
    } else if !isValidUUID(p.UserID) {
        errors["user_id"] = "User ID must be a valid UUID"
    }

    // Vérifier la présence de ContentType
    if p.ContentType == "" {
        errors["content_type"] = "Content type is required"
    }

    // Vérifier la présence de Content
    if p.Content == "" {
        errors["content"] = "Content is required"
    }

    return errors
}

func isValidUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
}
