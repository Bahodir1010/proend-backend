package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lukasjarosch/go-docx"
)

// DocService handles file operations for templates and documents.
type DocService struct {
	// The directory where master templates are stored.
	TemplateDir string
	// The directory where generated documents are stored (this will be the shared volume).
	DocDir string
}

// NewDocService creates a new DocService with paths configured for the shared volume.
func NewDocService() *DocService {
    // Check if we are setting a specific storage path (Production)
    rootPath := os.Getenv("STORAGE_PATH")
    
    // Default to local folder if variable is empty (Local development)
    if rootPath == "" {
        root, _ := os.Getwd()
        rootPath = root
    }

    return &DocService{
        TemplateDir: filepath.Join(rootPath, "templates"),
        DocDir:      filepath.Join(rootPath, "documents"),
    }
}

// CopyTemplate copies a template file to the documents folder with a specific new name.
// It now correctly returns only an error.
func (s *DocService) CopyTemplate(templateFilename, newDocFilename string) error {
	sourcePath := filepath.Join(s.TemplateDir, templateFilename)
	destinationPath := filepath.Join(s.DocDir, newDocFilename)

	// Ensure the destination directory exists
	if err := os.MkdirAll(s.DocDir, os.ModePerm); err != nil {
		return fmt.Errorf("could not create destination directory %s: %w", s.DocDir, err)
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("could not open source template %s: %w", sourcePath, err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("could not create destination document %s: %w", destinationPath, err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy template content: %w", err)
	}

	return nil
}

// ReplaceInDoc replaces placeholders by writing to a temporary file first, then replacing the original.
// It is updated to use the correct DocDir path.
func (s *DocService) ReplaceInDoc(docFilename string, replacements map[string]interface{}) error {
	docPath := filepath.Join(s.DocDir, docFilename) // Uses the correct DocDir

	doc, err := docx.Open(docPath)
	if err != nil {
		return fmt.Errorf("failed to open doc: %v", err)
	}
	defer doc.Close()

	err = doc.ReplaceAll(replacements)
	if err != nil {
		return fmt.Errorf("failed to replace: %v", err)
	}

	tempFile, err := os.CreateTemp(s.DocDir, "doc-*.docx")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	err = doc.Write(tempFile)
	if err != nil {
		return fmt.Errorf("failed to write to temp file: %v", err)
	}
	
	tempFile.Close()
	doc.Close()

	err = os.Rename(tempFile.Name(), docPath)
	if err != nil {
		return fmt.Errorf("failed to replace original file with temp file: %v", err)
	}

	return nil
}