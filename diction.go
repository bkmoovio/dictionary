package dictionary

type Diction interface {
	Read() error
	RoutingNumberSearch(s string) *Participant
	FinancialInstitutionSearch(s string) []*Participant
	GetParticipants() []*Participant
	GetIndexRoutingNumber() map[string]*Participant
	GetIndexCustomerName() map[string][]*Participant
}
