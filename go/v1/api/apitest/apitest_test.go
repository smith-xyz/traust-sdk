package apitest_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/traust-security/traust-sdk/go/v1/api"
	"github.com/traust-security/traust-sdk/go/v1/api/apitest"
)

func TestFakeCoreAnswersRecordsAndRefusesUnknownRoutes(t *testing.T) {
	core := apitest.New(t)
	core.Handle("POST /v1/orders/claim", apitest.Status(http.StatusNoContent))
	client := core.Client()
	order, err := client.Claim(context.Background(), api.ClaimRequest{LeaseSeconds: 30})
	if err != nil || order != nil {
		t.Fatalf("an empty claim: order=%v err=%v", order, err)
	}
	_, err = client.GetRepository(context.Background(), uuid.New())
	var problem *api.ProblemError
	if !errors.As(err, &problem) || problem.StatusCode != http.StatusNotFound {
		t.Fatalf("unhandled route: %v", err)
	}
	sent := core.Requests()
	if len(sent) != 2 || sent[0].Method != http.MethodPost || !strings.Contains(string(sent[0].Body), `"lease_seconds":30`) {
		t.Fatalf("recorded %+v", sent)
	}
}
