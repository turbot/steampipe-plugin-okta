package okta

import (
	"context"
	"slices"
	"strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
	oktav4 "github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
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
			KeyColumns: []*plugin.KeyColumn{
				{Name: "user_id", Require: plugin.Optional},
			},
		},
		Columns: commonColumns([]*plugin.Column{
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
			{Name: "profile", Type: proto.ColumnType_JSON, Description: "Specific attributes related to the Factor.", Transform: transform.FromField("Factor.Profile")},
			{Name: "embedded", Type: proto.ColumnType_JSON, Description: "The Group's Profile properties.", Transform: transform.FromField("Factor.Embedded")},
			{Name: "verify", Type: proto.ColumnType_JSON, Description: "List of all users that are a member of this Group.", Transform: transform.FromField("Factor.Verify")},

			// Steampipe Columns
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Factor.Id"), Description: titleDescription},
		}),
	}
}

type UserFactorInfo struct {
	UserId   string
	UserName string
	Factor   OktaFactor
}

type OktaFactor struct {
	oktav4.UserFactor
	Profile interface{}
}

//// LIST FUNCTION

func listOktaFactors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := ConnectV4(ctx, d)
	if err != nil {
		logger.Error("okta_factor.listOktaFactors", "connect_error", err)
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

	// Minimize the API call with the given user id
	if d.EqualsQuals["user_id"] != nil {
		if d.EqualsQualString("user_id") != "" {
			if d.EqualsQualString("user_id") != "" && d.EqualsQualString("user_id") != userId {
				return nil, nil
			}
		} else if len(getListValues(d.EqualsQuals["user_id"].GetListValue())) > 0 {
			if !slices.Contains(types.StringValueSlice(getListValues(d.EqualsQuals["user_id"].GetListValue())), userId) {
				return nil, nil
			}
		}
	}

	if userId == "" {
		return nil, nil
	}

	factorReq := client.UserFactorAPI.ListFactors(ctx, userId)

	factors, resp, err := factorReq.Execute()
	if err != nil {
		logger.Error("okta_factor.listOktaFactors", "api_error", err)
		if strings.Contains(err.Error(), "Not found") {
			return nil, nil
		}
		return nil, err
	}

	for _, factor := range factors {
		if factor.GetActualInstance() != nil {
			factorDetails := getFactorDetails(factor.GetActualInstance())
			d.StreamListItem(ctx, UserFactorInfo{
				UserId:   userId,
				UserName: userName,
				Factor:   factorDetails,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	// paging
	for resp.HasNextPage() {
		var nextFactorSet []oktav4.ListFactors200ResponseInner
		resp, err = resp.Next(&nextFactorSet)
		if err != nil {
			logger.Error("okta_factor.listOktaFactors", "api_paging_error", err)
			return nil, err
		}

		for _, factor := range nextFactorSet {
			if factor.GetActualInstance() != nil {
				f := getFactorDetails(factor.GetActualInstance())
				d.StreamListItem(ctx, UserFactorInfo{
					UserId:   userId,
					UserName: userName,
					Factor:   f,
				})

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getOktaFactor(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	userId := d.EqualsQuals["user_id"].GetStringValue()
	factorId := d.EqualsQuals["id"].GetStringValue()

	if userId == "" || factorId == "" {
		return nil, nil
	}

	client, err := ConnectV4(ctx, d)
	if err != nil {
		logger.Error("okta_factor.getOktaFactor", "connection_error", err)
		return nil, err
	}

	userReq := client.UserAPI.GetUser(ctx, userId)
	user, _, err := userReq.Execute()
	if err != nil {
		logger.Error("okta_factor.getOktaFactor", "GetUser", err)
		if strings.Contains(err.Error(), "Not found") {
			return nil, nil
		}
		return nil, err
	}

	userProfile := *user.Profile
	userName := userProfile.Login

	factorReq := client.UserFactorAPI.GetFactor(ctx, userId, factorId)
	result, _, err := factorReq.Execute()
	if err != nil {
		logger.Error("okta_factor.getOktaFactor", "api_error", err)
		return nil, err
	}

	if result.GetActualInstance() == nil {
		return nil, nil
	}
	f := getFactorDetails(result.GetActualInstance())

	return &UserFactorInfo{UserId: userId, UserName: *userName, Factor: f}, nil
}

//// UTILITY FUNCTION

func getFactorDetails(i interface{}) OktaFactor {
	f := OktaFactor{}

	switch item := i.(type) {
	case *oktav4.UserFactorCall:
		f = OktaFactor{
			item.UserFactor,
			item.Profile,
		}
	case *oktav4.UserFactorCustomHOTP:
		f = OktaFactor{
			item.UserFactor,
			item.Profile,
		}
	case *oktav4.UserFactorEmail:
		f = OktaFactor{
			item.UserFactor,
			item.Profile,
		}
	case *oktav4.UserFactorHardware:
		f = OktaFactor{
			item.UserFactor,
			item.Profile,
		}
	case *oktav4.UserFactorPush:
		f = OktaFactor{
			item.UserFactor,
			item.Profile,
		}
	case *oktav4.UserFactorSMS:
		f = OktaFactor{
			item.UserFactor,
			item.Profile,
		}
	case *oktav4.UserFactorSecurityQuestion:
		f = OktaFactor{
			item.UserFactor,
			item.Profile,
		}
	case *oktav4.UserFactorTOTP:
		f = OktaFactor{
			item.UserFactor,
			item.Profile,
		}
	case *oktav4.UserFactorToken:
		f = OktaFactor{
			item.UserFactor,
			item.Profile,
		}
	case *oktav4.UserFactorU2F:
		f = OktaFactor{
			item.UserFactor,
			item.Profile,
		}
	case *oktav4.UserFactorWeb:
		f = OktaFactor{
			item.UserFactor,
			item.Profile,
		}
	case *oktav4.UserFactorWebAuthn:
		f = OktaFactor{
			item.UserFactor,
			item.Profile,
		}
	}
	return f
}
