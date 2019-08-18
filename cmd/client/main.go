package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	v1 "github.com/utrescu/grpccolors/pkg/api/v1"
)

const (
	apiVersion = "v1"
)

func main() {
	// get configuration
	address := flag.String("servidor", "", "servidor gRPC = host:port")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No s'ha pogut connectar: %v", err)
	}
	defer conn.Close()

	c := v1.NewColorServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Prova la creació d'un color amb "Create"
	req1 := v1.CreateRequest{
		Api: apiVersion,
		Color: &v1.Color{
			Nom: "Vermell",
			Rgb: "FF0000",
		},
	}
	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res1)

	id := res1.Id

	// Read
	req2 := v1.ReadRequest{
		Api: apiVersion,
		Id:  id,
	}
	res2, err := c.Read(ctx, &req2)
	if err != nil {
		log.Fatalf("Ha fallat la lectura: %v", err)
	}
	log.Printf("Resultat: <%+v>\n\n", res2)

	// Call ReadAll
	req4 := v1.ReadAllRequest{
		Api: apiVersion,
	}
	res4, err := c.ReadAll(ctx, &req4)
	if err != nil {
		log.Fatalf("Ha fallat la recuperació de tots: %v", err)
	}
	log.Printf("Tots els colors: <%+v>\n\n", res4)

}
