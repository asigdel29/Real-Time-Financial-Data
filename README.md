# Go-Polling-Worker-and-Loggly
-By Anubhav Sigdel

The purpose of this assignment is to design and build a Go program that will periodically collect formatted data via API requests, display this data on the console, and report its results (or errors) to Loggly.
Go is a programming language developed by Google whose features make it an ideal for writing back-end server code. In this assignment, you will make use of the net/http library to make your API requests. Loggly is an online logging service with extensive tagging and search capabilities.
The data collected by an API request should be stored in a Go struct type that is appropriate for your data source. It should also be presented on the console such that the keys/values are clear.The worker should send messages to Loggly indicating success or failure of a given request and the amount of data collected. It should also make use of the built in tagging feature

-Extraced from James Early's CSC482 A06b
