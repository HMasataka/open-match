package main

import (
	"open-match.dev/open-match/pkg/pb"
)

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
	// modes := []string{"mode.demo", "mode.ctf", "mode.battleroyale", "mode.2v2"}
	// Implement your logic to return any two of the above game-modes randomly
	return []string{}
}
