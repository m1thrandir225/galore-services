package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

type TemplateData struct {
	Message string
}

func generateTemplateWithMessage(message string) string {
	htmlTemplate := `
<!--
* This email was built using Tabular.
* For more information, visit https://tabular.email
-->
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office" lang="en">
<head>
<title></title>
<meta charset="UTF-8" />
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
<!--[if !mso]>-->
<meta http-equiv="X-UA-Compatible" content="IE=edge" />
<!--<![endif]-->
<meta name="x-apple-disable-message-reformatting" content="" />
<meta content="target-densitydpi=device-dpi" name="viewport" />
<meta content="true" name="HandheldFriendly" />
<meta content="width=device-width" name="viewport" />
<meta name="format-detection" content="telephone=no, date=no, address=no, email=no, url=no" />
<style type="text/css">
table {
border-collapse: separate;
table-layout: fixed;
mso-table-lspace: 0pt;
mso-table-rspace: 0pt
}
table td {
border-collapse: collapse
}
.ExternalClass {
width: 100%
}
.ExternalClass,
.ExternalClass p,
.ExternalClass span,
.ExternalClass font,
.ExternalClass td,
.ExternalClass div {
line-height: 100%
}
body, a, li, p, h1, h2, h3 {
-ms-text-size-adjust: 100%;
-webkit-text-size-adjust: 100%;
}
html {
-webkit-text-size-adjust: none !important
}
body, #innerTable {
-webkit-font-smoothing: antialiased;
-moz-osx-font-smoothing: grayscale
}
#innerTable img+div {
display: none;
display: none !important
}
img {
Margin: 0;
padding: 0;
-ms-interpolation-mode: bicubic
}
h1, h2, h3, p, a {
line-height: inherit;
overflow-wrap: normal;
white-space: normal;
word-break: break-word
}
a {
text-decoration: none
}
h1, h2, h3, p {
min-width: 100%!important;
width: 100%!important;
max-width: 100%!important;
display: inline-block!important;
border: 0;
padding: 0;
margin: 0
}
a[x-apple-data-detectors] {
color: inherit !important;
text-decoration: none !important;
font-size: inherit !important;
font-family: inherit !important;
font-weight: inherit !important;
line-height: inherit !important
}
u + #body a {
color: inherit;
text-decoration: none;
font-size: inherit;
font-family: inherit;
font-weight: inherit;
line-height: inherit;
}
a[href^="mailto"],
a[href^="tel"],
a[href^="sms"] {
color: inherit;
text-decoration: none
}
</style>
<style type="text/css">
@media (min-width: 481px) {
.hd { display: none!important }
}
</style>
<style type="text/css">
@media (max-width: 480px) {
.hm { display: none!important }
}
</style>
<style type="text/css">
@media (max-width: 480px) {
.t15{mso-line-height-alt:0px!important;line-height:0!important;display:none!important}.t16{padding-left:30px!important;padding-bottom:40px!important;padding-right:30px!important}.t18,.t28{width:480px!important}.t6{padding-bottom:20px!important}.t13,.t24,.t8{width:420px!important}.t5{line-height:28px!important;font-size:26px!important;letter-spacing:-1.04px!important}.t26{padding:40px 30px!important}.t1{padding-bottom:50px!important}.t3{width:80px!important}
}
</style>
<!--[if !mso]>-->
<link href="https://fonts.googleapis.com/css2?family=Albert+Sans:wght@500;700;800&amp;display=swap" rel="stylesheet" type="text/css" />
<!--<![endif]-->
<!--[if mso]>
<xml>
<o:OfficeDocumentSettings>
<o:AllowPNG/>
<o:PixelsPerInch>96</o:PixelsPerInch>
</o:OfficeDocumentSettings>
</xml>
<![endif]-->
</head>
<body id="body" class="t32" style="min-width:100%;Margin:0px;padding:0px;background-color:#242424;"><div class="t31" style="background-color:#242424;"><table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" align="center"><tr><td class="t30" style="font-size:0;line-height:0;mso-line-height-rule:exactly;background-color:#242424;" valign="top" align="center">
<!--[if mso]>
<v:background xmlns:v="urn:schemas-microsoft-com:vml" fill="true" stroke="false">
<v:fill color="#242424"/>
</v:background>
<![endif]-->
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" align="center" id="innerTable"><tr><td><div class="t15" style="mso-line-height-rule:exactly;mso-line-height-alt:45px;line-height:45px;font-size:1px;display:block;">&nbsp;&nbsp;</div></td></tr><tr><td align="center">
<table class="t19" role="presentation" cellpadding="0" cellspacing="0" style="Margin-left:auto;Margin-right:auto;">
<tr>
<!--[if mso]>
<td width="600" class="t18" style="background-color:#F8F8F8;width:600px;">
<![endif]-->
<!--[if !mso]>-->
<td class="t18" style="background-color:#F8F8F8;width:600px;">
<!--<![endif]-->
<table class="t17" role="presentation" cellpadding="0" cellspacing="0" width="100%" style="width:100%;"><tr>
<td class="t16" style="padding:0 50px 60px 50px;"><table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="width:100% !important;"><tr><td align="left">
<table class="t4" role="presentation" cellpadding="0" cellspacing="0" style="Margin-right:auto;">
<tr>
<!--[if mso]>
<td width="130" class="t3" style="width:130px;">
<![endif]-->
<!--[if !mso]>-->
<td class="t3" style="width:130px;">
<!--<![endif]-->
<table class="t2" role="presentation" cellpadding="0" cellspacing="0" width="100%" style="width:100%;"><tr>
<td class="t1" style="padding:0 0 60px 0;"><div style="font-size:0px;"><img class="t0" style="display:block;border:0;height:auto;width:100%;Margin:0;max-width:100%;" width="130" height="130" alt="" src="https://dc0e5ccb-df52-44bb-9615-1eb7b77289db.b-cdn.net/e/5b5ecc22-9cd2-4ffa-86f8-f08253f67a3d/38b44de3-fb30-4723-a067-56e73b3d567b.png"/></div></td>
</tr></table>
</td>
</tr></table>
</td></tr><tr><td align="center">
<table class="t9" role="presentation" cellpadding="0" cellspacing="0" style="Margin-left:auto;Margin-right:auto;">
<tr>
<!--[if mso]>
<td width="500" class="t8" style="width:500px;">
<![endif]-->
<!--[if !mso]>-->
<td class="t8" style="width:500px;">
<!--<![endif]-->
<table class="t7" role="presentation" cellpadding="0" cellspacing="0" width="100%" style="width:100%;"><tr>
<td class="t6" style="padding:0 0 25px 0;"><h1 class="t5" style="margin:0;Margin:0;font-family:Albert Sans,BlinkMacSystemFont,Segoe UI,Helvetica Neue,Arial,sans-serif;line-height:41px;font-weight:800;font-style:normal;font-size:39px;text-decoration:none;text-transform:none;letter-spacing:-1.56px;direction:ltr;color:#191919;text-align:left;mso-line-height-rule:exactly;mso-text-raise:1px;">Password Reset Request</h1></td>
</tr></table>
</td>
</tr></table>
</td></tr><tr><td align="center">
<table class="t14" role="presentation" cellpadding="0" cellspacing="0" style="Margin-left:auto;Margin-right:auto;">
<tr>
<!--[if mso]>
<td width="500" class="t13" style="width:500px;">
<![endif]-->
<!--[if !mso]>-->
<td class="t13" style="width:500px;">
<!--<![endif]-->
<table class="t12" role="presentation" cellpadding="0" cellspacing="0" width="100%" style="width:100%;"><tr>
<td class="t11" style="padding:0 0 22px 0;"><p class="t10" style="margin:0;Margin:0;font-family:Albert Sans,BlinkMacSystemFont,Segoe UI,Helvetica Neue,Arial,sans-serif;line-height:22px;font-weight:500;font-style:normal;font-size:14px;text-decoration:none;text-transform:none;letter-spacing:-0.56px;direction:ltr;color:#333333;text-align:left;mso-line-height-rule:exactly;mso-text-raise:2px;">{{.Message}}</p></td>
</tr></table>
</td>
</tr></table>
</td></tr></table></td>
</tr></table>
</td>
</tr></table>
</td></tr><tr><td align="center">
<table class="t29" role="presentation" cellpadding="0" cellspacing="0" style="Margin-left:auto;Margin-right:auto;">
<tr>
<!--[if mso]>
<td width="600" class="t28" style="background-color:#242424;width:600px;">
<![endif]-->
<!--[if !mso]>-->
<td class="t28" style="background-color:#242424;width:600px;">
<!--<![endif]-->
<table class="t27" role="presentation" cellpadding="0" cellspacing="0" width="100%" style="width:100%;"><tr>
<td class="t26" style="padding:48px 50px 48px 50px;"><table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="width:100% !important;"><tr><td align="center">
<table class="t25" role="presentation" cellpadding="0" cellspacing="0" style="Margin-left:auto;Margin-right:auto;">
<tr>
<!--[if mso]>
<td width="500" class="t24" style="width:500px;">
<![endif]-->
<!--[if !mso]>-->
<td class="t24" style="width:500px;">
<!--<![endif]-->
<table class="t23" role="presentation" cellpadding="0" cellspacing="0" width="100%" style="width:100%;"><tr>
<td class="t22"><p class="t21" style="margin:0;Margin:0;font-family:Albert Sans,BlinkMacSystemFont,Segoe UI,Helvetica Neue,Arial,sans-serif;line-height:22px;font-weight:500;font-style:normal;font-size:12px;text-decoration:none;text-transform:none;direction:ltr;color:#888888;text-align:center;mso-line-height-rule:exactly;mso-text-raise:3px;"><a class="t20" href="https://tabular.email" style="margin:0;Margin:0;font-weight:700;font-style:normal;text-decoration:none;direction:ltr;color:#888888;mso-line-height-rule:exactly;" target="_blank">© All Rights Registered. Galore 2024</a></p></td>
</tr></table>
</td>
</tr></table>
</td></tr></table></td>
</tr></table>
</td>
</tr></table>
</td></tr></table></td></tr></table></div><div class="gmail-fix" style="display: none; white-space: nowrap; font: 15px courier; line-height: 0;">&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;</div></body>
</html>
`

	tmpl, err := template.New("generic-template").Parse(htmlTemplate)

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Create a PasswordResetData instance with the OTP value
	data := TemplateData{Message: message}

	// Create a buffer to store the rendered template
	var renderedHTML bytes.Buffer

	// Execute the template with the data
	err = tmpl.Execute(&renderedHTML, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	// Return the rendered HTML as a string
	return renderedHTML.String()
}

func GeneratePasswordResetSuccessfullyMail() string {
	message := fmt.Sprintf("You have successfully reset your password.")
	return generateTemplateWithMessage(message)
}

func GeneratePasswordOTPMail(otp string) string {
	message := fmt.Sprintf("You requested a password reset. Here is your OTP Code: %s", otp)

	return generateTemplateWithMessage(message)
}

func GenerateWelcomeMail(email string) string {
	message := fmt.Sprintf("Hi, %s. Thank you for registering to the Galore app. Hope you enjoy our services.", email)
	return generateTemplateWithMessage(message)
}
