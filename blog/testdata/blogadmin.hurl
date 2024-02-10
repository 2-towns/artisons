# Display the articles
GET {{host}}/admin/blog
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//input[@name='q'][@hx-get='/adminblog?page=1']" exists
xpath "//table/tbody/tr" count > 0

# Access the articles with huge page shoes empty message
GET {{host}}/admin/blog?page=100
HX-Request: true
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
body contains "No results found." 

# Search articles by id 
GET {{host}}/admin/blog?q=1
HX-Request: true
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//a[@href='/admin/blog/1/edit']" exists
xpath "//a[@hx-post='/admin/blog/1/delete']" exists

# Search blog by title 
GET {{host}}/admin/blog?q=ail
HX-Request: true
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//a[@href='/admin/blog/1/edit']" exists
xpath "//a[@hx-post='/admin/blog/1/delete']" exists

# Search articles by non existing keyword shows empty message 
GET {{host}}/admin/blog?q=totoisverydangerous
HX-Request: true
[Cookies]
wsid: 444444
HTTP 200
[Asserts]
body contains "No results found." 

