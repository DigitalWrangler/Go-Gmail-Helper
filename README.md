Gmail Helper - Mark Unread Emails as Read

📌 Overview

This Go application connects to your Gmail account using the Gmail API, retrieves all unread emails, prints their details (subject and received time), and marks them as read.

🚀 Setup Instructions

1️⃣ Get credentials.json from Google Cloud

Before running the application, you need to enable the Gmail API and download your OAuth 2.0 credentials.

Go to the Google Cloud Console.

Create a new project (or select an existing one).

Enable the Gmail API for your project.

Go to APIs & Services > Credentials.

Click Create Credentials > OAuth Client ID.

Select "Desktop App" and create the credentials.

Download the credentials.json file and place it in the project folder.

2️⃣ Install Dependencies

Before running the project, install all required dependencies by running:

 go mod tidy

3️⃣ Run the Application

Start the program by executing:

 go run main.go

4️⃣ Authenticate and Get API Token

When you run the program for the first time:

It will provide a Google authentication link.

Open the link in your browser.

Grant the necessary permissions.

Copy the authorization code displayed.

Paste the code into the terminal.

Once authenticated, the program will save your API token in token.json, and it will start processing unread emails.

5️⃣ Let the Script Run

The program will fetch all unread emails, print their subjects and received time.

It will mark them as read.

Finally, it will print a summary of how many emails were processed and how many remain unread.

🎯 Example Output

📩 Total unread emails: 25
📧 Subject: Welcome to Gmail - Received: 2025-02-08 14:05:32
✅ Marked as read!
📧 Subject: Your order has shipped - Received: 2025-02-08 12:47:10
✅ Marked as read!
...
📊 Summary:
🔹 Emails marked as read: 25
📩 Remaining unread emails: 0

🎉 Congratulations!