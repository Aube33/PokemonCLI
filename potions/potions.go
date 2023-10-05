package potions

var Potions = map[string]map[string]interface{} {
	"Potion": {
		"Description": "Heal 6 HP",
		"Heal": 6,
		"DropRate": 60,
	},
	"Super Potion": {
		"Description": "Heal 10 HP",
		"Heal": 10,
		"DropRate": 20,
	},
	"Hyper Potion": {
		"Description": "Heal 20 HP",
		"Heal": 20,
		"DropRate": 10,
	},
	"Full Restore": {
		"Description": "Restore all HP and remove any status",
		"Heal": 50,
		"DropRate": 5,
	},
	"Revive": {
		"Description": "Revive a Pokemon with half his life",
		"Heal": 0,
		"DropRate": 5,
	},
}