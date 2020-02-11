package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/ivch/dynasty/models/dto"
	"github.com/ivch/dynasty/models/entities"
	"github.com/ivch/dynasty/modules/users"
)

type Client struct {
	svc userService
}

func New(svc users.Service) *Client {
	return &Client{
		svc: svc,
	}
}

type userService interface {
	Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	UserByPhoneAndPassword(ctx context.Context, phone, password string) (*dto.UserAuthResponse, error)
	UserByID(ctx context.Context, id uint) (*dto.UserByIDResponse, error)
}

func (c *Client) UserByID(ctx context.Context, id uint) (*entities.User, error) {
	if id == 0 {
		return nil, errors.New("empty id")
	}

	res, err := c.svc.UserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	return &entities.User{
		ID:        res.ID,
		Role:      res.Role,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Phone:     res.Phone,
		Email:     res.Email,
	}, nil
}

func (c *Client) UserByPhoneAndPassword(ctx context.Context, phone, password string) (*entities.User, error) {
	if phone == "" {
		return nil, errors.New("empty phone")
	}

	res, err := c.svc.UserByPhoneAndPassword(ctx, phone, password)
	if err != nil {
		return nil, errors.New("user with give credentials not found")
	}
	return &entities.User{
		ID:        res.ID,
		Role:      res.Role,
		FirstName: res.FirstName,
		LastName:  res.LastName,
	}, nil
}
