package domain

type SMS struct {
	ID        int
	PhoneTo   PhoneNumber
	PhoneFrom string
	Text      string
}
