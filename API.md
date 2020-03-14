# API Summary

```
Version:     1.0.0
Title:       Example OpenAPI Version 2 Specification
Description: Example OpenAPI Version 2 Specification
Host:        api.example.com
BasePath:    /api
Consumes:    [application/json]
Produces:    [application/json]
```

<details>
<summary>/contacts: post</summary>


```
Create a new contact request entity.
```

`body parameter`
- name: `body`, type: `models.ContactRequest`
	- name: `address`, type: `object`
		- name: `city`, type: `string`
		- name: `state`, type: `string`
		- name: `street`, type: `string`
		- name: `zipCode`, type: `string`
	- name: `eMail`, type: `string`, format: `email`
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
<summary>/contacts: delete</summary>


```
DeleteBulk contact by id
```

`query parameters`
- name: `ids`, type: `string`


`responses`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/contacts/{id}: get</summary>


```
GetOne contact by id
```

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


```
Update a contact entity with provided data.
```

`path parameters`
- name: `id`, type: `integer`

`body parameter`
- name: `body`, type: `models.ContactRequest`
	- name: `address`, type: `object`
		- name: `city`, type: `string`
		- name: `state`, type: `string`
		- name: `street`, type: `string`
		- name: `zipCode`, type: `string`
	- name: `eMail`, type: `string`, format: `email`
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


```
Replace a contact entity completely.
```

`path parameters`
- name: `id`, type: `integer`

`body parameter`
- name: `body`, type: `models.ContactRequest`
	- name: `address`, type: `object`
		- name: `city`, type: `string`
		- name: `state`, type: `string`
		- name: `street`, type: `string`
		- name: `zipCode`, type: `string`
	- name: `eMail`, type: `string`, format: `email`
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


```
Get contacts list by date
```

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
<summary>/doodads: post</summary>


```
Create a new doodad entity.
```

`body parameter`
- name: `body`, type: `models.ThingRequest`
	- name: `name`, type: `string`

`responses`
- code: `200`, type: `models.ThingResponse`
	- name: `bool`, type: `boolean`
	- name: `createDate`, type: `string`, format: `date-time,2006-01-02`
	- name: `float32`, type: `number`, format: `float`
	- name: `float64`, type: `number`, format: `double`
	- name: `int`, type: `integer`
	- name: `int16`, type: `integer`, format: `int16`
	- name: `int32`, type: `integer`, format: `int32`
	- name: `int64`, type: `integer`, format: `int64`
	- name: `int8`, type: `integer`, format: `int8`
	- name: `name`, type: `string`
	- name: `uint`, type: `integer`
	- name: `uint16`, type: `integer`, format: `int16`
	- name: `uint32`, type: `integer`, format: `int32`
	- name: `uint64`, type: `integer`, format: `int64`
	- name: `uint8`, type: `integer`, format: `int8`
	- name: `updateDate`, type: `string`, format: `date-time`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/echo: post</summary>


```
Echo returns body with 'i's replaced with 'o's
```

`body parameter`
- name: `body`, type: `EchoRequest`
	- name: `input`, type: `string`

`responses`
- code: `200`, type: `EchoResponse`
	- name: `output`, type: `string`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/things: get</summary>


```
Get things by date range

@from `format:"date-time,2006-01-02"` date should be in Go time format
@to   `format:"date-time,2006-01-02"` date should be in Go time format
```

`query parameters`
- name: `from`, type: `string`, format: `date-time,2006-01-02`
- name: `to`, type: `string`, format: `date-time,2006-01-02`


`responses`
- code: `200`, type: `ThingListResponse`
	- name: `things`, type: `[]array`
		- name: `bool`, type: `boolean`
		- name: `createDate`, type: `string`, format: `date-time,2006-01-02`
		- name: `float32`, type: `number`, format: `float`
		- name: `float64`, type: `number`, format: `double`
		- name: `int`, type: `integer`
		- name: `int16`, type: `integer`, format: `int16`
		- name: `int32`, type: `integer`, format: `int32`
		- name: `int64`, type: `integer`, format: `int64`
		- name: `int8`, type: `integer`, format: `int8`
		- name: `name`, type: `string`
		- name: `uint`, type: `integer`
		- name: `uint16`, type: `integer`, format: `int16`
		- name: `uint32`, type: `integer`, format: `int32`
		- name: `uint64`, type: `integer`, format: `int64`
		- name: `uint8`, type: `integer`, format: `int8`
		- name: `updateDate`, type: `string`, format: `date-time`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/things: post</summary>


```
Create thing
```

`body parameter`
- name: `body`, type: `models.ThingRequest`
	- name: `name`, type: `string`

`responses`
- code: `200`, type: `models.ThingResponse`
	- name: `bool`, type: `boolean`
	- name: `createDate`, type: `string`, format: `date-time,2006-01-02`
	- name: `float32`, type: `number`, format: `float`
	- name: `float64`, type: `number`, format: `double`
	- name: `int`, type: `integer`
	- name: `int16`, type: `integer`, format: `int16`
	- name: `int32`, type: `integer`, format: `int32`
	- name: `int64`, type: `integer`, format: `int64`
	- name: `int8`, type: `integer`, format: `int8`
	- name: `name`, type: `string`
	- name: `uint`, type: `integer`
	- name: `uint16`, type: `integer`, format: `int16`
	- name: `uint32`, type: `integer`, format: `int32`
	- name: `uint64`, type: `integer`, format: `int64`
	- name: `uint8`, type: `integer`, format: `int8`
	- name: `updateDate`, type: `string`, format: `date-time`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/things/{category}: get</summary>


```
Get things by category and search query
```

`path parameters`
- name: `category`, type: `string`

`query parameters`
- name: `q`, type: `string`


`responses`
- code: `200`, type: `ThingListResponse`
	- name: `things`, type: `[]array`
		- name: `bool`, type: `boolean`
		- name: `createDate`, type: `string`, format: `date-time,2006-01-02`
		- name: `float32`, type: `number`, format: `float`
		- name: `float64`, type: `number`, format: `double`
		- name: `int`, type: `integer`
		- name: `int16`, type: `integer`, format: `int16`
		- name: `int32`, type: `integer`, format: `int32`
		- name: `int64`, type: `integer`, format: `int64`
		- name: `int8`, type: `integer`, format: `int8`
		- name: `name`, type: `string`
		- name: `uint`, type: `integer`
		- name: `uint16`, type: `integer`, format: `int16`
		- name: `uint32`, type: `integer`, format: `int32`
		- name: `uint64`, type: `integer`, format: `int64`
		- name: `uint8`, type: `integer`, format: `int8`
		- name: `updateDate`, type: `string`, format: `date-time`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

<details>
<summary>/things/{id}: delete</summary>


```
Delete thing by id
```

`path parameters`
- name: `id`, type: `integer`


`responses`
- `default`, type: `Error`
	- name: `code`, type: `integer`
	- name: `status`, type: `string`
</details>

