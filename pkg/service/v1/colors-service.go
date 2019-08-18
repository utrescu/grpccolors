package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes"
	v1 "github.com/utrescu/grpccolors/pkg/api/v1"
)

const (
	apiVersion = "v1"
)

// colorServiceServer és la implementació de la interfície del servidor
type colorServiceServer struct {
	lastID       int64
	llistaColors []v1.Color
}

// NewColorServiceServer crea el servei
func NewColorServiceServer() v1.ColorServiceServer {
	var llista []v1.Color

	return &colorServiceServer{
		lastID:       0,
		llistaColors: llista,
	}
}

func (s *colorServiceServer) incrementID() int64 {
	s.lastID = s.lastID + 1
	return s.lastID
}

func (s *colorServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

func (s *colorServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	creacio := ptypes.TimestampNow()

	var id = s.incrementID()

	// insert Color entity data
	color := v1.Color{
		Id:      id,
		Nom:     req.Color.Nom,
		Rgb:     req.Color.Rgb,
		Creacio: creacio,
	}

	s.llistaColors = append(s.llistaColors, color)
	// get ID of creates ToDo

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

// Read color
func (s *colorServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {

	var color v1.Color
	found := false

	for _, element := range s.llistaColors {
		if element.Id == req.Id {
			color = element
			found = true
			break
		}
	}

	if found == false {
		return nil, errors.New("Color no trobat")
	}
	return &v1.ReadResponse{
		Api:   apiVersion,
		Color: &color,
	}, nil
}

func (s *colorServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {

	var updated int64 = 1

	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: updated,
	}, nil
}

func (s *colorServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {

	var deleted int64 = 1

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: deleted,
	}, nil
}

func (s *colorServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	list := []*v1.Color{}

	for _, color := range s.llistaColors {
		list = append(list, &color)
	}

	return &v1.ReadAllResponse{
		Api:    apiVersion,
		Colors: list,
	}, nil
}
