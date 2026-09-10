Task tracker. Educational go project.

It can work with a JSON file or a database. See env.example.

### Init

go build ./

### Adding a new task
task-tracker add "Buy groceries"

### Updating and deleting tasks
task-tracker update 1 "Buy groceries and cook dinner"
task-tracker delete 1

### Marking a task as in progress or done
task-tracker mark-in-progress 1
task-tracker mark-done 1

### Listing all tasks
task-tracker list

### Listing tasks by status
task-tracker list done
task-tracker list todo
task-tracker list in-progress
