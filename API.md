# API Summary

```
Version:     1.0.0
Title:       Example OpenAPI Version 2 Specification
Description: Example OpenAPI Version 2 Specification
Host:        api.example.com
BasePath:    /api
Consumes:    [application/json]
Produces:    [application/json application/text]
```

<details>
<summary>/contacts: get</summary>


```
GetBunch of contacts by ids
```

`query parameters`
- ids: `array`


`responses`
- code: `200`, type: `services.ContactResponse`
	- id: `string`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/contacts: post</summary>


```
Create a new contact request entity.
```

`body parameter`
- body: `services.ContactRequest`
	- address: `object`
		- city: `string`
		- state: `string`
		- street: `string`
		- zipCode: `string`
	- eMail: `string`, format: `email`
	- firstName: `string`
	- id: `string`
	- lastName: `string`

`responses`
- code: `200`, type: `services.ContactResponse`
	- id: `string`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/contacts: delete</summary>


```
DeleteBulk contact by id
```

`query parameters`
- ids: `array`


`responses`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/contacts/{id}: get</summary>


```
GetOne contact by id
id contact id
```

`path parameters`
- id: `integer`


`responses`
- code: `200`, type: `services.ContactResponse`
	- id: `string`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/contacts/{id}: post</summary>


```
Update a contact entity with provided data.
```

`path parameters`
- id: `integer`

`body parameter`
- body: `services.ContactRequest`
	- address: `object`
		- city: `string`
		- state: `string`
		- street: `string`
		- zipCode: `string`
	- eMail: `string`, format: `email`
	- firstName: `string`
	- id: `string`
	- lastName: `string`

`responses`
- code: `200`, type: `services.ContactResponse`
	- id: `string`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/contacts/{id}: put</summary>


```
Replace a contact entity completely.
```

`path parameters`
- id: `integer`

`body parameter`
- body: `services.ContactRequest`
	- address: `object`
		- city: `string`
		- state: `string`
		- street: `string`
		- zipCode: `string`
	- eMail: `string`, format: `email`
	- firstName: `string`
	- id: `string`
	- lastName: `string`

`responses`
- code: `200`, type: `services.ContactResponse`
	- id: `string`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/contacts/{month}-{day}-{year}: get</summary>


```
Get contacts list by date
```

`path parameters`
- month: `string`
- day: `string`
- year: `string`


`responses`
- code: `200`, type: `services.ContactResponse`
	- id: `string`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/doodads: post</summary>


```
Create a new doodad entity.
```

`body parameter`
- body: `models.ThingRequest`
	- name: `string`

`responses`
- code: `200`, type: `models.ThingResponse`
	- bool: `boolean`
	- createDate: `string`, format: `date-time,2006-01-02`
	- float32: `number`, format: `float`
	- float64: `number`, format: `double`
	- int: `integer`
	- int16: `integer`, format: `int16`
	- int32: `integer`, format: `int32`
	- int64: `integer`, format: `int64`
	- int8: `integer`, format: `int8`
	- name: `string`
	- uint: `integer`
	- uint16: `integer`, format: `int16`
	- uint32: `integer`, format: `int32`
	- uint64: `integer`, format: `int64`
	- uint8: `integer`, format: `int8`
	- updateDate: `string`, format: `date-time`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/echo/{input}: get</summary>


```
Echo returns body with 'i's replaced with 'o's
```

`path parameters`
- input: `string`


`responses`
- code: `200`, type: `EchoResponse`
	- output: `string`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/ping: post</summary>


```
Ping returns body with 'i's replaced with 'o's
```

`body parameter`
- body: `services.PingRequest`
	- input: `string`

`responses`
- code: `200`, type: `services.PingResponse`
	- output: `string`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/things: get</summary>


```
Get things by date range
```

`query parameters`
- from: `string`, format: `date-time,2006-01-02`
- to: `string`, format: `date-time,2006-01-02`


`responses`
- code: `200`, type: `ThingListResponse`
	- things: `[]array`
		- bool: `boolean`
		- createDate: `string`, format: `date-time,2006-01-02`
		- float32: `number`, format: `float`
		- float64: `number`, format: `double`
		- int: `integer`
		- int16: `integer`, format: `int16`
		- int32: `integer`, format: `int32`
		- int64: `integer`, format: `int64`
		- int8: `integer`, format: `int8`
		- name: `string`
		- uint: `integer`
		- uint16: `integer`, format: `int16`
		- uint32: `integer`, format: `int32`
		- uint64: `integer`, format: `int64`
		- uint8: `integer`, format: `int8`
		- updateDate: `string`, format: `date-time`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/things: post</summary>


```
Create thing
```

`body parameter`
- body: `models.ThingRequest`
	- name: `string`

`responses`
- code: `200`, type: `models.ThingResponse`
	- bool: `boolean`
	- createDate: `string`, format: `date-time,2006-01-02`
	- float32: `number`, format: `float`
	- float64: `number`, format: `double`
	- int: `integer`
	- int16: `integer`, format: `int16`
	- int32: `integer`, format: `int32`
	- int64: `integer`, format: `int64`
	- int8: `integer`, format: `int8`
	- name: `string`
	- uint: `integer`
	- uint16: `integer`, format: `int16`
	- uint32: `integer`, format: `int32`
	- uint64: `integer`, format: `int64`
	- uint8: `integer`, format: `int8`
	- updateDate: `string`, format: `date-time`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/things/{category}: get</summary>


```
Get things by category and search query
```

`path parameters`
- category: `string`

`query parameters`
- q: `string`


`responses`
- code: `200`, type: `ThingListResponse`
	- things: `[]array`
		- bool: `boolean`
		- createDate: `string`, format: `date-time,2006-01-02`
		- float32: `number`, format: `float`
		- float64: `number`, format: `double`
		- int: `integer`
		- int16: `integer`, format: `int16`
		- int32: `integer`, format: `int32`
		- int64: `integer`, format: `int64`
		- int8: `integer`, format: `int8`
		- name: `string`
		- uint: `integer`
		- uint16: `integer`, format: `int16`
		- uint32: `integer`, format: `int32`
		- uint64: `integer`, format: `int64`
		- uint8: `integer`, format: `int8`
		- updateDate: `string`, format: `date-time`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

<details>
<summary>/things/{id}: delete</summary>


```
Delete thing by id
```

`path parameters`
- id: `integer`


`responses`
- `default`, type: `Error`
	- code: `integer`
	- status: `string`
</details>

