package calendar_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lynxflow/patter-go/pkg/calendar"
)

func TestCalendarAvailabilityAndBooking(t *testing.T) {
	svc := calendar.NewUnifiedCalendarService(nil)

	cfg := calendar.CalendarConfig{
		Provider: calendar.ProviderGoogleCalendar,
		TimeZone: "Europe/Paris",
	}

	req := calendar.AvailabilityRequest{
		StartDate: time.Now(),
		Duration:  30,
	}

	avail, err := svc.CheckAvailability(context.Background(), "tenant_1", cfg, req)
	if err != nil {
		t.Fatalf("Unexpected error checking availability: %v", err)
	}

	if len(avail.Slots) == 0 {
		t.Errorf("Expected available slots generated, got 0")
	}

	bookReq := calendar.BookingRequest{
		Slot: avail.Slots[0],
		Attendee: calendar.AttendeeInfo{
			Name:  "Jean Dupont",
			Email: "jean@example.com",
			Phone: "+33612345678",
		},
		Summary: "Rendez-vous Démonstration",
	}

	booking, err := svc.BookAppointment(context.Background(), "tenant_1", cfg, bookReq)
	if err != nil {
		t.Fatalf("Unexpected error booking appointment: %v", err)
	}

	if booking.Status != "CONFIRMED" {
		t.Errorf("Expected status CONFIRMED, got %s", booking.Status)
	}
	if booking.BookingID == "" {
		t.Errorf("Expected valid booking ID")
	}
}

func TestCalendarCustomWebhook(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(calendar.AvailabilityResponse{
			TimeZone: "Europe/Paris",
			Slots: []calendar.TimeSlot{
				{StartTime: time.Now(), EndTime: time.Now().Add(30 * time.Minute), Available: true},
			},
		})
	}))
	defer ts.Close()

	svc := calendar.NewUnifiedCalendarService(nil)

	cfg := calendar.CalendarConfig{
		Provider:   calendar.ProviderCustomWebhook,
		WebhookURL: ts.URL,
	}

	res, err := svc.CheckAvailability(context.Background(), "tenant_1", cfg, calendar.AvailabilityRequest{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(res.Slots) != 1 {
		t.Errorf("Expected 1 slot returned from webhook, got %d", len(res.Slots))
	}
}
