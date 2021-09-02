package okta

import (
	"context"
	"strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableOktaFactor() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_factor",
		Description: "Represents an Okta Factor.",
		Get: &plugin.GetConfig{
			Hydrate:           getOktaFactor,
			KeyColumns:        plugin.AllColumns([]string{"id", "user_id"}),
			ShouldIgnoreError: isNotFoundError([]string{"Not found", "Invalid Factor"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listOktaUsers,
			Hydrate:       listOktaFactors,
		},
		Columns: []*plugin.Column{
			// Top Columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique key for Group.", Transform: transform.FromField("Factor.Id")},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "Unique key for Group."},
			{Name: "user_name", Type: proto.ColumnType_STRING, Description: "Unique identifier for the user (username)."},
			{Name: "factor_type", Type: proto.ColumnType_STRING, Description: "Description of the Group.", Transform: transform.FromField("Factor.FactorType")},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when Group was created.", Transform: transform.FromField("Factor.Created")},

			// Other Columns
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the factor was last updated.", Transform: transform.FromField("Factor.LastUpdated")},
			{Name: "provider", Type: proto.ColumnType_STRING, Description: "The provider for the factor.", Transform: transform.FromField("Factor.Provider")},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The current status of the factor.", Transform: transform.FromField("Factor.Status")},

			// JSON Columns
			{Name: "embedded", Type: proto.ColumnType_JSON, Description: "The Group's Profile properties.", Transform: transform.FromField("Factor.Embedded")},
			{Name: "verify", Type: proto.ColumnType_JSON, Description: "List of all users that are a member of this Group.", Transform: transform.FromField("Factor.Verify")},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Factor.Id"), Description: titleDescription},
		},
	}
}

type UserFactorInfo struct {
	UserId   string
	UserName string
	Factor   okta.Factor
}

//// LIST FUNCTION

func listOktaFactors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("listOktaFactors", "connect", err)
		return nil, err
	}

	var userId string
	var userName string
	if h.Item != nil {
		userData := h.Item.(*okta.User)
		userId = userData.Id
		userProfile := *userData.Profile
		userName = userProfile["login"].(string)
	}

	if userId == "" {
		return nil, nil
	}

	factors, resp, err := client.UserFactor.ListFactors(ctx, userId)
	if err != nil {
		logger.Error("listOktaFactors", "Error", err)
		if strings.Contains(err.Error(), "Not found") {
			return nil, nil
		}
		return nil, err
	}

	for _, factor := range factors {
		d.StreamListItem(ctx, UserFactorInfo{
			UserId:   userId,
			UserName: userName,
			Factor:   factor,
		})
	}

	// paging
	for resp.HasNextPage() {
		var nextFactorSet []*okta.Factor
		resp, err = resp.Next(ctx, &nextFactorSet)
		if err != nil {
			logger.Error("listOktaFactors", "list factor paging", err)
			return nil, err
		}
		for _, factor := range nextFactorSet {
			d.StreamListItem(ctx, UserFactorInfo{
				UserId:   userId,
				UserName: userName,
				Factor:   *factor,
			})
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getOktaFactor(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("getOktaFactor")
	userId := d.KeyColumnQuals["user_id"].GetStringValue()
	factorId := d.KeyColumnQuals["id"].GetStringValue()

	if userId == "" || factorId == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getOktaFactor", "connect", err)
		return nil, err
	}

	user, _, err := client.User.GetUser(ctx, userId)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return nil, nil
		}
		return nil, err
	}

	userProfile := *user.Profile
	userName := userProfile["login"].(string)

	var factorInstance okta.Factor
	factor, _, err := client.UserFactor.GetFactor(ctx, userId, factorId, factorInstance)
	if err != nil {
		logger.Error("getOktaFactor", "get factor", err)
		return nil, err
	}

	return &UserFactorInfo{UserId: userId, UserName: userName, Factor: factor}, nil
}
