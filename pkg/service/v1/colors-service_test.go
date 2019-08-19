package v1

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	v1 "github.com/utrescu/grpccolors/pkg/api/v1"
)

func Test_toDoServiceServer_Create(t *testing.T) {
	ctx := context.Background()

	s := NewColorServiceServer()

	type args struct {
		ctx context.Context
		req *v1.CreateRequest
	}
	tests := []struct {
		name    string
		s       v1.ColorServiceServer
		args    args
		want    *v1.CreateResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1",
					Color: &v1.Color{
						Nom: "blanc",
						Rgb: "000000",
					},
				},
			},
			want: &v1.CreateResponse{
				Api: "v1",
				Id:  1,
			},
		},
		{
			name: "Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1000",
					Color: &v1.Color{
						Nom: "blanc",
						Rgb: "000000",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("colorServiceServer.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("colorServiceServer.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_colorServiceServer_Read(t *testing.T) {
	ctx := context.Background()

	tm := time.Now().In(time.UTC)
	creation, _ := ptypes.TimestampProto(tm)

	colors := []v1.Color{
		v1.Color{
			Id:      1,
			Nom:     "blanc",
			Rgb:     "000000",
			Creacio: creation,
		},
	}

	s := NewColorServiceServerWithValues(colors)

	type args struct {
		ctx context.Context
		req *v1.ReadRequest
	}
	tests := []struct {
		name    string
		s       v1.ColorServiceServer
		args    args
		want    *v1.ReadResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v1",
					Id:  colors[0].Id,
				},
			},

			want: &v1.ReadResponse{
				Api:   "v1",
				Color: &colors[0],
			},
		},
		{
			name: "Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v1000",
					Id:  colors[0].Id,
				},
			},
			wantErr: true,
		},
		{
			name: "Not found",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadRequest{
					Api: "v1",
					Id:  colors[0].Id + 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Read(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("toDoServiceServer.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toDoServiceServer.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toDoServiceServer_ReadAll(t *testing.T) {
	ctx := context.Background()

	tm := time.Now().In(time.UTC)
	creation, _ := ptypes.TimestampProto(tm)

	colors := []v1.Color{
		v1.Color{
			Id:      1,
			Nom:     "blanc",
			Rgb:     "000000",
			Creacio: creation,
		},
		v1.Color{
			Id:      2,
			Nom:     "negre",
			Rgb:     "FFFFFF",
			Creacio: creation,
		},
	}

	s := NewColorServiceServerWithValues(colors)

	type args struct {
		ctx context.Context
		req *v1.ReadAllRequest
	}
	tests := []struct {
		name    string
		s       v1.ColorServiceServer
		args    args
		mock    func()
		want    *v1.ReadAllResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadAllRequest{
					Api: "v1",
				},
			},
			want: &v1.ReadAllResponse{
				Api: "v1",
				Colors: []*v1.Color{
					{
						Id:      1,
						Nom:     "blanc",
						Rgb:     "000000",
						Creacio: creation,
					},
					{
						Id:      2,
						Nom:     "negre",
						Rgb:     "FFFFFF",
						Creacio: creation,
					},
				},
			},
		},
		{
			name: "Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadAllRequest{
					Api: "v100",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ReadAll(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("toDoServiceServer.ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toDoServiceServer.ReadAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toDoServiceServer_ReadAll_Empty(t *testing.T) {
	ctx := context.Background()

	s := NewColorServiceServer()

	type args struct {
		ctx context.Context
		req *v1.ReadAllRequest
	}
	tests := []struct {
		name    string
		s       v1.ColorServiceServer
		args    args
		mock    func()
		want    *v1.ReadAllResponse
		wantErr bool
	}{
		{
			name: "Empty",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadAllRequest{
					Api: "v1",
				},
			},
			want: &v1.ReadAllResponse{
				Api:    "v1",
				Colors: []*v1.Color{},
			},
		},
		{
			name: "Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.ReadAllRequest{
					Api: "v100",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ReadAll(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("toDoServiceServer.ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toDoServiceServer.ReadAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
