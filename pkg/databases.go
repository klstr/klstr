package klstr

type NewDatabase struct {
	Name      string
	DBType    string
	Instance  string
	Namespace string
}

func CreateNewDatabase(newdb NewDatabase) error {
	return nil
}
