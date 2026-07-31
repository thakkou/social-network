package fixtures

// ConversationInfo describes a seeded private conversation.
type ConversationInfo struct {
	ID      int
	User1ID int
	User2ID int
}

// MessageInfo describes a seeded message.
type MessageInfo struct {
	ID             int
	ConversationID int
	SenderID       int
	Text           string
}

// Named conversation constants.
var (
	ConvAliceBob = ConversationInfo{ID: 1, User1ID: Alice.ID, User2ID: Bob.ID}
	ConvAliceFarid = ConversationInfo{ID: 2, User1ID: Alice.ID, User2ID: Farid.ID}
)

// MessagesByConversation maps conversation IDs to their seeded messages.
var MessagesByConversation = map[int][]MessageInfo{
	1: {
		{ID: 1, ConversationID: 1, SenderID: Alice.ID, Text: "Hey Bob, are we still on for tomorrow?"},
		{ID: 2, ConversationID: 1, SenderID: Bob.ID, Text: "Yep, see you tomorrow!"},
	},
	2: {
		{ID: 3, ConversationID: 2, SenderID: Alice.ID, Text: "Can you send me the trail map?"},
		{ID: 4, ConversationID: 2, SenderID: Farid.ID, Text: "Sounds good, thanks!"},
	},
}

// AllConversations returns all seeded conversations.
func AllConversations() []ConversationInfo {
	return []ConversationInfo{ConvAliceBob, ConvAliceFarid}
}

// IsConversationParticipant returns true if userID is part of the conversation.
func IsConversationParticipant(convID, userID int) bool {
	for _, conv := range AllConversations() {
		if conv.ID == convID {
			return conv.User1ID == userID || conv.User2ID == userID
		}
	}
	return false
}
