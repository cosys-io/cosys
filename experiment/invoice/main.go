package main

import (
	"context"
	"fmt"
	proto "github.com/cosys-io/cosys/experiment/invoice/invoicer"
	"google.golang.org/grpc"
	"log"
	"net"
)

type myInvoiceServer struct {
	proto.UnimplementedInvoicerServer
}

func (receiver myInvoiceServer) Create(ctx context.Context, req *proto.CreateRequest) (*proto.CreateResponse, error) {
	return &proto.CreateResponse{
		Pdf:  []byte(req.From),
		Docx: []byte("Hello"),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	serverRegistrar := grpc.NewServer()
	service := &myInvoiceServer{}

	proto.RegisterInvoicerServer(serverRegistrar, service)
	fmt.Printf("serving on %s\n", lis.Addr().String())
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}
