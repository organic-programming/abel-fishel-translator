---
# Holon Identity v1
uuid: "31d98dfa-c08e-4db5-ab4f-180cecad50ef"
given_name: "Abel-Fishel"
family_name: "Translator"
motto: "One thought, every tongue."
composer: "B. ALTER"
clade: "probabilistic/generative"
status: draft
born: "2026-02-12"

# Lineage
parents: []
reproduction: "manual"

# Pinning
binary_path: null
binary_version: "0.1.0"
git_tag: null
git_commit: null
os: "darwin"
arch: "arm64"
dependencies: []

# Optional
aliases: ["translate", "babel"]
wrapped_license: null

# Metadata
generated_by: "sophia-who"
lang: "go"
proto_status: draft
---

# Abel-Fishel Translator

> *"One thought, every tongue."*

## Description

Abel-Fishel Translator is the multilingual holon. It translates Markdown
documents between languages while preserving structure, YAML frontmatter,
code blocks, and Cartouche metadata.

Named after Douglas Adams' Babel Fish — the small, yellow, leech-like
creature that provides instant translation when placed in the ear.

## Commands

```
translate <file> --to <lang>           — translate a document to a target language
translate <file> --to <lang> --from <lang> — explicit source language
translate check <file>                 — verify a translation is up to date with its origin
translate status                       — show translation coverage for all documents
```

## Contract

- Proto file: `babel_fish.proto`
- Service: `BabelFishService`
- RPCs: `Translate`, `CheckTranslation`, `TranslationStatus`
