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

// This tutorial implenents a basic Match Function that is hosted in the below
// configured port. You can also configure the Open Match QueryService endpoint
// with which the Match Function communicates to query the Tickets.

const (
	queryServiceAddress = "open-match-query.open-match.svc.cluster.local:50503" // Address of the QueryService endpoint.
	serverPort          = 50502                                                 // The port for hosting the Match Function.
)

func main() {
	Start(queryServiceAddress, serverPort)
}

// MatchFunctionService implements pb.MatchFunctionServer, the server generated
// by compiling the protobuf, by fulfilling the pb.MatchFunctionServer interface.
type MatchFunctionService struct {
	grpc               *grpc.Server
	queryServiceClient pb.QueryServiceClient
	port               int
}

// Start creates and starts the Match Function server and also connects to Open
// Match's queryService. This connection is used at runtime to fetch tickets
// for pools specified in MatchProfile.
func Start(queryServiceAddr string, serverPort int) {
	// Connect to QueryService.
	conn, err := grpc.Dial(queryServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Open Match, got %s", err.Error())
	}
	defer conn.Close()

	mmfService := MatchFunctionService{
		queryServiceClient: pb.NewQueryServiceClient(conn),
	}

	// Create and host a new gRPC service on the configured port.
	server := grpc.NewServer()
	pb.RegisterMatchFunctionServer(server, &mmfService)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		log.Fatalf("TCP net listener initialization failed for port %v, got %s", serverPort, err.Error())
	}

	log.Printf("TCP net listener initialized for port %v", serverPort)
	err = server.Serve(ln)
	if err != nil {
		log.Fatalf("gRPC serve failed, got %s", err.Error())
	}
}

// This match function fetches all the Tickets for all the pools specified in
// the profile. It uses a configured number of tickets from each pool to generate
// a Match Proposal. It continues to generate proposals till one of the pools
// runs out of Tickets.
const (
	matchName              = "basic-matchfunction"
	ticketsPerPoolPerMatch = 4
)

// Run is this match function's implementation of the gRPC call defined in api/matchfunction.proto.
func (s *MatchFunctionService) Run(req *pb.RunRequest, stream pb.MatchFunction_RunServer) error {
	// Fetch tickets for the pools specified in the Match Profile.
	log.Printf("Generating proposals for function %v", req.GetProfile().GetName())

	poolTickets, err := matchfunction.QueryPools(stream.Context(), s.queryServiceClient, req.GetProfile().GetPools())
	if err != nil {
		log.Printf("Failed to query tickets for the given pools, got %s", err.Error())
		return err
	}

	// Generate proposals.
	proposals, err := makeMatches(req.GetProfile(), poolTickets)
	if err != nil {
		log.Printf("Failed to generate matches, got %s", err.Error())
		return err
	}

	log.Printf("Streaming %v proposals to Open Match", len(proposals))
	// Stream the generated proposals back to Open Match.
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
				// This pool is completely drained out. Stop creating matches.
				insufficientTickets = true
				break
			}

			// Remove the Tickets from this pool and add to the match proposal.
			matchTickets = append(matchTickets, tickets[0:ticketsPerPoolPerMatch]...)
			poolTickets[pool] = tickets[ticketsPerPoolPerMatch:]
		}

		if insufficientTickets {
			break
		}

		// Add score calculation logic here.

		matches = append(matches, &pb.Match{
			MatchId:       fmt.Sprintf("profile-%v-time-%v-%v", p.GetName(), time.Now().Format("2006-01-02T15:04:05.00"), count),
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
	// Add your logic to compute the score for this Match here
	return matchScore
}
