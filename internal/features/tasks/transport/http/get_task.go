package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/nullphosphor/golang-todoapp/internal/core/logger"
	core_http_request "github.com/nullphosphor/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/nullphosphor/golang-todoapp/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

func (h *TaskHTTPHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTaskResponse(taskDomain)

	responseHandler.JSONResponse(response, http.StatusOK)
}
