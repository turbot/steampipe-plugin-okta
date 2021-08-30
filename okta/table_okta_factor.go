package okta

import (
	"context"
	"time"
	// "fmt"
	// "strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
	// "github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableOktaFactor() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_factor",
		Description: "Okta Factor Sequencing enables passwordless MFA by requiring end users to successfully pass all specified MFA factors.",
		Get: &plugin.GetConfig{
			Hydrate:           getOktaFactor,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Not found"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listOktaUsers,
			Hydrate: listOktaFactors,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func:           listGroupMembers,
				MaxConcurrency: 10,
			},
		},
		Columns: []*plugin.Column{
			// Top Columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique key for Group."},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "Unique key for Group."},
			{Name: "factor_type", Type: proto.ColumnType_STRING, Description: "Description of the Group."},
			{Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when Group was created."},

			// Other Columns
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when Group's profile was last updated."},
			{Name: "provider", Type: proto.ColumnType_STRING, Description: "Determines how a Group's Profile and memberships are managed. Can be one of OKTA_GROUP, APP_GROUP or BUILT_IN."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Timestamp when Group's memberships were last updated."},

			// JSON Columns
			{Name: "embedded", Type: proto.ColumnType_JSON, Description: "The Group's Profile properties."},
			{Name: "links", Type: proto.ColumnType_JSON, Description: "Determines the Group's profile."},
			{Name: "verify", Type: proto.ColumnType_JSON, Description: "List of all users that are a member of this Group."},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Id"), Description: titleDescription},
		},
	}
}

// type UserFactorInfo = struct {
// 	okta.Factor
// 	UserId string
// }

type UserFactorInfo struct {
	Embedded    interface{}
	Links       interface{}
	Created     *time.Time
	FactorType  string
	Id          string
	LastUpdated *time.Time
	Provider    string
	Status      string
	Verify      interface{}
	UserId      string
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
	if h.Item != nil {
		userId = h.Item.(*okta.User).Id
	}

	if userId == "" {
		return nil, nil
	}

	factors, resp, err := client.UserFactor.ListFactors(ctx, userId)
	if err != nil {
		logger.Error("listOktaFactors", "list factors", err)
		return nil, err
	}

	for _, factor := range factors {
		// var data interface{}
		// data = factor
		// d.StreamListItem(ctx, UserFactorInfo{
		// 	Embedded: data["Embedded"],
		// })
		// fmt.Print(factor, "arnab \n")
		logger.Trace("listOktaFactors list", "list factor paging", factor)
		// d.StreamListItem(ctx, factor)
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
			d.StreamListItem(ctx, factor)
			// d.StreamListItem(ctx, UserFactorInfo{*factor, userId})
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getOktaFactor(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("getOktaGroup")
	groupId := d.KeyColumnQuals["id"].GetStringValue()
	

	if groupId == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("getOktaGroup", "connect", err)
		return nil, err
	}

	group, _, err := client.Group.GetGroup(ctx, groupId)
	if err != nil {
		logger.Error("getOktaGroup", "get group", err)
		return nil, err
	}

	return group, nil
}
