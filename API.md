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
<summary>get   : /contacts/{month}-{day}-{year}</summary>
`path parameters`
- name: month, type: string
- name: day, type: string
- name: year, type: string

`query parameters`

`body parameter`

`responses`
- code: 200, type: models.ContactResponse
	- name: output, type: string
</details>

<details>
<summary>get   : /things/{category}</summary>
`path parameters`
- name: category, type: string

`query parameters`
- name: q, type: string

`body parameter`

`responses`
- code: 200, type: []web.ThingResponse
	- name: output, type: string
</details>

<details>
<summary>delete: /things/{id}</summary>
`path parameters`
- name: id, type: integer

`query parameters`

`body parameter`

`responses`
</details>

<details>
<summary>post  : /echo</summary>
`path parameters`

`query parameters`

`body parameter`
- name: body, type: web.EchoRequest
	- name: input, type: string

`responses`
- code: 200, type: web.EchoResponse
	- name: output, type: string
</details>

<details>
<summary>post  : /contacts</summary>
`path parameters`

`query parameters`

`body parameter`
- name: body, type: models.ContactRequest
	- name: input, type: string
	- name: firstName, type: string
	- name: lastName, type: string
	- name: address, type: object
		- name: street, type: string
		- name: city, type: string
		- name: state, type: string
		- name: zipCode, type: string

`responses`
- code: 200, type: models.ContactResponse
	- name: output, type: string
</details>

<details>
<summary>get   : /contacts/{id}</summary>
`path parameters`
- name: id, type: integer

`query parameters`

`body parameter`

`responses`
- code: 200, type: models.ContactResponse
	- name: output, type: string
</details>

<details>
<summary>post  : /contacts/{id}</summary>
`path parameters`
- name: id, type: integer

`query parameters`

`body parameter`
- name: body, type: models.ContactRequest
	- name: address, type: object
		- name: street, type: string
		- name: city, type: string
		- name: state, type: string
		- name: zipCode, type: string
	- name: input, type: string
	- name: firstName, type: string
	- name: lastName, type: string

`responses`
- code: 200, type: models.ContactResponse
	- name: output, type: string
</details>

<details>
<summary>put   : /contacts/{id}</summary>
`path parameters`
- name: id, type: integer

`query parameters`

`body parameter`
- name: body, type: models.ContactRequest
	- name: input, type: string
	- name: firstName, type: string
	- name: lastName, type: string
	- name: address, type: object
		- name: street, type: string
		- name: city, type: string
		- name: state, type: string
		- name: zipCode, type: string

`responses`
- code: 200, type: models.ContactResponse
	- name: output, type: string
</details>

