package actions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"onlyoffice/models"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

// ListTemplates returns a list of all available document templates.
func ListTemplates(c buffalo.Context) error {
	// SAFETY CHECK: Ensure the database is actually connected
	if models.DB == nil {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("fatal error: database connection (models.DB) is nil"))
	}

	templates := &models.Templates{}
	tx := models.DB 

	err := tx.Order("name asc, version desc").All(templates)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// BYPASSING BUFFALO RENDERER: Using standard Go JSON encoder
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(templates)
}

// CreateTemplate handles uploading a new template file.
func CreateTemplate(c buffalo.Context) error {
	templateName := c.Param("name")
	if templateName == "" {
		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(map[string]string{"error": "'name' field is required"})
	}

	file, err := c.File("template_file")
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	defer file.Close()

	fileID := uuid.Must(uuid.NewV4())
	filename := fileID.String() + ".docx"

	root, _ := os.Getwd()
	templatesDir := filepath.Join(root, "templates")
	
	// Ensure the directory exists
	if err := os.MkdirAll(templatesDir, os.ModePerm); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	filePath := filepath.Join(templatesDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if models.DB == nil {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("database not connected"))
	}
	tx := models.DB
	template := &models.Template{
		ID:      uuid.Must(uuid.NewV4()),
		Name:    templateName,
		FileID:  fileID,
		Version: 1,
	}

	if err := tx.Create(template); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// BYPASSING BUFFALO RENDERER
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(template)
}