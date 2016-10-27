package query

func Build(fields map[string]string) string {
	query := "SELECT * FROM acesso WHERE 1 = 1"
	for key, value := range fields {
		query = query + " AND "
		if key != "dataInicio" && key != "dataFim" {
			query = query + key + " = '" + value + "'"
		} else {
			if key == "dataInicio" {
				query = query + "date >= '" + value + "'"
			}
			if key == "dataFim" {
				query = query + "date <= '" + value + "'"
			}
		}
	}
	query = query + " ORDER BY date desc, time desc"
	return query
}
