// actions/documents.go
package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"onlyoffice/models"
	"onlyoffice/services"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

// FormData for updating document fields
type FormData struct {
	FIO       string `json:"fio"`
	Lavozim   string `json:"lavozim"`
	Oylik     string `json:"oylik"`
	Stavka    string `json:"stavka"`
	Username  string `json:"username"`
	OrderType string `json:"order_type"`
}

// Request struct for init
type InitDocumentRequest struct {
	TemplateID string `json:"template_id"`
}

// InitDocument creates a document copy.
func InitDocument(c buffalo.Context) error {
	reqData := &InitDocumentRequest{}
	// Standard JSON decoder
	if err := json.NewDecoder(c.Request().Body).Decode(reqData); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	// Manual check to ensure DB is connected
	if models.DB == nil { return c.Error(http.StatusInternalServerError, fmt.Errorf("DB nil")) }
	tx := models.DB

	template := &models.Template{}
	if err := tx.Find(template, reqData.TemplateID); err != nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("template not found"))
	}

	templateFilename := template.FileID.String() + ".docx"
	newDocID := uuid.Must(uuid.NewV4())
	docService := services.NewDocService()

	// New filename based on ID
	newFilename := newDocID.String() + ".docx"
	
	// Call CopyTemplate (matches your doc_service.go logic)
	if err := docService.CopyTemplate(templateFilename, newFilename); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	doc := &models.Document{
		ID:         newDocID,
		TemplateID: template.ID,
		Filename:   newFilename,
		Status:     "draft",
	}

	if err := tx.Create(doc); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	// Standard JSON response
	c.Response().Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response()).Encode(map[string]string{
		"document_id": newDocID.String(),
	})
}

// UpdateContent fills placeholders
func UpdateContent(c buffalo.Context) error {
	docID := c.Param("document_id")
	reqData := &FormData{}
	if err := json.NewDecoder(c.Request().Body).Decode(reqData); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	if models.DB == nil { return c.Error(http.StatusInternalServerError, fmt.Errorf("DB nil")) }
	tx := models.DB

	doc := &models.Document{}
	if err := tx.Find(doc, docID); err != nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("document not found"))
	}

	replacements := map[string]interface{}{
		"fio":        reqData.FIO,
		"lavozim":    reqData.Lavozim,
		"oylik":      reqData.Oylik,
		"stavka":     reqData.Stavka,
		"username":   reqData.Username,
		"order_type": reqData.OrderType,
	}
	docService := services.NewDocService()
	
	// Calls ReplaceInDoc in your doc_service.go
	if err := docService.ReplaceInDoc(doc.Filename, replacements); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	doc.FIO = reqData.FIO
	doc.Lavozim = reqData.Lavozim
	doc.Oylik = reqData.Oylik
	doc.Stavka = reqData.Stavka
	doc.Username = reqData.Username
	doc.OrderType = reqData.OrderType
	doc.Status = "filled"
	if err := tx.Update(doc); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response()).Encode(map[string]string{
		"status": "success",
	})
}

// GetDocumentDetails returns JSON data about doc
func GetDocumentDetails(c buffalo.Context) error {
	docID := c.Param("document_id")
	if models.DB == nil { return c.Error(http.StatusInternalServerError, fmt.Errorf("DB nil")) }
	tx := models.DB
	
	doc := &models.Document{}
	if err := tx.Find(doc, docID); err != nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("document not found"))
	}

	c.Response().Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response()).Encode(doc)
}

// =======================================================
// CRITICAL FUNCTION: Serves the file to ONLYOFFICE
// =======================================================
func GetDocumentFile(c buffalo.Context) error {
	docID := c.Param("document_id")
	if models.DB == nil { return c.Error(http.StatusInternalServerError, fmt.Errorf("DB nil")) }
	tx := models.DB

	doc := &models.Document{}
	if err := tx.Find(doc, docID); err != nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("document not found"))
	}

	docService := services.NewDocService()
	// Combines the shared_storage/documents path with filename
	filePath := filepath.Join(docService.DocDir, doc.Filename)
	
	file, err := os.Open(filePath)
	if err != nil {
		// Log this so we know if path is wrong
		fmt.Printf("File not found at: %s\n", filePath) 
		return c.Error(http.StatusNotFound, fmt.Errorf("file open error: %v", err))
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return c.Error(http.StatusInternalServerError, fmt.Errorf("file stat error: %v", err))
	}

	// This allows the browser/ONLYOFFICE to download it
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+doc.Filename)
	http.ServeContent(c.Response(), c.Request(), doc.Filename, fileInfo.ModTime(), file)
	return nil
}