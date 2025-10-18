# Run templ generation in watch mode
templ:
    templ generate --watch --proxy="http://localhost:8000" --open-browser=false

# Run air for Go hot reload
server:
    air \
    --build.cmd "go build -o tmp/bin/main ./cmd/webapp/main.go" \
    --build.bin "tmp/bin/main" \
    --build.delay "100" \
    --build.exclude_dir "node_modules" \
    --build.include_ext "go" \
    --build.stop_on_error "false" \
    --misc.clean_on_exit true

# Watch Tailwind CSS changes
tailwind:
    tailwindcss -i ./web/assets/css/input.css -o ./web/static/css/output.css --watch

#Local turso db
local_db:
    turso dev --db-file local.db

# Start development server with all watchers
dev:
    just local_db &
    just templ &
    just server
