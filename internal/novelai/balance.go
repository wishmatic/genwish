package novelai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wishmatic/genwish/internal/utils"
)

const userDataPath = "/user/data"

type Balance struct {
	Anlas int64
	Tier  string
}

type userDataResponse struct {
	Subscription struct {
		TrainingStepsLeft json.RawMessage `json:"trainingStepsLeft"`
		Tier              json.RawMessage `json:"tier"`
	} `json:"subscription"`
}

func (c *Client) CheckBalance(ctx context.Context) (*Balance, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+userDataPath, nil)
	if err != nil {
		return nil, fmt.Errorf("novelai: build balance request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("novelai: read balance: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, utils.FmtHTTPErr(http.MethodGet, userDataPath, resp)
	}

	var out userDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("novelai: decode balance: %w", err)
	}

	anlas, err := parseAnlas(out.Subscription.TrainingStepsLeft)
	if err != nil {
		return nil, err
	}

	tier, err := parseTier(out.Subscription.Tier)
	if err != nil {
		return nil, err
	}

	return &Balance{Anlas: anlas, Tier: tier}, nil
}

func parseTier(raw json.RawMessage) (string, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return "", nil
	}

	var name string
	if err := json.Unmarshal(raw, &name); err == nil {
		return name, nil
	}

	var level int
	if err := json.Unmarshal(raw, &level); err != nil {
		return "", fmt.Errorf("novelai: subscription.tier is neither a string nor a number")
	}

	switch level {
	case 0:
		return "paper", nil
	case 1:
		return "tablet", nil
	case 2:
		return "scroll", nil
	case 3:
		return "opus", nil
	default:
		return strconv.Itoa(level), nil
	}
}

func parseAnlas(raw json.RawMessage) (int64, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return 0, fmt.Errorf("novelai: balance response has no subscription.trainingStepsLeft")
	}

	var object struct {
		FixedTrainingStepsLeft int64 `json:"fixedTrainingStepsLeft"`
		PurchasedTrainingSteps int64 `json:"purchasedTrainingSteps"`
	}

	if err := json.Unmarshal(raw, &object); err == nil {
		return object.FixedTrainingStepsLeft + object.PurchasedTrainingSteps, nil
	}

	var integer int64
	if err := json.Unmarshal(raw, &integer); err != nil {
		return 0, fmt.Errorf(
			"novelai: subscription.trainingStepsLeft is neither a number nor an object",
		)
	}

	return integer, nil
}
