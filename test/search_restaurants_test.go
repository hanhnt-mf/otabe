package main

import (
	"context"
	"github.com/golang/mock/gomock"
	mock_v1 "otabe/test/mock_proto"
	v1 "otabe/pb"
	"testing"
	"time"
)

var (
	fullListRestaurants = mock_v1.ListRestaurants
)

func TestSearchRestaurantsNotSpecified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOTabeManagerClient := mock_v1.NewMockOTabeManagerClient(ctrl)
	mockOTabeManagerClient.EXPECT().ListRestaurantsByOptions(
		gomock.Any(),
		&v1.ListRestaurantsRequest{},
		).Return(fullListRestaurants, nil)
	testSearchRestaurantsNotSpecified(t, mockOTabeManagerClient)
}

func testSearchRestaurantsNotSpecified(t *testing.T, client v1.OTabeManagerClient) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.ListRestaurantsByOptions(ctx, &v1.ListRestaurantsRequest{})
	if err != nil {
		t.Errorf("ListRestaurantsByOptions: mocking failed - %v", err)
	}
	if len(res.Data) != len(fullListRestaurants.Data) {
		t.Errorf("ListRestaurantsByOptions: data length not match, expected: %v - received : %v",fullListRestaurants, res)
	}
	if res.Data[0].Restaurant.GetName() != fullListRestaurants.Data[0].Restaurant.GetName() {
		t.Errorf("ListRestaurantsByOptions: data restaurant not match,")
	}
}
