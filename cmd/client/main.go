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

type colorsACrear struct {
	Nom string
	Rgb string
}

func main() {

	// Colors que crearà en el servidor
	crearColors := []colorsACrear{
		colorsACrear{
			Nom: "vermell",
			Rgb: "FF0000",
		},
		colorsACrear{
			Nom: "blau",
			Rgb: "00FF00",
		},
		colorsACrear{
			Nom: "verd",
			Rgb: "0000FF",
		},
		colorsACrear{
			Nom: "blanc",
			Rgb: "FFFFFF",
		},
		colorsACrear{
			Nom: "negre",
			Rgb: "000000",
		},
	}

	// get configuration
	address := flag.String("servidor", "", "servidor gRPC = host:port")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No s'ha pogut connectar: %v", err)
	}
	defer conn.Close()

	// Crear el client
	c := v1.NewColorServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Creació dels colors
	for _, dadesColor := range crearColors {
		req1 := v1.CreateRequest{
			Api: apiVersion,
			Color: &v1.Color{
				Nom: dadesColor.Nom,
				Rgb: dadesColor.Rgb,
			},
		}
		res1, err := c.Create(ctx, &req1)
		if err != nil {
			log.Fatalf("Ha fallat la creació: %v", err)
		}
		log.Printf("Create: Id=%d\n\n", res1.Id)
	}

	// Recuperar el que té Id 1
	var id int64 = 2
	reqblau := v1.ReadRequest{
		Api: apiVersion,
		Id:  id,
	}
	res2, err := c.Read(ctx, &reqblau)
	if err != nil {
		log.Fatalf("Ha fallat la lectura: %v", err)
	}
	log.Printf("Read: %d, %s, %s\n\n", res2.Color.Id, res2.Color.Nom, res2.Color.Rgb)

	// Veure'ls tots
	req3 := v1.ReadAllRequest{
		Api: apiVersion,
	}
	res3, err := c.ReadAll(ctx, &req3)
	if err != nil {
		log.Fatalf("Ha fallat la recuperació de tots: %v", err)
	}

	log.Print("----- COLORS ----\n")
	for _, color := range res3.Colors {
		log.Printf("Color: %d %s : %s", color.Id, color.Nom, color.Rgb)
	}

}
