package indexer

const createFolderTable = `
  CREATE TABLE IF NOT EXISTS folder (
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    abs_path TEXT NOT NULL,
    documents_count INTEGER NOT NULL DEFAULT 0
  );
`

const createFolderTermsFrequencyTable = `
  CREATE TABLE IF NOT EXISTS folder_terms_frequency (
    id INTEGER NOT NULL PRIMARY KEY,
    folder_id INTEGER NOT NULL,
    token TEXT NOT NULL,
    frequency INTEGER NOT NULL DEFAULT 0,
    
    FOREIGN KEY(folder_id) REFERENCES folder(id)
  );
`

const createDocumentTermsFrequencyTable = `
  CREATE TABLE IF NOT EXISTS document_terms_frequency (
    id INTEGER NOT NULL PRIMARY KEY,
    abs_path TEXT NOT NULL,
    folder_id INTEGER NOT NULL,
    token TEXT NOT NULL,
    frequency INTEGER NOT NULL DEFAULT 0,

    FOREIGN KEY(folder_id) REFERENCES folder(id)
  );
`

const createFolderIndexes = `
  CREATE INDEX IF NOT EXISTS idx_folder_abs_path ON folder(abs_path);
`

const createFolderTermsFrequencyIndexes = `
  CREATE INDEX IF NOT EXISTS idx_folder_terms_frequency_token ON folder_terms_frequency(token);
`

const createDocumentTermsFrequencyIndexes = `
  CREATE INDEX IF NOT EXISTS idx_document_terms_frequency_token ON document_terms_frequency(token);
`

func (i *Indexer) DbSetup() error {
	if _, err := i.db.Exec(createFolderTable); err != nil {
		return err
	}

	if _, err := i.db.Exec(createFolderTermsFrequencyTable); err != nil {
		return err
	}

	if _, err := i.db.Exec(createDocumentTermsFrequencyTable); err != nil {
		return err
	}

	if _, err := i.db.Exec(createFolderIndexes); err != nil {
		return err
	}

	if _, err := i.db.Exec(createFolderTermsFrequencyIndexes); err != nil {
		return err
	}

	if _, err := i.db.Exec(createDocumentTermsFrequencyIndexes); err != nil {
		return err
	}

	return nil
}
