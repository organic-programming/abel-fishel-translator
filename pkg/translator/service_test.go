package translator

import (
	"context"
	"math"
	"testing"

	pb "github.com/organic-programming/abel-fishel-translator/gen/go/abel_fishel_translator/v1"
)

func TestTranslateAddsTargetPrefix(t *testing.T) {
	t.Parallel()

	svc := NewService()
	resp, err := svc.Translate(context.Background(), &pb.TranslateRequest{
		InputMarkdown: "Hello",
		ToLang:        "fr",
	})
	if err != nil {
		t.Fatalf("Translate returned error: %v", err)
	}
	if got, want := resp.GetTranslatedMarkdown(), "<fr> Hello"; got != want {
		t.Fatalf("TranslatedMarkdown = %q, want %q", got, want)
	}
	if got, want := resp.GetDetectedSourceLang(), "auto"; got != want {
		t.Fatalf("DetectedSourceLang = %q, want %q", got, want)
	}
}

func TestCheckTranslationDetectsMissingFields(t *testing.T) {
	t.Parallel()

	svc := NewService()
	resp, err := svc.CheckTranslation(context.Background(), &pb.CheckTranslationRequest{})
	if err != nil {
		t.Fatalf("CheckTranslation returned error: %v", err)
	}
	if resp.GetUpToDate() {
		t.Fatal("expected UpToDate = false")
	}
	if got, want := len(resp.GetReasons()), 2; got != want {
		t.Fatalf("len(Reasons) = %d, want %d", got, want)
	}
}

func TestTranslationStatusComputesCoverage(t *testing.T) {
	t.Parallel()

	svc := NewService()
	resp, err := svc.TranslationStatus(context.Background(), &pb.TranslationStatusRequest{
		Documents: []*pb.TranslationDocument{
			{Path: "a.md", Lang: "en", Exists: true},
			{Path: "b.fr.md", Lang: "fr", Exists: false},
		},
	})
	if err != nil {
		t.Fatalf("TranslationStatus returned error: %v", err)
	}
	if got, want := resp.GetTotalDocuments(), int32(2); got != want {
		t.Fatalf("TotalDocuments = %d, want %d", got, want)
	}
	if got, want := resp.GetTranslatedDocuments(), int32(1); got != want {
		t.Fatalf("TranslatedDocuments = %d, want %d", got, want)
	}
	if got, want := len(resp.GetMissingDocuments()), 1; got != want {
		t.Fatalf("len(MissingDocuments) = %d, want %d", got, want)
	}
	if got := resp.GetCoverage(); math.Abs(float64(got-0.5)) > 1e-6 {
		t.Fatalf("Coverage = %f, want 0.5", got)
	}
}
