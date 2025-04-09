package adapters

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	pb "movies/pkg/api"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// func GetAuthClient() pb.AuthClient {
// 	ctx, _ := context.WithTimeout(context.Background(), time.Second)

// 	conn, err := grpc.DialContext(ctx, "localhost:50053",
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithBlock(),
// 	)
// 	if err != nil {
// 		log.Fatalf("failed to connect AuthServer: %v", err)
// 	}

// 	client := pb.NewAuthClient(conn)

// 	return client
// }

// func GetNewConn() *grpc.ClientConn {
// 	for {
// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 		conn, err := grpc.DialContext(ctx, "localhost:50053",
// 			grpc.WithTransportCredentials(insecure.NewCredentials()),
// 			grpc.WithBlock(),
// 		)

// 		cancel()
// 		if err != nil {
// 			fmt.Printf("failed to connect ClientService: %v\n", err)

// 			time.Sleep(10 * time.Second)

// 			continue
// 		}

// 		return conn
// 	}
// }
