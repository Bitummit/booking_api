package authclient

import (
	"context"
	"fmt"

	"github.com/Bitummit/booking_api/internal/models"
	"github.com/Bitummit/booking_api/pkg/config"
	auth "github.com/Bitummit/booking_auth/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Client auth.AuthClient
	Cfg *config.Config
}

func New(cfg *config.Config) (*Client, error) {
	authClient := Client {
		Cfg: cfg,
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(cfg.GrpcAuthAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating grpc auth client: %w", err)
	}

	client := auth.NewAuthClient(conn)
	authClient.Client = client

	return &authClient, nil
}

func (c *Client) Registration(user models.User) (string, error) {
	request := &auth.RegistrationRequest {
		Username: user.Username,
		Email: user.Email,
		Password: user.Password,
		FirstName: user.FirstName,
		LastName: user.LastName,
	}
	res, err := c.Client.Registration(context.Background(), request)
	if err != nil {
		return "", fmt.Errorf("auth service error: %w", err)
	}
	
	return res.GetToken(), nil
}

func (c *Client) Login(user models.User) (string, error) {
	request := &auth.LoginRequest {
		Username: user.Username,
		Password: user.Password,
	}
	res, err := c.Client.Login(context.Background(), request)
	if err != nil {
		return "", fmt.Errorf("auth service error: %w", err)
	}
	
	return res.GetToken(), nil
}

func (c *Client) CheckIsADmin(token string) error {
	request := &auth.CheckTokenRequest {
		Token: token,
	}
	_, err := c.Client.IsAdmin(context.Background(), request)
	if err != nil {
		return fmt.Errorf("auth service error: %w", err)
	}
	
	return nil
}

func (c *Client) UpdateUserRole(role, username string) error {
	request := &auth.UpdateUserRoleRequest {
		Username: username,
		Role: role,
	}
	_, err := c.Client.UpdateUserRole(context.Background(), request)
	if err != nil {
		return fmt.Errorf("auth service error: %w", err)
	}
	
	return nil
}

func (c *Client) GetUser(token string) (*models.User, error) {
	req := &auth.GetUserRequest {
		Token: token,
	}
	resp, err := c.Client.GetUser(context.Background(), req)
	if err != nil {
		return nil,fmt.Errorf("auth service error: %w", err)
	}
	
	user := models.User{
		Id: resp.Id,
		Username: resp.Username,
		FirstName: resp.FirstName,
		LastName: resp.LastName,
		Email: resp.Email,
		Role: resp.Role,
	}
	return &user, nil
}