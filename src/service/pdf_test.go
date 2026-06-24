package service

import (
	"os"
	"testing"
	"time"

	"alc/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

// TestGenerateComplaintPDF verifica que el PDF se genera correctamente con las
// fuentes UTF-8 cargadas (tildes y ñ) y un encabezado de PDF válido.
func TestGenerateComplaintPDF(t *testing.T) {
	s := NewPDFService(os.DirFS("../assets"))
	if len(s.regularFont) == 0 || len(s.boldFont) == 0 {
		t.Fatalf("fuentes no cargadas: regular=%d bold=%d", len(s.regularFont), len(s.boldFont))
	}

	c := repository.Complaint{
		ComplaintID:     42,
		FullName:        "José Ramírez Ñández",
		DocumentNumber:  "12345678",
		Address:         pgtype.Text{String: "Av. España 123, Carabayllo", Valid: true},
		Phone:           pgtype.Text{String: "948846618", Valid: true},
		Email:           pgtype.Text{String: "jose@example.com", Valid: true},
		GoodType:        "servicio",
		GoodDescription: "Instalación de ventanas de aluminio con vidrio templado.",
		ClaimType:       "reclamo",
		Detail:          "La instalación presentó filtraciones de agua y el acabado no corresponde a lo acordado en la cotización.",
		Request:         "Solicito la reparación o reinstalación completa sin costo adicional.",
		CompanyNotes:    "",
		RegisteredAt:    pgtype.Date{Time: time.Date(2026, 6, 24, 0, 0, 0, 0, time.UTC), Valid: true},
	}

	pdf, err := s.GenerateComplaintPDF(c)
	if err != nil {
		t.Fatalf("error generando PDF: %v", err)
	}
	if len(pdf) < 1000 || string(pdf[:4]) != "%PDF" {
		t.Fatalf("PDF inválido: len=%d", len(pdf))
	}

	// Escribe el PDF en el directorio temporal del test para posible inspección.
	out := t.TempDir() + "/complaint.pdf"
	if err := os.WriteFile(out, pdf, 0o644); err != nil {
		t.Fatalf("error escribiendo PDF: %v", err)
	}
	t.Logf("PDF generado: %s (%d bytes)", out, len(pdf))
}
