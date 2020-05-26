package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	pb "user/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	Serve()
}

// User ..
type User interface {
	Louder(ctx context.Context, r *pb.UserRequest) (*pb.UserResponse, error)
}

// User ..
type user struct {
}

// NreUserServer ..
func NreUserServer() User {
	return &user{}
}

func (u user) Louder(ctx context.Context, r *pb.UserRequest) (*pb.UserResponse, error) {
	fmt.Println(r.Name, r.Password, r.Code, r.Codeid)
	return &pb.UserResponse{
		JWT: "test",
	}, nil
}

// GrpcHandlerFunc ..
func GrpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	if otherHandler == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			grpcServer.ServeHTTP(w, r)
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("--------1-")
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			fmt.Println("--------2-")
			grpcServer.ServeHTTP(w, r)
		} else {
			fmt.Println("--------3-")
			otherHandler.ServeHTTP(w, r)
		}
	})
}

// GetTLSConfig ..
func GetTLSConfig(certPemPath, certKeyPath string) *tls.Config {
	var certKeyPair *tls.Certificate
	cert, _ := ioutil.ReadFile(certPemPath)
	key, _ := ioutil.ReadFile(certKeyPath)

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Printf("TLS KeyPair err: %v\n", err)
	}

	certKeyPair = &pair

	return &tls.Config{
		Certificates: []tls.Certificate{*certKeyPair},
		NextProtos:   []string{http2.NextProtoTLS},
	}
}

// CertPemPath .
var CertPemPath = "./certs/server.pem"

// CertKeyPath .
var CertKeyPath = "./certs/server.key"

// CertName .
var CertName = "localhost"

// ServerPort .
var ServerPort = "8890"

// EndPoint .
var EndPoint string

// Serve ..
func Serve() (err error) {
	EndPoint = ":" + ServerPort
	conn, err := net.Listen("tcp", EndPoint)
	if err != nil {
		log.Printf("TCP Listen err:%v\n", err)
	}

	tlsConfig := GetTLSConfig(CertPemPath, CertKeyPath)
	srv := createInternalServer(conn, tlsConfig)

	log.Printf("gRPC and https listen on: %s\n", ServerPort)

	if err = srv.Serve(tls.NewListener(conn, tlsConfig)); err != nil {
		log.Printf("ListenAndServe: %v\n", err)
	}

	return err
}

func createInternalServer(conn net.Listener, tlsConfig *tls.Config) *http.Server {
	var opts []grpc.ServerOption

	// grpc server
	creds, err := credentials.NewServerTLSFromFile(CertPemPath, CertKeyPath)
	if err != nil {
		log.Panicf("Failed to create server TLS credentials %v", err)
	}

	opts = append(opts, grpc.Creds(creds))
	grpcServer := grpc.NewServer(opts...)

	// register grpc pb
	pb.RegisterUserServer(grpcServer, NreUserServer())

	// gw server
	ctx := context.Background()
	dcreds, err := credentials.NewClientTLSFromFile(CertPemPath, CertName)
	if err != nil {
		log.Panicf("Failed to create client TLS credentials %v", err)
	}
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	gwmux := runtime.NewServeMux()
	// register grpc-gateway pb
	if err := pb.RegisterUserHandlerFromEndpoint(ctx, gwmux, EndPoint, dopts); err != nil {
		log.Printf("Failed to register gw server: %v\n", err)
	}

	// https服务
	mux := http.NewServeMux()
	serveSwagger(mux)
	mux.Handle("/", gwmux)

	return &http.Server{
		Addr:      EndPoint,
		Handler:   GrpcHandlerFunc(grpcServer, mux),
		TLSConfig: tlsConfig,
	}
}

//swagger UI
func serveSwagger(mux *http.ServeMux) {
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, strings.NewReader(pb.AiasJSON))
	})

	// mime.AddExtensionType(".svg", "image/svg+xml")

	// Expose files in third_party/swagger-ui/ on <host>/swagger-ui
	// fileServer := http.FileServer(&assetfs.AssetFS{
	// 	Asset:    swagger.Asset,
	// 	AssetDir: swagger.AssetDir,
	// 	Prefix:   "third_party/swagger-ui",
	// })
	prefix := "/swagger-ui/"
	// mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
	mux.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir("./proto/third_party/swagger-ui"))))
}
