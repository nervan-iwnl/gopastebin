package storage

// Storage — общий интерфейс для любого хранилища паст.
type Storage interface {
	// Save кладёт content по path и возвращает фактический path (обычно тот же).
	Save(path string, content []byte) (string, error)
	// Delete удаляет контент по path.
	Delete(path string) error
	// Download возвращает содержимое по path.
	Download(path string) ([]byte, error)
}
