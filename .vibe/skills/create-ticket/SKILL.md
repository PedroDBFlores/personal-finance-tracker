---
name: create-ticket
description: Load when the user asks to create a new ticket/issue in the "Personal Finance Tracker" GitHub Project (e.g. "create a ticket for ..."). Classifies it as Task/Bug/Feature, labels bug/enhancement, adds it to the project, and optionally creates a PFTXXXXXX branch.
---

# Create GitHub Project ticket

Use this skill whenever the user requests a new ticket in GitHub Projects.

Targets (hardcoded for this repo):

- Repository: `PedroDBFlores/personal-finance-tracker` (node ID `R_kgDOVCuRiQ`)
- Project: **"Personal Finance Tracker"**, number 1, owner `PedroDBFlores`
  (node ID `PVT_kwHOACWgtM4BslwW`)
- Requires `gh` authenticated with `repo` and `project` scopes
  (`gh auth status`); if missing, tell the user to run
  `gh auth login -s repo -s read:project -s project`.

## Steps

1. **Classify from the request text.** Map to the issue type of the project item:
   - Reports something broken / misbehaving → type **Bug**, label **bug**
   - New capability or improvement → type **Feature**, label **enhancement**
   - Chore, refactor, docs, small piece of work → type **Task**, no label
   - **If the type or label is not clear from the context, ask the user before
     creating anything.** Do not guess.
2. **Ask about a branch:** "Create a branch for this ticket now?" (yes/no).
   If yes, the branch name is `PFTXXXXXX` where `XXXXXX` is the issue number
   zero-padded to 6 digits (issue #3 → `PFT000003`), created from `main`.
   No dash — branch names with dashes have caused issues, use the bare form.
3. **Create the issue** in the repo (title + body summarizing the request),
   then set its issue type and label, then **add it to the project** with
   Status = **Todo**.
4. **Report back**: issue number + URL, project item, and branch name if created.

## Implementation (gh / GraphQL)

Issue type IDs (this repo): Task `IT_kwDOACWgtM4Ci6co`, Bug `IT_kwDOACWgtM4Ci6cp`,
Feature `IT_kwDOACWgtM4Ci6cq`. Re-fetch via `repository.issueTypes` if unsure.

```bash
# 1. create the issue, capture its node id
ISSUE=$(gh api graphql -f query='mutation($repo:ID!,$title:String!,$body:String!){ createIssue(input:{repositoryId:$repo,title:$title,body:$body}){ issue{ id number url } } }' \
  -f repo=R_kgDOVCuRiQ -f title="..." -f body="...")

# 2. set issue type + add label (addLabelsToLabelable needs the issue node id)
# 3. add to project with Status=Todo:
gh api graphql -f query='mutation($project:ID!,$content:ID!){ addProjectV2ItemById(input:{projectId:$project, contentId:$content}){ item{ id } } }' \
  -f project=PVT_kwHOACWgtM4BslwW -f content=<issue node id>
# then updateProjectV2ItemFieldValue with Status field, option "Todo"

# 4. branch, only if the user said yes:
git checkout -b PFT000003 main   # zero-pad issue number to 6 digits, no dash
```

## Rules

- Never create a ticket with an unclear type/label — ask first.
- Always ask the branch question explicitly; do not create a branch unprompted.
- Branch format is exactly `PFTXXXXXX` (PFT + 6 zero-padded digits, no dash).
