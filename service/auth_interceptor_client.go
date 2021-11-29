package service

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	v1 "otabe/v1"
	"time"
)

type OTabeClient struct {
	service v1.OTabeManagerClient
	username string
	password string
}

func NewAuthClient(cc *grpc.ClientConn, username string, password string) *OTabeClient {
	service :=v1.NewOTabeManagerClient(cc)
	return &OTabeClient{service, username, password}
}

func (client *OTabeClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &v1.LoginRequest{
		Username: client.username,
		Password: client.password,
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}


type AuthInterceptorClient struct {
	authClient *OTabeClient
	authMethods map[string]bool
	accessToken string
}

func NewAuthInterceptorClient (
	authClient *OTabeClient,
	authMethods map[string]bool,
	refreshDuration time.Duration,
	) (*AuthInterceptorClient, error) {
	interceptor := &AuthInterceptorClient{
		authClient: authClient,
		authMethods: authMethods,
	}

	err := interceptor.scheduleRefreshToken(refreshDuration)
	if err != nil {
		return nil, err
	}

	return interceptor, nil
}

func (interceptor *AuthInterceptorClient) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req,
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
		) error {
		log.Printf("--> unary interceptor: %s", method)
		if interceptor.authMethods[method] {
			return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (interceptor *AuthInterceptorClient) attachToken(ctx context.Context)  context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}

func (interceptor *AuthInterceptorClient) scheduleRefreshToken(refreshDuration time.Duration) error {
	err := interceptor.refreshToken()
	if err != nil {
		return err
	}
	go func() {
		wait := refreshDuration
		for {
			time.Sleep(wait)
			err := interceptor.refreshToken()
			if err != nil {
				wait = time.Second
			} else {
				wait = refreshDuration
			}
		}
	}()

	return nil
}

func (interceptor *AuthInterceptorClient) refreshToken() error {
	accessToken, err := interceptor.authClient.Login()
	if err != nil {
		return err
	}

	interceptor.accessToken = accessToken
	log.Printf("token refreshed: %v", accessToken)

	return nil
}