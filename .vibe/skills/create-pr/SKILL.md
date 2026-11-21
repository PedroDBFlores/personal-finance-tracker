---
name: create-pr
description: Load when the user asks to create/open a pull request in the "Personal Finance Tracker" repo (e.g. "create a PR for this branch"). Fills .github/PULL_REQUEST_TEMPLATE.md (description, Fixes #issue link, Type of change checkboxes, test summary) and opens the PR with gh against main.
---

# Create a pull request

Use this skill whenever the user requests a pull request in this repo.

Targets (hardcoded for this repo):

- Repository: `PedroDBFlores/personal-finance-tracker`, base branch `main`
- PR template: `.github/PULL_REQUEST_TEMPLATE.md` (placeholders: `<!description>`, `<!issue>`, `<!tests>`)
- Requires `gh` authenticated with the `repo` scope (`gh auth status`)
- Branch naming convention: `PFTXXXXXX` where `XXXXXX` is the issue number zero-padded to 6 digits (issue #8 -> `PFT000008`). Derive the linked issue number from the current branch name when it matches this pattern; otherwise ask the user which issue the PR fixes.

## Steps

1. **Gather context**: current branch (`git branch --show-current`), commits vs main
   (`git log main..HEAD --oneline`), changed files (`git diff --stat main...HEAD`).
   If the current branch is `main`, ask the user which branch/issue this PR is for
   before doing anything else.
2. **Fill the template.** Read `.github/PULL_REQUEST_TEMPLATE.md` and produce the PR
   body by replacing the placeholders:
   - `<!description>` -> a concise summary of the changes (derived from the commits /
     diff) plus motivation and context. If the intent is unclear from the commits,
     ask the user.
   - `Fixes # <!issue>` -> `Fixes [#<N>](https://github.com/PedroDBFlores/personal-finance-tracker/issues/<N>)`
     with the issue number `<N>` (from the branch name or the user).
   - **Type of change**: keep only the relevant checkbox lines and mark the
     appropriate one(s) with `[x]`, deleting the rest:
     - reports something broken -> "Bug fix"
     - new capability or improvement -> "New feature"
     - changes existing behavior incompatibly -> "Breaking change"
     - docs-only -> "This change requires a documentation update"
     If the classification is not clear from the context, ask the user. Do not guess.
   - `<!tests>` -> a short description of *where* tests were added or touched
     (e.g. "investment usecases and REST handlers", "webapp AccountsPage component
     tests"), NOT the test code itself. If no tests were added, replace it with
     exactly `No tests added`.
3. **Push the branch** if it is not on the remote yet: `git push -u origin <branch>`.
4. **Create the PR** with the filled body written to a temp file (avoids shell quoting
   issues with the markdown):

   ```bash
   gh pr create --base main --head <branch> --title "<title>" --body-file /tmp/pr-body.md
   ```

   Title: the first commit subject, or the linked issue title. Ask the user if unclear.
5. **Report back**: PR number + URL.

## Rules

- Never guess the issue number, the type of change, or the test summary — ask first, unless the user provides the PR number. Like the PR number is X.
- Base is always `main`; head is the current feature branch.
