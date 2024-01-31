# The blog list returns article
GET {{host}}/blog.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//div[@class='blog-card']" exists

# The blog list show nothing when no artiche match the query
GET {{host}}/blog.html?q=xyz
User-Agent: {{ua}}
HTTP 200
[Asserts]
body contains "No results found."

# The blog article is displayed
GET {{host}}/blog/mangez-de-lail.html
User-Agent: {{ua}}
HTTP 200
[Asserts]
xpath "//p[text()=\"Mangez de l'ail\"]" exists

# The offline article returns 404
GET {{host}}/blog/non-active.html
User-Agent: {{ua}}
HTTP 404
[Asserts]
xpath "//h2[text()=\"Oops..You just found an error page.\"]" exists
