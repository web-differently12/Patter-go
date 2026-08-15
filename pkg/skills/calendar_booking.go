package skills

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/lynxflow/patter-go/pkg/calendar"
	"github.com/lynxflow/patter-go/pkg/core"
)

type CalendarBookingSkill struct {
	calService calendar.UnifiedCalendarService
	logger     *slog.Logger
}

func NewCalendarBookingSkill(svc calendar.UnifiedCalendarService, logger *slog.Logger) *CalendarBookingSkill {
	if logger == nil {
		logger = core.GetLogger()
	}
	if svc == nil {
		svc = calendar.NewUnifiedCalendarService(logger)
	}
	return &CalendarBookingSkill{
		calService: svc,
		logger:     logger.With("skill", "calendar_booking"),
	}
}

func (s *CalendarBookingSkill) GetToolDefinitions() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "check_calendar_availability",
				"description": "Vérifie les créneaux disponibles dans l'agenda du conseiller (Google Calendar, Outlook, Cal.com, Calendly).",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"start_date": map[string]interface{}{
							"type":        "string",
							"description": "Date de début souhaitée au format YYYY-MM-DD.",
						},
						"duration_minutes": map[string]interface{}{
							"type":        "integer",
							"description": "Durée du rendez-vous en minutes (défaut 30).",
						},
					},
					"required": []string{"start_date"},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "book_calendar_appointment",
				"description": "Réserve un créneau de rendez-vous avec confirmation des coordonnées du prospect.",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"start_time": map[string]interface{}{
							"type":        "string",
							"description": "Heure de début au format ISO RFC3339 (ex: '2026-08-16T10:00:00Z').",
						},
						"attendee_name": map[string]interface{}{
							"type":        "string",
							"description": "Nom complet du prospect.",
						},
						"attendee_email": map[string]interface{}{
							"type":        "string",
							"description": "Email du prospect.",
						},
						"attendee_phone": map[string]interface{}{
							"type":        "string",
							"description": "Téléphone E.164 du prospect.",
						},
						"summary": map[string]interface{}{
							"type":        "string",
							"description": "Titre / Objet du rendez-vous.",
						},
					},
					"required": []string{"start_time", "attendee_name", "attendee_phone"},
				},
			},
		},
	}
}

func (s *CalendarBookingSkill) CheckAvailability(ctx context.Context, tenantID, rawArgs string, cfg calendar.CalendarConfig) (string, error) {
	var args struct {
		StartDate string `json:"start_date"`
		Duration  int    `json:"duration_minutes"`
	}
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return "Erreur d'analyse des arguments", err
	}

	start, err := time.Parse("2006-01-02", args.StartDate)
	if err != nil {
		start = time.Now().Add(24 * time.Hour)
	}

	resp, err := s.calService.CheckAvailability(ctx, tenantID, cfg, calendar.AvailabilityRequest{
		StartDate: start,
		EndDate:   start.Add(48 * time.Hour),
		Duration:  args.Duration,
	})
	if err != nil {
		return "Impossible de vérifier la disponibilité de l'agenda", err
	}

	out, _ := json.Marshal(resp)
	return string(out), nil
}

func (s *CalendarBookingSkill) BookAppointment(ctx context.Context, tenantID, rawArgs string, cfg calendar.CalendarConfig) (string, error) {
	var args struct {
		StartTime     string `json:"start_time"`
		AttendeeName  string `json:"attendee_name"`
		AttendeeEmail string `json:"attendee_email"`
		AttendeePhone string `json:"attendee_phone"`
		Summary       string `json:"summary"`
	}
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return "Erreur d'analyse des arguments de réservation", err
	}

	st, err := time.Parse(time.RFC3339, args.StartTime)
	if err != nil {
		st = time.Now().Add(24 * time.Hour)
	}

	resp, err := s.calService.BookAppointment(ctx, tenantID, cfg, calendar.BookingRequest{
		Slot: calendar.TimeSlot{
			StartTime: st,
			EndTime:   st.Add(30 * time.Minute),
			Available: true,
		},
		Attendee: calendar.AttendeeInfo{
			Name:  args.AttendeeName,
			Email: args.AttendeeEmail,
			Phone: args.AttendeePhone,
		},
		Summary: args.Summary,
	})
	if err != nil {
		return "Échec de la réservation du rendez-vous", err
	}

	return fmt.Sprintf("RENDEZ_VOUS_CONFIRMÉ: Réservation ID [%s] enregistrée pour %s le %s.", resp.BookingID, args.AttendeeName, st.Format("02/01/2006 a 15:04")), nil
}
