# Guide: Refactoring Notifications to Database Triggers

## Step-by-Step Refactoring Process

### Step 1: Identify the Notification

Find where the notification is currently sent in application code.

**Example: Entry creation**
```go
// File: internal/app/mindwell-server/entries/entries.go
func postEntry(...) {
    // ... insert entry ...
    srv.Ntf.SendNewEntry(tx, entry)  // <-- This is what we're replacing
}
```

**Find all calls:**
```bash
grep -rn "SendNewEntry" internal/ --include="*.go"
```

### Step 2: Analyze the Notification Method

Understand what data the notification method needs.

**Example:**
```go
func (ntf *CompositeNotifier) SendNewEntry(tx *AutoTx, entry *models.Entry) {
    // Needs: entry.ID, entry.UserID, entry.TlogID, entry.Title
    // Queries for: watchers, user settings, privacy rules
    // Creates: notification records for each watcher
}
```

**Key questions:**
- What fields from the model are required?
- What additional data is queried?
- Who receives notifications? (complex logic = harder to move to triggers)
- Are there privacy/permission checks? (these stay in Go code)

### Step 3: Create Database Trigger Functions

Add trigger functions to `scripts/mindwell.sql` and `scripts/update.sql`.

**Location in schema:**
Find related triggers for the same table and add nearby. Use grep:
```bash
grep -n "CREATE TRIGGER.*entries" scripts/mindwell.sql
```

**Template:**
```sql
-- Trigger function for new entries
CREATE OR REPLACE FUNCTION mindwell.notify_new_entry() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('new_entry', json_build_object(
        'id', NEW.id,
        'user_id', NEW.user_id,
        'tlog_id', NEW.tlog_id,
        'author_id', NEW.author_id
        -- Include any fields needed by notification handler
    )::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for INSERT
CREATE TRIGGER entry_insert_trigger
    AFTER INSERT ON mindwell.entries
    FOR EACH ROW
    EXECUTE PROCEDURE mindwell.notify_new_entry();

-- Trigger for UPDATE (optional - only if notifications needed on update)
CREATE TRIGGER entry_update_trigger
    AFTER UPDATE ON mindwell.entries
    FOR EACH ROW
    WHEN (OLD.title IS DISTINCT FROM NEW.title)  -- Condition: only on changes
    EXECUTE PROCEDURE mindwell.notify_update_entry();

-- Trigger for DELETE (optional)
CREATE TRIGGER entry_delete_trigger
    AFTER DELETE ON mindwell.entries
    FOR EACH ROW
    EXECUTE PROCEDURE mindwell.notify_remove_entry();
```

**Important notes:**
- Use `AFTER` triggers (not `BEFORE`) so the data is committed
- Use `mindwell.` schema prefix to match existing triggers
- For UPDATE triggers, use `WHEN` clause to only fire on relevant changes
- For DELETE triggers, use `OLD` row data
- `pg_notify()` queues notification, sends AFTER transaction commits

### Step 4: Add PubSub Subscription

In `lib/notifications/composite.go`, subscribe to the new channel.

**Find subscription location:**
```go
func NewCompositeNotifier(srv ServerDependencies, ps *pubsub.PubSub) *CompositeNotifier {
    // ...
    ps.Subscribe("moved_entries", ntf.notifyMovedEntry)
    ps.Subscribe("user_badges", ntf.notifyBadge)
    ps.Subscribe("new_comment", ntf.notifyNewComment)

    // ADD NEW SUBSCRIPTIONS HERE:
    ps.Subscribe("new_entry", ntf.notifyNewEntry)  // <-- Add this

    return ntf
}
```

### Step 5: Create Notification Handler

Add handler function at the end of `lib/notifications/composite.go`.

**Template:**
```go
func (ntf *CompositeNotifier) notifyNewEntry(entryData []byte) {
    // Step 1: Unmarshal JSON payload from trigger
    var entryInfo struct {
        ID       int64 `json:"id"`
        UserID   int64 `json:"user_id"`
        TlogID   int64 `json:"tlog_id"`
        AuthorID int64 `json:"author_id"`
    }
    err := json.Unmarshal(entryData, &entryInfo)
    if err != nil {
        ntf.srv.LogSystem().Error(err.Error())
        return
    }

    // Step 2: Create new transaction for notification processing
    tx := database.NewAutoTx(ntf.srv.GetDB())
    defer tx.Finish()

    // Step 3: Load additional data needed by notification method
    const q = `
        SELECT title, content, privacy
        FROM entries
        WHERE id = $1
    `

    entry := &models.Entry{
        ID:       entryInfo.ID,
        User:     &models.User{ID: entryInfo.UserID},
        Tlog:     &models.Tlog{ID: entryInfo.TlogID},
        Author:   &models.User{ID: entryInfo.AuthorID},
    }

    tx.Query(q, entryInfo.ID).Scan(&entry.Title, &entry.Content, &entry.Privacy)

    // Step 4: Handle errors
    if tx.Error() != nil {
        return
    }

    // Step 5: Call the notification method
    ntf.sendNewEntry(tx, entry)
}
```

**Key points:**
- Use `defer tx.Finish()` for automatic commit/rollback
- Populate only the fields actually needed by the notification method
- Initialize nested structs (User, Tlog, etc.) to avoid nil pointer panics

### Step 6: Unexport the Notification Method

Make the notification method private since it's no longer called by application code.

**Before:**
```go
func (ntf *CompositeNotifier) SendNewEntry(tx *AutoTx, entry *Entry) {
    // ...
}
```

**After:**
```go
func (ntf *CompositeNotifier) sendNewEntry(tx *AutoTx, entry *Entry) {
    // ...
}
```

**Update all internal calls:**
```bash
# Find internal calls within composite.go
grep -n "SendNewEntry" lib/notifications/composite.go

# Change to lowercase in the same file
# ntf.SendNewEntry(tx, entry) → ntf.sendNewEntry(tx, entry)
```

### Step 7: Remove Application-Level Calls

Remove the synchronous notification calls from application code.

**Before:**
```go
func postEntry(srv *server.MindwellServer, tx *AutoTx, entry *Entry) {
    tx.Exec("INSERT INTO entries (...) VALUES (...)")
    // ...
    srv.Ntf.SendNewEntry(tx, entry)  // <-- Remove this
}
```

**After:**
```go
func postEntry(srv *server.MindwellServer, tx *AutoTx, entry *Entry) {
    tx.Exec("INSERT INTO entries (...) VALUES (...)")
    // Notification will be sent by trigger
}
```

**Find all calls to remove:**
```bash
grep -rn "srv.Ntf.SendNewEntry" internal/ --include="*.go"
```

### Step 8: Test and Verify

**Compile check:**
```bash
go build ./cmd/mindwell-server/
go build ./cmd/mindwell-images-server/
go build ./cmd/mindwell-helper/
```

**Database migration:**
```bash
# Apply triggers to development database
psql -d mindwell -f scripts/update.sql
```

**Run tests:**
```bash
# Run all tests
go test ./test/ --failfast

# Run specific notification tests
go test ./test/ -run TestNotification --failfast -v
```

**Manual testing:**
1. Create the resource (entry, comment, etc.)
2. Check that notifications are created
3. Verify email/Telegram/Centrifugo notifications sent
4. Test update operations
5. Test delete operations
6. Test privacy rules still work
