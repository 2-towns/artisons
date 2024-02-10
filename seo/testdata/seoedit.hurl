# Editing the seo without title shows an error 
POST {{host}}/admin/seo/home/edit
HX-Request: true
[Cookies]
wsid: 444444 
[FormParams]
title: 
url: /
description: Very nice home page
HTTP 200
[Asserts]
header "HX-Retarget" == "#title-error" 
header "HX-Reswap" == "innerHTML show:#title-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the seo without description shows an error 
POST {{host}}/admin/seo/home/edit
HX-Request: true
[Cookies]
wsid: 444444 
[FormParams]
title: Welcome !
url: /
description: 
HTTP 200
[Asserts]
header "HX-Retarget" == "#description-error" 
header "HX-Reswap" == "innerHTML show:#description-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the seo without url shows an error 
POST {{host}}/admin/seo/home/edit
HX-Request: true
[Cookies]
wsid: 444444 
[FormParams]
title: Welcome !
url: 
description: Very nice home page
HTTP 200
[Asserts]
header "HX-Retarget" == "#url-error" 
header "HX-Reswap" == "innerHTML show:#url-row:top" 
xpath "normalize-space(//div[@class='form-error']/text())" == "The data is invalid." 

# Editing the seo with correct data works
POST {{host}}/admin/seo/home/edit
HX-Request: true
[Cookies]
wsid: 444444 
[FormParams]
title: Welcome !
url: /
description: Very nice home page
HTTP 200
[Asserts]
header "HX-Redirect" == "/admin/seo"
cookie "flash" exists
cookie "flash[HttpOnly]" exists
cookie "flash[Max-Age]" != 0
cookie "flash[Path]" == "/"
