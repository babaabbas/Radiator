package product

import (
	"encoding/json"
	"fmt"
	"mma_api/internal/storage/postgres"
	"mma_api/internal/types"
	"mma_api/internal/utils/response"
	"net/http"
	"strconv"
	"strings"
)

func GetProductsHandler(storage *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := response.GeneralError(http.ErrNotSupported)
			_ = response.WriteJson(w, http.StatusMethodNotAllowed, resp)
			return
		}

		products, err := storage.GetProducts()
		if err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusInternalServerError, resp)
			return
		}

		_ = response.WriteJson(w, http.StatusOK, map[string]interface{}{
			"custom_status": response.Status_Ok,
			"data":          products,
		})
	}
}

func GetProductByIDHandler(storage *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := response.GeneralError(http.ErrNotSupported)
			_ = response.WriteJson(w, http.StatusMethodNotAllowed, resp)
			return
		}

		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			resp := response.GeneralError(fmt.Errorf("missing product ID"))
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		idStr := pathParts[3]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			resp := response.GeneralError(fmt.Errorf("invalid product ID: %w", err))
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		product, err := storage.GetProductById(id)
		if err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusNotFound, resp)
			return
		}

		_ = response.WriteJson(w, http.StatusOK, map[string]interface{}{
			"custom_status": response.Status_Ok,
			"data":          product,
		})
	}
}

func CreateProductHandler(storage *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := response.GeneralError(http.ErrNotSupported)
			_ = response.WriteJson(w, http.StatusMethodNotAllowed, resp)
			return
		}

		var product types.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		if product.Name == "" || product.Unit == "" {
			resp := response.GeneralError(fmt.Errorf("name and unit are required"))
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		newProduct, err := storage.CreateProduct(product.Name, product.Description, product.Category, product.Unit)
		if err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusInternalServerError, resp)
			return
		}

		_ = response.WriteJson(w, http.StatusCreated, map[string]interface{}{
			"custom_status": response.Status_Ok,
			"data":          newProduct,
		})
	}
}

type BoMCreateRequest struct {
	ComponentID   int     `json:"component_id"`
	Quantity      float64 `json:"quantity"`
	OperationName string  `json:"operation_name,omitempty"`
}

func CreateBoMHandler(storage *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := response.GeneralError(http.ErrNotSupported)
			_ = response.WriteJson(w, http.StatusMethodNotAllowed, resp)
			return
		}

		// URL: /api/products/{id}/bom
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 5 || pathParts[4] != "bom" {
			resp := response.GeneralError(fmt.Errorf("invalid URL"))
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		// Extract product ID
		productID, err := strconv.Atoi(pathParts[3])
		if err != nil {
			resp := response.GeneralError(fmt.Errorf("invalid product ID: %w", err))
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		// Decode request body
		var req struct {
			ComponentID   int     `json:"component_id"`
			Quantity      float64 `json:"quantity"`
			OperationName string  `json:"operation_name,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		if req.ComponentID == 0 || req.Quantity <= 0 {
			resp := response.GeneralError(fmt.Errorf("component_id and quantity are required"))
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		// Create BoM entry
		bom, err := storage.CreateBoM(productID, req.ComponentID, req.Quantity, req.OperationName)
		if err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusInternalServerError, resp)
			return
		}

		_ = response.WriteJson(w, http.StatusCreated, map[string]interface{}{
			"custom_status": response.Status_Ok,
			"data":          bom,
		})
	}
}

func GetBoMHandler(storage *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := response.GeneralError(http.ErrNotSupported)
			_ = response.WriteJson(w, http.StatusMethodNotAllowed, resp)
			return
		}

		// URL: /api/products/{id}/bom
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) != 5 || pathParts[4] != "bom" {
			resp := response.GeneralError(fmt.Errorf("invalid URL"))
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		// Extract product ID
		productID, err := strconv.Atoi(pathParts[3])
		if err != nil {
			resp := response.GeneralError(fmt.Errorf("invalid product ID: %w", err))
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		// Fetch BoM entries from DB
		boms, err := storage.GetBoM(productID)
		if err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusInternalServerError, resp)
			return
		}

		// Send JSON response
		_ = response.WriteJson(w, http.StatusOK, map[string]interface{}{
			"custom_status": response.Status_Ok,
			"data":          boms,
		})
	}
}
