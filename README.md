# Amazon SNS Subscriber

The application exposes an HTTP endpoint ready to process Amazon SNS messages. It implements subcription confirmation and signature verification.
The payloads received are stored in the plain txt file.

# Execution

The execution expects the name of the output file as a parameter:

`go run main.go ./out.txt`
