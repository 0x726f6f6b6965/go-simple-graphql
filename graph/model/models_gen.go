// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Blog struct {
	ID        string     `json:"id" bson:"_id,omitempty"`
	Title     string     `json:"title" bson:"title"`
	Content   string     `json:"content" bson:"content"`
	Author    *User      `json:"author,omitempty" bson:"author"`
	CreatedAt time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:"updatedAt"`
}

type DeleteBlog struct {
	BlogID string `json:"blogId" bson:"blogId"`
}

type EditBlog struct {
	BlogID  string `json:"blogId" bson:"blogId"`
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
}

type LoginInput struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type NewBlog struct {
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
}

type NewUser struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type User struct {
	ID        string     `json:"id" bson:"_id,omitempty"`
	Username  string     `json:"username" bson:"username"`
	Email     string     `json:"email" bson:"email"`
	Password  string     `json:"password" bson:"password"`
	CreatedAt time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" bson:"updatedAt"`
}
