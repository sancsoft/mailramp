# mailramp
A simple tool for warming up an email domain or sending IP address

## Why You Might Need MailRamp

Modern mail systems react to sudden changes in email volume by mail servers and domains. If a sudden change is detected in sending volume, spam filters will lock down and block delivery under the assumption that there is a problem.  Many orginizations refuse to alter their spam filtering configurations, even for legitimate reasons.

Consistently delivering email to popular email providers and specific enterprise domains builds email reputation, which improves deliverability when email delivery spikes are needed.

## What You Will Need

MailRamp sends messages using SMTP. You'll need an email server (the one you plan to send from in the future) and an account that can send emails. Next, you'll need a bundle of email addresses to receive messages. 

## Ramping

MailRamp sends the configured number of messages per hour by sending a batch of messages every minute.  It repeatedly cycles through the defined email addresses.  Use Ctrl-C to stop execution.

## Configuring MailRamp

### Configuration File

The **config.json** file contains the settings for the application. 

```
{
    "sender": "sender@email.com",
    "rate": 10,
    "server": {
        "host": "smtp.mailgun.org",
        "port": 587,
        "user": "smtpuser@domain.com",
        "password": "password"
    },
    "subject": "[MAILRAMP] Subject Text"
}
```

* sender = the sending email address
* subject = the subject line for the message
* rate = the number of messages to send per hour
* server = email server credentials
  * host = server host name or ip address
  * port = server port number 25/587/465
  * user = mail server account - often just the email address
  * password = credentials for the mail server account

### Message Body

The **body.txt** file is used for the contents of the email messsage.  Each message has the same content.

### Recipients

The **recipients.txt** file contains the list of email addresses for delivery.  Messages are sent sequentially according the list.  The program wraps around to the beginning of the list once it reaches the end.  

For general public distribution lists (ex: email newsletters), try to include at least 2 addresses for each major email provivder (gmail.com, outlook.com, yahoo.com, etc.) 

## Feedback

Feedback is welcome at support@sancsoft.com.  





