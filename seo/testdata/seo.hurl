# Display the seo
GET {{host}}/admin/seo
[Cookies]
wsid: 444444 
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//table/tbody/tr" count > 0

# Access the seo with huge page shoes empty message
GET {{host}}/admin/seo?page=100
HX-Request: true
[Cookies]
wsid: 444444 
HTTP 200
[Asserts]
body contains "No results found." 

