package internal

import "strconv"

func FindByPhrases(phrases []string) string {
	query := `SELECT * FROM phrases WHERE phrase IN (`

	for i := 0; i < len(phrases); i++ {
		query += `'` + phrases[i] + `'`
		if i != len(phrases)-1 {
			query += `,`
		}
	}

	query += `)`

	return query
}

func BulkInsertPhrases(phrases []string) string {
	query := `INSERT INTO phrases (phrase) VALUES `

	for i := 0; i < len(phrases); i++ {
		query += `('` + phrases[i] + `')`
		if i != len(phrases)-1 {
			query += `,`
		}
	}

	return query
}

func GetListPhrases(offset int) string {
	return `SELECT * FROM phrases LIMIT 100 OFFSET ` + strconv.Itoa(offset)
}
