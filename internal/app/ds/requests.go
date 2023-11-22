package ds

type GetSubstancesRequestBody struct {
	Name   string
	Status string
}
type OrderSynthesisRequestBody struct {
	User_id          int
	Substance_first  int
	Substance_second int
}
