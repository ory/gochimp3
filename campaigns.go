package gochimp3

import (
	"errors"
	"fmt"
)

const (
	campaigns_path       = "/campaigns"
	single_campaign_path = campaigns_path + "/%s"

	send_test_path = single_campaign_path + "/actions/test"
	send_path = single_campaign_path + "/actions/send"

	CAMPAIGN_TYPE_REGULAR = "regular"
	CAMPAIGN_TYPE_PLAINTEXT = "plaintext"
	CAMPAIGN_TYPE_ABSPLIT = "absplit" // deprecated by mailchimp
	CAMPAIGN_TYPE_RSS = "rss"
	CAMPAIGN_TYPE_VARIATE = "variate"

	CAMPAIGN_SEND_TYPE_HTML = "html"
	CAMPAIGN_SEND_TYPE_PLAINTEXT = "plaintext"

	CONDITION_MATCH_ANY = "any"
	CONDITION_MATCH_ALL = "all"

	CONDITION_TYPE_INTERESTS = "Interests"

	CONDITION_OP_CONTAINS = "interestcontains"
)


type CampaignQueryParams struct {
	ExtendedQueryParams

	Type string
	Status string
	BeforeSendTime string
	SinceSendTime string
	BeforeCreateTime string
	SinceCreateTime string
	ListId string
	FolderId string
	SortField string
	SortDir string
}

func (q CampaignQueryParams) Params() map[string]string {
	m := q.ExtendedQueryParams.Params()
	m["type"] = q.Type
	m["status"] = q.Status
	m["before_send_time"] = q.BeforeSendTime
	m["since_send_time"] = q.SinceSendTime
	m["before_create_time"] = q.BeforeCreateTime
	m["since_create_time"] = q.SinceCreateTime
	m["list_id"] = q.ListId
	m["folder_id"] = q.FolderId
	m["sort_field"] = q.SortField
	m["sort_dir"] = q.SortDir
	return m
}

type ListOfCampaigns struct {
	baseList
	Campaigns []CampaignResponse `json:"campaigns"`
}

type CampaignCreationRecipients struct {
	ListId string `json:"list_id"`
	SegmentOptions *CampaignCreationSegmentOptions `json:"segment_opts,omitempty"`
}

type CampaignCreationSegmentOptions struct {
	SavedSegmentId int `json:"saved_segment_id"`
	Match      string               `json:"match"`		// one of CONDITION_MATCH_*

	// this accepts various payloads. See http://developer.mailchimp.com/documentation/mailchimp/reference/campaigns/#create-post_campaigns
	Conditions interface{} `json:"conditions"`
}

type InterestsCondition struct {
	ConditionType string `json:"condition_type"`
	Field string `json:"field"`
	Op    string `json:"op"`
	Value []string `json:"value"`
}

type CampaignCreationSettings struct {
	SubjectLine     string `json:"subject_line"`
	PreviewText     string `json:"preview_text"`
	Title           string `json:"title"`
	FromName        string `json:"from_name"`
	ReplyTo         string `json:"reply_to"`
	UseConversation bool `json:"use_conversation"`
	ToName          string `json:"to_name"`
	FolderId        string `json:"folder_id"`
	Authenticate    bool `json:"authenticate"`
	AutoFooter      bool `json:"auto_footer"`
	InlineCss       bool `json:"inline_css"`
	AutoTweet       bool `json:"auto_tweet"`
	FbComments      bool `json:"fb_comments"`
	TemplateId      uint `json:"template_id"`
}

type CampaignCreationRequest struct {
	Type       string           `json:"type"` // must be one of the CAMPAIGN_TYPE_* consts
	Recipients CampaignCreationRecipients `json:"recipients"`
	Settings   CampaignCreationSettings `json:"settings"`
	// variate_settings not implemented
	Tracking          CampaignTracking `json:"tracking"`
	// rss_opts not implemented
	// social_card not implemented
}

type CampaignResponseRecipients struct {
	ListId string `json:"list_id"`
	ListName string `json:"list_name"`
	SegmentText string `json:"segment_text"`
	RecipientCount int `json:"recipient_count"`
}

type CampaignResponseSettings struct {
	SubjectLine     string `json:"subject_line"`
	PreviewText     string `json:"preview_text"`
	Title           string `json:"title"`
	FromName        string `json:"from_name"`
	ReplyTo         string `json:"reply_to"`
	UseConversation bool `json:"use_conversation"`
	ToName          string `json:"to_name"`
	FolderId        string `json:"folder_id"`
	Authenticate    bool `json:"authenticate"`
	AutoFooter      bool `json:"auto_footer"`
	InlineCss       bool `json:"inline_css"`
	AutoTweet       bool `json:"auto_tweet"`
	FbComments      bool `json:"fb_comments"`
	Timewarp        bool `json:"timewarp"`
	TemplateId      uint `json:"template_id"`
	DragAndDrop     bool `json:"drag_and_drop"`
}

type CampaignTracking struct {
	Opens           bool `json:"opens"`
	HtmlClicks      bool `json:"html_clicks"`
	TextClicks      bool `json:"text_clicks"`
	GoalTracking    bool `json:"goal_tracking"`
	Ecomm360        bool `json:"ecomm360"`
	GoogleAnalytics string `json:"google_analytics"`
	Clicktale       string `json:"clicktale"`
}

type CampaignEcommerce struct {
	TotalOrders  int `json:"total_orders"`
	TotalSpent   int `json:"total_spent"`
	TotalRevenue int `json:"total_revenue"`
}

type CampaignReportSummary struct {
	Opens            int `json:"opens"`
	UniqueOpens      int `json:"unique_opens"`
	OpenRate         float32 `json:"open_rate"`
	Clicks           int `json:"clicks"`
	SubscriberClicks int `json:"subscriber_clicks"`
	ClickRate        float32 `json:"click_rate"`
	Ecommerce        CampaignEcommerce `json:"ecommerce"`
}

type CampaignDeliveryStatus struct {
	Enabled  		 bool `json:"enabled"`
}

type CampaignResponse struct {
	withLinks

	ID                string `json:"id"`
	WebID             uint   `json:"web_id"`
	Type              string `json:"type"`
	CreateTime        string `json:"create_time"`
	ArchiveUrl        string `json:"archive_url"`
	LongArchiveUrl    string `json:"long_archive_url"`
	Status            string `json:"status"`
	EmailsSent        uint   `json:"emails_sent"`
	SendTime          string `json:"send_time"`
	ContentType       string `json:"content_type"`
	NeedsBlockRefresh bool   `json:"needs_block_refresh"`
	Recipients        CampaignResponseRecipients `json:"recipients"`
	Settings          CampaignResponseSettings `json:"settings"`
	Tracking          CampaignTracking `json:"tracking"`
	ReportSummary     CampaignReportSummary `json:"report_summary"`
	DeliveryStatus    CampaignDeliveryStatus `json:"delivery_status"`

	api *API
}

func (campaign CampaignResponse) CanMakeRequest() error {
	if campaign.ID == "" {
		return errors.New("No ID provided on campaign")
	}

	return nil
}

func (api API) GetCampaigns(params *CampaignQueryParams) (*ListOfCampaigns, error) {
	response := new(ListOfCampaigns)

	err := api.Request("GET", campaigns_path, params, nil, response)
	if err != nil {
		return nil, err
	}

	for _, l := range response.Campaigns {
		l.api = &api
	}

	return response, nil
}

func (api API) GetCampaign(id string, params *BasicQueryParams) (*CampaignResponse, error) {
	endpoint := fmt.Sprintf(single_campaign_path, id)

	response := new(CampaignResponse)
	response.api = &api

	return response, api.Request("GET", endpoint, params, nil, response)
}

func (api API) CreateCampaign(body *CampaignCreationRequest) (*CampaignResponse, error) {
	response := new(CampaignResponse)
	response.api = &api
	return response, api.Request("POST", campaigns_path, nil, body, response)
}

func (api API) UpdateCampaign(id string, body *CampaignCreationRequest) (*CampaignResponse, error) {
	endpoint := fmt.Sprintf(single_campaign_path, id)

	response := new(CampaignResponse)
	response.api = &api

	return response, api.Request("PATCH", endpoint, nil, body, response)
}

func (api API) DeleteCampaign(id string) (bool, error) {
	endpoint := fmt.Sprintf(single_campaign_path, id)
	return api.RequestOk("DELETE", endpoint)
}

// ------------------------------------------------------------------------------------------------
// Actions
// ------------------------------------------------------------------------------------------------

type TestEmailRequest struct {
	TestEmails []string `json:"test_emails"`
	SendType	string	`json:"send_type"`	// one of the CAMPAIGN_SEND_TYPE_* constants
}

type SendCampaignRequest struct {
	CampaignId	string	`json:"campaign_id"`
}

func (api API) SendTestEmail(id string, body *TestEmailRequest) (bool, error) {
	endpoint := fmt.Sprintf(send_test_path, id)
	err := api.Request("POST", endpoint, nil, body, nil)

	if err != nil {
		return false, err
	}
	return true, nil
}

func (api API) SendCampaign(id string, body *SendCampaignRequest) (bool, error) {
	endpoint := fmt.Sprintf(send_path, id)
	err := api.Request("POST", endpoint, nil, body, nil)

	if err != nil {
		return false, err
	}
	return true, nil
}
