package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"open-match.dev/open-match/pkg/matchfunction"
	"open-match.dev/open-match/pkg/pb"
)

const (
	queryServiceAddress = "open-match-query.open-match.svc.cluster.local:50503"
	serverPort          = 50502
)

func main() {
	Start(queryServiceAddress, serverPort)
}

type MatchFunctionService struct {
	queryServiceClient pb.QueryServiceClient
}

func Start(queryServiceAddr string, serverPort int) {
	conn, err := grpc.Dial(queryServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Open Match, got %s", err.Error())
	}
	defer conn.Close()

	mmfService := MatchFunctionService{
		queryServiceClient: pb.NewQueryServiceClient(conn),
	}

	server := grpc.NewServer()

	pb.RegisterMatchFunctionServer(server, &mmfService)

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		log.Fatalf("TCP net listener initialization failed for port %v, got %s", serverPort, err.Error())
	}

	log.Printf("TCP net listener initialized for port %v", serverPort)

	if err := server.Serve(ln); err != nil {
		log.Fatalf("gRPC serve failed, got %s", err.Error())
	}
}

const (
	matchName              = "basic-matchfunction"
	ticketsPerPoolPerMatch = 4
)

func (s *MatchFunctionService) Run(req *pb.RunRequest, stream pb.MatchFunction_RunServer) error {
	log.Printf("Generating proposals for function %v", req.GetProfile().GetName())

	poolTickets, err := matchfunction.QueryPools(stream.Context(), s.queryServiceClient, req.GetProfile().GetPools())
	if err != nil {
		log.Printf("Failed to query tickets for the given pools, got %s", err.Error())
		return err
	}

	proposals, err := makeMatches(req.GetProfile(), poolTickets)
	if err != nil {
		log.Printf("Failed to generate matches, got %s", err.Error())
		return err
	}

	log.Printf("Streaming %v proposals to Open Match", len(proposals))

	for _, proposal := range proposals {
		if err := stream.Send(&pb.RunResponse{Proposal: proposal}); err != nil {
			log.Printf("Failed to stream proposals to Open Match, got %s", err.Error())
			return err
		}
	}

	return nil
}

func makeMatches(p *pb.MatchProfile, poolTickets map[string][]*pb.Ticket) ([]*pb.Match, error) {
	var matches []*pb.Match

	count := 0

	for {
		insufficientTickets := false
		matchTickets := []*pb.Ticket{}

		for pool, tickets := range poolTickets {
			if len(tickets) < ticketsPerPoolPerMatch {
				insufficientTickets = true
				break
			}

			matchTickets = append(matchTickets, tickets[0:ticketsPerPoolPerMatch]...)
			poolTickets[pool] = tickets[ticketsPerPoolPerMatch:]
		}

		if insufficientTickets {
			break
		}

		// Add score calculation logic here.

		matches = append(matches, &pb.Match{
			MatchId:       fmt.Sprintf("profile-%v-time-%v-%v", p.GetName(), time.Now().Format(time.RFC3339Nano), count),
			MatchProfile:  p.GetName(),
			MatchFunction: matchName,
			Tickets:       matchTickets,
		})

		count++
	}

	return matches, nil
}

func scoreCalculator(tickets []*pb.Ticket) float64 {
	matchScore := 0.0
	return matchScore
}
