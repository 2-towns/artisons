# The blog list returns article
GET {{host}}/blog.html
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//div[@class='blog-card']" exists

# The blog list show nothing when no artiche match the query
GET {{host}}/blog.html?q=xyz
HX-Request: true
User-Agent: {{ua}}
HTTP 200
[Asserts]
body contains "No results found."

