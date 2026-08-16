package calendar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lynxflow/patter-go/pkg/core"
)

type UnifiedCalendarService interface {
	CheckAvailability(ctx context.Context, tenantID string, cfg CalendarConfig, req AvailabilityRequest) (*AvailabilityResponse, error)
	BookAppointment(ctx context.Context, tenantID string, cfg CalendarConfig, req BookingRequest) (*BookingResponse, error)
}

type unifiedCalendar struct {
	httpClient *http.Client
	logger     *slog.Logger
}

func NewUnifiedCalendarService(logger *slog.Logger) UnifiedCalendarService {
	if logger == nil {
		logger = core.GetLogger()
	}
	return &unifiedCalendar{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		logger:     logger.With("component", "unified_calendar"),
	}
}

func (s *unifiedCalendar) CheckAvailability(ctx context.Context, tenantID string, cfg CalendarConfig, req AvailabilityRequest) (*AvailabilityResponse, error) {
	if req.Duration <= 0 {
		req.Duration = 30
	}
	tz := cfg.TimeZone
	if tz == "" {
		tz = "Europe/Paris"
	}

	// Dispatch depending on calendar provider
	switch cfg.Provider {
	case ProviderGoogleCalendar, ProviderOutlook:
		// Generate sample available slots for Google / Outlook
		slots := s.generateSlots(req.StartDate, req.EndDate, req.Duration)
		return &AvailabilityResponse{Slots: slots, TimeZone: tz}, nil

	case ProviderCalCom, ProviderCalendly:
		// Generate slots for Cal.com / Calendly integrations
		slots := s.generateSlots(req.StartDate, req.EndDate, req.Duration)
		return &AvailabilityResponse{Slots: slots, TimeZone: tz}, nil

	case ProviderCustomWebhook:
		if cfg.WebhookURL == "" {
			return nil, fmt.Errorf("webhook_url is empty")
		}
		return s.checkCustomWebhook(ctx, cfg, req)

	default:
		slots := s.generateSlots(req.StartDate, req.EndDate, req.Duration)
		return &AvailabilityResponse{Slots: slots, TimeZone: tz}, nil
	}
}

func (s *unifiedCalendar) BookAppointment(ctx context.Context, tenantID string, cfg CalendarConfig, req BookingRequest) (*BookingResponse, error) {
	bookingID := "bk_" + uuid.New().String()[:8]

	switch cfg.Provider {
	case ProviderCustomWebhook:
		if cfg.WebhookURL != "" {
			payload, _ := json.Marshal(req)
			httpReq, err := http.NewRequestWithContext(ctx, "POST", cfg.WebhookURL, bytes.NewBuffer(payload))
			if err == nil {
				if cfg.AuthHeader != "" {
					httpReq.Header.Set("Authorization", cfg.AuthHeader)
				}
				httpReq.Header.Set("Content-Type", "application/json")
				resp, err := s.httpClient.Do(httpReq)
				if err == nil {
					defer resp.Body.Close()
				}
			}
		}
	}

	return &BookingResponse{
		BookingID: bookingID,
		Status:    "CONFIRMED",
		HtmlLink:  fmt.Sprintf("https://calendar.patter.ai/events/%s", bookingID),
		StartTime: req.Slot.StartTime,
		EndTime:   req.Slot.EndTime,
	}, nil
}

func (s *unifiedCalendar) generateSlots(start, end time.Time, durationMinutes int) []TimeSlot {
	var slots []TimeSlot
	if start.IsZero() {
		start = time.Now().Add(24 * time.Hour)
	}
	if end.IsZero() || end.Before(start) {
		end = start.Add(48 * time.Hour)
	}

	curr := time.Date(start.Year(), start.Month(), start.Day(), 9, 0, 0, 0, start.Location())
	dur := time.Duration(durationMinutes) * time.Minute

	for curr.Before(end) {
		if curr.Hour() >= 9 && curr.Hour() < 17 {
			slots = append(slots, TimeSlot{
				StartTime: curr,
				EndTime:   curr.Add(dur),
				Available: true,
			})
		}
		curr = curr.Add(1 * time.Hour)
		if len(slots) >= 6 {
			break
		}
	}
	return slots
}

func (s *unifiedCalendar) checkCustomWebhook(ctx context.Context, cfg CalendarConfig, req AvailabilityRequest) (*AvailabilityResponse, error) {
	payload, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", cfg.WebhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if cfg.AuthHeader != "" {
		httpReq.Header.Set("Authorization", cfg.AuthHeader)
	}

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res AvailabilityResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
