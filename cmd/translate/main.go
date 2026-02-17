package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/organic-programming/abel-fishel-translator/gen/go/abel_fishel_translator/v1"
	"github.com/organic-programming/abel-fishel-translator/pkg/translator"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "translate: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return usageErr()
	}

	svc := translator.NewService()

	switch args[0] {
	case "check":
		return runCheck(svc, args[1:])
	case "status":
		return runStatus(svc)
	default:
		return runTranslate(svc, args)
	}
}

func runTranslate(svc *translator.Service, args []string) error {
	if len(args) == 0 {
		return usageErr()
	}

	path := args[0]
	fs := flag.NewFlagSet("translate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	to := fs.String("to", "", "target language")
	from := fs.String("from", "", "source language")
	if err := fs.Parse(args[1:]); err != nil {
		return usageErr()
	}
	if strings.TrimSpace(*to) == "" {
		return errors.New("--to is required")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	resp, err := svc.Translate(context.Background(), &pb.TranslateRequest{
		InputMarkdown: string(content),
		ToLang:        *to,
		FromLang:      *from,
	})
	if err != nil {
		return err
	}

	out := resp.GetTranslatedMarkdown()
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	_, err = io.WriteString(os.Stdout, out)
	return err
}

func runCheck(svc *translator.Service, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: translate check <file>")
	}

	path := args[0]
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	resp, err := svc.CheckTranslation(context.Background(), &pb.CheckTranslationRequest{
		OriginMarkdown:     string(content),
		TranslatedMarkdown: string(content),
	})
	if err != nil {
		return err
	}

	status := "stale"
	if resp.GetUpToDate() {
		status = "up to date"
	}
	fmt.Printf("%s: %s\n", filepath.Base(path), status)
	return nil
}

func runStatus(svc *translator.Service) error {
	resp, err := svc.TranslationStatus(context.Background(), &pb.TranslationStatusRequest{})
	if err != nil {
		return err
	}
	fmt.Printf("coverage: %.2f%% (%d/%d)\n", resp.GetCoverage()*100, resp.GetTranslatedDocuments(), resp.GetTotalDocuments())
	return nil
}

func usageErr() error {
	return errors.New("usage: translate <file> --to <lang> [--from <lang>] | translate check <file> | translate status")
}
