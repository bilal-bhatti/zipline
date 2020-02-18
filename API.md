# API Summary

```
Version:     1.0.0
Title:       OpenAPI Version 2 Specification
Description: OpenAPI Version 2 Specification
Host:        api.example.com
BasePath:    /api
Consumes:    [application/json]
Produces:    [application/json]
```

<details>
<summary>/echo: post</summary>

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
<summary>/contacts: post</summary>

`body parameter`
- name: `body`, type: `models.ContactRequest`
	- name: `input`, type: `string`
	- name: `firstName`, type: `string`
	- name: `lastName`, type: `string`
	- name: `address`, type: `object`
		- name: `street`, type: `string`
		- name: `city`, type: `string`
		- name: `state`, type: `string`
		- name: `zipCode`, type: `string`

`responses`
- code: `200`, type: `models.ContactResponse`
	- name: `output`, type: `string`
- `default`, type: `Error`
	- name: `status`, type: `string`
	- name: `code`, type: `integer`
</details>

<details>
<summary>/contacts/{id}: get</summary>

`path parameters`
- name: `id`, type: `integer`


`responses`
- code: `200`, type: `models.ContactResponse`
	- name: `output`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/contacts/{id}: post</summary>

`path parameters`
- name: `id`, type: `integer`

`body parameter`
- name: `body`, type: `models.ContactRequest`
	- name: `lastName`, type: `string`
	- name: `address`, type: `object`
		- name: `street`, type: `string`
		- name: `city`, type: `string`
		- name: `state`, type: `string`
		- name: `zipCode`, type: `string`
	- name: `input`, type: `string`
	- name: `firstName`, type: `string`

`responses`
- code: `200`, type: `models.ContactResponse`
	- name: `output`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/contacts/{id}: put</summary>

`path parameters`
- name: `id`, type: `integer`

`body parameter`
- name: `body`, type: `models.ContactRequest`
	- name: `input`, type: `string`
	- name: `firstName`, type: `string`
	- name: `lastName`, type: `string`
	- name: `address`, type: `object`
		- name: `street`, type: `string`
		- name: `city`, type: `string`
		- name: `state`, type: `string`
		- name: `zipCode`, type: `string`

`responses`
- code: `200`, type: `models.ContactResponse`
	- name: `output`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/contacts/{month}-{day}-{year}: get</summary>

`path parameters`
- name: `month`, type: `string`
- name: `day`, type: `string`
- name: `year`, type: `string`


`responses`
- code: `200`, type: `models.ContactResponse`
	- name: `output`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/things/{category}: get</summary>

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

`path parameters`
- name: `id`, type: `integer`


`responses`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

