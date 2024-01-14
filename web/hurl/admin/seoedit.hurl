# Display otp form when the correct email is used
POST {{host}}/auth/otp.html
HX-Request: true
[FormParams]
email: hello@world.com
HTTP 200
[Asserts]
xpath "//form[@hx-post='/auth/login.html']" exists
[Captures]
glue: xpath "normalize-space(//input[@id='glue']/@value)"

# Redirect to dashboard when the otp is correct
POST {{host}}/auth/login.html
HX-Request: true
User-Agent: {{ua}}
[FormParams]
glue: {{glue}}
otp: 1
otp: 1
otp: 1 
otp: 1 
otp: 1
otp: 1 
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/index.html"

# Editing the seo without title shows an error 
POST {{host}}/admin/seo/home/edit.html
HX-Request: true
[FormParams]
title: 
url: /
description: Very nice home page
HTTP 200
[Asserts]
header "HX-Retarget" == "#title-error" 
header "HX-Reswap" == "innerHTML show:#title-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The title is invalid." 

# Editing the seo without description shows an error 
POST {{host}}/admin/seo/home/edit.html
HX-Request: true
[FormParams]
title: Welcome !
url: /
description: 
HTTP 200
[Asserts]
header "HX-Retarget" == "#description-error" 
header "HX-Reswap" == "innerHTML show:#description-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The description is invalid." 

# Editing the seo without url shows an error 
POST {{host}}/admin/seo/home/edit.html
HX-Request: true
[FormParams]
title: Welcome !
url: 
description: Very nice home page
HTTP 200
[Asserts]
header "HX-Retarget" == "#url-error" 
header "HX-Reswap" == "innerHTML show:#url-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "input_url_invalid" 

# Editing the seo with correct data works
POST {{host}}/admin/seo/home/edit.html
HX-Request: true
[FormParams]
title: Welcome !
url: /
description: Very nice home page
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/seo.html"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"
