package entity

type Notification struct {
	ID       int             `json:"id"`
	UUID     string          `json:"uuid"`
	Title    string          `json:"title"`
	Message  string          `json:"message"`
	Channels map[int]Channel `json:"channels"`
	Event    Event           `json:"event"`
	Retries  int64           `json:"retries"`
}

type Event struct {
	Name      string `json:"name"`
	Currency  string `json:"currency"`
	Requester string `json:"requester"`
	Receiver  string `json:"receiver"`
	Category  string `json:"category"`
	Timestamp int64  `json:"timestamp"`
	CostCents int64  `json:"cost_cents"`
}

type NotificationError struct {
	ID    int    `json:"id"`
	UUID  string `json:"uuid"`
	Body  []byte `json:"body"`
	Error string `json:"error"`
}

type Email struct {
	Recipient string `json:"recipient" validate:"required" example:"example@gmail.com"`
	Subject   string `json:"subject" validate:"required" example:"subject - create account"`
	Body      string `json:"body" validate:"required" example:"hello sr...e"`
}
