package model

// ComplaintForm contiene los valores del formulario del Libro de Reclamaciones.
// Se usa para repoblar el formulario y mostrar mensajes de éxito o error.
type ComplaintForm struct {
	FullName        string
	DocumentNumber  string
	Address         string
	Phone           string
	Email           string
	GoodType        string
	GoodDescription string
	ClaimType       string
	Detail          string
	Request         string
	RegisteredAt    string // formato YYYY-MM-DD

	Error   string // mensaje de error a mostrar (vacío si no hay)
	Success bool   // true si la reclamación se registró correctamente
}
