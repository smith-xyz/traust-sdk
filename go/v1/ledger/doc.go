// Package ledger provides typed writes and reads for the peer ledger service.
// It validates write artifacts and is transport-agnostic through Provider;
// NewHTTPClient supplies the built-in HTTP transport.
//
//	client := ledger.NewHTTPClient("https://ledger.example.com",
//		ledger.WithBearerToken(os.Getenv("LEDGER_TOKEN")),
//	)
//	resp, err := client.SubmitTriageReport(ctx, ledger.TriageReportInput{...})
//	findings, err := client.GetFindings(ctx, "repo-a")
package ledger
