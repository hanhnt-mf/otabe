package main


import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
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

func main() {
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	// WithBlock: ensure Dial() will not return value until the connection is made
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}
	defer conn.Close()
	client := pb.NewOTabeManagerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//getRestaurantDetails(client, ctx, &pb.GetRestaurantRequest{RestaurantId: 1})
	//restaurantName := "HaNoi & Hanoi"
	//nation := "Vietnamese"
	//itemName := "Banh mi nhan thit"
	//prefecture := "1-3-28 Shaiba Tokyo"
	//lat := 35.64489421538165
	//long := 139.74929356703967
	//distance := float64(100000)
	//location:= &pb.SearchLocationConditions{Long: &long, Lat: &lat, Distance: &distance}
	//paging := &pb.Paging{PageLimit: 10, PageNumber: 2}
	//sortedBy := "created_at"

	listRestaurantsByOptions(client, ctx, &pb.ListRestaurantsRequest{
		//RestaurantName: &restaurantName,
		//Nation: &nation,
		//ItemName: &itemName,
		//Prefecture: &prefecture,
		//Location: location,
		//Paging: paging,
		//SortedBy: &sortedBy,
	})
}
