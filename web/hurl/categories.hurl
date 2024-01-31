# The categories page is displayed
GET {{host}}/categories.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//li[@class='category']" exists


