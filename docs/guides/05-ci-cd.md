---
title: CI/CD Pipelines
description: Automate documentation builds in CI/CD pipelines
tags:
  - guides
  - ci-cd
  - automation
---

# CI/CD Pipelines

Automate documentation builds and deployments.

## GitHub Actions

### Basic Build

```yaml
# .github/workflows/docs.yml
name: Build Documentation

on:
  push:
    branches: [main]
    paths: ['docs/**']
  pull_request:
    paths: ['docs/**']

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Install MinimalDoc
        run: go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest

      - name: Build
        run: minimaldoc build ./docs --output dist

      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: docs
          path: dist/
```

### Build and Deploy

```yaml
name: Deploy Documentation

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Install MinimalDoc
        run: go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest

      - name: Build
        run: |
          minimaldoc build ./docs \
            --base-url "${{ vars.DOCS_URL }}" \
            --output dist

      - name: Deploy to S3
        run: aws s3 sync dist/ s3://${{ vars.S3_BUCKET }}/ --delete
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
```

### PR Preview

```yaml
name: Preview Documentation

on:
  pull_request:
    paths: ['docs/**']

jobs:
  preview:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Build Preview
        run: |
          go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest
          minimaldoc build ./docs \
            --base-url "https://preview.example.com/pr-${{ github.event.number }}" \
            --output dist

      - name: Deploy Preview
        id: deploy
        run: |
          # Deploy to preview environment
          echo "url=https://preview.example.com/pr-${{ github.event.number }}" >> $GITHUB_OUTPUT

      - name: Comment PR
        uses: actions/github-script@v7
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: '📚 Documentation preview: ${{ steps.deploy.outputs.url }}'
            })
```

## GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - build
  - deploy

variables:
  GO_VERSION: "1.24"

build-docs:
  stage: build
  image: golang:${GO_VERSION}
  script:
    - go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest
    - minimaldoc build ./docs --output dist
  artifacts:
    paths:
      - dist/
    expire_in: 1 week
  only:
    changes:
      - docs/**/*

deploy-docs:
  stage: deploy
  image: alpine:latest
  dependencies:
    - build-docs
  script:
    - apk add --no-cache aws-cli
    - aws s3 sync dist/ s3://${S3_BUCKET}/ --delete
  only:
    - main
  environment:
    name: production
    url: https://docs.example.com
```

## CircleCI

```yaml
# .circleci/config.yml
version: 2.1

jobs:
  build-docs:
    docker:
      - image: cimg/go:1.24
    steps:
      - checkout
      - run:
          name: Install MinimalDoc
          command: go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest
      - run:
          name: Build Documentation
          command: minimaldoc build ./docs --output dist
      - persist_to_workspace:
          root: .
          paths:
            - dist/

  deploy-docs:
    docker:
      - image: cimg/aws:2024.03
    steps:
      - attach_workspace:
          at: .
      - run:
          name: Deploy to S3
          command: aws s3 sync dist/ s3://${S3_BUCKET}/ --delete

workflows:
  build-and-deploy:
    jobs:
      - build-docs:
          filters:
            branches:
              only: main
      - deploy-docs:
          requires:
            - build-docs
          filters:
            branches:
              only: main
```

## Jenkins

```groovy
// Jenkinsfile
pipeline {
    agent any

    environment {
        GO_VERSION = '1.24'
        DOCS_URL = 'https://docs.example.com'
    }

    stages {
        stage('Setup') {
            steps {
                sh '''
                    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
                    tar -xzf go${GO_VERSION}.linux-amd64.tar.gz
                    export PATH=$PWD/go/bin:$PATH
                    go install github.com/studiowebux/minimaldoc/cmd/minimaldoc@latest
                '''
            }
        }

        stage('Build') {
            steps {
                sh '''
                    export PATH=$PWD/go/bin:$HOME/go/bin:$PATH
                    minimaldoc build ./docs --base-url ${DOCS_URL} --output dist
                '''
            }
        }

        stage('Deploy') {
            when {
                branch 'main'
            }
            steps {
                withAWS(credentials: 'aws-credentials') {
                    sh 'aws s3 sync dist/ s3://${S3_BUCKET}/ --delete'
                }
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: 'dist/**/*', fingerprint: true
        }
    }
}
```

## Caching

### GitHub Actions

```yaml
- name: Cache Go modules
  uses: actions/cache@v4
  with:
    path: |
      ~/go/pkg/mod
      ~/go/bin
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

### GitLab CI

```yaml
build-docs:
  cache:
    key: go-modules
    paths:
      - /go/pkg/mod/
```

## Environment Variables

### Configuration

```yaml
# GitHub Actions
env:
  DOCS_BASE_URL: ${{ vars.DOCS_BASE_URL }}
  DOCS_TITLE: ${{ vars.DOCS_TITLE }}

- name: Build
  run: |
    minimaldoc build ./docs \
      --base-url "$DOCS_BASE_URL" \
      --title "$DOCS_TITLE" \
      --output dist
```

### Secrets

```yaml
# Never log secrets
- name: Deploy
  run: aws s3 sync dist/ s3://$BUCKET/
  env:
    AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
    AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
```

## Multi-Environment

### Staging and Production

```yaml
name: Deploy Documentation

on:
  push:
    branches:
      - main
      - develop

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set Environment
        run: |
          if [ "${{ github.ref }}" == "refs/heads/main" ]; then
            echo "ENV=production" >> $GITHUB_ENV
            echo "BASE_URL=https://docs.example.com" >> $GITHUB_ENV
            echo "S3_BUCKET=docs-prod" >> $GITHUB_ENV
          else
            echo "ENV=staging" >> $GITHUB_ENV
            echo "BASE_URL=https://staging-docs.example.com" >> $GITHUB_ENV
            echo "S3_BUCKET=docs-staging" >> $GITHUB_ENV
          fi

      - name: Build
        run: |
          minimaldoc build ./docs \
            --base-url "$BASE_URL" \
            --output dist

      - name: Deploy
        run: aws s3 sync dist/ s3://$S3_BUCKET/
```

## Validation

### Link Checking

```yaml
- name: Check Links
  run: |
    npm install -g linkinator
    linkinator dist/ --recurse --skip "^(?!http)"
```

### HTML Validation

```yaml
- name: Validate HTML
  run: |
    npm install -g html-validate
    html-validate "dist/**/*.html"
```

## Notifications

### Slack

```yaml
- name: Notify Slack
  if: always()
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    fields: repo,message,commit,author
  env:
    SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK }}
```

### Discord

```yaml
- name: Notify Discord
  if: failure()
  run: |
    curl -H "Content-Type: application/json" \
      -d '{"content": "Documentation build failed!"}' \
      ${{ secrets.DISCORD_WEBHOOK }}
```

## Scheduled Builds

Rebuild docs periodically:

```yaml
on:
  schedule:
    - cron: '0 0 * * *'  # Daily at midnight
  push:
    branches: [main]
```

## Monorepo Support

Build multiple doc sites:

```yaml
jobs:
  build:
    strategy:
      matrix:
        project: [api, sdk, cli]
    steps:
      - name: Build ${{ matrix.project }}
        run: |
          minimaldoc build ./packages/${{ matrix.project }}/docs \
            --output dist/${{ matrix.project }}
```
