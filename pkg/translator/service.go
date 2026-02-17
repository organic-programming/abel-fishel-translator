package translator

import (
	"context"
	"strings"

	pb "github.com/organic-programming/abel-fishel-translator/gen/go/abel_fishel_translator/v1"
)

// Service provides a minimal implementation of the BabelFish gRPC contract.
type Service struct {
	pb.UnimplementedBabelFishServiceServer
}

// NewService builds a new translator service instance.
func NewService() *Service {
	return &Service{}
}

// Translate currently applies a deterministic placeholder transform.
func (s *Service) Translate(_ context.Context, req *pb.TranslateRequest) (*pb.TranslateResponse, error) {
	detected := req.GetFromLang()
	if detected == "" {
		detected = "auto"
	}

	translated := req.GetInputMarkdown()
	if target := strings.TrimSpace(req.GetToLang()); target != "" {
		translated = "<" + target + "> " + translated
	}

	return &pb.TranslateResponse{
		TranslatedMarkdown: translated,
		DetectedSourceLang: detected,
	}, nil
}

// CheckTranslation reports stale status for missing source or translated text.
func (s *Service) CheckTranslation(_ context.Context, req *pb.CheckTranslationRequest) (*pb.CheckTranslationResponse, error) {
	upToDate := true
	reasons := make([]string, 0, 2)

	if strings.TrimSpace(req.GetOriginMarkdown()) == "" {
		upToDate = false
		reasons = append(reasons, "origin_markdown is empty")
	}
	if strings.TrimSpace(req.GetTranslatedMarkdown()) == "" {
		upToDate = false
		reasons = append(reasons, "translated_markdown is empty")
	}

	return &pb.CheckTranslationResponse{
		UpToDate: upToDate,
		Reasons:  reasons,
	}, nil
}

// TranslationStatus computes a simple coverage ratio from document entries.
func (s *Service) TranslationStatus(_ context.Context, req *pb.TranslationStatusRequest) (*pb.TranslationStatusResponse, error) {
	total := len(req.GetDocuments())
	translated := 0
	missing := make([]*pb.TranslationDocument, 0, total)

	for _, doc := range req.GetDocuments() {
		if doc.GetExists() {
			translated++
			continue
		}
		missing = append(missing, doc)
	}

	coverage := float32(0)
	if total > 0 {
		coverage = float32(translated) / float32(total)
	}

	return &pb.TranslationStatusResponse{
		TotalDocuments:      int32(total),
		TranslatedDocuments: int32(translated),
		MissingDocuments:    missing,
		Coverage:            coverage,
	}, nil
}
