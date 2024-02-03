# Display the tags
GET {{host}}/admin/tags.html
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//table/tbody/tr" count > 0

# Access the tags with huge page shoes empty message
GET {{host}}/admin/tags.html?page=100
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

