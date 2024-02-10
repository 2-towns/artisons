# Display the tags
GET {{host}}/admin/tags
[Cookies]
wsid: 444444 
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//table/tbody/tr" count > 0

# Access the tags with huge page shoes empty message
GET {{host}}/admin/tags?page=100
[Cookies]
wsid: 444444 
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

