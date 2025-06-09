package resource

type SecurityPurchase struct {
	Security *Security `json:"security"`
	Count    int       `json:"count"`
}

func (sb *SecurityPurchase) TotalAmount() float64 {
	return sb.Security.Price * float64(sb.Count)
}

type PurchaseBalance struct {
	SecuritiesBalance []*SecurityPurchase `json:"securitiesBalance"`
	AmountInvested    float64             `json:"amountInvested"`
}

func (pb *PurchaseBalance) TotalCount() int {
	var total int = 0
	for _, sb := range pb.SecuritiesBalance {
		total += sb.Count
	}
	return total
}

func (pb *PurchaseBalance) AmountSpent() float64 {
	var total float64 = 0.0
	for _, sb := range pb.SecuritiesBalance {
		total += sb.TotalAmount()
	}
	return total
}

func (pb *PurchaseBalance) RemainingBalance() float64 {
	return pb.AmountInvested - pb.AmountSpent()
}
