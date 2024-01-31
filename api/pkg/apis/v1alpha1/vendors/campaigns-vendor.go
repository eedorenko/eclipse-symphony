/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 * SPDX-License-Identifier: MIT
 */

package vendors

import (
	"encoding/json"

	"github.com/eclipse-symphony/symphony/api/pkg/apis/v1alpha1/managers/campaigns"
	"github.com/eclipse-symphony/symphony/api/pkg/apis/v1alpha1/model"
	"github.com/eclipse-symphony/symphony/api/pkg/apis/v1alpha1/utils"
	"github.com/eclipse-symphony/symphony/coa/pkg/apis/v1alpha2"
	"github.com/eclipse-symphony/symphony/coa/pkg/apis/v1alpha2/managers"
	"github.com/eclipse-symphony/symphony/coa/pkg/apis/v1alpha2/observability"
	observ_utils "github.com/eclipse-symphony/symphony/coa/pkg/apis/v1alpha2/observability/utils"
	"github.com/eclipse-symphony/symphony/coa/pkg/apis/v1alpha2/providers"
	"github.com/eclipse-symphony/symphony/coa/pkg/apis/v1alpha2/providers/pubsub"
	"github.com/eclipse-symphony/symphony/coa/pkg/apis/v1alpha2/vendors"
	"github.com/eclipse-symphony/symphony/coa/pkg/logger"
	"github.com/valyala/fasthttp"
)

var cLog = logger.NewLogger("coa.runtime")

type CampaignsVendor struct {
	vendors.Vendor
	CampaignsManager *campaigns.CampaignsManager
}

func (o *CampaignsVendor) GetInfo() vendors.VendorInfo {
	return vendors.VendorInfo{
		Version:  o.Vendor.Version,
		Name:     "Campaigns",
		Producer: "Microsoft",
	}
}

func (e *CampaignsVendor) Init(config vendors.VendorConfig, factories []managers.IManagerFactroy, providers map[string]map[string]providers.IProvider, pubsubProvider pubsub.IPubSubProvider) error {
	err := e.Vendor.Init(config, factories, providers, pubsubProvider)
	if err != nil {
		return err
	}
	for _, m := range e.Managers {
		if c, ok := m.(*campaigns.CampaignsManager); ok {
			e.CampaignsManager = c
		}
	}
	if e.CampaignsManager == nil {
		return v1alpha2.NewCOAError(nil, "campaigns manager is not supplied", v1alpha2.MissingConfig)
	}
	return nil
}

func (o *CampaignsVendor) GetEndpoints() []v1alpha2.Endpoint {
	route := "campaigns"
	if o.Route != "" {
		route = o.Route
	}
	return []v1alpha2.Endpoint{
		{
			Methods:    []string{fasthttp.MethodGet, fasthttp.MethodPost, fasthttp.MethodDelete},
			Route:      route,
			Version:    o.Version,
			Handler:    o.onCampaigns,
			Parameters: []string{"name?"},
		},
	}
}

func (c *CampaignsVendor) onCampaigns(request v1alpha2.COARequest) v1alpha2.COAResponse {
	pCtx, span := observability.StartSpan("Campaigns Vendor", request.Context, &map[string]string{
		"method": "onCampaigns",
	})
	defer span.End()
	cLog.Infof("V (Campaigns): onCampaigns, method: %s, traceId: %s", string(request.Method), span.SpanContext().TraceID().String())

	switch request.Method {
	case fasthttp.MethodGet:
		ctx, span := observability.StartSpan("onCampaigns-GET", pCtx, nil)
		id := request.Parameters["__name"]
		var err error
		var state interface{}
		isArray := false
		if id == "" {
			state, err = c.CampaignsManager.ListSpec(ctx)
			isArray = true
		} else {
			state, err = c.CampaignsManager.GetSpec(ctx, id)
		}
		if err != nil {
			cLog.Infof("V (Campaigns): onCampaigns failed - %s, traceId: %s", err.Error(), span.SpanContext().TraceID().String())
			return observ_utils.CloseSpanWithCOAResponse(span, v1alpha2.COAResponse{
				State: v1alpha2.InternalError,
				Body:  []byte(err.Error()),
			})
		}
		jData, _ := utils.FormatObject(state, isArray, request.Parameters["path"], request.Parameters["doc-type"])
		resp := observ_utils.CloseSpanWithCOAResponse(span, v1alpha2.COAResponse{
			State:       v1alpha2.OK,
			Body:        jData,
			ContentType: "application/json",
		})
		if request.Parameters["doc-type"] == "yaml" {
			resp.ContentType = "application/text"
		}
		return resp
	case fasthttp.MethodPost:
		ctx, span := observability.StartSpan("onCampaigns-POST", pCtx, nil)
		id := request.Parameters["__name"]

		var campaign model.CampaignSpec

		err := json.Unmarshal(request.Body, &campaign)
		if err != nil {
			cLog.Infof("V (Campaigns): onCampaigns failed - %s, traceId: %s", err.Error(), span.SpanContext().TraceID().String())
			return observ_utils.CloseSpanWithCOAResponse(span, v1alpha2.COAResponse{
				State: v1alpha2.InternalError,
				Body:  []byte(err.Error()),
			})
		}

		err = c.CampaignsManager.UpsertSpec(ctx, id, campaign)
		if err != nil {
			cLog.Infof("V (Campaigns): onCampaigns failed - %s, traceId: %s", err.Error(), span.SpanContext().TraceID().String())
			return observ_utils.CloseSpanWithCOAResponse(span, v1alpha2.COAResponse{
				State: v1alpha2.InternalError,
				Body:  []byte(err.Error()),
			})
		}
		return observ_utils.CloseSpanWithCOAResponse(span, v1alpha2.COAResponse{
			State: v1alpha2.OK,
		})
	case fasthttp.MethodDelete:
		ctx, span := observability.StartSpan("onCampaigns-DELETE", pCtx, nil)
		id := request.Parameters["__name"]
		err := c.CampaignsManager.DeleteSpec(ctx, id)
		if err != nil {
			cLog.Infof("V (Campaigns): onCampaigns failed - %s, traceId: %s", err.Error(), span.SpanContext().TraceID().String())
			return observ_utils.CloseSpanWithCOAResponse(span, v1alpha2.COAResponse{
				State: v1alpha2.InternalError,
				Body:  []byte(err.Error()),
			})
		}
		return observ_utils.CloseSpanWithCOAResponse(span, v1alpha2.COAResponse{
			State: v1alpha2.OK,
		})
	}
	cLog.Infof("V (Campaigns): onCampaigns failed - 405 method not allowed, traceId: %s", span.SpanContext().TraceID().String())
	resp := v1alpha2.COAResponse{
		State:       v1alpha2.MethodNotAllowed,
		Body:        []byte("{\"result\":\"405 - method not allowed\"}"),
		ContentType: "application/json",
	}
	observ_utils.UpdateSpanStatusFromCOAResponse(span, resp)
	return resp
}
