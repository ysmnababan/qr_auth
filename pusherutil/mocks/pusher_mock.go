package mocks

import "sync"

type MockPush struct {
	Wg           *sync.WaitGroup
	TriggerCalls []TriggerCall
}

type TriggerCall struct {
	Channel   string
	EventName string
	Data      any
}

func NewMockPush() *MockPush {
	return &MockPush{}
}

func (m *MockPush) Trigger(channel string, eventName string, data any) error {
	if m.Wg != nil {
		defer m.Wg.Done()
	}
	m.TriggerCalls = append(m.TriggerCalls, TriggerCall{
		Channel:   channel,
		EventName: eventName,
		Data:      data,
	})
	return nil
}
