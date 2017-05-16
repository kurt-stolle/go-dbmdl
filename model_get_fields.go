package dbmdl


func (m *Model) GetFields() []string {
    var fields []string
    for i := 0; i < m.Type.NumField(); i++ {
        field := m.Type.Field(i) // Get the field at index i
        if field.Tag.Get("dbmdl") == "" {
            continue
        }

        for _, tag := range getTagParameters(field) {
            if tag == omit {
                continue
            }
        }

        fields = append(fields, field.Name)
    }

    return fields
}

