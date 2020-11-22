package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/ivch/dynasty/server/handlers/users"
	"github.com/ivch/dynasty/server/handlers/users/transport"
)

type Client struct {
	svc userService
}

func New(svc userService) *Client {
	return &Client{
		svc: svc,
	}
}

type userService interface {
	UserByPhoneAndPassword(ctx context.Context, phone, password string) (*users.User, error)
	UserByID(ctx context.Context, id uint) (*users.User, error)
}

func (c *Client) UserByID(ctx context.Context, id uint) (*transport.UserByIDResponse, error) {
	if id == 0 {
		return nil, errors.New("empty id")
	}

	res, err := c.svc.UserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	return &transport.UserByIDResponse{
		ID:        res.ID,
		Role:      res.Role,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Phone:     res.Phone,
		Email:     res.Email,
		Active:    res.Active,
	}, nil
}

func (c *Client) UserByPhoneAndPassword(ctx context.Context, phone, password string) (*transport.UserByIDResponse, error) {
	if phone == "" {
		return nil, errors.New("empty phone")
	}

	if password == "" {
		return nil, errors.New("empty phone")
	}

	res, err := c.svc.UserByPhoneAndPassword(ctx, phone, password)
	if err != nil {
		return nil, errors.New("user with given credentials not found")
	}

	return &transport.UserByIDResponse{
		ID:        res.ID,
		Apartment: res.Apartment,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Phone:     res.Phone,
		Email:     res.Email,
		Role:      res.Role,
		Building:  &res.Building,
		Entry:     &res.Entry,
		Active:    res.Active,
	}, nil
}
