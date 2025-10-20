package service

import (
	"alc/repository"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	// Register types for gob encoding (used by gorilla/sessions)
	gob.Register(time.Time{})
}

var (
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrUserNotActive      = errors.New("usuario no activo")
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrSessionNotFound    = errors.New("sesión no encontrada")
)

// SessionData contiene la información almacenada en la sesión
type SessionData struct {
	AdminID  int32
	Username string
	Email    string
}

// AuthService define la interfaz para el servicio de autenticación
type AuthService interface {
	// Login autentica un administrador con username/email y contraseña
	Login(ctx context.Context, identifier, password string) (*SessionData, error)

	// Logout invalida la sesión del administrador
	Logout(sessionID string) error

	// ValidateSession valida si existe una sesión activa y devuelve los datos
	ValidateSession(sessionID string) (*SessionData, error)

	// HashPassword genera un hash bcrypt de la contraseña
	HashPassword(password string) (string, error)

	// ComparePassword compara una contraseña con su hash
	ComparePassword(hashedPassword, password string) error
}

// SessionAuthService implementa AuthService usando sesiones de cookie
type SessionAuthService struct {
	queries      *repository.Queries
	sessionStore *sessions.CookieStore
}

// NewSessionAuthService crea una nueva instancia de SessionAuthService
func NewSessionAuthService(queries *repository.Queries, sessionSecret string) *SessionAuthService {
	store := sessions.NewCookieStore([]byte(sessionSecret))

	// Configurar opciones de cookie
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 días
		HttpOnly: true,
		Secure:   false, // Cambiar a true en producción con HTTPS
		SameSite: 2,     // Lax
	}

	return &SessionAuthService{
		queries:      queries,
		sessionStore: store,
	}
}

// Login autentica un administrador
func (s *SessionAuthService) Login(ctx context.Context, identifier, password string) (*SessionData, error) {
	// Intentar buscar por username primero
	admin, err := s.queries.GetAdminByUsername(ctx, identifier)
	if err != nil {
		// Si no se encuentra por username, intentar por email
		admin, err = s.queries.GetAdminByEmail(ctx, identifier)
		if err != nil {
			return nil, ErrInvalidCredentials
		}
	}

	// Verificar que el admin esté activo
	if !admin.IsActive.Valid || !admin.IsActive.Bool {
		return nil, ErrUserNotActive
	}

	// Comparar contraseñas
	if err := s.ComparePassword(admin.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Crear datos de sesión
	sessionData := &SessionData{
		AdminID:  admin.AdminID,
		Username: admin.Username,
		Email:    admin.Email,
	}

	return sessionData, nil
}

// Logout invalida la sesión (en este caso, el cliente debe borrar la cookie)
func (s *SessionAuthService) Logout(sessionID string) error {
	// En una implementación cookie-based simple, el logout se maneja en el cliente
	// Si necesitáramos invalidar sesiones del lado del servidor, usaríamos Redis o similar
	return nil
}

// ValidateSession valida una sesión (no usado en implementación cookie-based)
func (s *SessionAuthService) ValidateSession(sessionID string) (*SessionData, error) {
	// Esta validación se hará en el middleware usando la cookie directamente
	return nil, ErrSessionNotFound
}

// HashPassword genera un hash bcrypt de la contraseña
func (s *SessionAuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error al hashear contraseña: %w", err)
	}
	return string(hash), nil
}

// ComparePassword compara una contraseña con su hash
func (s *SessionAuthService) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GetSessionStore devuelve el store de sesiones (útil para handlers)
func (s *SessionAuthService) GetSessionStore() *sessions.CookieStore {
	return s.sessionStore
}

// Constantes para nombres de sesión y claves
const (
	SessionName     = "admin_session"
	SessionKeyID    = "admin_id"
	SessionKeyUser  = "username"
	SessionKeyEmail = "email"
	SessionKeyExpiry = "expires_at"
)

// GetSessionData extrae los datos de sesión de una sesión de Gorilla
func GetSessionData(session *sessions.Session) (*SessionData, error) {
	adminID, ok := session.Values[SessionKeyID].(int32)
	if !ok {
		return nil, ErrSessionNotFound
	}

	username, ok := session.Values[SessionKeyUser].(string)
	if !ok {
		return nil, ErrSessionNotFound
	}

	email, ok := session.Values[SessionKeyEmail].(string)
	if !ok {
		return nil, ErrSessionNotFound
	}

	// Verificar expiración
	expiresAt, ok := session.Values[SessionKeyExpiry].(time.Time)
	if !ok || time.Now().After(expiresAt) {
		return nil, ErrSessionNotFound
	}

	return &SessionData{
		AdminID:  adminID,
		Username: username,
		Email:    email,
	}, nil
}

// SetSessionData almacena los datos en la sesión de Gorilla
func SetSessionData(session *sessions.Session, data *SessionData) {
	session.Values[SessionKeyID] = data.AdminID
	session.Values[SessionKeyUser] = data.Username
	session.Values[SessionKeyEmail] = data.Email
	session.Values[SessionKeyExpiry] = time.Now().Add(7 * 24 * time.Hour) // 7 días
}
