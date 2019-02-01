package cloud

import (
	"bytes"
	"context"
	"database/sql"
	"github.com/gofunct/gotransport/html"
	"github.com/gorilla/mux"
	"gocloud.dev/blob"
	"gocloud.dev/health"
	"gocloud.dev/runtimevar"
	"gocloud.dev/server"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Runtime is the main server struct for Guestbook. It contains the state of
// the most recently read message of the day.
type Runtime struct {
	runvar    string
	env       Env
	bannerSrc string
	messages  []greeting
	srv       *server.Server
	db        *sql.DB
	bucket    *blob.Bucket
	checks    []health.Checker
	mu        sync.RWMutex
}

// newApplication creates a new Runtime struct based on the backends and the message
// of the day variable.
func NewRuntime(rvar string, srv *server.Server, db *sql.DB, bucket *blob.Bucket, runvar *runtimevar.Variable) *Runtime {
	hthz, _ := AppHealthChecks(db)

	app := &Runtime{
		runvar:   rvar,
		env:      SelectEnv(),
		messages: []greeting{},
		srv:      srv,
		db:       db,
		bucket:   bucket,
		checks:   hthz,
		mu:       sync.RWMutex{},
	}

	go app.WatchRuntimeVar(runvar)
	return app
}

// watchMOTDVar listens for changes in v and updates the app's message of the
// day. It is run in a separate goroutine.
func (r *Runtime) WatchRuntimeVar(v *runtimevar.Variable) {
	ctx := context.Background()
	for {
		snap, err := v.Watch(ctx)
		if err != nil {
			log.Printf("watch runvar variable: %v", err)
			continue
		}
		log.Println("updated runvar to", snap.Value)
		r.mu.Lock()
		r.runvar = snap.Value.(string)
		r.mu.Unlock()
	}
}

// index serves the server's landing page. It lists the 100 most recent
// greetings, shows a cloud environment banner, and displays the message of the
// day.

func GuestBookIndex(run *Runtime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		run.mu.RLock()
		run.mu.RUnlock()

		const query = "SELECT content FROM (SELECT content, post_date FROM greetings ORDER BY post_date DESC LIMIT 100) AS recent_greetings ORDER BY post_date ASC;"
		q, err := run.db.QueryContext(r.Context(), query)
		if err != nil {
			log.Println("main page SQL error:", err)
			http.Error(w, "could not load greetings", http.StatusInternalServerError)
			return
		}
		defer q.Close()
		for q.Next() {
			var g greeting
			if err := q.Scan(&g.Content); err != nil {
				log.Println("main page SQL error:", err)
				http.Error(w, "could not load greetings", http.StatusInternalServerError)
				return
			}
			run.messages = append(run.messages, g)
		}
		if err := q.Err(); err != nil {
			log.Println("main page SQL error:", err)
			http.Error(w, "could not load messages", http.StatusInternalServerError)
			return
		}
		buf := new(bytes.Buffer)
		if err := html.GuestBookTmpl.Execute(buf, run); err != nil {
			log.Println("template error:", err)
			http.Error(w, "could not render page", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
		if _, err := w.Write(buf.Bytes()); err != nil {
			log.Println("writing response:", err)
		}
	}
}

type greeting struct {
	Content string
}

// sign is a form action handler for adding a greeting.
func Sign(run *Runtime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("Allow", "POST")
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}
		content := r.FormValue("content")
		if content == "" {
			http.Error(w, "content must not be empty", http.StatusBadRequest)
			return
		}
		const sqlStmt = "INSERT INTO greetings (content) VALUES (?);"
		_, err := run.db.ExecContext(r.Context(), sqlStmt, content)
		if err != nil {
			log.Println("sign SQL error:", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// serveBlob handles a request for a static asset by retrieving it from a bucket.
func (app *Runtime) ServeBlob(run *Runtime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]
		blobRead, err := run.bucket.NewReader(r.Context(), key, nil)
		if err != nil {
			// TODO(light): Distinguish 404.
			// https://github.com/google/go-cloud/issues/2
			log.Println("serve blob:", err)
			http.Error(w, "blob read error", http.StatusInternalServerError)
			return
		}
		// TODO(light): Get content type from blob storage.
		// https://github.com/google/go-cloud/issues/9
		switch {
		case strings.HasSuffix(key, ".png"):
			w.Header().Set("Content-Type", "image/png")
		case strings.HasSuffix(key, ".jpg"):
			w.Header().Set("Content-Type", "image/jpeg")
		default:
			w.Header().Set("Content-Type", "Runtime/octet-stream")
		}
		w.Header().Set("Content-Length", strconv.FormatInt(blobRead.Size(), 10))
		if _, err = io.Copy(w, blobRead); err != nil {
			log.Println("Copying blob:", err)
		}
	}
}
