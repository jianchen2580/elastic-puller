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
		<div id="index-select">
		<select id="index">
		<option value="prism" selected="selected">logstash-us</option>
		<option value="vcs">us-vcs</option>
		</select>
		</div>
    <div id="messages"></div>
    <form id="prism" action="/search" method="get"> 
		<input type="hidden" name="index" value="prism">
    Date Start: <input id="date_gte" name="date_gte" value=""></input> 
    Date End: <input id="date_lte" name="date_lte"></input> 
    Account ID: <input id="account_id" name="account_id" value=""></input> 
    App ID: <input id="app_id" name="app_id" value=""></input> 
    Session ID: <input id="session_id" name="session_id" value=""></input> 
    <input type="submit" value="Submit" /> 
    </form>

    <form id="vcs" action="/search" method="get"> 
		<input type="hidden" name="index" value="vcs">
    Date Start: <input id="date_gte" name="date_gte" value=""></input> 
    Date End: <input id="date_lte" name="date_lte"></input> 
    Account ID: <input id="account_id" name="account_id" value=""></input> 
    App ID: <input id="app_id" name="app_id" value=""></input> 
    Session ID: <input id="session_id" name="session_id" value=""></input> 
    <input type="submit" value="Submit" /> 
    </form>

    </body>
		<script>
		$(document).ready(function(){
			$("#vcs").hide()
		  $("#index").change(function() {
				$("#prism").hide()
			  $("#vcs").show()
			});
		});
		</script>
</html>
`))
