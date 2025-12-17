// Filename: actions/app.go
package actions

import (
	"sync"

	"onlyoffice/locales" // Your project's locales package
	"onlyoffice/models"  // Your project's models package

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/middleware/contenttype"
	"github.com/gobuffalo/middleware/forcessl"
	"github.com/gobuffalo/middleware/i18n"
	"github.com/gobuffalo/middleware/paramlogger"
	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
	"github.com/unrolled/secure"

	// --- SWAGGER IMPORTS ---
	swagger "github.com/swaggo/buffalo-swagger"
	"github.com/swaggo/buffalo-swagger/swaggerFiles"
	_ "onlyoffice/docs" // Points to the generated docs folder
	// -----------------------
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app     *buffalo.App
	appOnce sync.Once
	T       *i18n.Translator
)

// --- SWAGGER MAIN ANNOTATION ---
// @title Document Generation API
// @version 1.0
// @description This is an API for creating and managing documents from templates.
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @host localhost:3000
// @BasePath /api
// @schemes http
// -------------------------------
func App() *buffalo.App {
	appOnce.Do(func() {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			// --- UPDATED CORS CONFIGURATION FOR DEPLOYMENT ---
			PreWares: []buffalo.PreWare{
				cors.New(cors.Options{
					AllowedOrigins:   []string{"*"}, // Allows the frontend to connect from anywhere (e.g., Render)
					AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
					AllowedHeaders:   []string{"*"},
					AllowCredentials: true,
				}).Handler,
			},
			// -------------------------------------------------
			SessionName: "_onlyoffice_session",
		})

		// Automatically redirect to SSL in production
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		// Wraps each request in a database transaction.
		app.Use(popmw.Transaction(models.DB))

		// --- Setup Application Routes ---

		// Route for the home page (if you have one)
		app.GET("/", HomeHandler)

		// Route to serve the Swagger UI documentation
		app.GET("/swagger/{doc:.*}", swagger.WrapHandler(swaggerFiles.Handler))

		// Group all API endpoints under the "/api" prefix
		api := app.Group("/api")

		// === TEMPLATES API ===
		// GET /api/templates - Get a list of all available templates.
		api.GET("/templates", ListTemplates)
		// POST /api/templates - Upload a new template file.
		api.POST("/templates", CreateTemplate)

		// === DOCUMENTS API ===
		// POST /api/documents/init - Create a new document instance from a template.
		api.POST("/documents/init", InitDocument)
		// POST /api/documents/{document_id}/generate - Fill a document with data and save it.
		api.POST("/documents/{document_id}/generate", UpdateContent)
		// GET /api/documents/{document_id} - Get the saved JSON data for an existing document.
		api.GET("/documents/{document_id}", GetDocumentDetails)
		// GET /api/documents/{document_id}/file - Serve the physical .docx file.
		api.GET("/documents/{document_id}/file", GetDocumentFile)

	})

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS.
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}