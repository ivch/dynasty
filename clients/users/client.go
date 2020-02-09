package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/ivch/dynasty/models"
	"github.com/ivch/dynasty/modules/users"
)

type Client struct {
	svc users.Service
}

func New(svc users.Service) *Client {
	return &Client{
		svc: svc,
	}
}

func (c *Client) UserByID(ctx context.Context, id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("empty id")
	}

	res, err := c.svc.UserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	return &models.User{
		ID:        res.ID,
		Role:      res.Role,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Phone:     res.Phone,
		Email:     res.Email,
	}, nil
}

func (c *Client) UserByPhoneAndPassword(ctx context.Context, phone, password string) (*models.User, error) {
	if phone == "" {
		return nil, errors.New("empty phone")
	}

	res, err := c.svc.UserByPhoneAndPassword(ctx, phone, password)
	if err != nil {
		return nil, errors.New("user with give credentials not found")
	}
	return &models.User{
		ID:        res.ID,
		Role:      res.Role,
		FirstName: res.FirstName,
		LastName:  res.LastName,
	}, nil
}
