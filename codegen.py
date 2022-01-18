import csv

f=open('cards.csv','r')
r = csv.reader(f, delimiter=',')
next(r)
print("""package gwent

var AllCards = []Card{""")
for line in r:
    name, displayname, strength, abilities, faction, row, picture = line

    struct = f"""Card{{
	Faction:     {faction},
	Row:         {row},
	Name:        "{name}",
	DisplayName: "{displayname}",
	Strength:     {strength},
	Effects:     []Effect{{{abilities}}},
	Image:       "{picture}",
}},"""
    print(struct)

print("}")
