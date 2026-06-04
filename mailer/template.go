package mailer

import (
	"fmt"
	"time"
)

func MailTemplate(code string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Verification Code</title>

<style>
body{
    margin:0;
    padding:20px;
    background:#f5f7fb;
    font-family:Arial,Helvetica,sans-serif;
}

.container{
    max-width:600px;
    margin:auto;
    background:#ffffff;
    border-radius:16px;
    overflow:hidden;
    box-shadow:0 4px 20px rgba(0,0,0,.08);
}

.header{
    background:#0f172a;
    color:#fff;
    padding:32px;
    text-align:center;
}

.logo{
    font-size:28px;
    font-weight:bold;
}

.content{
    padding:40px 30px;
}

.title{
    font-size:26px;
    color:#111827;
    margin-bottom:20px;
}

.text{
    color:#4b5563;
    line-height:1.7;
    margin-bottom:24px;
}

.code-box{
    background:#eff6ff;
    border:2px solid #3b82f6;
    border-radius:12px;
    padding:24px;
    text-align:center;
    margin:30px 0;
}

.code{
    font-size:42px;
    font-weight:700;
    letter-spacing:10px;
    color:#2563eb;
}

.info{
    background:#f8fafc;
    border-left:4px solid #3b82f6;
    padding:16px;
    border-radius:8px;
    color:#475569;
    margin-top:20px;
}

.footer{
    padding:24px;
    text-align:center;
    font-size:13px;
    color:#94a3b8;
    border-top:1px solid #e5e7eb;
}
</style>
</head>

<body>

<div class="container">

    <div class="header">
        <div class="logo">📊 VisitTracker</div>
    </div>

    <div class="content">

        <h1 class="title">Verify Your Account</h1>

        <p class="text">
            Use the verification code below to complete your registration
            and access your VisitTracker dashboard.
        </p>

        <div class="code-box">
            <div class="code">%s</div>
        </div>

        <div class="info">
            ⏱ This code expires in <strong>15 minutes</strong>.<br>
            🔒 Never share this code with anyone.
        </div>

        <p class="text" style="margin-top:30px;">
            If you didn't request this verification code, you can safely
            ignore this email.
        </p>

    </div>

    <div class="footer">
        © %d VisitTracker<br>
        Real-Time Visitor Analytics & Monitoring
    </div>

</div>

</body>
</html>`, code, time.Now().Year())
}
