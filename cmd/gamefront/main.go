package main

import (
	"context"
	"log"
	"time"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"open-match.dev/open-match/pkg/pb"
)

const (
	omFrontendEndpoint = "open-match-frontend.open-match.svc.cluster.local:50504"
	ticketsPerIter     = 20
)

func main() {
	// Connect to Open Match Frontend.
	conn, err := grpc.Dial(omFrontendEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Open Match, got %v", err)
	}

	defer conn.Close()
	fe := pb.NewFrontendServiceClient(conn)
	for range time.Tick(time.Second * 2) {
		for i := 0; i <= ticketsPerIter; i++ {
			req := &pb.CreateTicketRequest{
				Ticket: makeTicket(),
			}

			resp, err := fe.CreateTicket(context.Background(), req)
			if err != nil {
				log.Printf("Failed to Create Ticket, got %s", err.Error())
				continue
			}

			log.Println("Ticket created successfully, id:", resp.Id)
			go deleteOnAssign(fe, resp)
		}
	}
}

// deleteOnAssign fetches the Ticket state periodically and deletes the Ticket
// once it has an assignment.
func deleteOnAssign(fe pb.FrontendServiceClient, t *pb.Ticket) {
	for {
		got, err := fe.GetTicket(context.Background(), &pb.GetTicketRequest{TicketId: t.GetId()})
		if err != nil {
			log.Fatalf("Failed to Get Ticket %v, got %s", t.GetId(), err.Error())
		}

		if got.GetAssignment() != nil {
			log.Printf("Ticket %v got assignment %v", got.GetId(), got.GetAssignment())
			break
		}

		time.Sleep(time.Second * 1)
	}

	_, err := fe.DeleteTicket(context.Background(), &pb.DeleteTicketRequest{TicketId: t.GetId()})
	if err != nil {
		log.Fatalf("Failed to Delete Ticket %v, got %s", t.GetId(), err.Error())
	}
}

func makeTicket() *pb.Ticket {
	ticket := &pb.Ticket{
		SearchFields: &pb.SearchFields{
			// Tags can support multiple values but for simplicity, the demo function
			// assumes only single mode selection per Ticket.
			Tags: gameModes(),
		},
	}

	return ticket
}

func enterQueueTime() float64 {
	// Implement your logic to return a random time interval.
	return 0.0
}

func gameModes() []string {
	modes := []string{"mode.demo", "mode.ctf", "mode.battleroyale"}

	return []string{lo.Sample(modes)}
}
