# dmarcanalyze

## Description
DMARC Analyze is a set of tools that enables you to download DMARC reports from an IMAP account and put that data into a SQL database.
Currently only SQLite, MySQL/MariaDB and Postgres are supported.

## dmarcfetch
The DMARC Fetch tool connects to an IMAP server, changes into the directory and reads the messages. It then downloads and unpacks the attachments and stores the data into the SQL server.

There are various options that make it easy to run this tool autonomously (either as a cron-job or as a process that pauses after each run and then updates the database again).

This assumes you have a legacy authentication IMAP connection available. If this is not the case, then you can forward the reports to a different server that does have legacy authentication (i.e. without MFA)

## dmarcsqltoxls
The tool DMARC SQL to XLS reads a SQL database and generates a spreadsheet out of it. This is for people that can do data analysis with Excel better than with SQL.
Also, this helps if you want to make nice graphs.

You can have a look at the provided template.xlsx to see how to make the graph in advance in the template and have this tool fill in the data.

## Installation and configuration
You only need to copy the binaries and their config.yml in a convenient spot and modify the config file.

If you do not feel comfortable storing credentials in the config file, you can use environment variables to provide the data. All environment variables are mentioned in the comments of the config file.