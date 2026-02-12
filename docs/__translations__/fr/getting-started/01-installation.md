---
title: Installation
description: Installer MinimalDoc sur votre systeme
tags:
  - demarrage
  - installation
---

# Installation

## Prerequis

- Go 1.24 ou version ulterieure
- Git (optionnel, pour le clonage)

## Installer depuis les sources

```bash
git clone https://github.com/studiowebux/minimaldoc.git
cd minimaldoc
go build -o minimaldoc ./cmd/minimaldoc
```

Deplacez le binaire dans votre PATH:

```bash
# macOS/Linux
sudo mv minimaldoc /usr/local/bin/

# Ou ajoutez au repertoire utilisateur
mv minimaldoc ~/.local/bin/
```

## Go Install

```bash
go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest
```

Assurez-vous que `$GOPATH/bin` est dans votre PATH.

## Binaires pre-compiles

Telechargez les binaires pre-compiles depuis la [page des versions](https://github.com/studiowebux/minimaldoc/releases).

Disponible pour:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

## Verifier l'installation

Apres l'installation, verifiez que MinimalDoc fonctionne:

```bash
minimaldoc --version
```
