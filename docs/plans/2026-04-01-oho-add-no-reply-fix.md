# oho add --no-reply Bug Fix Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Fix the `oho add --no-reply` command so that it properly verifies message delivery and provides clear feedback to users.

**Architecture:** The issue is that `--no-reply` mode sends to `/session/{id}/message` with `noReply=true`, but treats empty server response as success without verification. The fix adds explicit response validation and clear messaging about async message status.

**Tech Stack:** Go, Cobra CLI, HTTP Client

---

## Root Cause Summary

When `oho add "message" --no-reply` is executed:
1. Session is created successfully
2. Message is sent to `/session/{id}/message` with `noReply=true`
3. Server returns empty body `[]byte{}` (async processing)
4. Client returns `(messageID="", err=nil)` - **false success**
5. User sees "Message sent successfully" but message may not be queued

**The Fix:** Add verification for no-reply mode by using a different approach:
- Option A: Use `/session/{id}/prompt_async` endpoint (dedicated async endpoint)
- Option B: Verify message was stored by checking session messages list
- Option C: Return the session ID so user can verify manually

**Recommended Fix:** Option A - Use proper async endpoint

---

## Atomic Commit Strategy

| Commit | Description | Files |
|--------|-------------|-------|
| 1 | Add verification test for no-reply mode | `oho/cmd/add/add_test.go` |
| 2 | Fix sendMessage to use prompt_async for no-reply mode | `oho/cmd/add/add.go` |
| 3 | Update error message for no-reply mode to be clearer | `oho/cmd/add/add.go` |

---

## Task 1: Write Failing Test for No-Reply Mode Verification

**Files:**
- Modify: `oho/cmd/add/add_test.go`

**Step 1: Add test case for no-reply verification**

Add a new test case that verifies when `noReply=true`:
- The mock should be called with the correct endpoint
- The response handling should be correct

```go
// Add to TestSendMessage test cases:
{
    name:      "no-reply mode uses correct endpoint",
    sessionID: "ses_test123",
    message:   "Hello",
    noReply:   true,
    mockResp:  []byte{},  // Empty response for no-reply
    mockErr:   nil,
    wantErr:   false,
    wantMsgID: "",  // No message ID returned
    verifyPath: "/session/ses_test123/prompt_async",  // Should use async endpoint
},
```

Wait - the issue is the endpoint should be `/prompt_async` not `/message`. Let me verify current behavior first.

**Step 2: Run existing tests to understand current behavior**

Run: `cd /mnt/d/fe/opencode_cli && go test -v ./oho/cmd/add/... -run TestSendMessage`
Expected: All tests pass (baseline)

---

## Task 2: Investigate Correct API Behavior

**Files:**
- None (investigation)

**Step 1: Check OpenCode API documentation**

The API mapping shows:
- `oho message prompt-async` → `/session/:id/prompt_async` (POST)

But `oho add --no-reply` uses `/session/:id/message` with `noReply=true`.

**Step 2: Determine correct approach**

The issue is:
- `/session/{id}/message` with `noReply=true` should return a message ID (even if empty body initially)
- But server returns empty body, causing ambiguity

**Recommended fix:** When `noReply=true`, use `/session/{id}/prompt_async` endpoint instead, which is designed for async messaging.

---

## Task 3: Implement the Fix

**Files:**
- Modify: `oho/cmd/add/add.go:197-261` (sendMessage function)

**Step 1: Update sendMessage to use prompt_async for no-reply mode**

```go
// sendMessage sends a message to the session and returns the message ID
func sendMessage(c client.ClientInterface, ctx context.Context, sessionID, message, agent, model string, noReply bool, system string, tools, files []string) (string, error) {
    // Build message parts
    var parts []types.Part

    // Add text part
    text := message
    parts = append(parts, types.Part{
        Type: "text",
        Text: &text,
    })

    // Add file parts
    for _, filePath := range files {
        // ... file handling code (lines 209-233)
    }

    // Build message request
    msgReq := types.MessageRequest{
        Model:   convertModel(model),
        Agent:   agent,
        NoReply: noReply,
        System:  system,
        Tools:   tools,
        Parts:   parts,
    }

    // Determine endpoint based on noReply flag
    endpoint := fmt.Sprintf("/session/%s/message", sessionID)
    if noReply {
        endpoint = fmt.Sprintf("/session/%s/prompt_async", sessionID)
        // For async endpoint, always set noReply to false (server handles async internally)
        msgReq.NoReply = false
    }

    resp, err := c.Post(ctx, endpoint, msgReq)
    if err != nil {
        return "", fmt.Errorf("API request failed: %w", err)
    }

    // Handle empty response (no-reply mode)
    if len(resp) == 0 {
        // For no-reply mode, we can't verify message ID, but we can confirm it was accepted
        if noReply {
            // Return session ID as reference since message ID is not available
            return sessionID, nil
        }
        return "", nil
    }

    var result types.MessageWithParts
    if err := json.Unmarshal(resp, &result); err != nil {
        return "", fmt.Errorf("failed to parse response: %w", err)
    }

    return result.Info.ID, nil
}
```

**Step 2: Update output message in runAdd to clarify async behavior**

In `runAdd` function (around line 136-142):

```go
} else {
    fmt.Printf("Session created: %s\n", sessionID)
    if messageID != "" {
        fmt.Printf("Message sent: %s\n", messageID)
    } else if addNoReply {
        // For no-reply mode, message ID won't be returned, but session was created
        fmt.Printf("Message sent (async mode - session: %s)\n", sessionID)
    } else {
        fmt.Println("Message sent successfully")
    }
}
```

**Step 3: Run tests to verify fix**

Run: `cd /mnt/d/fe/opencode_cli && go test -v ./oho/cmd/add/... -run TestSendMessage`
Expected: Tests should pass (including new no-reply test)

---

## Task 4: Add Integration Test for No-Reply Mode

**Files:**
- Modify: `oho/cmd/add/add_test.go`

**Step 1: Add integration test case for no-reply with session verification**

```go
func TestNoReplyModeVerification(t *testing.T) {
    // Test that no-reply mode properly creates session and accepts message
    mock := &client.MockClient{
        PostWithQueryFunc: func(ctx context.Context, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
            return testutil.MockSessionResponse(), nil
        },
        PostFunc: func(ctx context.Context, path string, body interface{}) ([]byte, error) {
            // Verify the endpoint used
            if path != "/session/session1/prompt_async" {
                t.Errorf("Expected /session/session1/prompt_async, got %s", path)
            }
            return []byte{}, nil  // Empty response for async
        },
    }

    ctx := context.Background()

    // Create session
    sessionID, err := createSession(mock, ctx, "Test", "", "/test")
    if err != nil {
        t.Fatalf("Failed to create session: %v", err)
    }

    // Send message with noReply=true
    msgID, err := sendMessage(mock, ctx, sessionID, "Test message", "", "", true, "", nil, nil)
    if err != nil {
        t.Errorf("Unexpected error in no-reply mode: %v", err)
    }

    // For no-reply mode, msgID should be sessionID (as fallback reference)
    t.Logf("No-reply mode returned session ID as reference: %s", msgID)
}
```

---

## Task 5: Update Test Mock Client

**Files:**
- Modify: `oho/internal/client/client_mock.go`

**Step 1: Verify MockClient supports PostFunc properly**

Check that the mock can track which endpoint was called.

```go
// client_mock.go should have:
type MockClient struct {
    PostWithQueryFunc func(ctx context.Context, path string, queryParams map[string]string, body interface{}) ([]byte, error)
    PostFunc          func(ctx context.Context, path string, body interface{}) ([]byte, error)
    // Add for tracking:
    LastPostPath      string  // Track last POST path for verification
}
```

---

## Task 6: Run Full Test Suite

**Files:**
- None

**Step 1: Run all add command tests**

Run: `cd /mnt/d/fe/opencode_cli && go test -v ./oho/cmd/add/...`
Expected: All tests pass

**Step 2: Run all oho tests**

Run: `cd /mnt/d/fe/opencode_cli && go test ./oho/...`
Expected: All tests pass

**Step 3: Verify build**

Run: `cd /mnt/d/fe/opencode_cli/oho && go build -o oho_test ./cmd/oho/`
Expected: Build succeeds

---

## Task 7: Commit Changes

**Files:**
- `oho/cmd/add/add.go`
- `oho/cmd/add/add_test.go`

**Step 1: Stage and commit**

```bash
git add oho/cmd/add/add.go oho/cmd/add/add_test.go
git commit -m "fix: use /prompt_async endpoint for --no-reply mode

Previously, oho add --no-reply sent to /session/{id}/message with
noReply=true, which returned an empty response treated as success
without verification.

Now, --no-reply mode uses the dedicated /session/{id}/prompt_async
endpoint which is designed for async message handling.

Fixes: Session created but message not properly queued issue.
"
```

---

## Verification Checklist

- [ ] `go test ./oho/cmd/add/...` passes
- [ ] `go test ./oho/...` passes
- [ ] `go build ./cmd/oho/...` succeeds
- [ ] Manual test: `oho add "test" --no-reply` works correctly
- [ ] Output clearly indicates async message was accepted
