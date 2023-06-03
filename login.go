package main

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleLogin(c *gin.Context) {
	htmlTemplate := `
<!DOCTYPE html>
<html>

<head>
        <meta charset="utf-8" />
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <title>Go OAuth2 Example</title>
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <!-- <link rel="stylesheet" type="text/css" media="screen" href="main.css" /> -->
        <!-- <script src="main.js"></script> -->
</head>

<body>
        <a href="https://github.com/login/oauth/authorize?client_id=03178089f0ff2ea0356d&redirect_uri=http://localhost:8080/oauth/redirect">Login with github</a>
</body>

</html>	`

	tmpl := template.Must(template.New("login.html").Parse(htmlTemplate))
	err := tmpl.Execute(c.Writer, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}
