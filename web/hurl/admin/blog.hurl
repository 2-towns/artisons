# Display the articles
GET {{host}}/admin/blog.html
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//input[@name='q'][@hx-get='/admin/blog.html?page=1']" exists
xpath "//table/tbody/tr" count > 0

# Access the articles with huge page shoes empty message
GET {{host}}/admin/blog.html?page=100
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

# Search articles by id 
GET {{host}}/admin/blog.html?q=1
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//a[@href='/admin/blog/1/edit.html']" exists
xpath "//a[@hx-post='/admin/blog/1/delete.html']" exists

# Search blog by title 
GET {{host}}/admin/blog.html?q=ail
HX-Request: true
HTTP 200
[Asserts]
body not contains "No results found." 
xpath "//a[@href='/admin/blog/1/edit.html']" exists
xpath "//a[@hx-post='/admin/blog/1/delete.html']" exists

# Search articles by non existing keyword shows empty message 
GET {{host}}/admin/blog.html?q=totoisverydangerous
HX-Request: true
HTTP 200
[Asserts]
body contains "No results found." 

