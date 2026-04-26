package domain

type ProductQuery struct {
	Category string
	Sort     string
	Limit    int
	Offset   int
}

type OrderQuery struct {
	Phone  string
	Status string
	Limit  int
	Offset int
}
