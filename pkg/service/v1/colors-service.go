package v1

import (
	"context"

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

// NewColorServiceServerWithValues inicialitzat amb valors
func NewColorServiceServerWithValues(colors []v1.Color) v1.ColorServiceServer {
	var max int64
	for _, c := range colors {
		if c.Id > max {
			max = c.Id
		}
	}

	return &colorServiceServer{
		lastID:       max,
		llistaColors: colors,
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

func (s *colorServiceServer) checkRGB(rgb string) error {

	for _, element := range s.llistaColors {
		if element.Rgb == rgb {
			return status.Errorf(codes.AlreadyExists,
				"el color RGB `%s` ja està entre les dades:'%s'", rgb, element.Nom)
		}
	}
	return nil
}

func (s *colorServiceServer) locateColor(id int64) (int, error) {

	for index, element := range s.llistaColors {
		if element.Id == id {
			return index, nil
		}
	}
	return -1, status.Errorf(codes.NotFound,
		"Id no trobat:'%d'", id)
}

func (s *colorServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// Check if rgb already exists
	if err := s.checkRGB(req.Color.Rgb); err != nil {
		return nil, err
	}

	if req.Color == nil || req.Color.Nom == "" || req.Color.Rgb == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			"S'han d'especificar les noves dades: color:{nom=%s,rgb=%s}", req.Color.Nom, req.Color.Rgb)
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

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

// Read color
func (s *colorServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	posicio, err := s.locateColor(req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.ReadResponse{
		Api:   apiVersion,
		Color: &s.llistaColors[posicio],
	}, nil
}

func (s *colorServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	if req.Color == nil || req.Color.Nom == "" || req.Color.Rgb == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			"S'han d'especificar les noves dades: color:{nom=%s,rgb=%s}", req.Color.Nom, req.Color.Rgb)
	}

	var updated int64

	posicio, err := s.locateColor(req.Color.Id)
	if err == nil {
		s.llistaColors[posicio].Nom = req.Color.Nom
		s.llistaColors[posicio].Rgb = req.Color.Rgb
		updated = 1
	}

	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: updated,
	}, nil
}

func (s *colorServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	var deleted int64

	posicio, err := s.locateColor(req.Id)
	if err == nil {
		s.llistaColors = append(s.llistaColors[:posicio], s.llistaColors[posicio+1:]...)
		deleted = 1
	}

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: deleted,
	}, nil
}

func (s *colorServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	llista := []*v1.Color{}

	for index := range s.llistaColors {
		llista = append(llista, &s.llistaColors[index])
	}

	return &v1.ReadAllResponse{
		Api:    apiVersion,
		Colors: llista,
	}, nil
}
