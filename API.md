# API Summary

```
Version:     1.0.0
Title:       Collabs Reporting Service 
Description: Collabs Reporting Service 
Host:        api.host.com
BasePath:    /
Consumes:    [application/json]
Produces:    [application/json]
```

<details>
<summary>/v1/reports/summary/{campaignID}: get</summary>

`path parameters`
- name: `campaignID`, type: `string`

`query parameters`
- name: `t`, type: `string`


`responses`
- code: `200`, type: `service.SummaryMetrics`
	- name: `brand_id`, type: `string`
	- name: `campaign_id`, type: `integer`
	- name: `date_range`, type: `object`
		- name: `start_date`, type: `string`
		- name: `end_date`, type: `string`
	- name: `gross_sale`, type: `object`
		- name: `sale_count`, type: `integer`
		- name: `total`, type: `number`
		- name: `channel_type`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/v1/reports/campaigns/{campaignID}: get</summary>

`path parameters`
- name: `campaignID`, type: `string`

`query parameters`
- name: `channelFilter`, type: `string`


`responses`
- code: `200`, type: `[]json.RawMessage`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

