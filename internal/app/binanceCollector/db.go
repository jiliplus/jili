package binancecollector

func save(trades []*trade) {
	for _, t := range trades {
		db.Create(t)
	}
}
