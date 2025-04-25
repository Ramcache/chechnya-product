package handlers

import (
	"chechnya-product/internal/services"
	"chechnya-product/internal/utils"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AnnouncementHandlerInterface interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type AnnouncementHandler struct {
	service services.AnnouncementServiceInterface
	logger  *zap.Logger
}

func NewAnnouncementHandler(service services.AnnouncementServiceInterface, logger *zap.Logger) *AnnouncementHandler {
	return &AnnouncementHandler{service: service, logger: logger}
}

// Update
// @Summary Обновить объявление
// @Tags Объявления
// @Security BearerAuth
// @Param id path int true "ID"
// @Param input body map[string]string true "title, content"
// @Success 200 {string} string "Updated"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/admin/announcements/{id} [put]
func (h *AnnouncementHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var body struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" || body.Content == "" {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	if err := h.service.Update(id, body.Title, body.Content); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to update announcement")
		return
	}

	utils.JSONResponse(w, http.StatusOK, "Updated", nil)
}

// GetByID
// @Summary Получить объявление по ID
// @Tags Объявления
// @Param id path int true "ID"
// @Success 200 {object} models.Announcement
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/announcements/{id} [get]
func (h *AnnouncementHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	ann, err := h.service.GetByID(id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusNotFound, "Announcement not found")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "Announcement fetched", ann)
}

// GetAll
// @Summary Получить все объявления
// @Tags Объявления
// @Success 200 {array} models.Announcement
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/announcements [get]
func (h *AnnouncementHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	anns, err := h.service.GetAll()
	if err != nil {
		h.logger.Error("get announcements failed", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to get announcements")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "Announcements fetched", anns)
}

// Create
// @Summary Создать объявление
// @Tags Объявления
// @Security BearerAuth
// @Accept json
// @Param input body map[string]string true "title, content"
// @Success 201 {object} models.Announcement
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/admin/announcements [post]
func (h *AnnouncementHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" || body.Content == "" {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	ann, err := h.service.Create(body.Title, body.Content)
	if err != nil {
		h.logger.Error("create announcement failed", zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to create announcement")
		return
	}
	utils.JSONResponse(w, http.StatusCreated, "Announcement created", ann)
}

// Delete
// @Summary Удалить объявление
// @Tags Объявления
// @Security BearerAuth
// @Param id path int true "ID"
// @Success 200 {string} string "Deleted"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/admin/announcements/{id} [delete]
func (h *AnnouncementHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		h.logger.Error("delete announcement failed", zap.Int("id", id), zap.Error(err))
		utils.ErrorJSON(w, http.StatusInternalServerError, "Failed to delete")
		return
	}
	utils.JSONResponse(w, http.StatusOK, "Deleted", nil)
}
