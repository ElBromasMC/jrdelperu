package service

import (
	"alc/config"
	"alc/repository"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/h2non/bimg"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvalidFileType = errors.New("tipo de archivo no permitido")
	ErrFileTooLarge    = errors.New("archivo demasiado grande")
	ErrFileNotFound    = errors.New("archivo no encontrado")
)

const (
	MaxImageSize = 10 * 1024 * 1024  // 10 MB
	MaxPDFSize   = 20 * 1024 * 1024  // 20 MB
	UploadDir    = "uploads"         // Directorio base para archivos

	// Configuración de optimización de imágenes
	MaxImageWidth    = 2000 // Ancho máximo en píxeles
	MaxImageHeight   = 2000 // Alto máximo en píxeles
	WebPQuality      = 80   // Calidad WebP (1-100) - 80 ofrece excelente compresión
	StripMetadata    = true // Eliminar datos EXIF
)

// FileService maneja la subida, almacenamiento y eliminación de archivos
type FileService struct {
	queries   *repository.Queries
	uploadDir string
}

// NewFileService crea una nueva instancia de FileService
func NewFileService(queries *repository.Queries) *FileService {
	return &FileService{
		queries:   queries,
		uploadDir: UploadDir,
	}
}

// FileUploadResult contiene la información del archivo subido
type FileUploadResult struct {
	FileID      int32
	FileName    string
	FilePath    string
	FileType    string
	MimeType    string
	FileSize    int64
	DisplayName string
}

// generateDisplayName generates a clean display name from filename
func generateDisplayName(filename string) string {
	// Remove extension
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// Convert to slug-like format
	displayName := nameWithoutExt

	// Replace common separators with spaces
	displayName = strings.ReplaceAll(displayName, "_", " ")
	displayName = strings.ReplaceAll(displayName, "-", " ")

	// Normalize unicode (remove accents)
	displayName = strings.Map(func(r rune) rune {
		switch r {
		case 'á', 'à', 'ä', 'â': return 'a'
		case 'é', 'è', 'ë', 'ê': return 'e'
		case 'í', 'ì', 'ï', 'î': return 'i'
		case 'ó', 'ò', 'ö', 'ô': return 'o'
		case 'ú', 'ù', 'ü', 'û': return 'u'
		case 'ñ': return 'n'
		case 'Á', 'À', 'Ä', 'Â': return 'A'
		case 'É', 'È', 'Ë', 'Ê': return 'E'
		case 'Í', 'Ì', 'Ï', 'Î': return 'I'
		case 'Ó', 'Ò', 'Ö', 'Ô': return 'O'
		case 'Ú', 'Ù', 'Ü', 'Û': return 'U'
		case 'Ñ': return 'N'
		default: return r
		}
	}, displayName)

	// Title case words
	words := strings.Fields(displayName)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// inferFileType infers file type from extension and MIME type
func inferFileType(filename, mimeType string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	// Check by extension first
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg":
		return "image", nil
	case ".pdf":
		return "pdf", nil
	}

	// Fallback to MIME type
	if strings.HasPrefix(mimeType, "image/") {
		return "image", nil
	}
	if mimeType == "application/pdf" {
		return "pdf", nil
	}

	return "", ErrInvalidFileType
}

// optimizeImage optimiza una imagen para la web usando bimg/libvips
// Convierte a WebP, redimensiona, comprime y elimina metadatos
// Retorna el nuevo nombre de archivo (con .webp) y el tamaño optimizado
func (s *FileService) optimizeImage(filePath string) (string, int64, error) {
	// Leer la imagen
	buffer, err := bimg.Read(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("error al leer imagen: %w", err)
	}

	// Obtener información de la imagen
	img := bimg.NewImage(buffer)
	size, err := img.Size()
	if err != nil {
		return "", 0, fmt.Errorf("error al obtener dimensiones: %w", err)
	}

	// Determinar si necesita redimensionamiento
	needsResize := size.Width > MaxImageWidth || size.Height > MaxImageHeight

	// Configurar opciones de optimización - SIEMPRE convertir a WebP
	options := bimg.Options{
		Type:          bimg.WEBP, // Convertir a WebP para mejor compresión
		Quality:       WebPQuality,
		StripMetadata: StripMetadata,
		Enlarge:       false, // No agrandar imágenes más pequeñas
	}

	// Si necesita redimensionamiento, calcular nuevas dimensiones manteniendo aspecto
	if needsResize {
		// Calcular escala para mantener aspect ratio
		scaleW := float64(MaxImageWidth) / float64(size.Width)
		scaleH := float64(MaxImageHeight) / float64(size.Height)
		scale := scaleW
		if scaleH < scaleW {
			scale = scaleH
		}

		options.Width = int(float64(size.Width) * scale)
		options.Height = int(float64(size.Height) * scale)
	}

	// Procesar la imagen
	optimized, err := img.Process(options)
	if err != nil {
		return "", 0, fmt.Errorf("error al procesar imagen: %w", err)
	}

	// Cambiar extensión del archivo a .webp
	newFilePath := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".webp"

	// Eliminar archivo original si tiene extensión diferente
	if filePath != newFilePath {
		os.Remove(filePath)
	}

	// Guardar la imagen optimizada en formato WebP
	err = bimg.Write(newFilePath, optimized)
	if err != nil {
		return "", 0, fmt.Errorf("error al guardar imagen optimizada: %w", err)
	}

	// Retornar el nuevo nombre de archivo y tamaño
	fileInfo, err := os.Stat(newFilePath)
	if err != nil {
		return "", 0, fmt.Errorf("error al obtener tamaño de archivo: %w", err)
	}

	// Retornar solo el nombre base del archivo (sin la ruta completa)
	newFileName := filepath.Base(newFilePath)
	return newFileName, fileInfo.Size(), nil
}

// UploadFile sube un archivo al filesystem y registra en la base de datos
func (s *FileService) UploadFile(ctx context.Context, file *multipart.FileHeader, customDisplayName string) (*FileUploadResult, error) {
	// Get MIME type
	mimeType := file.Header.Get("Content-Type")

	// Infer file type from extension and MIME type
	fileType, err := inferFileType(file.Filename, mimeType)
	if err != nil {
		return nil, err
	}

	// Validar tamaño
	maxSize := MaxImageSize
	if fileType == "pdf" {
		maxSize = MaxPDFSize
	}
	if file.Size > int64(maxSize) {
		return nil, ErrFileTooLarge
	}

	// Validar MIME type
	if !s.isValidMimeType(mimeType, fileType) {
		return nil, ErrInvalidFileType
	}

	// Abrir archivo
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("error al abrir archivo: %w", err)
	}
	defer src.Close()

	// Generar nombre único basado en hash y timestamp
	uniqueName, err := s.generateUniqueFileName(src, file.Filename)
	if err != nil {
		return nil, fmt.Errorf("error al generar nombre de archivo: %w", err)
	}

	// Reset reader position
	src.Seek(0, 0)

	// Crear directorio si no existe
	uploadPath := filepath.Join(s.uploadDir, fileType+"s") // "images" o "pdfs"
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return nil, fmt.Errorf("error al crear directorio: %w", err)
	}

	// Ruta completa del archivo
	filePath := filepath.Join(uploadPath, uniqueName)

	// Crear archivo en el filesystem
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al crear archivo: %w", err)
	}
	defer dst.Close()

	// Copiar contenido
	fileSize, err := io.Copy(dst, src)
	if err != nil {
		os.Remove(filePath) // Limpiar en caso de error
		return nil, fmt.Errorf("error al copiar archivo: %w", err)
	}

	// Cerrar archivo antes de optimizar
	dst.Close()

	// Optimizar imagen si corresponde (convierte a WebP)
	if fileType == "image" {
		newFileName, optimizedSize, err := s.optimizeImage(filePath)
		if err != nil {
			// Log error pero no fallar - mantener imagen original
			fmt.Printf("Advertencia: no se pudo optimizar imagen: %v\n", err)
		} else {
			// Actualizar con el nuevo nombre de archivo (.webp) y tamaño
			uniqueName = newFileName
			fileSize = optimizedSize
			mimeType = "image/webp"
		}
	}

	// Use custom display name if provided, otherwise auto-generate from filename
	displayName := customDisplayName
	if displayName == "" {
		displayName = generateDisplayName(file.Filename)
	}

	// Registrar en base de datos
	staticFile, err := s.queries.CreateStaticFile(ctx, repository.CreateStaticFileParams{
		FileName:      uniqueName,
		FileType:      pgtype.Text{String: fileType, Valid: true},
		MimeType:      pgtype.Text{String: mimeType, Valid: true},
		FileSizeBytes: pgtype.Int8{Int64: fileSize, Valid: true},
		DisplayName:   pgtype.Text{String: displayName, Valid: true},
	})
	if err != nil {
		os.Remove(filePath) // Limpiar en caso de error
		return nil, fmt.Errorf("error al registrar archivo en BD: %w", err)
	}

	return &FileUploadResult{
		FileID:      staticFile.FileID,
		FileName:    uniqueName,
		FilePath:    filePath,
		FileType:    fileType,
		MimeType:    mimeType,
		FileSize:    fileSize,
		DisplayName: displayName,
	}, nil
}

// DeleteFile elimina un archivo del filesystem y de la base de datos
func (s *FileService) DeleteFile(ctx context.Context, fileID int32) error {
	// Obtener información del archivo
	file, err := s.queries.GetStaticFile(ctx, fileID)
	if err != nil {
		return ErrFileNotFound
	}

	// Determinar directorio basado en tipo
	fileType := "others"
	if file.FileType.Valid {
		fileType = file.FileType.String + "s"
	}

	// Construir ruta del archivo
	filePath := filepath.Join(s.uploadDir, fileType, file.FileName)

	// Eliminar archivo del filesystem
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error al eliminar archivo: %w", err)
	}

	// Eliminar registro de la base de datos
	if err := s.queries.DeleteStaticFile(ctx, fileID); err != nil {
		return fmt.Errorf("error al eliminar registro de BD: %w", err)
	}

	return nil
}

// GetFilePath devuelve la ruta del archivo en el filesystem
func (s *FileService) GetFilePath(ctx context.Context, fileID int32) (string, error) {
	file, err := s.queries.GetStaticFile(ctx, fileID)
	if err != nil {
		return "", ErrFileNotFound
	}

	fileType := "others"
	if file.FileType.Valid {
		fileType = file.FileType.String + "s"
	}

	return filepath.Join(s.uploadDir, fileType, file.FileName), nil
}

// generateUniqueFileName genera un nombre único para el archivo
func (s *FileService) generateUniqueFileName(reader io.Reader, originalName string) (string, error) {
	// Generar hash del contenido
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	hashStr := fmt.Sprintf("%x", hash.Sum(nil))[:16]

	// Obtener extensión original
	ext := filepath.Ext(originalName)
	if ext == "" {
		ext = ".bin"
	}

	// Timestamp
	timestamp := time.Now().Unix()

	// Nombre único: timestamp_hash.ext
	return fmt.Sprintf("%d_%s%s", timestamp, hashStr, ext), nil
}

// isValidMimeType valida si el MIME type es permitido para el tipo de archivo
func (s *FileService) isValidMimeType(mimeType, fileType string) bool {
	validImageTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
		"image/gif",
	}

	validPDFTypes := []string{
		"application/pdf",
	}

	if fileType == "image" {
		for _, valid := range validImageTypes {
			if strings.EqualFold(mimeType, valid) {
				return true
			}
		}
		return false
	}

	if fileType == "pdf" {
		for _, valid := range validPDFTypes {
			if strings.EqualFold(mimeType, valid) {
				return true
			}
		}
		return false
	}

	return false
}

// GetFileURL devuelve la URL pública del archivo
func (s *FileService) GetFileURL(ctx context.Context, fileID int32) (string, error) {
	file, err := s.queries.GetStaticFile(ctx, fileID)
	if err != nil {
		return "", ErrFileNotFound
	}

	// Get file type
	if !file.FileType.Valid {
		return "", fmt.Errorf("tipo de archivo no válido")
	}

	// Build URL based on file type
	var fileURL string
	switch file.FileType.String {
	case "image":
		fileURL = filepath.Join(config.IMAGES_PATH, file.FileName)
	case "pdf":
		fileURL = filepath.Join(config.PDFS_PATH, file.FileName)
	default:
		return "", fmt.Errorf("tipo de archivo no soportado")
	}

	// Normalize path separators to forward slashes for URLs
	return filepath.ToSlash(fileURL), nil
}
