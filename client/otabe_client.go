package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"otabe/service"
	pb "otabe/v1"
	"time"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:8080", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)



func getRestaurantDetails(client pb.OTabeManagerClient, ctx context.Context, request *pb.GetRestaurantRequest) {
	res, err := client.GetRestaurantDetails(ctx, request)
	if err != nil {
		log.Fatalf("%v.GetRestaurantDetails(_) = _, %v: ", client, err)
	}
	resJson, _ := json.Marshal(res)
	fmt.Printf(`%s`,resJson)
}

func listRestaurantsByOptions(client pb.OTabeManagerClient, ctx context.Context, req *pb.ListRestaurantsRequest) {
	res, err := client.ListRestaurantsByOptions(ctx, req)
	if err != nil {
		log.Fatalf("%v.ListReqtaurantsByOptions(_) = _, %v: ", client, err)
	}
	resJson, _ := json.Marshal(res)
	fmt.Printf(`%s`,resJson)
}

func createNewRestaurant(client pb.OTabeManagerClient, ctx context.Context, req *pb.CreateRestaurantRequest) {
	res, err := client.CreateNewRestaurant(ctx, req)
	if err != nil {
		log.Fatalf("%v.CreateNewRestaurant(_) = _, %v: ", client, err)
	}
	resJson, _ := json.Marshal(res)
	fmt.Printf(`%s`,resJson)
}

//func updateRestaurant(client pb.OTabeManagerClient, ctx context.Context, req *pb.UpdateRestaurantRequest) {
//
//}
//


const (
	username = "HaNa"
	password = "Hana123"
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	const restaurantsPath = "/v1.OTabeManager/"
	return map[string]bool{
		restaurantsPath + "CreateNewRestaurant" : true,
		restaurantsPath + "UpdateNewRestaurant" : true,
		//restaurantsPath + "GetRestaurantDetails" : {"admin", "user"},
		//restaurantsPath + "ListRestaurantsByOptions" : {"admin", "user"},
	}
}

func main() {
	conn1, err := grpc.Dial(
		*serverAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	// WithBlock: ensure Dial() will not return value until the connection is made
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}
	authClient := service.NewAuthClient(conn1, username, password)
	interceptor, err := service.NewAuthInterceptorClient(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	conn, err := grpc.Dial(
		*serverAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
	)
	// WithBlock: ensure Dial() will not return value until the connection is made
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}
	defer conn.Close()

	client := pb.NewOTabeManagerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	restaurantName := flag.String("restaurantName","Ai banh my" , "name for restaurant")
	//nation := flag.String("restaurantNation", "", "user's nation voted items for restaurant")
	//itemName := flag.String("itemName", "", "name of item in restaurant")
	prefecture := flag.String("prefecture", "Tokyo", "name of item in restaurant")
	long := flag.Float64("long", 53.34542, "longitude of search area")
	lat := flag.Float64("lat", 435.2352345, "latitude of search area")
	//distance := flag.Float64("distance", 0, "radius of search area")
	//pageLimit := flag.Uint64("pageLimit", 10, "limit of results in 1 page")
	//pageNumber := flag.Uint64("pageNumber", 1, "page number")
	//location := &pb.SearchLocationConditions{Long: long, Lat: lat, Distance: distance}
	//paging := &pb.Paging{PageLimit: *pageLimit, PageNumber: *pageNumber}
	//sortedBy := flag.String("sortedBy", "created_at", "sorted results by column")
	//
	//if *long == 0 || *lat == 0 || *distance == 0 {
	//	location = nil
	//}
	//flag.Parse()
	//listRestaurantsByOptions(client, ctx, &pb.ListRestaurantsRequest{
	//	RestaurantName: restaurantName,
	//	Nation: nation,
	//	ItemName: itemName,
	//	Prefecture: prefecture,
	//	Location: location,
	//	Paging: paging,
	//	SortedBy: sortedBy,
	//})

	//getRestaurantDetails(client, ctx, &pb.GetRestaurantRequest{})

	restaurant := &pb.RestaurantRequest{
		Name: *restaurantName,
		Website: "aibanhmy.com",
		Phone: "023345878",
		Description: "ngon banh my",
		PostalCode: "1080023",
		Address: *prefecture,
		Geo: &pb.Geo{Long: *long, Lat: *lat},
	}

	menuItems := make([]*pb.MenuItemsRequest, 0)
	menuItems = append(menuItems, &pb.MenuItemsRequest{ItemName: "Banh my tuyet", Description: "ngon lam", Price: 2341})
	menus := make([]*pb.MenuRequest, 0)
	menus = append(menus, &pb.MenuRequest{Name: "First", MenuItems: menuItems})
	createNewRestaurant(client, ctx, &pb.CreateRestaurantRequest{Restaurant: restaurant, Menus: menus})

	//res, err := client.Login(ctx, &pb.LoginRequest{Username: "HaNa", Password: "Hana123"})
	//if err != nil {
	//	log.Fatalf("%v.Login(_) = _, %v: ", client, err)
	//}
	//resJson, _ := json.Marshal(res)
	//fmt.Printf(`%s`,resJson)
}