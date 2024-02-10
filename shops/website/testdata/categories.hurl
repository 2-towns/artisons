# The categories page is displayed
GET {{host}}/categories
HTTP 200
[Asserts]
xpath "//li[@class='category']" exists


