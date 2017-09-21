package calendar

import (
	"google.golang.org/grpc"
	"log"
	"context"
	pb "github.com/integraal/chat-ops-calendar/calendar"
)

type Config struct {
	Address string `json:"address"`
}

var client pb.CalendarClient

func Initialize(config Config) {
	conn, err := grpc.Dial(config.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	client = pb.NewCalendarClient(conn)
}

func GetEvents(email string) ([]*pb.Event, error) {
	collection, err := client.GetEvents(context.Background(), &pb.EventRequest{
		Email: email,
		Start: "-2days",
		End:   "+2days",
	})
	if err != nil {
		return nil, err
	}
	return collection.Items, nil
}
