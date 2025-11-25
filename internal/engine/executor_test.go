package engine_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/grantcarthew/start/internal/domain"
	"github.com/grantcarthew/start/internal/engine"
	"github.com/grantcarthew/start/test/assert"
	"github.com/grantcarthew/start/test/mocks"
)

func TestExecutor_Execute_Success(t *testing.T) {
	mockRunner := &mocks.MockRunner{
		Outputs: map[string]string{
			"smith --model test-model 'hello world'": "Response from smith",
		},
	}

	resolver := engine.NewPlaceholderResolver()
	executor := engine.NewExecutor(mockRunner, resolver)

	agent := domain.Agent{
		Name:    "smith",
		Bin:     "smith",
		Command: "{bin} --model {model} '{prompt}'",
	}

	ctx := context.Background()
	err := executor.Execute(ctx, agent, "test-model", "hello world")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockRunner.CalledWith))
	assert.Equal(t, "bash", mockRunner.CalledWith[0].Shell)
	assert.Equal(t, "smith --model test-model 'hello world'", mockRunner.CalledWith[0].Command)
	assert.Equal(t, 2*time.Minute, mockRunner.CalledWith[0].Timeout)
}

func TestExecutor_Execute_PlaceholderResolution(t *testing.T) {
	mockRunner := &mocks.MockRunner{}
	resolver := engine.NewPlaceholderResolver()
	executor := engine.NewExecutor(mockRunner, resolver)

	agent := domain.Agent{
		Name:    "test-agent",
		Bin:     "test-bin",
		Command: "{bin} --model {model} --prompt '{prompt}'",
	}

	ctx := context.Background()
	err := executor.Execute(ctx, agent, "my-model", "test prompt")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockRunner.CalledWith))
	assert.Equal(t, "test-bin --model my-model --prompt 'test prompt'", mockRunner.CalledWith[0].Command)
}

func TestExecutor_Execute_Error(t *testing.T) {
	mockRunner := &mocks.MockRunner{
		ShouldError:  true,
		ErrorMessage: "command failed",
	}

	resolver := engine.NewPlaceholderResolver()
	executor := engine.NewExecutor(mockRunner, resolver)

	agent := domain.Agent{
		Name:    "smith",
		Bin:     "smith",
		Command: "{bin} --model {model} '{prompt}'",
	}

	ctx := context.Background()
	err := executor.Execute(ctx, agent, "test-model", "hello")

	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "command failed"), "Error should contain 'command failed'")
}

func TestExecutor_Execute_DatePlaceholder(t *testing.T) {
	mockRunner := &mocks.MockRunner{}
	resolver := engine.NewPlaceholderResolver()
	executor := engine.NewExecutor(mockRunner, resolver)

	agent := domain.Agent{
		Name:    "smith",
		Bin:     "smith",
		Command: "{bin} --date {date} '{prompt}'",
	}

	ctx := context.Background()
	err := executor.Execute(ctx, agent, "test-model", "hello")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(mockRunner.CalledWith))
	assert.Contains(t, mockRunner.CalledWith[0].Command, "smith --date ")
	assert.NotContains(t, mockRunner.CalledWith[0].Command, "{date}")
}
