package storage

func GetTestConfigFixtures() []ItemRecord {
	return []ItemRecord{
		{
			Code:  "A12T-4GH7-QPL9-3N4M",
			Name:  "Lettuce",
			Price: 3_46,
		},
		{
			Code:  "E5T6-9UI3-TH15-QR88",
			Name:  "Peach",
			Price: 2_99,
		},
		{
			Code:  "YRT6-72AS-K736-L4AR",
			Name:  "Green Pepper",
			Price: 79,
		},
		{
			Code:  "TQ4C-VV6T-75ZX-1RMR",
			Name:  "Gala Apple",
			Price: 3_59,
		},
	}
}
