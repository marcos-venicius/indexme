package indexer

import "database/sql"

type FolderTable struct {
	id             int
	name           string
	abspath        string
	documentsCount int
}

type FolderTermsFrequencyTable struct {
	id       int
	folderId int
	token    string
	frequeny int
}

type DocumentTermsFrequencyTable struct {
	id       int
	abspath  string
	folderId int
	token    string
	frequeny int
}

func (i *Indexer) CreateFolder(name, abspath string) error {
	const createFolderQuery = `INSERT INTO folder (name, abs_path) VALUES(?,?)`

	if _, err := i.db.Exec(createFolderQuery, name, abspath); err != nil {
		return err
	}

	return nil
}

func (i *Indexer) GetFolderByAbsPath(path string) *FolderTable {
	const getFolderByAbsPathQuery = `SELECT id, name, abs_path, documents_count FROM folder where abs_path = ?;`

	row := i.db.QueryRow(getFolderByAbsPathQuery, path)

	folder := &FolderTable{}

	err := row.Scan(&folder.id, &folder.name, &folder.abspath, &folder.documentsCount)

	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}

	return folder
}

func (i *Indexer) GetFolderTermsFrequency(abspath string) map[string]int {
	frequency := make(map[string]int)

	folder := i.GetFolderByAbsPath(abspath)

	const getFolderTermsFrequencyQuery = `SELECT * FROM folder_terms_frequency WHERE folder_id = ?;`

	rows, err := i.db.Query(getFolderTermsFrequencyQuery, folder.id)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		termsFrequency := FolderTermsFrequencyTable{}

		err := rows.Scan(&termsFrequency.id, &termsFrequency.folderId, &termsFrequency.token, &termsFrequency.frequeny)

		if err != nil {
			panic(err)
		}

		if count, ok := frequency[termsFrequency.token]; ok {
			frequency[termsFrequency.token] = count + termsFrequency.frequeny
		} else {
			frequency[termsFrequency.token] = termsFrequency.frequeny
		}
	}

	return frequency
}

func (i *Indexer) GetDocumentTermsFrequency(folderId int) map[string]map[string]int {
	frequency := make(map[string]map[string]int)

	const getDocumentTermsFrequencyQuery = `
    SELECT id, abs_path, folder_id, token, frequency FROM
    document_terms_frequency WHERE folder_id = ?;`

	rows, err := i.db.Query(getDocumentTermsFrequencyQuery, folderId)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		termsFrequency := DocumentTermsFrequencyTable{}

		err := rows.Scan(
			&termsFrequency.id,
			&termsFrequency.abspath,
			&termsFrequency.folderId,
			&termsFrequency.token,
			&termsFrequency.frequeny,
		)

		if err != nil {
			panic(err)
		}

		if _, ok := frequency[termsFrequency.abspath]; !ok {
			frequency[termsFrequency.abspath] = make(map[string]int)
		}

		if count, ok := frequency[termsFrequency.abspath][termsFrequency.token]; ok {
			frequency[termsFrequency.abspath][termsFrequency.token] = count + termsFrequency.frequeny
		} else {
			frequency[termsFrequency.abspath][termsFrequency.token] = termsFrequency.frequeny
		}
	}

	return frequency
}

func (i *Indexer) AddFolderTermsFrequency(folderId int, token string, frequency int) error {
	const createDocumentTermsFrequencyQuery = `
    INSERT INTO folder_terms_frequency (folder_id, token, frequency)
    VALUES (?, ?, ?);
  `

	if _, err := i.db.Exec(createDocumentTermsFrequencyQuery, folderId, token, frequency); err != nil {
		return err
	}

	return nil
}

func (i *Indexer) AddDocumentTermsFrequency(absPath string, folderId int, token string, frequency int) error {
	const createDocumentTermsFrequencyQuery = `
    INSERT INTO document_terms_frequency (abs_path, folder_id, token, frequency)
    VALUES (?, ?, ?, ?);
  `

	if _, err := i.db.Exec(createDocumentTermsFrequencyQuery, absPath, folderId, token, frequency); err != nil {
		return err
	}

	return nil
}
