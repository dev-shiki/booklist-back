package controllers

import (
	"booklist-back/database"
	"booklist-back/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		log.Printf("Error binding JSON: %v", err)
		log.Printf("Received data: %+v", c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pubDate, err := time.Parse("2006-01-02", book.PublicationDate)
	if err != nil {
		errorMessage := fmt.Sprintf("Invalid date format. Expected YYYY-MM-DD, but received: %s", book.PublicationDate)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	result, err := database.DB.Exec("INSERT INTO books (title, author, publication_date, publisher, number_of_pages, category_id) VALUES (?, ?, ?, ?, ?, ?)",
		book.Title, book.Author, pubDate, book.Publisher, book.NumberOfPages, book.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	book.ID = int(id)
	c.JSON(http.StatusCreated, book)
}

func ListBooks(c *gin.Context) {
	query := "SELECT * FROM books WHERE 1=1"
	var args []interface{}

	// Filter by search text (title, author, or publisher)
	if search := c.Query("search"); search != "" {
		query += " AND (title LIKE ? OR author LIKE ? OR publisher LIKE ?)"
		searchParam := "%" + search + "%"
		args = append(args, searchParam, searchParam, searchParam)
	}

	// Filter by category
	if category := c.Query("category_id"); category != "" {
		query += " AND category_id = ?"
		args = append(args, category)
	}

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		query += " AND publication_date >= ?"
		args = append(args, startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query += " AND publication_date <= ?"
		args = append(args, endDate)
	}

	log.Printf("Executing query: %s with args: %v", query, args)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Database query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		var pubDate []uint8
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &pubDate, &book.Publisher, &book.NumberOfPages, &book.CategoryID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		book.PublicationDate = string(pubDate)
		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}

func UpdateBook(c *gin.Context) {
	id := c.Param("id")

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	pubDate, err := time.Parse("2006-01-02", book.PublicationDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	_, err = database.DB.Exec("UPDATE books SET title=?, author=?, publication_date=?, publisher=?, number_of_pages=?, category_id=? WHERE id=?",
		book.Title, book.Author, pubDate, book.Publisher, book.NumberOfPages, book.CategoryID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	_, err := database.DB.Exec("DELETE FROM books WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
