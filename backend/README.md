<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fwww.almrsal.com%2Fwp-content%2Fuploads%2F2015%2F12%2FEnron-Corporation-was-an-American-energy-commodities-and-services-company-based-in-Houston.jpg&f=1&nofb=1&ipt=7d291f71e280fc04c928387d0f0f199f056c6e7a2c4aabdd17289b045038898f&ipo=images" alt="Enron logo"></a>
</p>

<h3 align="center">Enron Corp - Backend</h3>

## API Documentation

### 1. GET /emails

Retrieves a paginated list of emails with optional filters and sorting.

#### Request

- **URL:** `/emails`
- **Method:** `GET`
- **Headers:**
  - `Content-Type: application/json`

#### Query Parameters

- `page` (integer, optional): The page number to retrieve. Default is `1`.
- `size` (integer, optional): The number of emails per page. Default is `10`.
- `filter` (string, optional): A filter string to apply to the email list. This can filter text included in`content`.
- `sort` (string, optional): The field by which to sort the emails. Possible values are `date`, `from`, `to`, `subject`. Default is `date`.
- `order` (string, optional): The order in which to sort the emails. Possible values are `asc` (ascending) and `desc` (descending). Default is `desc`.

#### Response

- **Status Code:** `200 OK`
- **Body:**

```json
{
  "page": 1,
  "size": 10,
  "total_elements": 100,
  "total_pages": 10,
  "emails": [
    {
      "id": 1,
      "date": "2024-08-01T12:00:00Z",
      "from": "example@example.com",
      "to": "recipient@example.com",
      "subject": "Subject Example",
      "content": "Email content goes here."
    },
    ...
  ]
}
```

### 2. GET /emails/{id}

Retrieves the details of a specific email by its ID.

#### Request

- **URL:** `/emails/{id}`
- **Method:** `GET`
- **Headers:**
  - `Content-Type: application/json`

#### Path Parameters

- `id` (integer, required): The ID of the email to retrieve.

#### Query Parameters

- `filter` (string, optional): A filter string to apply to the `content` to highlight text.

#### Response

- **Status Code:** `200 OK`
- **Body:**

```json
{
  "id": 1,
  "date": "2024-08-01T12:00:00Z",
  "from": "example@example.com",
  "to": "recipient@example.com",
  "subject": "Subject Example",
  "content": "Email content goes here.",
  "path": "path/in/directory",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Test

You can test the APIs running them in Postman. You will need to update the credentials, which can be found in the `docker-compose.yml` file.

[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://app.getpostman.com/run-collection/14923491-a3b09f92-a416-4efe-8464-960649c9376d?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D14923491-a3b09f92-a416-4efe-8464-960649c9376d%26entityType%3Dcollection%26workspaceId%3Dfc072b4c-ff0b-4789-b8cf-025441bf4132)
