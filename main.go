package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "sync"
)

type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

type Store struct {
    mu     sync.RWMutex
    nextID int
    items  map[int]Todo
}

func NewStore() *Store {
    return &Store{nextID: 1, items: make(map[int]Todo)}
}

func (s *Store) listTodos(w http.ResponseWriter, r *http.Request) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    out := make([]Todo, 0, len(s.items))
    for _, t := range s.items {
        out = append(out, t)
    }
    writeJSON(w, http.StatusOK, out)
}

func (s *Store) createTodo(w http.ResponseWriter, r *http.Request) {
    var in struct {
        Title string `json:"title"`
        Done  bool   `json:"done"`
    }
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON body")
        return
    }
    if strings.TrimSpace(in.Title) == "" {
        writeError(w, http.StatusBadRequest, "title is required")
        return
    }
    s.mu.Lock()
    id := s.nextID
    s.nextID++
    todo := Todo{ID: id, Title: in.Title, Done: in.Done}
    s.items[id] = todo
    s.mu.Unlock()
    writeJSON(w, http.StatusCreated, todo)
}

func (s *Store) getTodo(w http.ResponseWriter, _ *http.Request, id int) {
    s.mu.RLock()
    todo, ok := s.items[id]
    s.mu.RUnlock()
    if !ok {
        writeError(w, http.StatusNotFound, "todo not found")
        return
    }
    writeJSON(w, http.StatusOK, todo)
}

func (s *Store) updateTodo(w http.ResponseWriter, r *http.Request, id int) {
    var in struct {
        Title *string `json:"title"`
        Done  *bool   `json:"done"`
    }
    if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON body")
        return
    }
    s.mu.Lock()
    todo, ok := s.items[id]
    if !ok {
        s.mu.Unlock()
        writeError(w, http.StatusNotFound, "todo not found")
        return
    }
    if in.Title != nil {
        title := strings.TrimSpace(*in.Title)
        if title == "" {
            s.mu.Unlock()
            writeError(w, http.StatusBadRequest, "title cannot be empty")
            return
        }
        todo.Title = title
    }
    if in.Done != nil {
        todo.Done = *in.Done
    }
    s.items[id] = todo
    s.mu.Unlock()
    writeJSON(w, http.StatusOK, todo)
}

func (s *Store) deleteTodo(w http.ResponseWriter, _ *http.Request, id int) {
    s.mu.Lock()
    if _, ok := s.items[id]; !ok {
        s.mu.Unlock()
        writeError(w, http.StatusNotFound, "todo not found")
        return
    }
    delete(s.items, id)
    s.mu.Unlock()
    w.WriteHeader(http.StatusNoContent)
}

func (s *Store) todosHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        s.listTodos(w, r)
    case http.MethodPost:
        s.createTodo(w, r)
    default:
        writeError(w, http.StatusMethodNotAllowed, "method not allowed")
    }
}

func (s *Store) todoItemHandler(w http.ResponseWriter, r *http.Request) {
    // Expect path: /todos/{id}
    idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        writeError(w, http.StatusBadRequest, "invalid id")
        return
    }
    switch r.Method {
    case http.MethodGet:
        s.getTodo(w, r, id)
    case http.MethodPatch, http.MethodPut:
        s.updateTodo(w, r, id)
    case http.MethodDelete:
        s.deleteTodo(w, r, id)
    default:
        writeError(w, http.StatusMethodNotAllowed, "method not allowed")
    }
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
    writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if strings.TrimSpace(name) == "" {
        name = "world"
    }
    writeJSON(w, http.StatusOK, map[string]string{"message": "Hello, " + name + "!"})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    if r.Method != http.MethodGet {
        writeError(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }
    writeJSON(w, http.StatusOK, map[string]string{"message": "welcome to todos simple API"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(v); err != nil {
        log.Printf("encode response: %v", err)
    }
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]string{"error": msg})
}

func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

// cors wraps a handler with simple CORS support.
// Allowed origins can be configured via env `CORS_ALLOWED_ORIGINS` (comma-separated, or "*").
func cors(next http.Handler) http.Handler {
    // Load allowed origins from env once.
    allowed := strings.Split(strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")), ",")
    useWildcard := false
    list := make([]string, 0, len(allowed))
    for _, o := range allowed {
        o = strings.TrimSpace(o)
        if o == "*" || o == "" {
            useWildcard = true
            continue
        }
        list = append(list, o)
    }

    allowOrigin := func(origin string) string {
        if useWildcard {
            return "*"
        }
        for _, o := range list {
            if strings.EqualFold(o, origin) {
                return o
            }
        }
        return ""
    }

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        allowedOrigin := allowOrigin(origin)
        if allowedOrigin != "" && origin != "" {
            w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
        }
        // Manage vary for caches
        w.Header().Add("Vary", "Origin")
        w.Header().Add("Vary", "Access-Control-Request-Method")
        w.Header().Add("Vary", "Access-Control-Request-Headers")

        if r.Method == http.MethodOptions {
            // Preflight
            reqHeaders := r.Header.Get("Access-Control-Request-Headers")
            if reqHeaders == "" {
                reqHeaders = "Content-Type, Authorization"
            }
            w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,PUT,DELETE,OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
            // If wildcard not used and origin matched, reflect it
            if allowedOrigin == "" && origin != "" {
                // No match -> block preflight
                w.WriteHeader(http.StatusNoContent)
                return
            }
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main() {
    store := NewStore()

    mux := http.NewServeMux()
    mux.HandleFunc("/", rootHandler)
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/hello", helloHandler)
    mux.HandleFunc("/todos", store.todosHandler)
    mux.HandleFunc("/todos/", store.todoItemHandler)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    addr := ":" + port
    log.Printf("listening on %s", addr)
    // Order: cors -> logging -> mux, so CORS headers apply to all responses.
    if err := http.ListenAndServe(addr, cors(logging(mux))); err != nil {
        log.Fatal(err)
    }
}
