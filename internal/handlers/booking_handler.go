package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/Lzrb0x/smartBookingGoApi/internal/dtos"
	"github.com/Lzrb0x/smartBookingGoApi/internal/models"
	"github.com/Lzrb0x/smartBookingGoApi/internal/repositories"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingRepo             *repositories.BookingRepository
	employeeRepo            *repositories.EmployeeRepository
	ownerRepo               *repositories.OwnerRepository
	barbershopRepo          *repositories.BarbershopRepository
	barbershopServiceRepo   *repositories.BarbershopServiceRepository
	serviceEmployeeRepo     *repositories.ServiceEmployeeRepository
	workingHourRepo         *repositories.EmployeeWorkingHourRepository
	workingHourOverrideRepo *repositories.EmployeeWorkingHourOverrideRepository
}

func NewBookingHandler(
	bookingRepo *repositories.BookingRepository,
	employeeRepo *repositories.EmployeeRepository,
	ownerRepo *repositories.OwnerRepository,
	barbershopRepo *repositories.BarbershopRepository,
	barbershopServiceRepo *repositories.BarbershopServiceRepository,
	serviceEmployeeRepo *repositories.ServiceEmployeeRepository,
	workingHourRepo *repositories.EmployeeWorkingHourRepository,
	workingHourOverrideRepo *repositories.EmployeeWorkingHourOverrideRepository,
) *BookingHandler {
	return &BookingHandler{
		bookingRepo:             bookingRepo,
		employeeRepo:            employeeRepo,
		ownerRepo:               ownerRepo,
		barbershopRepo:          barbershopRepo,
		barbershopServiceRepo:   barbershopServiceRepo,
		serviceEmployeeRepo:     serviceEmployeeRepo,
		workingHourRepo:         workingHourRepo,
		workingHourOverrideRepo: workingHourOverrideRepo,
	}
}

func (h *BookingHandler) GetAvailability(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	employeeID, err := strconv.ParseInt(c.Param("employeeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id inválido"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date é obrigatório (formato: YYYY-MM-DD)"})
		return
	}

	serviceIDStr := c.Query("barbershop_service_id")
	if serviceIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_service_id é obrigatório"})
		return
	}
	serviceID, err := strconv.ParseInt(serviceIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_service_id inválido"})
		return
	}

	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date inválido (formato: YYYY-MM-DD)"})
		return
	}

	// Verify employee exists in barbershop
	exists, err := h.employeeRepo.ExistsInBarbershop(c.Request.Context(), employeeID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não encontrado nesta barbearia"})
		return
	}

	service, err := h.barbershopServiceRepo.FindByID(c.Request.Context(), barbershopID, serviceID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "serviço não encontrado para esta barbearia"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	assigned, err := h.serviceEmployeeRepo.IsEmployeeAssignedToService(c.Request.Context(), employeeID, serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !assigned {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não está atribuído a este serviço"})
		return
	}

	filterStartStr := c.Query("start_time")
	filterEndStr := c.Query("end_time")
	var filterStartSec, filterEndSec int
	if filterStartStr != "" || filterEndStr != "" {
		if filterStartStr == "" || filterEndStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_time e end_time devem ser informados juntos"})
			return
		}
		filterStart, err := time.ParseInLocation("15:04:05", filterStartStr, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_time inválido (formato: HH:MM:SS)"})
			return
		}
		filterEnd, err := time.ParseInLocation("15:04:05", filterEndStr, time.Local)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_time inválido (formato: HH:MM:SS)"})
			return
		}
		filterStartSec = timeToSeconds(filterStart)
		filterEndSec = timeToSeconds(filterEnd)
		if filterEndSec <= filterStartSec {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_time deve ser após start_time"})
			return
		}
	}

	windows, err := h.buildWindowsForDate(c.Request.Context(), employeeID, barbershopID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if filterStartStr != "" {
		windows = intersectWindows(windows, filterStartSec, filterEndSec)
	}

	if isSameDate(date, time.Now()) {
		nowSec := nowSeconds()
		windows = intersectWindows(windows, nowSec, 24*60*60)
	}

	sort.Slice(windows, func(i, j int) bool {
		return windows[i].startSec < windows[j].startSec
	})

	bookings, err := h.bookingRepo.ListByEmployeeAndDate(c.Request.Context(), employeeID, date)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookingIntervals := make([]timeWindow, 0, len(bookings))
	for _, booking := range bookings {
		bookingIntervals = append(bookingIntervals, timeWindow{
			startSec: timeToSeconds(booking.StartTime),
			endSec:   timeToSeconds(booking.EndTime),
		})
	}

	durationSec := service.Duration * 60
	intervalSec := 15 * 60
	slots := make([]string, 0)

	for _, window := range windows {
		start := alignUp(window.startSec, intervalSec)
		for start+durationSec <= window.endSec {
			end := start + durationSec
			if !overlapsAny(start, end, bookingIntervals) {
				slots = append(slots, secondsToClock(start))
			}
			start += intervalSec
		}
	}

	c.JSON(http.StatusOK, dtos.AvailabilityResponse{
		Date:                date.Format("2006-01-02"),
		EmployeeID:          employeeID,
		BarbershopServiceID: serviceID,
		ServiceDuration:     service.Duration,
		Slots:               slots,
	})
}

func (h *BookingHandler) GetAllByBarbershop(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	bookings, err := h.bookingRepo.ListByBarbershop(c.Request.Context(), barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookingsToResponse(bookings))
}

func (h *BookingHandler) GetAllByEmployee(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	employeeID, err := strconv.ParseInt(c.Param("employeeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id inválido"})
		return
	}

	exists, err := h.employeeRepo.ExistsInBarbershop(c.Request.Context(), employeeID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não encontrado nesta barbearia"})
		return
	}

	bookings, err := h.bookingRepo.ListByEmployeeAndBarbershop(c.Request.Context(), employeeID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookingsToResponse(bookings))
}

func (h *BookingHandler) GetProfessionalDashboard(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	authenticatedUserID, ok := authenticatedUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}
	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date inválido (formato: YYYY-MM-DD)"})
		return
	}

	filterEmployeeID, ok := h.authorizedProfessionalEmployeeFilter(c, authenticatedUserID, barbershopID)
	if !ok {
		return
	}

	bookings, err := h.bookingRepo.ListProfessionalDashboard(c.Request.Context(), barbershopID, date, filterEmployeeID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var estimatedRevenue float64
	var totalDuration int
	for _, booking := range bookings {
		estimatedRevenue += booking.ServicePrice
		totalDuration += booking.ServiceDuration
	}

	c.JSON(http.StatusOK, dtos.ProfessionalBookingDashboardResponse{
		Date:             date.Format("2006-01-02"),
		EmployeeID:       filterEmployeeID,
		EstimatedRevenue: estimatedRevenue,
		AppointmentCount: len(bookings),
		TotalDuration:    totalDuration,
		Bookings:         dtos.FromProfessionalBookingItemModels(bookings),
	})
}

func (h *BookingHandler) authorizedProfessionalEmployeeFilter(c *gin.Context, userID, barbershopID int64) (*int64, bool) {
	barbershop, err := h.barbershopRepo.FindByID(c.Request.Context(), barbershopID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "barbearia não encontrada"})
		return nil, false
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, false
	}

	isOwner := false
	owner, err := h.ownerRepo.FindByUserID(c.Request.Context(), userID)
	if err == nil && owner.ID == barbershop.OwnerID {
		isOwner = true
	} else if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, false
	}

	var ownEmployeeID *int64
	employee, err := h.employeeRepo.FindByUserAndBarbershop(c.Request.Context(), userID, barbershopID)
	if err == nil {
		ownEmployeeID = &employee.ID
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, false
	}

	if !isOwner && ownEmployeeID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "você não tem acesso profissional a esta barbearia"})
		return nil, false
	}

	employeeFilter := c.Query("employee_id")
	if employeeFilter == "" || employeeFilter == "all" {
		if isOwner {
			return nil, true
		}
		return ownEmployeeID, true
	}

	requestedEmployeeID, err := strconv.ParseInt(employeeFilter, 10, 64)
	if err != nil || requestedEmployeeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id inválido"})
		return nil, false
	}

	exists, err := h.employeeRepo.ExistsInBarbershop(c.Request.Context(), requestedEmployeeID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, false
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não encontrado nesta barbearia"})
		return nil, false
	}

	if !isOwner && (ownEmployeeID == nil || requestedEmployeeID != *ownEmployeeID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "funcionário só pode ver os próprios agendamentos"})
		return nil, false
	}

	return &requestedEmployeeID, true
}

func (h *BookingHandler) GetUserDashboard(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}

	currentBooking, err := h.bookingRepo.FindCurrentByCustomer(c.Request.Context(), customerID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	recentBookings, err := h.bookingRepo.ListRecentByCustomer(c.Request.Context(), customerID, 5)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	recentBarbershops, err := h.bookingRepo.ListRecentBarbershopsByCustomer(c.Request.Context(), customerID, 3)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var currentResponse *dtos.BookingDashboardItemResponse
	if currentBooking != nil {
		mapped := dtos.FromBookingDashboardItemModel(currentBooking)
		currentResponse = &mapped
	}

	c.JSON(http.StatusOK, dtos.UserDashboardResponse{
		CurrentBooking:    currentResponse,
		RecentBookings:    dtos.FromBookingDashboardItemModels(recentBookings),
		RecentBarbershops: dtos.FromRecentBarbershopModels(recentBarbershops),
	})
}

func (h *BookingHandler) GetRecentByCustomer(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}

	bookings, err := h.bookingRepo.ListRecentByCustomer(c.Request.Context(), customerID, parseLimit(c, 5, 20))
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.FromBookingDashboardItemModels(bookings))
}

func (h *BookingHandler) GetCurrentByCustomer(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}

	booking, err := h.bookingRepo.FindCurrentByCustomer(c.Request.Context(), customerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, nil)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.FromBookingDashboardItemModel(booking))
}

func (h *BookingHandler) GetRecentBarbershopsByCustomer(c *gin.Context) {
	customerID, ok := parseCustomerID(c)
	if !ok {
		return
	}

	barbershops, err := h.bookingRepo.ListRecentBarbershopsByCustomer(c.Request.Context(), customerID, parseLimit(c, 3, 10))
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dtos.FromRecentBarbershopModels(barbershops))
}

func (h *BookingHandler) Create(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	var req dtos.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.ParseInLocation("2006-01-02", req.Date, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date inválido (formato: YYYY-MM-DD)"})
		return
	}

	startTime, err := time.ParseInLocation("15:04:05", req.StartTime, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_time inválido (formato: HH:MM:SS)"})
		return
	}

	if !isAlignedToInterval(startTime, 15) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_time deve estar em intervalos de 15 minutos"})
		return
	}

	// Verify employee exists in barbershop
	exists, err := h.employeeRepo.ExistsInBarbershop(c.Request.Context(), req.EmployeeID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não encontrado nesta barbearia"})
		return
	}

	service, err := h.barbershopServiceRepo.FindByID(c.Request.Context(), barbershopID, req.BarbershopServiceID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "serviço não encontrado para esta barbearia"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	assigned, err := h.serviceEmployeeRepo.IsEmployeeAssignedToService(c.Request.Context(), req.EmployeeID, req.BarbershopServiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !assigned {
		c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não está atribuído a este serviço"})
		return
	}

	startSec := timeToSeconds(startTime)
	endSec := startSec + (service.Duration * 60)
	if endSec > 24*60*60 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "horário inválido para a duração do serviço"})
		return
	}

	if isSameDate(date, time.Now()) {
		if startSec < nowSeconds() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_time não pode estar no passado"})
			return
		}
	}
	if date.Before(beginningOfDay(time.Now())) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date não pode estar no passado"})
		return
	}

	windows, err := h.buildWindowsForDate(c.Request.Context(), req.EmployeeID, barbershopID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sort.Slice(windows, func(i, j int) bool {
		return windows[i].startSec < windows[j].startSec
	})
	if !slotFitsWindows(startSec, endSec, windows) {
		c.JSON(http.StatusConflict, gin.H{"error": "horário fora do expediente"})
		return
	}

	booking := &models.Booking{
		CustomerID:          req.CustomerID,
		EmployeeID:          req.EmployeeID,
		BarbershopID:        barbershopID,
		BarbershopServiceID: req.BarbershopServiceID,
		Date:                date,
		StartTime:           startTime,
		EndTime:             secondsToTime(endSec, startTime),
	}

	tx, err := h.bookingRepo.BeginTx(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	if err := h.bookingRepo.LockEmployee(c.Request.Context(), tx, req.EmployeeID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "funcionário não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	overlap, err := h.bookingRepo.HasOverlapTx(c.Request.Context(), tx, req.EmployeeID, date, booking.StartTime, booking.EndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if overlap {
		c.JSON(http.StatusConflict, gin.H{"error": "horário indisponível"})
		return
	}

	if err := h.bookingRepo.CreateTx(c.Request.Context(), tx, booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dtos.FromBookingModel(booking))
}

func (h *BookingHandler) Cancel(c *gin.Context) {
	barbershopID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barbershop_id inválido"})
		return
	}

	bookingID, err := strconv.ParseInt(c.Param("bookingId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking_id inválido"})
		return
	}

	authenticatedUserID, ok := authenticatedUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	details, err := h.bookingRepo.FindCancellationDetails(c.Request.Context(), bookingID, barbershopID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "agendamento não encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if details.CustomerID != authenticatedUserID &&
		details.EmployeeUserID != authenticatedUserID &&
		details.OwnerUserID != authenticatedUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "você não tem permissão para cancelar este agendamento"})
		return
	}

	deleted, err := h.bookingRepo.Delete(c.Request.Context(), bookingID, barbershopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "agendamento não encontrado"})
		return
	}

	c.Status(http.StatusNoContent)
}

func authenticatedUserIDFromContext(c *gin.Context) (int64, bool) {
	value, exists := c.Get("authenticated_user_id")
	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)
	return userID, ok && userID > 0
}

func parseCustomerID(c *gin.Context) (int64, bool) {
	customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id inválido"})
		return 0, false
	}
	return customerID, true
}

func parseLimit(c *gin.Context, defaultLimit, maxLimit int) int {
	limitStr := c.Query("limit")
	if limitStr == "" {
		return defaultLimit
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

type timeWindow struct {
	startSec int
	endSec   int
}

func (h *BookingHandler) buildWindowsForDate(ctx context.Context, employeeID, barbershopID int64, date time.Time) ([]timeWindow, error) {
	override, err := h.workingHourOverrideRepo.FindByEmployeeAndDate(ctx, employeeID, date)
	if err == nil {
		if override.IsDayOff {
			return []timeWindow{}, nil
		}
		if override.StartTime == nil || override.EndTime == nil {
			return []timeWindow{}, nil
		}
		startSec := timeToSeconds(*override.StartTime)
		endSec := timeToSeconds(*override.EndTime)
		if endSec <= startSec {
			return []timeWindow{}, nil
		}
		return []timeWindow{{startSec: startSec, endSec: endSec}}, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	whs, err := h.workingHourRepo.FindByEmployeeAndBarbershop(ctx, employeeID, barbershopID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	dayOfWeek := int(date.Weekday())
	windows := make([]timeWindow, 0)
	hasDayOff := false
	for _, wh := range whs {
		if wh.DayOfWeek != dayOfWeek {
			continue
		}
		if wh.IsDayOff {
			hasDayOff = true
			continue
		}
		if wh.StartTime == nil || wh.EndTime == nil {
			continue
		}
		startSec := timeToSeconds(*wh.StartTime)
		endSec := timeToSeconds(*wh.EndTime)
		if endSec <= startSec {
			continue
		}
		windows = append(windows, timeWindow{startSec: startSec, endSec: endSec})
	}

	if len(windows) == 0 && hasDayOff {
		return []timeWindow{}, nil
	}

	return windows, nil
}

func intersectWindows(windows []timeWindow, startSec, endSec int) []timeWindow {
	if len(windows) == 0 {
		return windows
	}
	bounded := make([]timeWindow, 0, len(windows))
	for _, window := range windows {
		start := window.startSec
		end := window.endSec
		if start < startSec {
			start = startSec
		}
		if end > endSec {
			end = endSec
		}
		if end <= start {
			continue
		}
		bounded = append(bounded, timeWindow{startSec: start, endSec: end})
	}
	return bounded
}

func slotFitsWindows(startSec, endSec int, windows []timeWindow) bool {
	for _, window := range windows {
		if startSec >= window.startSec && endSec <= window.endSec {
			return true
		}
	}
	return false
}

func overlapsAny(startSec, endSec int, windows []timeWindow) bool {
	for _, window := range windows {
		if startSec < window.endSec && endSec > window.startSec {
			return true
		}
	}
	return false
}

func timeToSeconds(t time.Time) int {
	return t.Hour()*3600 + t.Minute()*60 + t.Second()
}

func secondsToClock(totalSec int) string {
	hours := totalSec / 3600
	minutes := (totalSec % 3600) / 60
	seconds := totalSec % 60
	return formatClock(hours, minutes, seconds)
}

func secondsToTime(totalSec int, base time.Time) time.Time {
	hours := totalSec / 3600
	minutes := (totalSec % 3600) / 60
	seconds := totalSec % 60
	location := base.Location()
	return time.Date(base.Year(), base.Month(), base.Day(), hours, minutes, seconds, 0, location)
}

func formatClock(hours, minutes, seconds int) string {
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func alignUp(value, interval int) int {
	if value%interval == 0 {
		return value
	}
	return value + (interval - (value % interval))
}

func nowSeconds() int {
	now := time.Now()
	return now.Hour()*3600 + now.Minute()*60 + now.Second()
}

func isAlignedToInterval(t time.Time, intervalMinutes int) bool {
	if t.Second() != 0 {
		return false
	}
	return (t.Minute()%(intervalMinutes) == 0)
}

func isSameDate(a time.Time, b time.Time) bool {
	a = a.In(time.Local)
	b = b.In(time.Local)
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

func beginningOfDay(t time.Time) time.Time {
	local := t.In(time.Local)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.Local)
}

func bookingsToResponse(bookings []models.Booking) []dtos.BookingResponse {
	response := make([]dtos.BookingResponse, 0, len(bookings))
	for i := range bookings {
		response = append(response, dtos.FromBookingModel(&bookings[i]))
	}
	return response
}
