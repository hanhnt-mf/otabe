package main

import (
	"context"
	"github.com/golang/mock/gomock"
	mock_v1 "otabe/test/mock_proto"
	v1 "otabe/pb"
	"testing"
	"time"
)

func TestGetRestaurantIdSpecified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOTabeManagerClient := mock_v1.NewMockOTabeManagerClient(ctrl)
	req := &v1.GetRestaurantRequest{RestaurantId: 1}
	mockOTabeManagerClient.EXPECT().GetRestaurantDetails(
		gomock.Any(),
		req,
		).Return(mock_v1.RestaurantDetails, nil)
	testGetRestaurantIdSpecified(t, mockOTabeManagerClient)
}

func testGetRestaurantIdSpecified(t *testing.T, client v1.OTabeManagerClient) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.GetRestaurantDetails(ctx, &v1.GetRestaurantRequest{RestaurantId: 1})
	if err != nil || res.Restaurant.GetName() != mock_v1.RestaurantDetails.Restaurant.GetName() {
		t.Errorf("GetRestaurantDetails: mocking failed")
	}
	t.Log("Reply: ", res.Restaurant)
}

func TestGetRestaurantIdNotSpecified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOTabeManagerClient := mock_v1.NewMockOTabeManagerClient(ctrl)
	req := &v1.GetRestaurantRequest{RestaurantId: 1}
	mockOTabeManagerClient.EXPECT().GetRestaurantDetails(
		gomock.Any(),
		req,
	).Return(mock_v1.RestaurantDetails, nil)
	testGetRestaurantIdNotSpecified(t, mockOTabeManagerClient)
}

func testGetRestaurantIdNotSpecified(t *testing.T, client v1.OTabeManagerClient) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.GetRestaurantDetails(ctx, &v1.GetRestaurantRequest{RestaurantId: 1})
	if err != nil || res.Restaurant.GetName() != mock_v1.RestaurantDetails.Restaurant.GetName() {
		t.Errorf("GetRestaurantDetails: mocking failed %v - %v", res.Restaurant.GetName(), mock_v1.RestaurantDetails.Restaurant.GetName())
	}
	t.Log("Reply: ", res.Restaurant)
}