package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	routes = append(routes, Route{"offert", "GET", "/offert/{id}", offertHandler})
}

func offertHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars

	type Items struct {
		Text   string
		Amount string
		Cost   string
	}
	data := struct {
		Title     string
		Name      string
		RegNo     string
		CarType   string
		Mileage   string
		OrderNo   string
		Items     []Items
		TotalCost string
	}{
		Title:     "Offert Gruffman",
		Name:      "Fredrik Grufman",
		RegNo:     "ABC123",
		CarType:   "Ford Mustang 2015",
		Mileage:   "1234",
		OrderNo:   "234567",
		Items:     []Items{},
		TotalCost: "70",
	}
	data.Items = append(data.Items, Items{Text: "Artikel 1", Amount: "2", Cost: "20"})
	data.Items = append(data.Items, Items{Text: "Artikel 2", Amount: "1", Cost: "30"})

	t, _ := template.New("webpage").Parse(offerHtml)
	t.Execute(w, data)
	//fmt.Fprint(w, vars["id"])
}

const offerHtml = `
<!DOCTYPE html>
<html>
	<head>
		<title>{{.Title}}</title>
	</head>
	<body>
		<div>
			{{.Name}}<br/>
			{{.RegNo}}<br/>
			{{.CarType}}<br/>
			Mätarställning: {{.Mileage}} km<br/>
			Ordernr: {{.OrderNo}}<br/>
		</div>
		<div>
		<table>
			<thead>
				<tr>
					<th>Artikeltext</th>
					<th>Antal</th>
					<th>radsumma_inkmoms<th>
				</tr>
			<thead>
			<tfoot>
				<tr>
					<td colspan="3">Totalsumma inkl moms: Att betala {{.TotalCost}} SEK</td>
				</tr>
			</tfoot>
			<tbody>
				{{range .Items}}
				<tr>
					<td>{{.Text}}</td>
					<td>{{.Amount}}</td>
					<td>{{.Cost}}</td>
				</tr>
				{{else}}
				<tr>
					<td colspan="3">Inga artikelrader på denna offert</td>
				</tr>
				{{end}}
			</tbody>
		</table>
		</div>
		<div>
			
		</div>
	</body>
</html>
`
