# API Summary

```
Version:     1.0.0
Title:       OpenAPI Version 2 Specification
Description: OpenAPI Version 2 Specification
Host:        api.host.com
BasePath:    /api
Consumes:    [application/json]
Produces:    [application/json]
```

<details>
<summary>/contacts: post</summary>
Create a new contact request entity.

`body parameter`
- name: `body`, type: `models.ContactRequest`
	- name: `address`, type: `object`
		- name: `city`, type: `string`
		- name: `state`, type: `string`
		- name: `street`, type: `string`
		- name: `zipCode`, type: `string`
	- name: `firstName`, type: `string`
	- name: `id`, type: `string`
	- name: `lastName`, type: `string`

`responses`
- code: `200`, type: `models.ContactResponse`
	- name: `id`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/contacts/{id}: get</summary>
GetOne contact by id

`path parameters`
- name: `id`, type: `integer`


`responses`
- code: `200`, type: `models.ContactResponse`
	- name: `id`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/contacts/{id}: post</summary>
Update a contact entity with provided data.

`path parameters`
- name: `id`, type: `integer`

`body parameter`
- name: `body`, type: `models.ContactRequest`
	- name: `address`, type: `object`
		- name: `city`, type: `string`
		- name: `state`, type: `string`
		- name: `street`, type: `string`
		- name: `zipCode`, type: `string`
	- name: `firstName`, type: `string`
	- name: `id`, type: `string`
	- name: `lastName`, type: `string`

`responses`
- code: `200`, type: `models.ContactResponse`
	- name: `id`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/contacts/{id}: put</summary>
Replace a contact entity completely.

`path parameters`
- name: `id`, type: `integer`

`body parameter`
- name: `body`, type: `models.ContactRequest`
	- name: `address`, type: `object`
		- name: `city`, type: `string`
		- name: `state`, type: `string`
		- name: `street`, type: `string`
		- name: `zipCode`, type: `string`
	- name: `firstName`, type: `string`
	- name: `id`, type: `string`
	- name: `lastName`, type: `string`

`responses`
- code: `200`, type: `models.ContactResponse`
	- name: `id`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/contacts/{month}-{day}-{year}: get</summary>
Get contacts list by date

`path parameters`
- name: `month`, type: `string`
- name: `day`, type: `string`
- name: `year`, type: `string`


`responses`
- code: `200`, type: `models.ContactResponse`
	- name: `id`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/echo: post</summary>
Echo returns body with 'i's replaced with 'o's

`body parameter`
- name: `body`, type: `web.EchoRequest`
	- name: `input`, type: `string`

`responses`
- code: `200`, type: `web.EchoResponse`
	- name: `output`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/things/{category}: get</summary>
Get things by category and search query

`path parameters`
- name: `category`, type: `string`

`query parameters`
- name: `q`, type: `string`


`responses`
- code: `200`, type: `[]web.ThingResponse`
	- name: `output`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/things/{id}: delete</summary>
Delete thing by id

`path parameters`
- name: `id`, type: `integer`


`responses`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

