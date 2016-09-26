package main

import "html/template"

var html = template.Must(template.New("ES_Puller").Parse(`
<html> 
<head> 
    <title>"Welcome to ES Puller"</title>
    <link rel="stylesheet" type="text/css" href="http://meyerweb.com/eric/tools/css/reset/reset.css"/>
    <script src="http://ajax.googleapis.com/ajax/libs/jquery/1.7/jquery.js"></script> 
    <script src="http://malsup.github.com/jquery.form.js"></script> 
    </head>
    <body>
    <h1>Welcome to ES Puller</h1>
    <div id="messages"></div>
    <form id="myForm" action="/search" method="post"> 
    Date Start: <input id="date_gte" name="date_gte" value=""></input> 
    Date End: <input id="date_lte" name="date_lte"></input> 
    Session Number: <input id="session_number" name="session_number" value=""></input> 
    <input type="submit" value="Submit" /> 
    </form>
</body>
</html>
`))
