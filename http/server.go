package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/willdobbins/notes"
	"fmt"
)

//Essentially just a controller which has a NoteManager attached.
type Server struct {
	Service notes.Service
}

//Simple health check
func (s Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

//index endpoint
func (s Server) ListNotes(c *gin.Context) {
	set, err := s.Service.All()
	if err != nil {
		log.Print(err)
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{"results": set})
}

//single note view handler
func (s Server) GetNote(c *gin.Context) {
	id := c.Param("id")
	idNumber, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Print(err)
	}

	note, err := s.Service.One(idNumber)
	if err != nil {
		log.Print(err)
	}
	c.HTML(http.StatusOK, "single.tmpl", note)
}

//notes POST handler
func (s Server) CreateNote(c *gin.Context) {
	var n notes.Note

	err := c.Bind(&n) //binding & form tags in the Notes struct tell it what form params to expect.
	if err != nil {
		log.Print(err)
	}
	_, err = s.Service.CreateNote(&n)
	if err != nil {
		log.Print(err)
	}

	c.Redirect(http.StatusMovedPermanently, "/notes/")
}

//notes DELETE handler
func (s Server) DeleteNote(c *gin.Context) {
	id := c.Param("id")
	idNumber, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Print(err)
	}

	err = s.Service.DeleteNote(idNumber)
	if err != nil {
		log.Print(err)
	}

	c.Redirect(http.StatusMovedPermanently, "/notes/")
}

func (s Server) UpdateNote(c *gin.Context) {
	id := c.Param("id")
	idNumber, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Print(err)
	}

	var n = new(notes.Note)
	err = c.Bind(n)
	if err != nil {
		log.Print(err)
	}

	n, err = s.Service.UpdateNote(idNumber, n)
	if err != nil {
		log.Print(err)
	}

	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/notes/%d", idNumber))
}
