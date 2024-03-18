package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertDOTToMermaid(t *testing.T) {
	dotInput := `
    digraph MessageArchitecture {
      messageClient
      messageQueue[shape=rarrow]
      messageBackend[shape=rectangle]
      messageDB[shape=cylinder]
      userService[shape=rectangle]
      userDB[shape=cylinder]
      pushNotifications[shape=octagon]
      messageNotifier[shape=rectangle]

      messageClient -> messageBackend[label=sendMessage]
      messageBackend -> userService
      userService -> userDB
      messageBackend -> messageDB
      messageDB -> messageQueue[label="change data capture"]
      messageQueue -> messageNotifier
      messageNotifier -> pushNotifications
    }`

	// expectedMermaid := `graph TD
	// messageClient -->|sendMessage| messageBackend
	// messageBackend --> userService
	// userService --> userDB[(userDB)]
	// messageBackend --> messageDB
	// messageQueue>messageQueue]
	// messageDB[(messageDB)] -->|change data capture| messageQueue
	// messageQueue --> messageNotifier
	// messageNotifier --> pushNotifications{{pushNotifications}}`

	expectedMermaid := `graph TB
    messageClient(["messageClient"])
    messageQueue>"messageQueue"]
    messageBackend["messageBackend"]
    messageDB[("messageDB")]
    userService["userService"]
    userDB[("userDB")]
    pushNotifications{{"pushNotifications"}}
    messageNotifier["messageNotifier"]
    messageClient -->|"sendMessage"| messageBackend
    messageBackend --> userService
    userService --> userDB
    messageBackend --> messageDB
    messageDB -->|"change data capture"| messageQueue
    messageQueue --> messageNotifier
    messageNotifier --> pushNotifications`

	actualMermaid, err := ConvertDOTToMermaid(dotInput)

	// Используем require для контроля ошибок, так как если есть ошибка, нет смысла продолжать тест
	require.NoError(t, err)

	// Убираем лишние пробелы и переносы строк для более простого сравнения строк
	normalizedActual := normalizeWhitespace(actualMermaid)
	normalizedExpected := normalizeWhitespace(expectedMermaid)

	// Используем assert для сравнения ожидаемого и фактического результатов
	assert.Equal(t, normalizedExpected, normalizedActual, "The converted MermaidJS diagram does not match the expected output")
}

// normalizeWhitespace удаляет все пробелы и переносы строк из строки для упрощения сравнения.
func normalizeWhitespace(input string) string {
	return strings.Join(strings.Fields(input), "; ")
}
