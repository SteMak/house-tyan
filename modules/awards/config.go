package awards

type config struct {
	BankirAPIToken   string `json:"bankir_api_token,omitempty"`
	GdHouseID        string `json:"gd_house_id,omitempty"`
	ChForRequestID   string `json:"ch_for_request_id,omitempty"`
	ChForLogsID      string `json:"ch_for_logs_id,omitempty"`
	ChForBumpSiupID  string `json:"ch_for_bump_siup_id,omitempty"`
	UsConfirmatorID  string `json:"us_confirmator_id,omitempty"`
	UsSiupID         string `json:"us_siup_id,omitempty"`
	UsBumpID         string `json:"us_bump_id,omitempty"`
	RoRequestMakerID string `json:"ro_request_maker_id,omitempty"`
}
