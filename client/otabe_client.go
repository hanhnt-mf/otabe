package main


import (
	"context"
	"flag"
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

	res, err := client.GetRestaurantDetails(ctx, &pb.GetRestaurantRequest{RestaurantId: 1})
	if err != nil {
	 log.Fatalf("%v.GetRestaurantDetails(_) = _, %v: ", client, err)
	}
	log.Printf(`Restaurant Details: %v`, res)
}
