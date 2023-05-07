package extractors

import (
	"context"
	"fmt"
	"github.com/danielchalef/zep/internal"
	"github.com/danielchalef/zep/pkg/llms"
	"github.com/danielchalef/zep/pkg/models"
	"strings"
)

const SummaryMaxOutputTokens = 512

// Force compiler to validate that RedisMemoryStore implements the MemoryStore interface.
var _ models.Extractor = &SummaryExtractor{}

type SummaryExtractor struct {
	models.BaseExtractor
}

// Extract gets a list of messages created since the last SummaryPoint,
// determines if the message count exceeds the configured message window, and if
// so:
// - determines the new SummaryPoint index, which will one message older than
// message_window / 2
// - summarizes the messages from this new SummaryPoint to the
// oldest message not yet Summarized.
//
// When summarizing, it adds context from these messages to an existing summary
// if there is one.
func (se *SummaryExtractor) Extract(
	ctx context.Context,
	appState *models.AppState,
	messageEvent *models.MessageEvent,
) error {
	if appState.Config.Memory.MessageWindow == 0 {
		return NewExtractorError("SummaryExtractor message window is 0", nil)
	}

	// GetMemory will return the Messages up to the last SummaryPoint, which is the
	// last message that was summarized, and last the Summary.
	messagesSummary, err := appState.MemoryStore.GetMemory(ctx, appState, messageEvent.SessionID, 0)
	if err != nil {
		return NewExtractorError("SummaryExtractor get memory failed", err)
	}

	messages := messagesSummary.Messages
	if messages == nil {
		return NewExtractorError("SummaryExtractor messages is nil", nil)
	}
	// If we're still under the message window, we don't need to summarize.
	if len(messages) < appState.Config.Memory.MessageWindow {
		return nil
	}

	newSummary, err := summarize(
		ctx, appState, appState.Config.Memory.MessageWindow, messages, messagesSummary.Summary, 0,
	)
	if err != nil {
		return NewExtractorError("SummaryExtractor summarize failed", err)
	}

	err = appState.MemoryStore.PutSummary(
		ctx,
		appState,
		messageEvent.SessionID,
		newSummary,
	)
	if err != nil {
		return NewExtractorError("SummaryExtractor put summary failed", err)
	}

	return nil
}

func (se *SummaryExtractor) Notify(
	ctx context.Context,
	appState *models.AppState,
	messageEvents *models.MessageEvent,
) error {
	log.Debugf("SummaryExtractor notify: %v", messageEvents)
	if messageEvents == nil {
		return NewExtractorError(
			"SummaryExtractor message events is nil at Notify",
			nil,
		)
	}
	go func() {
		err := se.Extract(ctx, appState, messageEvents)
		if err != nil {
			log.Error(fmt.Sprintf("SummaryExtractor extract failed: %v", err))
		}
	}()
	return nil
}

func NewSummaryExtractor() *SummaryExtractor {
	return &SummaryExtractor{}
}

// summarize takes a slice of messages and a summary and returns a slice of messages that,
// if larger than the window size, results in the messages slice being halved. If the slice of messages is larger than
// the window size, the summary is updated to reflect the oldest messages that are removed. Expects messages to be in
// chronological order, with the oldest first.
func summarize(
	ctx context.Context,
	appState *models.AppState,
	windowSize int,
	messages []models.Message,
	summary *models.Summary,
	promptTokens int,
) (*models.Summary, error) {
	var currentSummaryContent string
	if summary != nil {
		currentSummaryContent = summary.Content
	}

	// New messages reduced to Half the windowSize to minimize the need to summarize new messages in the future.
	newMessageCount := windowSize / 2

	// Oldest messages that are over the newMessageCount
	messagesToSummarize := messages[:len(messages)-newMessageCount]

	modelName, err := llms.GetLLMModelName(appState.Config)
	if err != nil {
		return &models.Summary{}, err
	}
	maxTokens, ok := llms.MaxLLMTokensMap[modelName]
	if !ok {
		return &models.Summary{}, fmt.Errorf("model name not found in MaxLLMTokensMap")
	}

	if promptTokens == 0 {
		// rough calculation of tokes for current prompt, plus some headroom
		promptTokens = 250
	}

	// We use this to determine how many tokens we can use for the incremental summarization
	// loop. We add more messages to a summarization loop until we hit this.
	summarizerMaxInputTokens := maxTokens - SummaryMaxOutputTokens - promptTokens

	// Take the oldest messages that are over newMessageCount and summarize them.
	newSummary, err := processOverLimitMessages(
		ctx,
		appState,
		messagesToSummarize,
		summarizerMaxInputTokens,
		currentSummaryContent,
	)
	if err != nil {
		return &models.Summary{}, err
	}

	if newSummary.Content == "" {
		fmt.Println(newSummary)
		return &models.Summary{}, fmt.Errorf(
			"no summary found after summarization",
		)
	}

	return newSummary, nil
}

// processOverLimitMessages takes a slice of messages and a summary and enriches
// the summary with the messages content. Summary can an empty string. Returns a
// Summary model with enriched summary and the number of tokens in the summary.
func processOverLimitMessages(
	ctx context.Context,
	appState *models.AppState,
	messages []models.Message,
	summarizerMaxInputTokens int,
	summary string,
) (*models.Summary, error) {
	var tempMessageText []string //nolint:prealloc
	var newSummary string
	var newSummaryTokens int

	var err error
	totalTokensTemp := 0

	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages to summarize")
	}

	newSummaryPointUUID := messages[len(messages)-1].UUID

	processSummary := func() error {
		newSummary, newSummaryTokens, err = incrementalSummarizer(
			ctx,
			appState,
			summary,
			tempMessageText,
			SummaryMaxOutputTokens,
		)
		if err != nil {
			return err
		}
		tempMessageText = []string{}
		totalTokensTemp = 0
		return nil
	}

	for _, message := range messages {
		messageText := fmt.Sprintf("%s: %s", message.Role, message.Content)
		messageTokens, err := llms.GetTokenCount(messageText)
		if err != nil {
			return nil, err
		}

		if totalTokensTemp+messageTokens > summarizerMaxInputTokens {
			err = processSummary()
			if err != nil {
				return nil, err
			}
		}

		tempMessageText = append(tempMessageText, messageText)
		totalTokensTemp += messageTokens
	}

	if len(tempMessageText) > 0 {
		err = processSummary()
		if err != nil {
			return nil, err
		}
	}

	return &models.Summary{
		Content:          newSummary,
		TokenCount:       newSummaryTokens,
		SummaryPointUUID: newSummaryPointUUID,
	}, nil
}

// incrementalSummarizer takes a slice of messages and a summary, calls the LLM,
// and returns a new summary enriched with the messages content. Summary can be
// an empty string. Returns a string with the new summary and the number of
// tokens in the summary.
func incrementalSummarizer(
	ctx context.Context,
	appState *models.AppState,
	currentSummary string,
	messages []string,
	summaryMaxTokens int,
) (string, int, error) {
	if len(messages) < 1 {
		return "", 0, NewExtractorError("No messages provided", nil)
	}

	messagesJoined := strings.Join(messages, "\n")
	prevSummary := ""
	if currentSummary != "" {
		prevSummary = currentSummary
	}

	promptData := SummaryPromptTemplateData{
		PrevSummary:    prevSummary,
		MessagesJoined: messagesJoined,
	}

	progressivePrompt, err := internal.ParsePrompt(summaryPromptTemplate, promptData)
	if err != nil {
		return "", 0, err
	}

	resp, err := llms.RunChatCompletion(ctx, appState, summaryMaxTokens, progressivePrompt)
	if err != nil {
		return "", 0, err
	}

	completion := resp.Choices[0].Message.Content
	tokensUsed := resp.Usage.TotalTokens

	return completion, tokensUsed, nil
}
