package web

// import (
	// "fmt"
	// "net/http"

	// "github.com/dirodriguezm/xmatch/service/internal/validator"
	// "github.com/gin-gonic/gin"
// )

// func (w *web) search(c *gin.context) {
	// ctx := c.request.context()

	// q := c.query("q")
	// v := validator.validator{}

	// v.checkfield(validator.matches(q, validator.searchrx), "q", "must be a valid search query")
	// if !v.valid() {
		// data := w.newtemplatedata(c)
		// data.validator = v
		// tmpl := "search.tmpl.html"
		// if c.request.referer() == w.getenv("base_url") {
			// tmpl = "home.tmpl.html"
		// }
		// w.render(c, http.statusunprocessableentity, tmpl, data)
		// return
	// }

	// data := w.newtemplatedata(ctx)
	// if err := w.render(c, http.statusok, "search.tmpl.html", data); err != nil {
		// w.servererror(c, fmt.errorf("failed to render search template: %w", err))
	// }
// }
