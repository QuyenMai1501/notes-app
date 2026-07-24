INSERT INTO notes (id, title, content, tags, created_at, updated_at) VALUES
(
  gen_random_uuid(),
  'Getting Started with Go Gin',
  'Gin is a lightweight web framework for Go that provides a martini-like API with better performance. It includes routing, middleware support, JSON validation, and error management. Perfect for building REST APIs quickly.',
  ARRAY['go', 'tutorial'],
  NOW() - INTERVAL '14 days',
  NOW() - INTERVAL '10 days'
),
(
  gen_random_uuid(),
  'React Hooks Best Practices',
  'Use useReducer for complex state logic instead of multiple useState calls. Keep effects focused on synchronization with external systems. Avoid calling setState synchronously inside effects — prefer async callbacks or useReducer dispatch.',
  ARRAY['react', 'tutorial'],
  NOW() - INTERVAL '12 days',
  NOW() - INTERVAL '8 days'
),
(
  gen_random_uuid(),
  'Docker Compose for PostgreSQL 17',
  'A minimal docker-compose.yml to run PostgreSQL 17 Alpine: define POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB environment variables, map port 5432, and mount a named volume for data persistence. Use pgdata volume to survive container restarts.',
  ARRAY['docker', 'database'],
  NOW() - INTERVAL '10 days',
  NOW() - INTERVAL '7 days'
),
(
  gen_random_uuid(),
  'CSS Modules in React',
  'CSS Modules scope styles locally by default. Create a .module.css file alongside your component, import it as styles, and use styles.className in JSX. No class name collisions, no global leakage. Works out of the box with Vite.',
  ARRAY['react', 'css'],
  NOW() - INTERVAL '9 days',
  NOW() - INTERVAL '6 days'
),
(
  gen_random_uuid(),
  'PostgreSQL Array Type Tricks',
  'PostgreSQL supports TEXT[] columns natively. Use ANY() to check membership, unnest() to expand arrays, and array_append() to modify. Index with GIN for efficient searches on array elements.',
  ARRAY['database', 'tutorial'],
  NOW() - INTERVAL '8 days',
  NOW() - INTERVAL '5 days'
),
(
  gen_random_uuid(),
  'Graceful Shutdown in Go',
  'Use signal.Notify with SIGINT/SIGTERM and an http.Server.Shutdown with a timeout context. Start the server in a goroutine, block on the signal channel, then initiate graceful shutdown. This ensures in-flight requests complete before the process exits.',
  ARRAY['go', 'devops'],
  NOW() - INTERVAL '7 days',
  NOW() - INTERVAL '4 days'
),
(
  gen_random_uuid(),
  'TypeScript 6 — erasableSyntaxOnly',
  'TypeScript 6 introduced erasableSyntaxOnly which bans enums, namespaces, and parameter properties. These features require runtime emit that conflicts with the erasable-type-only philosophy. Use union types and const objects instead of enums.',
  ARRAY['react', 'typescript'],
  NOW() - INTERVAL '6 days',
  NOW() - INTERVAL '3 days'
),
(
  gen_random_uuid(),
  'Vite 8 Dev Server Setup',
  'Vite 8 uses rolldown for bundling. Configure dev server port via server.port in vite.config.ts. The server: { port: 3000 } setting makes the dev server listen on localhost:3000 instead of the default 5173.',
  ARRAY['react', 'tutorial'],
  NOW() - INTERVAL '5 days',
  NOW() - INTERVAL '2 days'
),
(
  gen_random_uuid(),
  'pgx Connection Pooling',
  'Use pgxpool.NewWithConfig for connection pooling. Set MaxConns (defaults to number of CPUs) and MinConns for baseline connections. Ping on startup to validate credentials. Use a single pool instance across all handlers via dependency injection.',
  ARRAY['go', 'database'],
  NOW() - INTERVAL '4 days',
  NOW() - INTERVAL '1 days'
),
(
  gen_random_uuid(),
  'CORS Middleware Without a Library',
  'Instead of importing gin-contrib/cors, you can write a 10-line middleware that sets Access-Control-Allow-Origin, Allow-Methods, Allow-Headers headers and handles OPTIONS preflight with a 204 response. Lighter and more transparent.',
  ARRAY['go', 'tutorial'],
  NOW() - INTERVAL '3 days',
  NOW()
),
(
  gen_random_uuid(),
  'Server-Side Search with ILIKE',
  'PostgreSQL ILIKE provides case-insensitive pattern matching. Use % wildcards around the search term for substring matching. Combine with GIN index on to_tsvector for full-text search on larger datasets. Keep the search parameterized to prevent SQL injection.',
  ARRAY['database', 'tutorial'],
  NOW() - INTERVAL '2 days',
  NOW()
),
(
  gen_random_uuid(),
  'Monitoring API Health Endpoints',
  'A /api/health endpoint returning {"status":"ok"} is the simplest way to integrate with monitoring tools and load balancers. The handler should validate database connectivity and return 5xx if the DB is unreachable, so alerts trigger on real outages.',
  ARRAY['devops', 'go'],
  NOW() - INTERVAL '1 days',
  NOW()
),
(
  gen_random_uuid(),
  'Go Module Layout for REST APIs',
  'Use cmd/server as the entrypoint and internal/ for private packages. Models, handlers, database, and router are separate packages under internal/. This prevents external imports and keeps the dependency graph clean. Migrations live alongside models.',
  ARRAY['go'],
  NOW() - INTERVAL '15 days',
  NOW() - INTERVAL '11 days'
),
(
  gen_random_uuid(),
  'UUID Primary Keys in PostgreSQL',
  'Use UUID primary keys instead of SERIAL for distributed systems. gen_random_uuid() generates v4 UUIDs with no ordering guarantee. Add a created_at default of NOW() for chronological queries. UUIDs are 16 bytes — slightly larger than integer but worth the distribution benefits.',
  ARRAY['database'],
  NOW() - INTERVAL '13 days',
  NOW() - INTERVAL '9 days'
),
(
  gen_random_uuid(),
  'TypeScript verbatimModuleSyntax',
  'With verbatimModuleSyntax: true, TypeScript preserves import statements as-is and does not elide type-only imports. You must explicitly use import type for any import that is only used as a type. This enables bundlers to tree-shake more effectively.',
  ARRAY['typescript', 'react'],
  NOW() - INTERVAL '11 days',
  NOW() - INTERVAL '7 days'
),
(
  gen_random_uuid(),
  'Docker Volume Persistence',
  'Named volumes (volumes: pgdata:) survive container recreation. Bind mounts tie data to host paths. For production PostgreSQL, always use named volumes or a remote volume driver. Anonymous volumes are discarded when the container is removed with --rm.',
  ARRAY['docker', 'devops'],
  NOW() - INTERVAL '10 days',
  NOW() - INTERVAL '6 days'
),
(
  gen_random_uuid(),
  'Go Gin Request Validation',
  'Use ShouldBindJSON with struct binding tags (binding:"required") for request validation. Gin returns 400 automatically on validation failure. For custom validation logic, validate inside the handler and return structured errors.',
  ARRAY['go', 'tutorial'],
  NOW() - INTERVAL '8 days',
  NOW() - INTERVAL '4 days'
),
(
  gen_random_uuid(),
  'React Fetch Pattern with useReducer',
  'A fetchReducer with start/success/error action types avoids the react-hooks/set-state-in-effect lint error. Dispatch start synchronously in the effect, then dispatch success or error in async callbacks. Use an ignore flag to prevent state updates after unmount.',
  ARRAY['react', 'tutorial'],
  NOW() - INTERVAL '6 days',
  NOW() - INTERVAL '2 days'
),
(
  gen_random_uuid(),
  'Environment Variables in Go',
  'Use os.Getenv to read configuration like DATABASE_URL and PORT. Provide sensible defaults for local development. Never hardcode credentials. For production, set env vars in the container runtime or orchestration layer.',
  ARRAY['go', 'devops'],
  NOW() - INTERVAL '4 days',
  NOW() - INTERVAL '1 days'
),
(
  gen_random_uuid(),
  'NPM Workspaces vs Multi-Repo',
  'This project uses separate package.json files in frontend/ and server/ without npm workspaces. Each directory is independent. Run npm commands from the frontend/ directory and go commands from the server/ directory.',
  ARRAY['devops'],
  NOW() - INTERVAL '2 days',
  NOW()
);
