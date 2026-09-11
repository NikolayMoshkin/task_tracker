Task tracker. Educational go project.

It can work with a JSON file or a database. The application runs in CLI or HTTP mode. See `env.example`.

### Configuration

Copy `env.example` to `.env` and adjust values:

| Variable     | Values       | Description                                      |
|--------------|--------------|--------------------------------------------------|
| `MODE`       | `cli`, `http` | How to run the application                      |
| `STORE`      | `file`, `db` | Where tasks are stored                           |
| `DB_CONNECT` | connection string | PostgreSQL DSN (required when `STORE=db`) |

### Init

```bash
go build -o task-tracker ./cmd
cp env.example .env
```

### CLI mode

Set `MODE=cli` in `.env`, then run commands:

```bash
# Adding a new task
task-tracker add "Buy groceries"

# Updating and deleting tasks
task-tracker update 1 "Buy groceries and cook dinner"
task-tracker delete 1

# Marking a task as in progress or done
task-tracker mark-in-progress 1
task-tracker mark-done 1

# Listing all tasks
task-tracker list

# Listing tasks by status
task-tracker list done
task-tracker list todo
task-tracker list in-progress
```

### HTTP mode

Set `MODE=http` in `.env` and start the server:

```bash
task-tracker
```

The server listens on `:8080`.

#### Endpoints

| Method   | Path                      | Description                    |
|----------|---------------------------|--------------------------------|
| `GET`    | `/health`                 | Health check                   |
| `GET`    | `/list`                   | List all tasks                 |
| `GET`    | `/list?status={status}`   | List tasks by status           |
| `POST`   | `/add`                    | Create a task                  |
| `POST`   | `/update`                 | Update a task                  |
| `DELETE` | `/delete/{id}`            | Delete a task                  |
| `POST`   | `/mark-done/{id}`         | Mark a task as done            |
| `POST`   | `/mark-in-progress/{id}`  | Mark a task as in progress     |

Status values for `/list?status=`: `todo`, `in-progress`, `done`.

#### Examples

```bash
curl http://localhost:8080/health

curl http://localhost:8080/list
curl "http://localhost:8080/list?status=done"

curl -X POST http://localhost:8080/add \
  -H "Content-Type: application/json" \
  -d '{"name": "Buy groceries"}'

curl -X POST http://localhost:8080/update \
  -H "Content-Type: application/json" \
  -d '{"id": 1, "name": "Buy groceries and cook dinner"}'

curl -X DELETE http://localhost:8080/delete/1

curl -X POST http://localhost:8080/mark-in-progress/1
curl -X POST http://localhost:8080/mark-done/1
```
