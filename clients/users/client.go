package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/ivch/dynasty/modules/users"
)

type User struct {
	ID        uint   `json:"id"`
	Role      uint   `json:"role,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type Client interface {
	UserByID(ctx context.Context, id uint) (*User, error)
	UserByPhoneAndPassword(ctx context.Context, phone, password string) (*User, error)
}

type client struct {
	svc users.Service
}

func New(svc users.Service) Client {
	return &client{
		svc: svc,
	}
}

func (c *client) UserByID(ctx context.Context, id uint) (*User, error) {
	if id == 0 {
		return nil, errors.New("empty id")
	}

	res, err := c.svc.UserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	return &User{
		ID:        res.ID,
		Role:      res.Role,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Phone:     res.Phone,
		Email:     res.Email,
	}, nil
}

func (c *client) UserByPhoneAndPassword(ctx context.Context, phone, password string) (*User, error) {
	if phone == "" {
		return nil, errors.New("empty phone")
	}

	res, err := c.svc.UserByPhoneAndPassword(ctx, phone, password)
	if err != nil {
		return nil, errors.New("user with give credentials not found")
	}
	return &User{
		ID:        res.ID,
		Role:      res.Role,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Phone:     res.Phone,
		Email:     res.Email,
	}, nil
}
