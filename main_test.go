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
      messageBackend[shape=rectanble]
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

	expectedMermaid := `graph TD;
    n3("messageBackend");
    n1("messageClient");
    n4("messageDB");
    n8("messageNotifier");
    n2("messageQueue");
    n7("pushNotifications");
    n6("userDB");
    n5("userService");
    n3-->|" "|n5;
    n3-->|" "|n4;
    n1-->|"sendMessage"|n3;
    n4-->|"&#34;change data capture&#34;"|n2;
    n8-->|" "|n7;
    n2-->|" "|n8;
    n5-->|" "|n6;`

	actualMermaid, err := СonvertDOTToMermaid(dotInput)

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
	return strings.Join(strings.Fields(input), " ")
}
