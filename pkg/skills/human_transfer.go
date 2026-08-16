package skills

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/lynxflow/patter-go/pkg/core"
)

// HumanTransferArgs sont les arguments passés par le LLM lors du Tool Call
type HumanTransferArgs struct {
	Reason             string `json:"reason"`               // Pourquoi on transfère (ex: "Question hors contrat")
	UserConsented      bool   `json:"user_consented"`       // Est-ce que l'utilisateur a validé oralement ?
	ProposedDepartment string `json:"proposed_department"` // "TECH", "SALES", "SUPPORT"
}

type CallDispatcher interface {
	TransferCall(ctx context.Context, tenantID, callSID, destinationNumber, reason string) error
}

type HumanTransferSkill struct {
	dispatcher CallDispatcher
	logger     *slog.Logger
}

func NewHumanTransferSkill(dispatcher CallDispatcher, logger *slog.Logger) *HumanTransferSkill {
	if logger == nil {
		logger = core.GetLogger()
	}
	return &HumanTransferSkill{
		dispatcher: dispatcher,
		logger:     logger.With("skill", "human_transfer"),
	}
}

// GetToolDefinition renvoie le schéma JSON Tool Call pour le LLM
func (s *HumanTransferSkill) GetToolDefinition() map[string]interface{} {
	return map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "request_human_transfer",
			"description": "Demande le transfert de l'appel vers un conseiller humain après avoir obtenu le consentement oral de l'utilisateur.",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"reason": map[string]interface{}{
						"type":        "string",
						"description": "Raison exacte du transfert (ex: 'Information tarifaire complexe introuvable dans la doc').",
					},
					"user_consented": map[string]interface{}{
						"type":        "boolean",
						"description": "Mettre à true UNIQUEMENT si l'utilisateur a explicitement dit 'Oui' ou donné son accord pour être transféré.",
					},
					"proposed_department": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"SALES", "SUPPORT", "TECHNICAL"},
						"description": "Service concerné par le transfert.",
					},
				},
				"required": []string{"reason", "user_consented"},
			},
		},
	}
}

// Execute exécute l'action de transfert dès que user_consented == true
func (s *HumanTransferSkill) Execute(ctx context.Context, tenantID, callSID, rawArgs string, destinationNumber string) (string, error) {
	var args HumanTransferArgs
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return "Erreur d'analyse des arguments du skill", err
	}

	// 🛑 Si l'utilisateur n'a pas encore consenti, le skill rappelle au LLM de demander
	if !args.UserConsented {
		return "CONSENTEMENT_REQUIS: Demande d'abord à l'utilisateur : 'Souhaitez-vous que je vous mette en relation avec mon collègue spécialiste ?'", nil
	}

	s.logger.Info("Executing human transfer skill",
		"call_sid", callSID,
		"destination", destinationNumber,
		"reason", args.Reason,
	)

	if s.dispatcher != nil {
		err := s.dispatcher.TransferCall(ctx, tenantID, callSID, destinationNumber, args.Reason)
		if err != nil {
			return "Échec technique du transfert vers le conseiller.", err
		}
	}

	return "TRANSFERT_INITIÉ: Informe l'utilisateur qu'il va être mis en relation et qu'il doit patienter quelques instants.", nil
}
