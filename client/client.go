/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	pb "example.com/nampkh-grpc/helloworld"

	"google.golang.org/grpc"
)

var (
	host        = "localhost"
	port        = "50051"
	defaultName = "world"
)

func greetMePlease(ctx context.Context, client pb.GreeterClient, request *pb.HelloRequest) error {
	stream, err := client.LotsOfReplies(ctx, request)

	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return err
	}

	for {
		greeting, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
			return err
		}

		log.Println(greeting.GetMessage())
	}
	return nil
}

var connectionCount = 1
var loopCount = 10

func connect(ws *sync.WaitGroup) {
	defer ws.Done()
	// Set up a connection to the server.
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 2 {
		name = os.Args[2]
	}

	ctx := context.Background()

	for i := 0; i < loopCount; i++ {
		err = greetMePlease(ctx, c, &pb.HelloRequest{Name: name + "-" + strconv.Itoa(i)})

		if err != nil {
			fmt.Println("in connect" + err.Error())
			break
		}
	}
}

func main() {
	if len(os.Args) > 1 {
		host = os.Args[1]
	}

	if len(os.Args) > 3 {
		var err error
		connectionCount, err = strconv.Atoi(os.Args[3])

		if err != nil {
			fmt.Println("supplied connectionCount is not a valid number")
			return
		}
	}

	if len(os.Args) > 4 {
		var err error
		loopCount, err = strconv.Atoi(os.Args[4])

		if err != nil {
			fmt.Println("supplied connectionCount is not a valid number")
			return
		}
	}

	var ws sync.WaitGroup
	ws.Add(connectionCount)

	startTime := time.Now()

	for i := 0; i < connectionCount; i++ {
		go connect(&ws)
	}

	ws.Wait()

	fmt.Println("duration:", time.Now().Sub(startTime))
	return
}
