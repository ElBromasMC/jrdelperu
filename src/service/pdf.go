package service

import (
	"alc/repository"
	"fmt"
	"io/fs"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/core/entity"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

// Paleta de colores del diseño (coincide con tailwind.config.cjs).
var (
	pdfNavy      = &props.Color{Red: 29, Green: 39, Blue: 71}    // #1d2747
	pdfApple     = &props.Color{Red: 171, Green: 36, Blue: 34}   // #ab2422
	pdfWhite     = &props.Color{Red: 255, Green: 255, Blue: 255} // #ffffff
	pdfGray      = &props.Color{Red: 88, Green: 91, Blue: 94}    // #585b5e (livid)
	pdfDark      = &props.Color{Red: 31, Green: 41, Blue: 55}    // texto principal
	pdfLightGray = &props.Color{Red: 245, Green: 245, Blue: 245} // fondo de tarjetas
	pdfBorder    = &props.Color{Red: 221, Green: 221, Blue: 221} // bordes suaves
)

const pdfFontFamily = "dejavu"

// PDFService genera documentos PDF (ej: el Libro de Reclamaciones).
type PDFService struct {
	regularFont []byte
	boldFont    []byte
	logo        []byte
}

// NewPDFService crea el servicio leyendo las fuentes y el logo del sistema de
// archivos embebido (assets). Si algún recurso falta, se omite sin fallar.
func NewPDFService(assets fs.FS) *PDFService {
	s := &PDFService{}
	if b, err := fs.ReadFile(assets, "static/fonts/DejaVuSans.ttf"); err == nil {
		s.regularFont = b
	}
	if b, err := fs.ReadFile(assets, "static/fonts/DejaVuSans-Bold.ttf"); err == nil {
		s.boldFont = b
	}
	if b, err := fs.ReadFile(assets, "static/img/logo.png"); err == nil {
		s.logo = b
	}
	return s
}

// GenerateComplaintPDF construye el PDF de una reclamación y devuelve sus bytes.
func (s *PDFService) GenerateComplaintPDF(c repository.Complaint) ([]byte, error) {
	builder := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithLeftMargin(15).
		WithRightMargin(15).
		WithTopMargin(12).
		WithBottomMargin(12)

	// Registrar las fuentes UTF-8 para soportar tildes y ñ.
	if len(s.regularFont) > 0 && len(s.boldFont) > 0 {
		builder = builder.
			WithCustomFonts([]*entity.CustomFont{
				{Family: pdfFontFamily, Style: fontstyle.Normal, Bytes: s.regularFont},
				{Family: pdfFontFamily, Style: fontstyle.Bold, Bytes: s.boldFont},
			}).
			WithDefaultFont(&props.Font{Family: pdfFontFamily, Size: 9, Color: pdfDark})
	}

	m := maroto.New(builder.Build())

	s.addHeader(m)
	s.addCompanyBar(m)
	s.addIntro(m)

	// 1. Datos del consumidor reclamante
	addSectionHeader(m, "1", "DATOS DEL CONSUMIDOR RECLAMANTE")
	addField(m, "Nombre Completo:", c.FullName)
	addField(m, "DNI / CE:", c.DocumentNumber)
	addField(m, "Dirección:", textOrDash(c.Address))
	addField(m, "Teléfono:", textOrDash(c.Phone))
	addField(m, "Correo Electrónico:", textOrDash(c.Email))
	addSpacer(m)

	// 2. Identificación del bien contratado
	addSectionHeader(m, "2", "IDENTIFICACIÓN DEL BIEN CONTRATADO")
	addField(m, "Tipo:", titleCase(goodTypeWord(c.GoodType)))
	addBlock(m, "Descripción del Producto o Servicio:", c.GoodDescription)
	addSpacer(m)

	// 3. Datos de la reclamación
	addSectionHeader(m, "3", "DATOS DE LA RECLAMACIÓN")
	addField(m, "Tipo:", titleCase(claimTypeWord(c.ClaimType)))
	addClaimInfoCards(m)
	addSpacer(m)

	// 4. Detalle de la reclamación
	addSectionHeader(m, "4", "DETALLE DE LA RECLAMACIÓN")
	addBlock(m, "Hechos que motivan el reclamo o queja:", c.Detail)
	addSpacer(m)

	// 5. Pedido del consumidor
	addSectionHeader(m, "5", "PEDIDO DEL CONSUMIDOR")
	addBlock(m, "Solución que espera recibir:", c.Request)
	addSpacer(m)

	// 6. Observaciones de la empresa
	addSectionHeader(m, "6", "OBSERVACIONES DE LA EMPRESA")
	addBlock(m, "Espacio reservado para uso interno:", dashIfEmpty(c.CompanyNotes))
	addSpacer(m)

	// Fecha de registro y firma
	addFooterFields(m, c)
	addLegalNotes(m)

	doc, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("error al generar PDF: %w", err)
	}
	return doc.GetBytes(), nil
}

// addHeader agrega el logo y el título principal.
func (s *PDFService) addHeader(m core.Maroto) {
	title := col.New(10).Add(
		text.New("LIBRO DE RECLAMACIONES", props.Text{
			Align: align.Center, Style: fontstyle.Bold, Size: 20, Color: pdfNavy, Top: 3,
		}),
		text.New("V I R T U A L", props.Text{
			Align: align.Center, Style: fontstyle.Bold, Size: 9, Color: pdfGray, Top: 13,
		}),
	)

	if len(s.logo) > 0 {
		logo := image.NewFromBytesCol(2, s.logo, extension.Png, props.Rect{Center: true, Percent: 90})
		m.AddRow(20, logo, title)
	} else {
		m.AddRow(20, title)
	}
}

// addCompanyBar agrega la barra azul con el nombre y RUC de la empresa.
func (s *PDFService) addCompanyBar(m core.Maroto) {
	bar := col.New(12).Add(
		text.New("CORPORACIÓN JR DEL PERÚ S.A.C.", props.Text{
			Align: align.Center, Style: fontstyle.Bold, Size: 12, Color: pdfWhite, Top: 2,
		}),
		text.New("RUC: 20604874221", props.Text{
			Align: align.Center, Size: 9, Color: pdfWhite, Top: 8,
		}),
	).WithStyle(&props.Cell{BackgroundColor: pdfNavy})
	m.AddRow(13, bar)
	addSpacer(m)
}

// addIntro agrega el texto introductorio.
func (s *PDFService) addIntro(m core.Maroto) {
	m.AddAutoRow(text.NewCol(12,
		"De conformidad con lo establecido en el Código de Protección y Defensa del Consumidor, "+
			"ponemos a disposición de nuestros clientes el presente Libro de Reclamaciones Virtual.",
		props.Text{Align: align.Center, Size: 8, Color: pdfGray},
	))
	addSpacer(m)
}

// addSectionHeader agrega un encabezado de sección numerado (rojo + azul).
func addSectionHeader(m core.Maroto, number, title string) {
	numCol := text.NewCol(1, number, props.Text{
		Align: align.Center, Style: fontstyle.Bold, Size: 11, Color: pdfWhite, Top: 1.5,
	}).WithStyle(&props.Cell{BackgroundColor: pdfApple})

	titleCol := text.NewCol(11, title, props.Text{
		Align: align.Left, Style: fontstyle.Bold, Size: 10, Color: pdfWhite, Left: 3, Top: 2,
	}).WithStyle(&props.Cell{BackgroundColor: pdfNavy})

	m.AddRow(8, numCol, titleCol)
}

// addField agrega una fila con etiqueta y valor en la misma línea.
func addField(m core.Maroto, label, value string) {
	labelCol := text.NewCol(4, label, props.Text{
		Style: fontstyle.Bold, Size: 9, Color: pdfNavy, Top: 1, Left: 1,
	})
	valueCol := text.NewCol(8, value, props.Text{
		Size: 9, Color: pdfDark, Top: 1,
	})
	m.AddRow(6, labelCol, valueCol).WithStyle(&props.Cell{
		BorderType: border.Bottom, BorderColor: pdfBorder, BorderThickness: 0.1,
	})
}

// addBlock agrega una etiqueta y, debajo, un texto largo que se ajusta solo.
func addBlock(m core.Maroto, label, value string) {
	m.AddRow(5, text.NewCol(12, label, props.Text{
		Style: fontstyle.Bold, Size: 9, Color: pdfNavy, Left: 1, Top: 1,
	}))
	m.AddAutoRow(text.NewCol(12, value, props.Text{
		Size: 9, Color: pdfDark, Left: 1, Top: 1, Bottom: 1,
	}))
}

// addClaimInfoCards agrega las dos tarjetas informativas (Reclamo / Queja).
func addClaimInfoCards(m core.Maroto) {
	cardStyle := &props.Cell{
		BackgroundColor: pdfLightGray,
		BorderType:      border.Full,
		BorderColor:     pdfBorder,
		BorderThickness: 0.1,
	}
	reclamo := col.New(6).Add(
		text.New("RECLAMO:", props.Text{Style: fontstyle.Bold, Size: 8, Color: pdfNavy, Top: 1.5, Left: 2}),
		text.New("Disconformidad relacionada con los productos o servicios brindados por la empresa.",
			props.Text{Size: 7, Color: pdfGray, Top: 5, Left: 2, Right: 2}),
	).WithStyle(cardStyle)
	queja := col.New(6).Add(
		text.New("QUEJA:", props.Text{Style: fontstyle.Bold, Size: 8, Color: pdfNavy, Top: 1.5, Left: 2}),
		text.New("Malestar o descontento respecto a la atención recibida por parte de la empresa.",
			props.Text{Size: 7, Color: pdfGray, Top: 5, Left: 2, Right: 2}),
	).WithStyle(cardStyle)
	m.AddRow(16, reclamo, queja)
}

// addFooterFields agrega la fecha de registro y el espacio de firma.
func addFooterFields(m core.Maroto, c repository.Complaint) {
	fecha := col.New(6).Add(
		text.New("FECHA DE REGISTRO:", props.Text{Style: fontstyle.Bold, Size: 9, Color: pdfNavy, Top: 2, Left: 2}),
		text.New(dateString(c.RegisteredAt), props.Text{Size: 9, Color: pdfDark, Top: 8, Left: 2}),
	).WithStyle(&props.Cell{BackgroundColor: pdfLightGray, BorderType: border.Full, BorderColor: pdfBorder, BorderThickness: 0.1})
	firma := col.New(6).Add(
		text.New("_______________________________", props.Text{Align: align.Center, Size: 9, Color: pdfGray, Top: 8}),
		text.New("FIRMA DEL CONSUMIDOR (Opcional)", props.Text{Align: align.Center, Style: fontstyle.Bold, Size: 8, Color: pdfNavy, Top: 12}),
	)
	m.AddRow(18, fecha, firma)
	addSpacer(m)
}

// addLegalNotes agrega las notas legales en la barra azul inferior.
func addLegalNotes(m core.Maroto) {
	notes := col.New(12).Add(
		text.New("La formulación del reclamo o queja no impide acudir a otras vías de solución de controversias "+
			"ni constituye una denuncia ante la autoridad competente.",
			props.Text{Size: 7, Color: pdfWhite, Top: 2, Left: 3, Right: 3}),
		text.New("La empresa dará respuesta dentro del plazo establecido por la normativa vigente.",
			props.Text{Size: 7, Color: pdfWhite, Top: 9, Left: 3, Right: 3}),
	).WithStyle(&props.Cell{BackgroundColor: pdfNavy})
	m.AddRow(15, notes)
}

// addSpacer agrega una fila vacía como separación vertical.
func addSpacer(m core.Maroto) {
	m.AddRow(3, col.New(12))
}

// --- Helpers de formato ---

func textOrDash(t pgtype.Text) string {
	if t.Valid && t.String != "" {
		return t.String
	}
	return "—"
}

func dashIfEmpty(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}

func dateString(d pgtype.Date) string {
	if !d.Valid {
		return "—"
	}
	return d.Time.Format("02/01/2006")
}

func claimTypeWord(t string) string {
	if t == "queja" {
		return "queja"
	}
	return "reclamo"
}

func goodTypeWord(t string) string {
	if t == "servicio" {
		return "servicio"
	}
	return "producto"
}

func titleCase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
