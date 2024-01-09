# dmarcanalyze

## Description
DMARC Analyze is a set of tools that enables you to download DMARC reports from an IMAP account and put that data into a SQL database.
Currently only SQLite, MySQL/MariaDB and Postgres are supported.

### Installing and usage
The binaries are fully self contained and have no other external dependencies. Just unzip the release somewhere convenient.
The tools are command line tools so you need a shell to run them.

The configuration is done via ```config.yml``` which must be placed in the same dir as the binary.
All of the items should be self-explanatory and otherwise I encourage you to experiment.

If you do not feel comfortable storing credentials in the config file, you can use environment variables to provide the data. All environment variables are mentioned in the comments of the config file.

### Example running
If you do not feel comfortale using credentials in a config file, you can use environment variables, like so:

```bash
DMARCANALYZE_IMAP_PASSWORD=bestestpassword DMARCANALYZE_IMAP_USERNAME="admin@mydomain.tld" DMARCANALYZE_IMAP_SERVER_ADDRESS=10.0.0.143 ./dmarcfetch
```
Be aware, however, that running it this way might find it's way into the history file of your shell.

If you do feel comfortable using credentials in a config file (make sure permissions are correct):
```bash
./dmarcfetch
```

## dmarcfetch
The DMARC Fetch tool connects to an IMAP server, changes into the directory and reads the messages. It then downloads and unpacks the attachments and stores the data into the SQL server.

There are various options that make it easy to run this tool autonomously (either as a cron-job or as a process that pauses after each run and then updates the database again).

This assumes you have a legacy authentication IMAP connection available. If this is not the case, then you can forward the reports to a different server that does have legacy authentication (i.e. without MFA) or you can enable legacy auth for just this single account.

## dmarcsqltoxls
The tool DMARC SQL to XLS reads a SQL database and generates a spreadsheet out of it. This is for people that can do data analysis with Excel better than with SQL.
Also, this helps if you want to make nice graphs.

You can have a look at the provided template.xlsx to see how to make the graph in advance in the template and have this tool fill in the data.

If you do not provide a template (or configure the filename to be empty) then the data is saved in an empty spreadsheet.


