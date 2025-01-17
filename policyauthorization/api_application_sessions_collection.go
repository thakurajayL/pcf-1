// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0
//

/*
 * Npcf_PolicyAuthorization Service API
 *
 * This is the Policy Authorization Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package policyauthorization

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/omec-project/http_wrapper"
	"github.com/omec-project/openapi"
	"github.com/omec-project/openapi/models"
	"github.com/omec-project/pcf/logger"
	"github.com/omec-project/pcf/producer"
	"github.com/omec-project/pcf/util"
)

// HTTPPostAppSessions - Creates a new Individual Application Session Context resource
func HTTPPostAppSessions(c *gin.Context) {
	var appSessionContext models.AppSessionContext

	requestBody, err := c.GetRawData()
	if err != nil {
		problemDetail := models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		logger.PolicyAuthorizationlog.Errorf("Get Request Body error: %+v", err)
		c.JSON(http.StatusInternalServerError, problemDetail)
		return
	}

	err = openapi.Deserialize(&appSessionContext, requestBody, "application/json")
	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.PolicyAuthorizationlog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	ascReqData := appSessionContext.AscReqData
	if ascReqData == nil || ascReqData.SuppFeat == "" || ascReqData.NotifUri == "" {
		// Check Mandatory IEs
		rsp := util.GetProblemDetail("Errorneous/Missing Mandotory IE", util.ERROR_INITIAL_PARAMETERS)
		logger.PolicyAuthorizationlog.Errorln(rsp.Detail)
		c.JSON(int(rsp.Status), rsp)
		return
	}

	req := http_wrapper.NewRequest(c.Request, appSessionContext)
	rsp := producer.HandlePostAppSessionsContext(req)

	for key, val := range rsp.Header {
		c.Header(key, val[0])
	}
	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.PolicyAuthorizationlog.Errorln(err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	} else {
		c.Data(rsp.Status, "application/json", responseBody)
	}
}
