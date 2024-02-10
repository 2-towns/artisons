# Display the tags
GET {{host}}/admin/filters
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//table/tbody/tr" count > 0

# Access the tags with huge page shoes empty message
GET {{host}}/admin/filters?page=100
HX-Request: true
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
body contains "No results found." 

