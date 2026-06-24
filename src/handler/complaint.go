package handler

import (
	"alc/model"
	"alc/repository"
	"alc/service"
	"alc/view"
	"context"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// complaintsPageSize es el número de reclamaciones que se listan por página en el admin.
const complaintsPageSize = 50

// HandleLibroReclamacionesShow muestra el formulario público del Libro de Reclamaciones.
func (h *Handler) HandleLibroReclamacionesShow(c echo.Context) error {
	// Fecha de registro por defecto: hoy
	form := model.ComplaintForm{RegisteredAt: time.Now().Format("2006-01-02")}
	return renderOK(c, view.LibroReclamaciones(form, h.recaptchaService.GetSiteKey()))
}

// HandleLibroReclamacionesSubmit procesa el envío del Libro de Reclamaciones.
func (h *Handler) HandleLibroReclamacionesSubmit(c echo.Context) error {
	ctx := c.Request().Context()

	// Verificar reCAPTCHA antes de procesar el formulario.
	if valid, err := h.recaptchaService.VerifyToken(c.FormValue("g-recaptcha-response")); err != nil || !valid {
		form := model.ComplaintForm{}
		if errors.Is(err, service.ErrRecaptchaNotFound) {
			form.Error = "Por favor, complete la verificación de seguridad."
		} else {
			form.Error = "No pudimos verificar que no es un robot. Intente nuevamente."
		}
		return h.renderFormView(c, form)
	}

	form := model.ComplaintForm{
		FullName:        strings.TrimSpace(c.FormValue("full_name")),
		DocumentNumber:  strings.TrimSpace(c.FormValue("document_number")),
		Address:         strings.TrimSpace(c.FormValue("address")),
		Phone:           strings.TrimSpace(c.FormValue("phone")),
		Email:           strings.TrimSpace(c.FormValue("email")),
		GoodType:        c.FormValue("good_type"),
		GoodDescription: strings.TrimSpace(c.FormValue("good_description")),
		ClaimType:       c.FormValue("claim_type"),
		Detail:          strings.TrimSpace(c.FormValue("detail")),
		Request:         strings.TrimSpace(c.FormValue("request")),
		RegisteredAt:    c.FormValue("registered_at"),
	}

	// Validación de campos requeridos
	if form.FullName == "" || form.DocumentNumber == "" || form.GoodDescription == "" ||
		form.Detail == "" || form.Request == "" {
		form.Error = "Por favor complete todos los campos obligatorios."
		return h.renderFormView(c, form)
	}
	if form.GoodType != "producto" && form.GoodType != "servicio" {
		form.Error = "Seleccione si su reclamación es sobre un producto o un servicio."
		return h.renderFormView(c, form)
	}
	if form.ClaimType != "reclamo" && form.ClaimType != "queja" {
		form.Error = "Seleccione si desea registrar un reclamo o una queja."
		return h.renderFormView(c, form)
	}

	// Parsear la fecha de registro
	registeredAt, err := time.Parse("2006-01-02", form.RegisteredAt)
	if err != nil {
		form.Error = "La fecha de registro no es válida."
		return h.renderFormView(c, form)
	}

	complaint, err := h.queries.CreateComplaint(ctx, repository.CreateComplaintParams{
		FullName:        form.FullName,
		DocumentNumber:  form.DocumentNumber,
		Address:         optionalText(form.Address),
		Phone:           optionalText(form.Phone),
		Email:           optionalText(form.Email),
		GoodType:        form.GoodType,
		GoodDescription: form.GoodDescription,
		ClaimType:       form.ClaimType,
		Detail:          form.Detail,
		Request:         form.Request,
		RegisteredAt:    pgtype.Date{Time: registeredAt, Valid: true},
	})
	if err != nil {
		form.Error = "Ocurrió un error al registrar su reclamación. Intente nuevamente."
		return h.renderFormView(c, form)
	}

	// Notificar por correo (sin bloquear la respuesta al usuario).
	h.sendComplaintNotifications(complaint)

	return renderOK(c, view.LibroReclamaciones(model.ComplaintForm{Success: true}, ""))
}

// renderForm vuelve a mostrar el formulario público (con valores y/o errores).
func (h *Handler) renderFormView(c echo.Context, form model.ComplaintForm) error {
	return Render(c, http.StatusOK, view.LibroReclamaciones(form, h.recaptchaService.GetSiteKey()))
}

// --- Administración ---

// HandleComplaintsIndex muestra la gestión de reclamaciones en el admin.
func (h *Handler) HandleComplaintsIndex(c echo.Context) error {
	ctx := c.Request().Context()
	sessionData := c.Get("session").(*service.SessionData)

	search := strings.TrimSpace(c.QueryParam("q"))
	complaints := h.listComplaints(ctx, search)

	unresolved, _ := h.queries.CountUnresolvedComplaints(ctx)

	return renderOK(c, view.AdminComplaintsPage(sessionData.Username, complaints, search, unresolved))
}

// HandleComplaintsSearch devuelve solo el listado de reclamaciones filtrado (HTMX).
func (h *Handler) HandleComplaintsSearch(c echo.Context) error {
	ctx := c.Request().Context()
	search := strings.TrimSpace(c.QueryParam("q"))
	complaints := h.listComplaints(ctx, search)
	return renderOK(c, view.AdminComplaintsGrid(complaints))
}

// HandleComplaintUpdate actualiza las observaciones de la empresa y el estado.
func (h *Handler) HandleComplaintUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de reclamación inválido")
	}

	err = h.queries.UpdateComplaintCompanyNotes(ctx, repository.UpdateComplaintCompanyNotesParams{
		ComplaintID:  int32(id),
		CompanyNotes: strings.TrimSpace(c.FormValue("company_notes")),
		IsResolved:   c.FormValue("is_resolved") == "on",
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error al actualizar la reclamación")
	}

	complaint, err := h.queries.GetComplaint(ctx, int32(id))
	if err != nil {
		return c.String(http.StatusNotFound, "Reclamación no encontrada")
	}

	return renderOK(c, view.AdminComplaintCard(complaint, true))
}

// HandleComplaintDelete elimina una reclamación.
func (h *Handler) HandleComplaintDelete(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "ID de reclamación inválido")
	}

	if err := h.queries.DeleteComplaint(ctx, int32(id)); err != nil {
		return c.String(http.StatusInternalServerError, "Error al eliminar la reclamación")
	}

	complaints := h.listComplaints(ctx, "")
	return renderOK(c, view.AdminComplaintsGrid(complaints))
}

// listComplaints devuelve las reclamaciones, filtradas por búsqueda si se indica.
func (h *Handler) listComplaints(ctx context.Context, search string) []repository.Complaint {
	var complaints []repository.Complaint
	var err error
	if search != "" {
		complaints, err = h.queries.SearchComplaints(ctx, repository.SearchComplaintsParams{
			Search: search,
			Limit:  complaintsPageSize,
			Offset: 0,
		})
	} else {
		complaints, err = h.queries.ListComplaints(ctx, repository.ListComplaintsParams{
			Limit:  complaintsPageSize,
			Offset: 0,
		})
	}
	if err != nil {
		return []repository.Complaint{}
	}
	return complaints
}

// optionalText convierte un string en pgtype.Text (NULL si está vacío).
func optionalText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

// --- Notificaciones por correo ---

// sendComplaintNotifications envía, en segundo plano, una notificación a la
// empresa (SMTP_TO_EMAIL) y una confirmación al consumidor (si dejó correo).
func (h *Handler) sendComplaintNotifications(complaint repository.Complaint) {
	if !h.emailService.IsConfigured() {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		detailsHTML := complaintEmailHTML(complaint)
		detailsText := complaintEmailText(complaint)
		claimWord := complaintClaimWord(complaint.ClaimType)

		// Generar el PDF de la reclamación para adjuntarlo.
		var attachments []service.EmailAttachment
		if pdfBytes, err := h.pdfService.GenerateComplaintPDF(complaint); err != nil {
			log.Printf("error generando PDF de reclamación #%d: %v", complaint.ComplaintID, err)
		} else {
			attachments = []service.EmailAttachment{{
				Filename:    fmt.Sprintf("reclamacion-%d.pdf", complaint.ComplaintID),
				Content:     pdfBytes,
				ContentType: "application/pdf",
			}}
		}

		// Notificación a la empresa.
		if to := h.emailService.CompanyEmail(); to != "" {
			err := h.emailService.Send(ctx, service.Email{
				To:      []string{to},
				Subject: fmt.Sprintf("Nueva %s #%d - %s", claimWord, complaint.ComplaintID, complaint.FullName),
				HTMLBody: "<p>Se ha registrado una nueva " + claimWord +
					" en el Libro de Reclamaciones Virtual. Se adjunta el documento en PDF.</p>" + detailsHTML,
				TextBody:    "Se ha registrado una nueva " + claimWord + ".\n\n" + detailsText,
				Attachments: attachments,
			})
			if err != nil {
				log.Printf("error notificando reclamación #%d a la empresa: %v", complaint.ComplaintID, err)
			}
		}

		// Confirmación al consumidor.
		if complaint.Email.Valid && complaint.Email.String != "" {
			err := h.emailService.Send(ctx, service.Email{
				To:      []string{complaint.Email.String},
				Subject: fmt.Sprintf("Hemos recibido su %s #%d", claimWord, complaint.ComplaintID),
				HTMLBody: "<p>Estimado(a) " + html.EscapeString(complaint.FullName) + ",</p>" +
					"<p>Hemos recibido su " + claimWord + " correctamente. La empresa dará respuesta " +
					"dentro del plazo establecido por la normativa vigente. Adjuntamos una copia en PDF " +
					"y, a continuación, el detalle registrado:</p>" +
					detailsHTML,
				TextBody:    "Hemos recibido su " + claimWord + " correctamente. Detalle registrado:\n\n" + detailsText,
				Attachments: attachments,
			})
			if err != nil {
				log.Printf("error enviando confirmación de reclamación #%d al consumidor: %v", complaint.ComplaintID, err)
			}
		}
	}()
}

// complaintClaimWord devuelve "reclamo" o "queja" en minúscula.
func complaintClaimWord(t string) string {
	if t == "queja" {
		return "queja"
	}
	return "reclamo"
}

// complaintGoodWord devuelve "producto" o "servicio".
func complaintGoodWord(t string) string {
	if t == "servicio" {
		return "servicio"
	}
	return "producto"
}

// complaintEmailHTML genera el detalle de la reclamación como tabla HTML.
func complaintEmailHTML(c repository.Complaint) string {
	row := func(label, value string) string {
		if value == "" {
			value = "—"
		}
		return "<tr>" +
			"<td style=\"padding:6px 10px;border:1px solid #ddd;background:#f5f5f5;font-weight:bold;vertical-align:top\">" +
			html.EscapeString(label) + "</td>" +
			"<td style=\"padding:6px 10px;border:1px solid #ddd\">" + html.EscapeString(value) + "</td>" +
			"</tr>"
	}

	var b strings.Builder
	b.WriteString("<table style=\"border-collapse:collapse;font-family:Arial,sans-serif;font-size:14px\">")
	b.WriteString(row("N° de reclamación", fmt.Sprint(c.ComplaintID)))
	b.WriteString(row("Fecha de registro", complaintDateString(c.RegisteredAt)))
	b.WriteString(row("Tipo", capitalize(complaintClaimWord(c.ClaimType))))
	b.WriteString(row("Nombre completo", c.FullName))
	b.WriteString(row("DNI / CE", c.DocumentNumber))
	b.WriteString(row("Dirección", optionalString(c.Address)))
	b.WriteString(row("Teléfono", optionalString(c.Phone)))
	b.WriteString(row("Correo electrónico", optionalString(c.Email)))
	b.WriteString(row("Bien contratado", capitalize(complaintGoodWord(c.GoodType))))
	b.WriteString(row("Descripción del bien", c.GoodDescription))
	b.WriteString(row("Detalle de la reclamación", c.Detail))
	b.WriteString(row("Pedido del consumidor", c.Request))
	b.WriteString("</table>")
	return b.String()
}

// complaintEmailText genera el detalle de la reclamación como texto plano.
func complaintEmailText(c repository.Complaint) string {
	line := func(label, value string) string {
		if value == "" {
			value = "—"
		}
		return label + ": " + value + "\n"
	}

	var b strings.Builder
	b.WriteString(line("N° de reclamación", fmt.Sprint(c.ComplaintID)))
	b.WriteString(line("Fecha de registro", complaintDateString(c.RegisteredAt)))
	b.WriteString(line("Tipo", capitalize(complaintClaimWord(c.ClaimType))))
	b.WriteString(line("Nombre completo", c.FullName))
	b.WriteString(line("DNI / CE", c.DocumentNumber))
	b.WriteString(line("Dirección", optionalString(c.Address)))
	b.WriteString(line("Teléfono", optionalString(c.Phone)))
	b.WriteString(line("Correo electrónico", optionalString(c.Email)))
	b.WriteString(line("Bien contratado", capitalize(complaintGoodWord(c.GoodType))))
	b.WriteString(line("Descripción del bien", c.GoodDescription))
	b.WriteString(line("Detalle de la reclamación", c.Detail))
	b.WriteString(line("Pedido del consumidor", c.Request))
	return b.String()
}

// complaintDateString formatea una fecha pgtype.Date como dd/mm/aaaa.
func complaintDateString(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("02/01/2006")
}

// optionalString extrae el valor de un pgtype.Text (vacío si es NULL).
func optionalString(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}

// capitalize pone en mayúscula la primera letra de una palabra.
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
