package engine

import (
	"fmt"

	"github.com/grantcarthew/start/internal/domain"
)

// TaskResolver resolves task names and aliases
type TaskResolver struct{}

// NewTaskResolver creates a new task resolver
func NewTaskResolver() *TaskResolver {
	return &TaskResolver{}
}

// ResolveResult represents the result of task resolution
type ResolveResult struct {
	Task  domain.Task
	Found bool
}

// Resolve resolves a task name or alias
// Resolution order: local task name → local alias → global task name → global alias
// Returns the resolved task or error if not found
func (r *TaskResolver) Resolve(
	input string,
	localTasks map[string]domain.Task,
	globalTasks map[string]domain.Task,
) (domain.Task, error) {
	// 1. Check local task name (exact match)
	if task, exists := localTasks[input]; exists {
		task.Name = input
		return task, nil
	}

	// 2. Check local task alias
	for name, task := range localTasks {
		if task.Alias == input {
			task.Name = name
			return task, nil
		}
	}

	// 3. Check global task name (exact match)
	if task, exists := globalTasks[input]; exists {
		task.Name = input
		return task, nil
	}

	// 4. Check global task alias
	for name, task := range globalTasks {
		if task.Alias == input {
			task.Name = name
			return task, nil
		}
	}

	// Task not found
	return domain.Task{}, fmt.Errorf("task %q not found", input)
}

// ListAllTasks returns all tasks from local and global configs
// Local tasks override global tasks with the same name
func (r *TaskResolver) ListAllTasks(
	localTasks map[string]domain.Task,
	globalTasks map[string]domain.Task,
) map[string]domain.Task {
	result := make(map[string]domain.Task)

	// Start with global tasks
	for name, task := range globalTasks {
		task.Name = name
		result[name] = task
	}

	// Override with local tasks
	for name, task := range localTasks {
		task.Name = name
		result[name] = task
	}

	return result
}
