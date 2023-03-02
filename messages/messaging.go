package messages

const MAKE_SELECTION_MSG string = "Make Selection from list..."

const (
	DASH    string = "-"
	ONE_SET int    = 1
)

type Trivia struct {
	QuestionID string   `json:"questionid"`
	Question   string   `json:"question"`
	Category   string   `json:"category"`
	Answer     string   `json:"answer"`
	Choices    []string `json:"choices"`
	Timestamp  string   `json:"timestamp"`
}

type TriviaTable struct {
	Question string `json:"question"`
	Category string `json:"category"`
	Answer   string `json:"answer"`
}

// QuestionResponse Request-Response messaging
type QuestionResponse struct {
	QuestionID string   `json:"questionid"`
	Question   string   `json:"question"`
	Category   string   `json:"category"`
	Choices    []string `json:"choices"`
	Timestamp  string   `json:"timestamp"`
	Warning    string   `json:"warning,omitempty"`
	Error      string   `json:"error,omitempty"`
}

type AnswerRequest struct {
	QuestionID string `json:"questionid"`
	Response   string `json:"response"`
}

type AnswerResponse struct {
	Question  string `json:"question"`
	Timestamp string `json:"timestamp"`
	Category  string `json:"category"`
	Response  string `json:"response"`
	Answer    string `json:"answer"`
	Correct   bool   `json:"correct"`
	Message   string `json:"message,omitempty"`
	Warning   string `json:"warning,omitempty"`
	Error     string `json:"error,omitempty"`
}
